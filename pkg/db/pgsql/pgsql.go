package pgsql

import (
	"context"

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
