package _mysql

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Conn(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}

type Dsn struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Query    string
}

func DsnBuild(dsn Dsn) string {
	dsnSrc := strings.Trim(dsn.Username, " ") + ":" + strings.Trim(dsn.Password, " ") + "@tcp(" + strings.Trim(dsn.Host, " ") + ":" + strconv.Itoa(dsn.Port) + ")/" + strings.Trim(dsn.Database, " ")

	query := strings.Trim(dsn.Query, " ")

	if len(query) > 0 {
		dsnSrc += "?" + strings.Trim(dsn.Query, " ")
	}

	return dsnSrc
}
