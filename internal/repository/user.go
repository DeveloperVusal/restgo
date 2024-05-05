package repository

import (
	"context"

	"apibgo/internal/domain"
	"apibgo/internal/storage/pgsql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRI interface {
	InsertUser(ctx context.Context, args []interface{}) (pgconn.CommandTag, error)
	UpdateUser(ctx context.Context, id int, args []interface{}) (pgconn.CommandTag, error)
	DeleteUser(ctx context.Context, id int) (pgconn.CommandTag, error)
	GetUserData(ctx context.Context, dto domain.UserDto) pgx.Row
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

func (ar *UserRepo) GetUserData(ctx context.Context, dto domain.UserDto) (domain.User, error) {
	sql := `SELECT * FROM users WHERE email = $1 or id = $2 LIMIT 1`
	args := []interface{}{dto.Email, dto.Id}

	rows, err := ar.db.Query(ctx, sql, args...)

	if err != nil {
		return domain.User{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ar *UserRepo) DeleteUser(ctx context.Context, id int) (pgconn.CommandTag, error) {
	sql := `DELETE FROM users WHERE id = $1`

	return ar.db.Exec(ctx, sql, id)
}

func (ar *UserRepo) InsertUser(ctx context.Context, args []interface{}) (int, error) {
	sql := `INSERT INTO users (email, password, name, surname, confirm_code, confirm_status, token_secret_key, confirmed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()::timestamp, NOW()::timestamp) RETURNING id;`
	id := 0

	err := ar.db.QueryRow(ctx, sql, args...).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
