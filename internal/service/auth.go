package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	domainAuth "apibgo/internal/domain/auth"
	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/repository"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/templates/mails"
	"apibgo/internal/utils/auth/generate"
	"apibgo/internal/utils/response"
	"apibgo/pkg/auth/ajwt"
	"apibgo/pkg/auth/device"
	"apibgo/pkg/auth/pswd"
	"apibgo/pkg/mail"

	"github.com/jackc/pgx/v5"
)

type Auths interface {
	Login(ctx context.Context, dto domainAuth.LoginDto) (*response.Response, error)
	Registration(ctx context.Context, dto domainAuth.RegistrationDto) (*response.Response, error)
	Activation(ctx context.Context, dto domainAuth.ActivationDto) (*response.Response, error)
	Logout(ctx context.Context, header_auth []string) (*response.Response, error)
	// Recover(ctx context.Context)
	VerifyToken(ctx context.Context, header_auth []string) (bool, error)
	Refresh(ctx context.Context, tokenCookie *http.Cookie, dto domainAuth.LoginDto)
	// RecoverPasswordCheckToken(ctx context.Context)
}

type AuthService struct {
	db *pgsql.Storage
}

func NewAuthService(store *pgsql.Storage) *AuthService {
	return &AuthService{
		db: store,
	}

}

func (ar *AuthService) Login(ctx context.Context, dto domainAuth.LoginDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ar.db)
	repoAuth := repository.NewAuthRepo(ar.db)
	dto.Device = strings.ToLower(device.DetectDevice(dto.UserAgent))
	user, err := repoUser.GetUser(ctx, domainAuth.UserDto{Email: dto.Email})

	if err != nil {
		return nil, err
	}

	// If valid data
	if user.Id > 0 && pswd.CheckPasswordHash(dto.Password, user.Password) {
		// If don't activated the user
		if !user.Activation {
			return &response.Response{
				Code:    response.ErrorAccountActivate,
				Status:  response.StatusError,
				Message: "account not activated",
			}, nil
		}

		// Checking exist already authentication a user
		auth, err := repoAuth.GetAuth(ctx, domainAuth.RefreshDto{
			Device:    dto.Device,
			Ip:        dto.Ip,
			UserAgent: dto.UserAgent,
		})

		if err != nil {
			return nil, err
		}

		// If exists, then we delete the record
		if auth.Id > 0 {
			cmdtag, err := repoAuth.DeleteAuth(ctx, domainAuth.DestroyDto{Id: int(auth.Id)})

			if err != nil {
				return nil, err
			}

			if cmdtag.RowsAffected() <= 0 {
				return nil, err
			}
		}

		// Creating pair tokens of jwt
		myjwt := ajwt.JWT{
			Secret: os.Getenv("APP_JWT_SECRET"),
			UserId: int(user.Id),
		}
		access, refresh := myjwt.NewPairTokens()

		// Inserting in sessions
		args := []interface{}{user.Id, access, refresh, dto.Ip, dto.Device, dto.UserAgent}
		cmdtag, err := repoAuth.InsertAuth(ctx, args)

		if err != nil {
			return nil, err
		}

		// If successfully, then we return the Response
		if cmdtag.RowsAffected() > 0 {
			var _cookies []*http.Cookie

			_cookies = append(_cookies, &http.Cookie{
				Name:     "refresh_token",
				Value:    refresh,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().AddDate(0, 1, 0),
			})

			// Prepare message for send to mailbox
			// Get template message
			subject, text := mails.Login(map[string]string{
				"email":         user.Email,
				"device":        device.DetectDevice(dto.UserAgent),
				"device_detail": strings.Join([]string{device.DetectOS(dto.UserAgent), device.DetectBrowser(dto.UserAgent), dto.Ip}, ","),
				"time":          time.Now().Format("02 Jan, 15:04"),
			})

			// TODO: recommendation use RabbitMQ
			go func() {
				smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				m := mail.Mailer{
					SmtpHost:     os.Getenv("SMTP_HOST"),
					SmtpPort:     smtpPort,
					SmtpUser:     os.Getenv("SMTP_USER"),
					SmtpPassword: os.Getenv("SMTP_PASSWORD"),
				}
				// Sends message to emails address
				m.SendMail([]string{user.Email}, subject, text)
			}()

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "Data is got",
				Result: map[string]interface{}{
					"access_token":  access,
					"refresh_token": refresh,
				},
				Cookies: _cookies,
			}, nil
		}
	}

	return nil, nil
}

func (ar *AuthService) Registration(ctx context.Context, dto domainAuth.RegistrationDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ar.db)
	user, err := repoUser.GetUser(ctx, domainAuth.UserDto{Email: dto.Email})

	if err != nil {
		return nil, err
	}

	// If valid data
	if user.Id <= 0 {
		// If don't match passwords
		if dto.Password != dto.ConfirmPassword {
			return &response.Response{
				Code:    response.ErrorAccountConfirmPassword,
				Status:  response.StatusError,
				Message: "don't match passwords",
			}, nil
		}

		// Generate codes and strings
		pwd_hash, err := pswd.HashPassword(dto.Password)
		confirmCode := generate.RandomNumbers(6)
		tokenSecret, err2 := generate.RandomStringBytes(32)

		if err != nil {
			return nil, err
		}

		if err2 != nil {
			return nil, err2
		}

		// start transaction
		tx, err := ar.db.Db.BeginTx(ctx, pgx.TxOptions{})

		if err != nil {
			return nil, err
		}

		defer func() {
			if err != nil {
				tx.Rollback(ctx)
			}
		}()

		// Inserting in users
		args := []interface{}{dto.Email, pwd_hash, dto.Name, dto.Surname, confirmCode, domainUser.ConfirmStatus_WAIT, tokenSecret}
		user, err := repoUser.InsertUser(ctx, args)

		if err != nil {
			tx.Rollback(ctx)

			return nil, err
		}

		// If successfully, then we return the Response
		if user.Id > 0 {
			// Prepare message for send to mailbox
			// Get template message
			subject, text := mails.Registration(map[string]string{
				"confirmCode":  confirmCode,
				"confirmed_at": user.ConfirmedAt.Time.Format("02-01-2006 15:04:05"),
			})

			// TODO: recommendation use RabbitMQ
			go func() {
				smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				m := mail.Mailer{
					SmtpHost:     os.Getenv("SMTP_HOST"),
					SmtpPort:     smtpPort,
					SmtpUser:     os.Getenv("SMTP_USER"),
					SmtpPassword: os.Getenv("SMTP_PASSWORD"),
				}
				// Sends message to emails address
				m.SendMail([]string{user.Email}, subject, text)
			}()

			// generate key for activation
			key := ar.generateToken(user.TokenSecretKey, user.Email)

			tx.Commit(ctx)

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "Data is got",
				Result: map[string]interface{}{
					"email": user.Email,
					"key":   key,
				},
				HttpCode: http.StatusCreated,
			}, nil
		} else {
			tx.Rollback(ctx)

			return &response.Response{
				Code:    response.ErrorAccountNotCreated,
				Status:  response.StatusError,
				Message: "failed, the user was not created",
			}, nil
		}
	}

	return &response.Response{
		Code:    response.ErrorAccountExists,
		Status:  response.StatusError,
		Message: "user with this email address already exists",
	}, nil
}

func (ar *AuthService) Logout(ctx context.Context, header_auth []string) (*response.Response, error) {
	// Parse header Authorization and get token
	split := strings.Split(header_auth[0], " ")
	token := split[1]

	// Checking on correct JWT
	if err := ajwt.IsJWT(token, os.Getenv("APP_JWT_SECRET")); err != nil {
		return nil, err
	}

	// Deleting session
	repoAuth := repository.NewAuthRepo(ar.db)
	cmdtag, err := repoAuth.DeleteAuth(ctx, domainAuth.DestroyDto{Token: token})

	if err != nil {
		return nil, err
	}

	// If whole successfully, then resets refresh token
	if cmdtag.RowsAffected() <= 0 {
		return nil, err
	} else {
		var _cookies []*http.Cookie

		_cookies = append(_cookies, &http.Cookie{
			Name:     "refresh_token",
			Value:    "empty",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		return &response.Response{
			Code:    response.ErrorEmpty,
			Status:  response.StatusSuccess,
			Message: "Session successfully destroyed",
			Result:  nil,
			Cookies: _cookies,
		}, nil
	}
}

func (ar *AuthService) Refresh(ctx context.Context, tokenCookie *http.Cookie, dto domainAuth.LoginDto) (*response.Response, error) {
	// Prepare data
	dto.Device = strings.ToLower(device.DetectDevice(dto.UserAgent))

	// Checking on correct JWT
	if err := ajwt.IsJWT(tokenCookie.Value, os.Getenv("APP_JWT_SECRET")); err != nil {
		return nil, err
	}

	// Get session
	repoAuth := repository.NewAuthRepo(ar.db)
	auth, err := repoAuth.GetAuth(ctx, domainAuth.RefreshDto{
		Refresh: tokenCookie.Value,
	})

	if err != nil {
		return nil, err
	}

	tx, err := ar.db.Db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Deleting session
	cmdtag, err := repoAuth.DeleteAuth(ctx, domainAuth.DestroyDto{Id: int(auth.Id)})

	// If whole successfully, then resets refresh token
	if cmdtag.RowsAffected() <= 0 {
		tx.Rollback(ctx)

		return nil, err
	} else {
		// Creating pair tokens of jwt
		myjwt := ajwt.JWT{
			Secret: os.Getenv("APP_JWT_SECRET"),
			UserId: int(auth.UserId),
		}
		access, refresh := myjwt.NewPairTokens()

		// Inserting in sessions
		args := []interface{}{auth.UserId, access, refresh, dto.Ip, dto.Device, dto.UserAgent}
		cmdtag, err := repoAuth.InsertAuth(ctx, args)

		if err != nil {
			tx.Rollback(ctx)

			return nil, err
		}

		// If successfully, then we return the Response
		if cmdtag.RowsAffected() > 0 {
			var _cookies []*http.Cookie

			_cookies = append(_cookies, &http.Cookie{
				Name:     "refresh_token",
				Value:    refresh,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().AddDate(0, 1, 0),
			})

			tx.Commit(ctx)

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "Data is got",
				Result: map[string]interface{}{
					"access_token":  access,
					"refresh_token": refresh,
				},
				Cookies: _cookies,
			}, nil
		} else {
			tx.Rollback(ctx)

			return nil, nil
		}
	}
}

func (ar *AuthService) VerifyToken(ctx context.Context, header_auth []string) (bool, error) {
	// Parse header Authorization and get token
	split := strings.Split(header_auth[0], " ")
	token := split[1]

	// Checking on correct JWT
	if err := ajwt.IsJWT(token, os.Getenv("APP_JWT_SECRET")); err != nil {
		return false, err
	}

	return true, nil
}

func (ar *AuthService) Activation(ctx context.Context, dto domainAuth.ActivationDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ar.db)
	user, err := repoUser.GetUser(ctx, domainAuth.UserDto{Email: dto.Email})

	if err != nil {
		return nil, err
	}

	// If valid data
	if user.Id > 0 && ar.checkToken(user.TokenSecretKey, dto.Email, dto.Key) {
		// If a user activated
		if user.Activation {
			return &response.Response{
				Code:    response.ErrorAccountAlreadyActivate,
				Status:  response.StatusError,
				Message: "account already activated",
			}, nil
		}

		// if confirmed_at valid
		if user.ConfirmedAt.Valid {
			confirmExpired := time.Now().Unix() - user.ConfirmedAt.Time.Unix()
			limitExprTime, _ := strconv.ParseInt(os.Getenv("APP_CONFIRM_TIME"), 10, 64)

			// if it hasn't been 5 minutes
			if confirmExpired < limitExprTime {
				// if don't match of confirm codes
				if strconv.Itoa(dto.Code) != user.ConfirmCode.String {
					return &response.Response{
						Code:    response.ErrorAccountInvalidCode,
						Status:  response.StatusError,
						Message: "don't match of confirm code",
					}, nil
				}

				// Begin transaction
				tx, err := ar.db.Db.BeginTx(ctx, pgx.TxOptions{})

				if err != nil {
					return nil, err
				}

				defer func() {
					if err != nil {
						tx.Rollback(ctx)
					}
				}()

				// Activating account
				cmdtag, err := repoUser.UpdateUser(ctx, int(user.Id), &domainUser.User{
					Activation: true,
				})

				if err != nil {
					tx.Rollback(ctx)

					return nil, err
				}

				if cmdtag.RowsAffected() <= 0 {
					tx.Rollback(ctx)

					return nil, err
				} else {
					tx.Commit(ctx)

					// Prepare message for send to mailbox
					// Get template message
					subject, text := mails.Activation(map[string]string{})

					// TODO: recommendation use RabbitMQ
					go func() {
						smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
						m := mail.Mailer{
							SmtpHost:     os.Getenv("SMTP_HOST"),
							SmtpPort:     smtpPort,
							SmtpUser:     os.Getenv("SMTP_USER"),
							SmtpPassword: os.Getenv("SMTP_PASSWORD"),
						}
						// Sends message to emails address
						m.SendMail([]string{user.Email}, subject, text)
					}()

					return &response.Response{
						Code:    response.ErrorEmpty,
						Status:  response.StatusSuccess,
						Message: "your account successfully activated",
					}, nil
				}
			} else {
				return &response.Response{
					Code:    response.ErrorAccountActivateTimeout,
					Status:  response.StatusError,
					Message: "this confirm code time out",
				}, nil
			}
		}
	}

	return nil, nil
}

func (ar *AuthService) generateToken(secret string, email string) string {
	h := sha256.New()
	h.Write([]byte(secret + `::` + email))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func (ar *AuthService) checkToken(secret string, email string, token string) bool {
	h := sha256.New()
	h.Write([]byte(secret + `::` + email))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs) == token
}
