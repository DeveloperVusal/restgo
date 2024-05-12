package repository

import (
	"context"
	"errors"
	"fmt"

	domainAuth "apibgo/internal/domain/auth"
	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/storage/pgsql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRI interface {
	GetUser(ctx context.Context, dto domainAuth.UserDto) (domainUser.User, error)
	InsertUser(ctx context.Context, args []interface{}) (pgconn.CommandTag, error)
	UpdateUser(ctx context.Context, id int, args []interface{}) (pgconn.CommandTag, error)
	DeleteUser(ctx context.Context, id int) (pgconn.CommandTag, error)
}

type UserRepo struct {
	db    *pgx.Conn
	store *pgsql.Storage
}

func NewUserRepo(store *pgsql.Storage) *UserRepo {
	return &UserRepo{
		db:    store.Db,
		store: store,
	}
}

func (ar *UserRepo) GetUser(ctx context.Context, dto domainAuth.UserDto) (domainUser.User, error) {
	var user domainUser.User

	args := []interface{}{}
	cond := ""

	if dto.Email == "" {
		cond = `id = $1`
		args = append(args, dto.Id)
	} else {
		cond = `email = $1`
		args = append(args, dto.Email)
	}

	sql := `SELECT * FROM ` + user.TableName() + ` WHERE ` + cond + ` LIMIT 1`

	err := ar.db.QueryRow(ctx, sql, args...).Scan(
		&user.Id, &user.Email, &user.Password, &user.Activation,
		&user.Name, &user.Surname, &user.TokenSecretKey,
		&user.UpdatedAt, &user.CreatedAt, &user.ConfirmCode,
		&user.ConfirmedAt, &user.ConfirmStatus,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainUser.User{}, nil
		}

		return domainUser.User{}, err
	}

	return user, nil
}

func (ar *UserRepo) DeleteUser(ctx context.Context, id int) (pgconn.CommandTag, error) {
	sql := `DELETE FROM users WHERE id = $1`

	return ar.db.Exec(ctx, sql, id)
}

func (ar *UserRepo) InsertUser(ctx context.Context, args []interface{}) (domainUser.User, error) {
	sql := `INSERT INTO users (email, password, name, surname, confirm_code, confirm_status, token_secret_key, confirmed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()::timestamp, NOW()::timestamp) RETURNING id;`
	id := 0

	err := ar.db.QueryRow(ctx, sql, args...).Scan(&id)

	if err != nil {
		return domainUser.User{}, err
	}

	user, err := ar.GetUser(ctx, domainAuth.UserDto{Id: id})

	if err != nil {
		fmt.Println("insert 2 <--")
		return domainUser.User{}, err
	}

	return user, nil
}
