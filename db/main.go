package main

import (
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

	// ルーティング
	http.HandleFunc("/user", userController.Handler)

	log.Println("Listening on :8080...")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// 終了待ち
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Server shutting down...")
}
