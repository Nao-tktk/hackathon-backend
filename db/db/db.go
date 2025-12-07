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
	dbName := os.Getenv("MYSQL_DATABASE")

	if user == "" || pwd == "" || dbName == "" {
		return nil, fmt.Errorf("env MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE must be set")
	}

	dsn := fmt.Sprintf("%s:%s@(localhost:3306)/%s", user, pwd, dbName)
	return sql.Open("mysql", dsn)
}
