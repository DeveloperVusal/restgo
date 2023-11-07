package pgsql

import (
	"apibgo/internal/storage"
	"apibgo/pkg/db/pgsql"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
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

	return &Storage{db: db}, nil
}
