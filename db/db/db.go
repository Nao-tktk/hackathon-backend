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

	dsn := fmt.Sprintf("%s:%s@%s/%s", user, pwd, mysqlHost, Database)
	return sql.Open("mysql", dsn)
}
