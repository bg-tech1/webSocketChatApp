package adapter

import (
	"database/sql"
	"fmt"
	"log"
	"messageApp/app/util"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type DatabaseManager struct {
	Con *sql.DB
}

var (
	instance *DatabaseManager
	once     sync.Once
)

func NewDatabaseManager() (*DatabaseManager, error) {
	var err error
	once.Do(func() {
		dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		con, err := sql.Open("postgres", dbInfo)
		if err != nil {
			log.Println("Error opening database:", err)
			return
		}
		instance = &DatabaseManager{Con: con}
		_, err = con.Exec("CREATE TABLE IF NOT EXISTS userinfo (id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, username VARCHAR(50) NOT NULL, password VARCHAR(50) NOT NULL, accesstoken VARCHAR(50) NOT NULL)")
		if err != nil {
			log.Println("Error creating table:", err)
			return
		}
		log.Println("Create userinfo table")
	})
	return instance, err
}

func (dm *DatabaseManager) RegisterUser(name string, pass string) error {
	token, _ := util.GenerateRandomString(32)
	_, err := dm.Con.Exec("INSERT INTO userinfo (username,password,accesstoken) VALUES ($1,$2,$3)", name, pass, token)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

func (dm *DatabaseManager) AuthenticateUser(name string, pass string) bool {
	var exists bool
	err := dm.Con.QueryRow("SELECT EXISTS(SELECT 1 FROM userinfo WHERE username=$1 AND password=$2)", name, pass).Scan(&exists)
	log.Println("exists:", exists)
	if err != nil {
		log.Println("Error querying user:", err)
		return false
	}
	return exists
}

func (dm *DatabaseManager) IsAccessTokenValid(token string) bool {
	var exists bool
	err := dm.Con.QueryRow("SELECT EXISTS(SELECT 1 FROM userinfo WHERE accesstoken=$1)", token).Scan(&exists)
	if err != nil {
		log.Println("Error querying user by access token:", err)
		return false
	}
	return exists
}

func (dm *DatabaseManager) SelectAccessToken(username string, pass string) string {
	var token string
	err := dm.Con.QueryRow("SELECT accesstoken FROM userinfo WHERE username=$1 AND password=$2", username, pass).Scan(&token)
	if err != nil {
		log.Println("Error querying access token:", err)
		return ""
	}
	return token
}
