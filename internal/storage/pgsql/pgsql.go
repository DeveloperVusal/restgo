package pgsql

import (
	"apibgo/internal/storage"
	"apibgo/pkg/db/pgsql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

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
	dsnstr := pgsql.Dsn{
		Host:     dsn.Host,
		Port:     cfgPort,
		Username: dsn.Username,
		Password: dsn.Password,
		Database: dsn.Database,
	}
	db, err := pgsql.Conn(pgsql.DsnBuild(dsnstr))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
