package repository

import (
	"context"
	"errors"

	domainAuth "apibgo/internal/domain/auth"
	"apibgo/internal/storage/pgsql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRI interface {
	GetAuth(ctx context.Context, dto domainAuth.RefreshDto) (domainAuth.Auth, error)
	DeleteAuth(ctx context.Context, id int) (pgconn.CommandTag, error)
	InsertAuth(ctx context.Context, args []interface{}) (pgconn.CommandTag, error)
	GetSessions(ctx context.Context, dto domainAuth.SessionDto) ([]domainAuth.Auth, error)
	GetCountSessions(ctx context.Context) (int, error)
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

func (ar *AuthRepo) GetAuth(ctx context.Context, dto domainAuth.RefreshDto) (domainAuth.Auth, error) {
	var auth domainAuth.Auth

	args := []interface{}{}
	cond := ""

	if dto.Refresh == "" {
		cond = `user_agent = $1 AND ip = $2 AND device = $3`
		args = append(args, dto.UserAgent, dto.Ip, dto.Device)
	} else {
		cond = `refresh_token = $1`
		args = append(args, dto.Refresh)
	}

	sql := `SELECT * FROM auths WHERE ` + cond + ` LIMIT 1`

	err := ar.db.QueryRow(ctx, sql, args...).Scan(
		&auth.Id, &auth.UserId, &auth.AccessToken, &auth.RefreshToken,
		&auth.Ip, &auth.Device, &auth.UserAgent, &auth.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainAuth.Auth{}, nil
		}

		return domainAuth.Auth{}, err
	}

	return auth, nil
}

func (ar *AuthRepo) GetSessions(ctx context.Context, dto domainAuth.SessionDto) ([]domainAuth.Auth, error) {
	var auths []domainAuth.Auth

	authModel := domainAuth.Auth{}

	sql := `SELECT * FROM ` + authModel.TableName() + ` WHERE user_id = $1`
	rows, err := ar.db.Query(ctx, sql, dto.Id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		auth := domainAuth.Auth{}
		err = rows.Scan(
			&auth.Id, &auth.UserId, &auth.AccessToken, &auth.RefreshToken,
			&auth.Ip, &auth.Device, &auth.UserAgent, &auth.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		auths = append(auths, auth)
	}

	return auths, nil
}

func (ar *AuthRepo) GetCountSessions(ctx context.Context, dto domainAuth.SessionDto) (int, error) {
	var auth domainAuth.Auth
	var count int

	sql := `SELECT COUNT(id) FROM ` + auth.TableName() + ` WHERE user_id = $1`
	err := ar.db.QueryRow(ctx, sql, dto.Id).Scan(&count)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (ar *AuthRepo) DeleteAuth(ctx context.Context, dto domainAuth.DestroyDto) (pgconn.CommandTag, error) {
	sql := `DELETE FROM auths WHERE access_token = $1 OR id = $2`

	return ar.db.Exec(ctx, sql, dto.Token, dto.Id)
}

func (ar *AuthRepo) InsertAuth(ctx context.Context, args []interface{}) (pgconn.CommandTag, error) {
	sql := `INSERT INTO auths (user_id, access_token, refresh_token, ip, device, user_agent, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()::timestamp)`

	return ar.db.Exec(ctx, sql, args...)
}
