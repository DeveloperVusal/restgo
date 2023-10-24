package _mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conn(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}
