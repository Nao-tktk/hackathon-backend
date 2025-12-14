package main

import (
	"fmt" // 追加
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"db/controller"
	"db/dao"
	"db/db"
	"db/usecase"
)

func main() {
	// DB接続
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// 組み立て (DI)
	userDao := dao.NewUserDao(dbConn)
	userUsecase := usecase.NewUserUsecase(userDao)
	userController := controller.NewUserController(userUsecase)

	itemDao := dao.NewItemDao(dbConn)
	itemUsecase := usecase.NewItemUsecase(itemDao)
	itemController := controller.NewItemController(itemUsecase)

	txDao := dao.NewTransactionDao(dbConn)
	txUsecase := usecase.NewTransactionUsecase(txDao)
	txController := controller.NewTransactionController(txUsecase)

	// ルーティング
	http.HandleFunc("/user", userController.Handler)
	http.HandleFunc("/items", itemController.Handler)
	http.HandleFunc("/purchace", txController.Handler)
	http.HandleFunc("/db-check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1. そもそも繋がっているか (PING)
		if err := dbConn.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status": "NG", "message": "DB接続自体が失敗しています: %v"}`, err)
			return
		}

		// 2. テーブルは何があるか (SHOW TABLES)
		rows, err := dbConn.Query("SHOW TABLES")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status": "NG", "message": "接続できたけどテーブル一覧が見れません: %v"}`, err)
			return
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var t string
			if err := rows.Scan(&t); err == nil {
				tables = append(tables, t)
			}
		}

		// 結果を表示
		fmt.Fprintf(w, `{"status": "OK", "message": "接続成功", "tables_found": %q}`, tables)
	})
	// ▲▲▲▲▲ 追加ここまで ▲▲▲▲▲

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %s...", addr)

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// 終了待ち
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Server shutting down...")
}
