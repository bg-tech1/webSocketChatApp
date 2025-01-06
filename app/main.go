package main

import (
	"log"
	"messageApp/app/db"
	"messageApp/app/handler"
	"net/http"
)

func main() {
	// データベースの初期化
	db.InitDatabase()
	// 各種ハンドラ
	http.HandleFunc("/ws", handler.HandleWebsocket)
	http.HandleFunc("/chat/", handler.HandleChatPage)
	http.HandleFunc("/home/", handler.MakeHandler("home"))
	http.HandleFunc("/register/", handler.MakeHandler("registerForm"))
	http.HandleFunc("/register/user/", handler.HandleRegisterUser)
	http.HandleFunc("/auth/login/", handler.HandleAuth)
	http.HandleFunc("/notfound/", handler.MakeHandler("notfound"))
	http.HandleFunc("/badrequest/", handler.MakeHandler("badRequest"))
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app/"))))
	// メッセージの処理
	go handler.HandleMessages()
	// サーバーの起動
	log.Println("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
