package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"

	domainAuth "apibgo/internal/domain/auth"
	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/storage/pgsql"
	"apibgo/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRI interface {
	GetUser(ctx context.Context, dto domainAuth.UserDto) (domainUser.User, error)
	InsertUser(ctx context.Context, args []interface{}) (domainUser.User, error)
	UpdateUser(ctx context.Context, id int, args []interface{}) (domainUser.User, error)
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
		&user.ConfirmedAt, &user.ConfirmStatus, &user.ConfirmAction,
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
	sql := `INSERT INTO users (email, password, name, surname, confirm_code, confirm_action, confirm_status, token_secret_key, confirmed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()::timestamp, NOW()::timestamp) RETURNING id;`
	id := 0

	err := ar.db.QueryRow(ctx, sql, args...).Scan(&id)

	if err != nil {
		return domainUser.User{}, err
	}

	user, err := ar.GetUser(ctx, domainAuth.UserDto{Id: id})

	if err != nil {
		return domainUser.User{}, err
	}

	return user, nil
}

func (ar *UserRepo) UpdateUser(ctx context.Context, id int, user *domainUser.User) (domainUser.User, pgconn.CommandTag, error) {
	// SQL-query
	sql := "UPDATE users SET "

	var args []interface{} = []interface{}{id}

	// Checking every field in struct, and if don't empty then add to SQL-query
	if utils.IsFieldInitialized(user, "Email") {
		sql += "email = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.Email)
	}
	if utils.IsFieldInitialized(user, "Password") {
		sql += "password = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.Password)
	}
	if utils.IsFieldInitialized(user, "Name") {
		sql += "name = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.Name.String)
	}
	if utils.IsFieldInitialized(user, "Surname") {
		sql += "surname = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.Surname.String)
	}
	if utils.IsFieldInitialized(user, "Activation") {
		sql += "activation = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.Activation)
	}
	if utils.IsFieldInitialized(user, "TokenSecretKey") {
		sql += "token_secret_key = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.TokenSecretKey)
	}
	if utils.IsFieldInitialized(user, "ConfirmCode") {
		sql += "confirm_code = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.ConfirmCode.String)
	}
	if utils.IsFieldInitialized(user, "ConfirmedAt") {
		sql += "confirmed_at = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.ConfirmedAt.Time.Format("2006-01-02 15:04:05"))
	}
	if utils.IsFieldInitialized(user, "ConfirmStatus") {
		sql += "confirm_status = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.ConfirmStatus)
	}
	if utils.IsFieldInitialized(user, "ConfirmAction") {
		sql += "confirm_action = $" + strconv.Itoa(len(args)+1) + ", "
		args = append(args, user.ConfirmAction.String)
	}

	// Removing extra the comma and space in SQL-query ends
	sql = strings.TrimSuffix(sql, ", ")
	sql += ", updated_at = NOW()"

	// Adding condition of WHERE for to update a specific record
	sql += " WHERE id = $1"

	// Execute SQL-query
	commandTag, err := ar.db.Exec(ctx, sql, args...)

	if err != nil {
		return domainUser.User{}, pgconn.CommandTag{}, err
	}

	updUser, err := ar.GetUser(ctx, domainAuth.UserDto{Id: id})

	if err != nil {
		return domainUser.User{}, pgconn.CommandTag{}, err
	}

	return updUser, commandTag, nil
}
