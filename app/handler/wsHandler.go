package handler

import (
	"encoding/json"
	"log"
	"messageApp/app/db"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }}
	Users     = make(map[User]bool)
	broadcast = make(chan []byte)
)

type AuthResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	AcessToken string `json:"access_token"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"pass"`
}

type User struct {
	client *websocket.Conn
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	var userInfo UserInfo
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
		})
	}
	username := userInfo.Username
	password := userInfo.Password
	result := db.AuthenticateUser(username, password)
	if result {
		log.Println("User authenticated")
		token := db.SelectAccessToken(username, password)
		if token != "" {
			json.NewEncoder(w).Encode(AuthResponse{
				StatusCode: http.StatusOK,
				Message:    "User authenticated",
				AcessToken: token,
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(AuthResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "User not authenticated",
		})
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("app/view/" + tmpl + ".html")
	if err != nil {
		log.Println("Error parsing template:", err)
		return
	}
	if data != nil {
		t.Execute(w, data)
	} else {
		t.Execute(w, nil)
	}
}

func HandleMessages() {
	for {
		message := <-broadcast
		for user := range Users {
			err := user.client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				user.client.Close()
				delete(Users, user)
			}
		}
	}
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	user := User{client: conn}
	Users[user] = true

	defer conn.Close()
	log.Println("Client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		broadcast <- msg
	}
}

func HandleChatPage(w http.ResponseWriter, r *http.Request) {
	//不正ないアクセスを防ぐためにアクセストークンの認証をする
	if !db.IsAccessTokenValid(r.URL.Query().Get("access_token")) {
		http.Redirect(w, r, "/badrequest/", http.StatusSeeOther)
	}
	renderTemplate(w, "chat", nil)
}

func MakeHandler(title string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, title, nil)
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userInfo UserInfo
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
		})
	}
	username := userInfo.Username
	password := userInfo.Password
	db.RegisterUser(username, password)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		StatusCode: http.StatusCreated,
		Message:    "User created",
		AcessToken: db.SelectAccessToken(username, password),
	})
}
