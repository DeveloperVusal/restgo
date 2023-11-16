package repository

import (
	"context"

	"apibgo/internal/domain"
	"apibgo/internal/storage/pgsql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Auths interface {
	Create(ctx context.Context, au domain.Auth) (pgconn.CommandTag, error)
	Save(ctx context.Context, au domain.Auth, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Delete(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Find(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	First(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Last(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type AuthRepo struct {
	db *pgx.Conn
}

func NewAuthRepo(store *pgsql.Storage) *AuthRepo {
	return &AuthRepo{
		db: store.Db,
	}
}

func (ar *AuthRepo) Create(ctx context.Context, au domain.Auth) (pgconn.CommandTag, error) {
	ds := goqu.Insert(au.TableName()).Rows(au)
	sql, args, _ := ds.ToSQL()
	cmmtag, err := ar.db.Exec(ctx, sql, args...)

	return cmmtag, err
}

func (ar *AuthRepo) Save(ctx context.Context, au domain.Auth, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ds := goqu.Update(au.TableName()).Set(au).Where(goqu.L(sql, args...))
	_sql, _args, _ := ds.ToSQL()
	cmmtag, err := ar.db.Exec(ctx, _sql, _args...)

	return cmmtag, err
}

func (ar *AuthRepo) Delete(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	au := domain.Auth{}
	ds := goqu.Delete(au.TableName()).Where(goqu.L(sql, args...))
	_sql, _args, _ := ds.ToSQL()
	cmmtag, err := ar.db.Exec(ctx, _sql, _args...)

	return cmmtag, err
}

func (ar *AuthRepo) Find(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	au := domain.Auth{}
	ds := goqu.From(au.TableName()).Where(goqu.L(sql, args...))
	_sql, _args, _ := ds.ToSQL()
	rows, err := ar.db.Query(ctx, _sql, _args...)

	return rows, err
}

func (ar *AuthRepo) First(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	au := domain.Auth{}
	ds := goqu.From(au.TableName()).
		Where(goqu.L(sql, args...)).
		Limit(1).
		Order(goqu.I("id").Asc())
	_sql, _args, _ := ds.ToSQL()
	rows, err := ar.db.Query(ctx, _sql, _args...)

	return rows, err
}

func (ar *AuthRepo) Last(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	au := domain.Auth{}
	ds := goqu.From(au.TableName()).
		Where(goqu.L(sql, args...)).
		Limit(1).
		Order(goqu.I("id").Desc())
	_sql, _args, _ := ds.ToSQL()
	rows, err := ar.db.Query(ctx, _sql, _args...)

	return rows, err
}
