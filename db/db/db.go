package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*sql.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	Database := os.Getenv("MYSQL_DATABASE")

	if user == "" || pwd == "" || Database == "" {
		return nil, fmt.Errorf("env MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE must be set")
	}

	var dsn string
	if len(mysqlHost) > 0 && mysqlHost[0] == '/' {
		// 【本番用】 Unixドメインソケット接続 (Cloud Run)
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true", user, pwd, mysqlHost, Database)
	} else {
		// 【ローカル用】 TCP接続
		dsn = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, pwd, mysqlHost, Database)
	}

	return sql.Open("mysql", dsn)
}
