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

	// ルートとハンドラのマッピング
	routes := map[string]http.HandlerFunc{
		"/ws":             handler.HandleWebsocket,
		"/chat/":          handler.HandleChatPage,
		"/home/":          handler.MakeHandler("home"),
		"/register/":      handler.MakeHandler("registerForm"),
		"/register/user/": handler.HandleRegisterUser,
		"/auth/login/":    handler.HandleAuth,
		"/notfound/":      handler.MakeHandler("notfound"),
		"/unauthorized/":  handler.MakeHandler("unauthorized"),
	}
	// ルートとハンドラのマッピングを登録
	for route, handler := range routes {
		http.HandleFunc(route, handler)
	}

	// 静的ファイルのサーブ
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app/"))))

	// メッセージの処理
	go handler.HandleMessages()

	// サーバーの起動
	log.Println("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
