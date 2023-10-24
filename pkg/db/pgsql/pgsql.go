package pgsql

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Conn(dsn string) (context.Context, *pgx.Conn, error) {
	config, _ := pgx.ParseConfig(dsn)
	// config.BuildStatementCache = func(conn *pgconn.PgConn) stmtcache.Cache {
	// 	return stmtcache.New(conn, stmtcache.ModeDescribe, 1024)
	// }

	_db, _err := pgx.ConnectConfig(context.Background(), config)

	return context.Background(), _db, _err
}

type Dsn struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func DsnBuild(dsn Dsn) string {
	return "postgres://" + strings.Trim(dsn.Username, " ") + ":" + strings.Trim(dsn.Password, " ") + "@" + strings.Trim(dsn.Host, " ") + ":" + strconv.Itoa(dsn.Port) + "/" + strings.Trim(dsn.Database, " ")
}
