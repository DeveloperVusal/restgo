package pgsql

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Conn(dsn string) (*pgx.Conn, error) {
	conf, err := pgx.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	conf.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	db, err := pgx.ConnectConfig(context.Background(), conf)

	return db, err
}

type Dsn struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func DsnBuild(dsn Dsn) string {
	return "postgres://" + strings.Trim(dsn.Username, " ") + ":" + url.QueryEscape(strings.Trim(dsn.Password, " ")) + "@" + strings.Trim(dsn.Host, " ") + ":" + strconv.Itoa(dsn.Port) + "/" + strings.Trim(dsn.Database, " ")
}
