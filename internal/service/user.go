package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	domainAuth "apibgo/internal/domain/auth"
	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/repository"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/utils/auth/generate"
	"apibgo/internal/utils/response"
	"apibgo/pkg/auth/ajwt"
	"apibgo/pkg/auth/device"
	"apibgo/pkg/auth/pswd"

	"github.com/jackc/pgx/v5"
)

type Users interface {
	GetUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error)
	GetUsers(ctx context.Context) (*response.Response, error)
	CreateUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error)
	UpdateUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error)
	DeleteUser(ctx context.Context, user_id int) (*response.Response, error)
	Sessions(ctx context.Context, header_auth []string) (*response.Response, error)
}

type UserService struct {
	db *pgsql.Storage
}

func NewUserService(store *pgsql.Storage) *UserService {
	return &UserService{
		db: store,
	}

}

func (ur *UserService) GetUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ur.db)
	user, err := repoUser.GetUser(ctx, domainUser.UserDto{Id: dto.Id})

	if err != nil {
		return nil, err
	}

	if user.Id > 0 {
		return &response.Response{
			Code:    response.ErrorEmpty,
			Status:  response.StatusSuccess,
			Message: "data is got",
			Result: map[string]interface{}{
				"email":      user.Email,
				"name":       user.Name.String,
				"surname":    user.Surname.String,
				"activation": user.Activation,
				"status":     user.ConfirmStatus,
			},
		}, nil
	}

	return nil, nil
}

func (ur *UserService) GetUsers(ctx context.Context) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ur.db)
	users, err := repoUser.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	count, err := repoUser.GetCountUsers(ctx)

	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		var respUsers []map[string]interface{}

		for _, user := range users {
			respUsers = append(respUsers, map[string]interface{}{
				"id":         user.Id,
				"email":      user.Email,
				"name":       user.Name.String,
				"surname":    user.Surname.String,
				"activation": user.Activation,
				"status":     user.ConfirmStatus,
			})
		}

		return &response.Response{
			Code:    response.ErrorEmpty,
			Status:  response.StatusSuccess,
			Message: "data is got",
			Result: map[string]interface{}{
				"count": count,
				"data":  respUsers,
			},
		}, nil
	}

	return nil, nil
}

func (ur *UserService) CreateUser(ctx context.Context, dto domainUser.CreateUserDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ur.db)
	user, err := repoUser.GetUser(ctx, domainUser.UserDto{Email: dto.Email})

	if err != nil {
		return nil, err
	}

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
		tokenSecret, err2 := generate.RandomStringBytes(32)

		if err != nil {
			return nil, err
		}

		if err2 != nil {
			return nil, err2
		}

		// start transaction
		tx, err := ur.db.Db.BeginTx(ctx, pgx.TxOptions{})

		if err != nil {
			return nil, err
		}

		defer func() {
			if err != nil {
				tx.Rollback(ctx)
			}
		}()

		// Inserting in users
		args := []interface{}{dto.Email, pwd_hash, dto.Name, dto.Surname, "", "", domainUser.ConfirmStatusEnum(dto.ConfirmStatus), tokenSecret}
		user, err := repoUser.InsertUser(ctx, args)

		if err != nil {
			tx.Rollback(ctx)

			return nil, err
		}

		// If successfully, then we return the Response
		if user.Id > 0 {
			tx.Commit(ctx)

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "user created successfully",
				Result: map[string]interface{}{
					"id":         user.Id,
					"email":      user.Email,
					"name":       user.Name.String,
					"surname":    user.Surname.String,
					"activation": user.Activation,
					"status":     user.ConfirmStatus,
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

func (ur *UserService) UpdateUser(ctx context.Context, dto domainUser.UpdateUserDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ur.db)

	if dto.Email != "" {
		user, err := repoUser.GetUser(ctx, domainUser.UserDto{Email: dto.Email})

		if err != nil {
			return nil, err
		}

		if user.Id > 0 && user.Id != uint(dto.Id) {
			return &response.Response{
				Code:    response.ErrorAccountExists,
				Status:  response.StatusError,
				Message: "user with this email address already exists",
			}, nil
		}
	}

	user, err := repoUser.GetUser(ctx, domainUser.UserDto{Id: dto.Id})

	if err != nil {
		return nil, err
	}

	if user.Id > 0 {
		var pwd_hash string

		if dto.Password != "" && dto.ConfirmPassword != "" {
			// If don't match passwords
			if dto.Password != dto.ConfirmPassword {
				return &response.Response{
					Code:    response.ErrorAccountConfirmPassword,
					Status:  response.StatusError,
					Message: "don't match passwords",
				}, nil
			}

			// Generate password
			pwd_hash, err = pswd.HashPassword(dto.Password)

			if err != nil {
				return nil, err
			}
		}

		// start transaction
		tx, err := ur.db.Db.BeginTx(ctx, pgx.TxOptions{})

		if err != nil {
			return nil, err
		}

		defer func() {
			if err != nil {
				tx.Rollback(ctx)
			}
		}()

		modelUser := &domainUser.User{
			Activation: dto.Activation,
		}

		if dto.Email != "" {
			modelUser.Email = dto.Email
		}

		if len(pwd_hash) > 5 {
			modelUser.Password = pwd_hash
		}

		if dto.Name != "" {
			modelUser.Name = sql.NullString{String: dto.Name}
		}

		if dto.Surname != "" {
			modelUser.Surname = sql.NullString{String: dto.Surname}
		}

		if dto.ConfirmStatus != "" {
			modelUser.ConfirmStatus = domainUser.ConfirmStatusEnum(dto.ConfirmStatus)
		}

		updUser, cmdtag, err := repoUser.UpdateUser(ctx, int(user.Id), modelUser)

		if err != nil {
			tx.Rollback(ctx)

			return nil, err
		}

		if cmdtag.RowsAffected() <= 0 {
			tx.Rollback(ctx)

			return nil, err
		} else {
			tx.Commit(ctx)

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "user data updated successfully",
				Result: map[string]interface{}{
					// "id":         updUser.Id,
					"email":      updUser.Email,
					"name":       updUser.Name.String,
					"surname":    updUser.Surname.String,
					"activation": updUser.Activation,
					"status":     updUser.ConfirmStatus,
				},
			}, nil
		}
	}

	return nil, nil
}

func (ur *UserService) DeleteUser(ctx context.Context, user_id int) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ur.db)
	user, err := repoUser.GetUser(ctx, domainUser.UserDto{Id: user_id})

	if err != nil {
		return nil, err
	}

	if user.Id > 0 {
		// start transaction
		tx, err := ur.db.Db.BeginTx(ctx, pgx.TxOptions{})

		if err != nil {
			return nil, err
		}

		defer func() {
			if err != nil {
				tx.Rollback(ctx)
			}
		}()

		cmdtag, err := repoUser.DeleteUser(ctx, int(user.Id))

		if err != nil {
			tx.Rollback(ctx)

			return nil, err
		}

		if cmdtag.RowsAffected() <= 0 {
			tx.Rollback(ctx)

			return nil, err
		} else {
			tx.Commit(ctx)

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "user deleted successfully",
			}, nil
		}
	}

	return nil, nil
}

func (ur *UserService) Sessions(ctx context.Context, header_auth []string) (*response.Response, error) {
	// Parse header Authorization and get token
	split := strings.Split(header_auth[0], " ")
	token := split[1]

	// Checking on correct JWT
	if err := ajwt.IsJWT(token, os.Getenv("APP_JWT_SECRET")); err != nil {
		return nil, err
	}

	payload, err := ajwt.GetClaims(token, os.Getenv("APP_JWT_SECRET"))

	if err != nil {
		return nil, err
	}

	user_id, _ := strconv.Atoi(fmt.Sprintf("%v", payload["user_id"]))

	if user_id > 0 {
		repoAuth := repository.NewAuthRepo(ur.db)

		auths, err := repoAuth.GetSessions(context.Background(), domainAuth.SessionDto{
			Id: user_id,
		})

		if err != nil {
			return nil, err
		}

		count, err := repoAuth.GetCountSessions(context.Background(), domainAuth.SessionDto{
			Id: user_id,
		})

		if err != nil {
			return nil, err
		}

		if len(auths) > 0 {
			var respAuths []map[string]interface{}

			for _, auth := range auths {
				respAuths = append(respAuths, map[string]interface{}{
					"id": auth.Id,
					"ip": auth.Ip,
					"device": map[string]string{
						"name": auth.Device,
						"info": strings.Join([]string{device.DetectOS(auth.UserAgent), device.DetectBrowser(auth.UserAgent)}, ","),
					},
					"time": auth.CreatedAt.Format("02-01-2006 15:04:05"),
				})
			}

			return &response.Response{
				Code:    response.ErrorEmpty,
				Status:  response.StatusSuccess,
				Message: "data is got",
				Result: map[string]interface{}{
					"count": count,
					"data":  respAuths,
				},
			}, nil
		}
	}

	return nil, nil
}
