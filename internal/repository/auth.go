package repository

import (
	"context"

	"apibgo/internal/domain"
	"apibgo/internal/storage/pgsql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRI interface {
	GetUser(ctx context.Context, dto domain.UserDto) pgx.Row
	GetUserToEmail(ctx context.Context, dto domain.UserDto) pgx.Row
	GetAuth(ctx context.Context, dto domain.LoginDto) pgx.Row
	DeleteAuth(ctx context.Context, id int) (pgconn.CommandTag, error)
	InsertAuth(ctx context.Context, args []interface{}) (pgconn.CommandTag, error)
	InsertUser(ctx context.Context, args []interface{}) (pgconn.CommandTag, error)
}

type AuthRepo struct {
	db    *pgx.Conn
	store *pgsql.Storage
}

func NewAuthRepo(store *pgsql.Storage) *AuthRepo {
	return &AuthRepo{
		db:    store.Db,
		store: store,
	}
}

func (ar *AuthRepo) GetUser(ctx context.Context, dto domain.UserDto) pgx.Row {
	sql := `SELECT id, password, activation FROM users WHERE email = $1 LIMIT 1`
	args := []interface{}{dto.Email}

	return ar.db.QueryRow(ctx, sql, args...)
}

func (ar *AuthRepo) GetAuth(ctx context.Context, dto domain.LoginDto) pgx.Row {
	sql := `SELECT id FROM auths WHERE user_agent = $1 AND ip = $2 AND device = $3 LIMIT 1`
	args := []interface{}{dto.UserAgent, dto.Ip, dto.Device}

	return ar.db.QueryRow(ctx, sql, args...)
}

func (ar *AuthRepo) DeleteAuth(ctx context.Context, id int) (pgconn.CommandTag, error) {
	sql := `DELETE FROM auths WHERE id = $1`

	return ar.db.Exec(ctx, sql, id)
}

func (ar *AuthRepo) InsertAuth(ctx context.Context, args []interface{}) (pgconn.CommandTag, error) {
	sql := `INSERT INTO auths (user_id, access_token, refresh_token, ip, device, user_agent, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()::timestamp)`

	return ar.db.Exec(ctx, sql, args...)
}

func (ar *AuthRepo) InsertUser(ctx context.Context, args []interface{}) (int, error) {
	sql := `INSERT INTO users (email, password, name, surname, confirm_code, confirm_status, token_secret_key, confirmed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()::timestamp, NOW()::timestamp) RETURNING id;`
	id := 0

	err := ar.db.QueryRow(ctx, sql, args...).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ar *AuthRepo) GetUserToEmail(ctx context.Context, dto domain.UserDto) pgx.Row {
	sql := `SELECT id, to_char(confirmed_at, 'DD-MM-YYYY HH24:MI:SS') AS confirmed_at, email FROM users WHERE id = $1 LIMIT 1`
	args := []interface{}{dto.Id}

	return ar.db.QueryRow(ctx, sql, args...)
}
