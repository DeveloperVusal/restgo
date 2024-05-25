package service

import (
	"context"
	"net/http"

	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/repository"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/utils/auth/generate"
	"apibgo/internal/utils/response"
	"apibgo/pkg/auth/pswd"

	"github.com/jackc/pgx/v5"
)

type Users interface {
	GetUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error)
	GetUsers(ctx context.Context) (*response.Response, error)
	CreateUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error)
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
		if dto.ConfirmStatus == "" {
			dto.ConfirmStatus = domainUser.ConfirmStatus_UNKNOWN
		}

		args := []interface{}{dto.Email, pwd_hash, dto.Name, dto.Surname, "", "", dto.ConfirmStatus, tokenSecret}
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
				Message: "Data is got",
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
