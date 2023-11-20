package repository

import (
	"context"

	"apibgo/internal/domain"
	"apibgo/internal/storage/pgsql"

	"github.com/jackc/pgx/v5"
)

type AuthRI interface {
	GetUser(ctx context.Context, dto domain.AuthDto) (pgx.Rows, error)
	GetAuth(ctx context.Context, dto domain.AuthDto) (pgx.Rows, error)
	DeleteAuth(ctx context.Context, id int) (pgx.Rows, error)
	InsertAuth(ctx context.Context, args []interface{}) (pgx.Rows, error)
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

func (ar *AuthRepo) GetUser(ctx context.Context, dto domain.AuthDto) (pgx.Rows, error) {
	sql := `SELECT id, password, activation FROM users WHERE email = $1 LIMIT 1`
	args := []interface{}{dto.Email}

	return ar.db.Query(ctx, sql, args...)
}

func (ar *AuthRepo) GetAuth(ctx context.Context, dto domain.AuthDto) (pgx.Rows, error) {
	sql := `SELECT id FROM auths WHERE user_agent = $1 AND ip = $2 AND device = $3 LIMIT 1`
	args := []interface{}{dto.UserAgent, dto.Ip, dto.Device}

	return ar.db.Query(ctx, sql, args...)
}

func (ar *AuthRepo) DeleteAuth(ctx context.Context, id int) (pgx.Rows, error) {
	sql := `DELETE FROM auths WHERE id = $1`

	return ar.db.Query(ctx, sql, id)
}

func (ar *AuthRepo) InsertAuth(ctx context.Context, args []interface{}) (pgx.Rows, error) {
	sql := `INSERT INTO auths (user_id, access_toen, refresh_toen, ip, device, user_agent) VALUES ($1, $2, $3, $4, $5, $6)`

	return ar.db.Query(ctx, sql, args...)
}
