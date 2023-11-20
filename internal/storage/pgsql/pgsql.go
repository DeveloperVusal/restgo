package pgsql

import (
	"apibgo/internal/domain"
	"apibgo/internal/storage"
	"apibgo/pkg/db/pgsql"
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	Db *pgx.Conn
}

func New(cfg *storage.Config, conn string) (*Storage, error) {
	var dsn storage.Cluster
	const op = "storage.pgsql.New()"

	pgcfg := cfg.PgSql
	reftt := reflect.TypeOf(pgcfg)
	re := regexp.MustCompile(`^yaml:"([a-zA-Z0-9_-]+)"`)

	for i := 0; i < reftt.NumField(); i++ {
		field := reftt.Field(i)
		fieldTag := field.Tag
		matches := re.FindAllStringSubmatch(string(fieldTag), -1)
		connName := ""

		if len(matches) == 1 {
			connName = matches[0][1]
		}

		if connName == conn {
			vv := reflect.ValueOf(pgcfg).Field(i)
			dsn = vv.Interface().(storage.Cluster)
			break
		}
	}

	cfgPort, _ := strconv.Atoi(dsn.Port)
	dsnobj := pgsql.Dsn{
		Host:     dsn.Host,
		Port:     cfgPort,
		Username: dsn.Username,
		Password: dsn.Password,
		Database: dsn.Database,
	}
	dsnstr := pgsql.DsnBuild(dsnobj)
	db, err := pgsql.Conn(dsnstr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if dsn.Migrate {
		mig, err := migrate.New("file://"+os.Getenv("MIGRATIONS_PATH"), dsnstr+"?sslmode=disable")

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = mig.Up()

		if err != nil && err != migrate.ErrNoChange {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Storage{Db: db}, nil
}

func (s *Storage) Create(ctx context.Context, au domain.Auth) (pgconn.CommandTag, error) {
	ds := goqu.Insert(au.TableName()).Rows(au)
	sql, args, _ := ds.ToSQL()
	cmmtag, err := s.Db.Exec(ctx, sql, args...)

	return cmmtag, err
}

func (s *Storage) Save(ctx context.Context, au domain.Auth, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ds := goqu.Update(au.TableName()).Set(au).Where(goqu.L(sql, args...))
	_sql, _args, _ := ds.ToSQL()
	cmmtag, err := s.Db.Exec(ctx, _sql, _args...)

	return cmmtag, err
}

func (s *Storage) Delete(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	au := domain.Auth{}
	ds := goqu.Delete(au.TableName()).Where(goqu.L(sql, args...))
	_sql, _args, _ := ds.ToSQL()
	cmmtag, err := s.Db.Exec(ctx, _sql, _args...)

	return cmmtag, err
}

func (s *Storage) Find(ctx context.Context, table string, sql string, args ...interface{}) (pgx.Rows, error) {
	ds := goqu.From(table).Where(goqu.L(sql, args...))
	_sql, _args, _ := ds.ToSQL()
	rows, err := s.Db.Query(ctx, _sql, _args...)

	return rows, err
}

func (s *Storage) First(ctx context.Context, table string, sql string, args ...interface{}) (pgx.Rows, error) {
	ds := goqu.From(table).
		Where(goqu.L(sql, args...)).
		Limit(1).
		Order(goqu.I("id").Asc())
	_sql, _args, _ := ds.ToSQL()
	rows, err := s.Db.Query(ctx, _sql, _args...)

	return rows, err
}

func (s *Storage) Last(ctx context.Context, table string, sql string, args ...interface{}) (pgx.Rows, error) {
	ds := goqu.From(table).
		Where(goqu.L(sql, args...)).
		Limit(1).
		Order(goqu.I("id").Desc())
	_sql, _args, _ := ds.ToSQL()
	rows, err := s.Db.Query(ctx, _sql, _args...)

	return rows, err
}
