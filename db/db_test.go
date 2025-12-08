package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// GoLandなら、この関数の左に出る「緑の三角マーク」を押すだけで実行できます
func TestDBConnection(t *testing.T) {
	// 1. 環境変数の取得 (GoLandの実行構成で設定するか、ここに直接書いてテスト後は消す)
	// 面倒なら一時的に直接書いてしまってもテストファイルなのでOKです
	dbUser := os.Getenv("MYSQL_USER")
	dbPwd := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DATABASE")

	// 接続情報の組み立て
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPwd, dbHost, dbName)

	// 2. 接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("接続設定エラー: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("DB接続失敗: %v", err)
	}
	t.Log("✅ DB接続成功")

	// 3. テーブル作成 (まだ作ってない場合)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		t.Fatalf("テーブル作成失敗: %v", err)
	}

	// 4. INSERTテスト
	res, err := db.Exec("INSERT INTO users (name) VALUES (?)", "Test User via Go")
	if err != nil {
		t.Fatalf("INSERT失敗: %v", err)
	}
	id, _ := res.LastInsertId()
	t.Logf("✅ INSERT成功 ID: %d", id)

	// 5. SELECTテスト
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		t.Fatalf("SELECT失敗: %v", err)
	}
	t.Logf("✅ SELECT成功 Name: %s", name)
}
