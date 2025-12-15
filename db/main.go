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

	//以下調整用
	categorySQL := `
    INSERT IGNORE INTO categories (id, name) VALUES 
    (1, '本・雑誌'),
    (2, '家電・スマホ'),
    (3, 'ファッション');
    `
	if _, err := dbConn.Exec(categorySQL); err != nil {
		log.Printf("カテゴリー追加エラー(無視可): %v", err)
	} else {
		log.Println("カテゴリーデータの準備完了")
	}

	// 2. テストユーザー (ID=1)
	userSQL := `
    INSERT IGNORE INTO users (id, name, password, created_at) VALUES 
    (1, 'テスト太郎', 'pass1234', NOW());
    `
	if _, err := dbConn.Exec(userSQL); err != nil {
		log.Printf("ユーザー追加エラー(無視可): %v", err)
	} else {
		log.Println("テストユーザー(ID=1)の準備完了")
	}
	////////////////

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
	http.HandleFunc("/login", userController.LoginHandler)
	http.HandleFunc("/items", itemController.Handler)
	http.HandleFunc("/purchase", txController.Handler)

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
