package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var con *sql.DB

func init() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	var err error
	for i := 0; i < 10; i++ {
		con, err = sql.Open("postgres", dbInfo)
		if err == nil {
			err = con.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Error connecting to database: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}
	con, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println("Error connecting to database:", err)
	}
	_, err = con.Exec("CREATE TABLE IF NOT EXISTS userinfo (id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(50) NOT NULL, accesstoken VARCHAR(50) NOT NULL)")
	if err != nil {
		log.Println("Error creating table:", err)
	}
	log.Println("Connected to database")
}

func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func RegisterUser(name string, pass string) {
	token, _ := GenerateRandomString(32)
	_, err := con.Exec("INSERT INTO userinfo (username,password,accesstoken) VALUES ($1,$2,$3)", name, pass, token)
	if err != nil {
		log.Println("Error inserting user:", err)
	}
}

func AuthenticateUser(name string, pass string) bool {
	var exists bool
	err := con.QueryRow("SELECT EXISTS(SELECT 1 FROM userinfo WHERE username=$1 AND password=$2)", name, pass).Scan(&exists)
	if err != nil {
		log.Println("Error querying user:", err)
		return false
	}
	return exists
}

func IsAccessTokenValid(token string) bool {
	var exists bool
	err := con.QueryRow("SELECT EXISTS(SELECT 1 FROM userinfo WHERE accesstoken=$1)", token).Scan(&exists)
	if err != nil {
		log.Println("Error querying user by access token:", err)
		return false
	}
	return exists
}

func SelectAccessToken(username string, pass string) string {
	var token string
	err := con.QueryRow("SELECT accesstoken FROM userinfo WHERE username=$1 AND password=$2", username, pass).Scan(&token)
	if err != nil {
		log.Println("Error querying access token:", err)
		return ""
	}
	return token
}
