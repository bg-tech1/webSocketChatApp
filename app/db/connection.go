package db

import (
	"database/sql"
	"fmt"
	"log"
	"messageApp/app/util"
	"os"

	_ "github.com/lib/pq"
)

var con *sql.DB

type dbInfo struct {
	host     string
	port     string
	name     string
	user     string
	password string
}

type DatabaseManager struct {
	Con  *sql.DB
	info dbInfo
}

func NewDatabaseManager() (*DatabaseManager, error) {
	info := dbInfo{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		name:     os.Getenv("DB_NAME"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
	}
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		info.host, info.port, info.user, info.password, info.name)
	con, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}
	return &DatabaseManager{
		Con:  con,
		info: info}, nil
}

func InitDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	con, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	// userinfoテーブル作成
	_, err = con.Exec("CREATE TABLE IF NOT EXISTS userinfo (id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(50) NOT NULL, accesstoken VARCHAR(50) NOT NULL)")
	if err != nil {
		log.Println("Error creating table:", err)
		return
	}
	log.Println("Connected to database")
}

func (dm *DatabaseManager) RegisterUser(name string, pass string) {
	token, _ := util.GenerateRandomString(32)
	_, err := con.Exec("INSERT INTO userinfo (username,password,accesstoken) VALUES ($1,$2,$3)", name, pass, token)
	if err != nil {
		log.Println("Error inserting user:", err)
	}
}

func (dm *DatabaseManager) AuthenticateUser(name string, pass string) bool {
	var exists bool
	err := con.QueryRow("SELECT EXISTS(SELECT 1 FROM userinfo WHERE username=$1 AND password=$2)", name, pass).Scan(&exists)
	if err != nil {
		log.Println("Error querying user:", err)
		return false
	}
	return exists
}

func (dm *DatabaseManager) IsAccessTokenValid(token string) bool {
	var exists bool
	err := con.QueryRow("SELECT EXISTS(SELECT 1 FROM userinfo WHERE accesstoken=$1)", token).Scan(&exists)
	if err != nil {
		log.Println("Error querying user by access token:", err)
		return false
	}
	return exists
}

func (dm *DatabaseManager) SelectAccessToken(username string, pass string) string {
	var token string
	err := con.QueryRow("SELECT accesstoken FROM userinfo WHERE username=$1 AND password=$2", username, pass).Scan(&token)
	if err != nil {
		log.Println("Error querying access token:", err)
		return ""
	}
	return token
}
