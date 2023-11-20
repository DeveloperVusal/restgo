package service

import (
	"context"
	"net/http"
	"os"
	"time"

	"apibgo/internal/domain"
	"apibgo/internal/repository"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/utils/response"
	"apibgo/pkg/auth/ajwt"
	"apibgo/pkg/auth/pswd"
)

type Auths interface {
	Login(ctx context.Context, dto domain.AuthDto) (*response.Response, error)
	Registration(ctx context.Context)
	Activation(ctx context.Context)
	Logout(ctx context.Context)
	Recover(ctx context.Context)
	VerifyToken(ctx context.Context)
	Refresh(ctx context.Context)
	RecoverPasswordCheckToken(ctx context.Context)
}

type AuthService struct {
	db *pgsql.Storage
}

func NewAuthService(store *pgsql.Storage) *AuthService {
	return &AuthService{
		db: store,
	}

}

func (ar *AuthService) Login(ctx context.Context, dto domain.AuthDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoAuth := repository.NewAuthRepo(ar.db)
	rows, err := repoAuth.GetUser(ctx, dto)

	if err != nil {
		return nil, err
	}

	var user_id int
	var user_password string
	var user_activation bool

	// Get above columns from row result
	err = rows.Scan(&user_id, &user_password, &user_activation)

	if err != nil {
		return nil, err
	}

	// If valid data
	if user_id > 0 && pswd.CheckPasswordHash(dto.Password, user_password) {
		// If don't activated the user
		if !user_activation {
			return &response.Response{
				Code:    response.ErrorAccountActivate,
				Status:  response.Error,
				Message: "account not activated",
			}, nil
		}

		// Checking exist already authentication a user
		rows, err := repoAuth.GetAuth(ctx, dto)

		if err != nil {
			return nil, err
		}

		var auth_id int

		// Get above columns from row result
		err = rows.Scan(&auth_id)

		if err != nil {
			return nil, err
		}

		// If exists, then we delete the record
		if auth_id > 0 {
			rows, err := repoAuth.DeleteAuth(ctx, auth_id)

			if err != nil {
				return nil, err
			}

			if rows.CommandTag().RowsAffected() <= 0 {
				return nil, err
			}
		}

		// Creating pair tokens of jwt
		myjwt := ajwt.JWT{
			Secret: os.Getenv("APP_JWT_SECRET"),
			UserId: user_id,
		}
		access, refresh := myjwt.NewPairTokens()

		// Inserting in sessions
		args := []interface{}{user_id, access, refresh, dto.Ip, dto.Device, dto.UserAgent}
		rows, err = repoAuth.InsertAuth(ctx, args)

		if err != nil {
			return nil, err
		}

		// If successfully, then we return the Response
		if rows.CommandTag().RowsAffected() > 0 {
			var _cookies []*http.Cookie

			_cookies = append(_cookies, &http.Cookie{
				Name:     "refresh_token",
				Value:    refresh,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().AddDate(0, 1, 0),
			})

			return &response.Response{
				Code:    0,
				Status:  response.Success,
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
