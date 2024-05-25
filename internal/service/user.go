package service

import (
	"context"

	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/repository"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/utils/response"
)

type Users interface {
	GetUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error)
	GetUsers(ctx context.Context) (*response.Response, error)
}

type UserService struct {
	db *pgsql.Storage
}

func NewUserService(store *pgsql.Storage) *UserService {
	return &UserService{
		db: store,
	}

}

func (ar *UserService) GetUser(ctx context.Context, dto domainUser.UserDto) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ar.db)
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

func (ar *UserService) GetUsers(ctx context.Context) (*response.Response, error) {
	// Trying find a user in the users table
	repoUser := repository.NewUserRepo(ar.db)
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
