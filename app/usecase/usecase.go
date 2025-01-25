package usecase

import (
	"log"
	"messageApp/app/adapter"
	"messageApp/app/domain"
)

type UserUsecase struct {
	reepository domain.UserRepository
}

func NewUserUsecase() *UserUsecase {
	reepository, err := adapter.NewDatabaseManager()
	if err != nil {
		log.Println("Cannot create database manager", err)
		return nil
	}
	return &UserUsecase{reepository: reepository}
}

func (uc *UserUsecase) RegisterUser(name string, pass string) error {
	err := uc.reepository.RegisterUser(name, pass)
	if err != nil {
		log.Println("Cannot register user", err)
		return err
	}
	return nil
}

func (uc *UserUsecase) AuthenticateUser(name string, pass string) bool {
	result := uc.reepository.AuthenticateUser(name, pass)
	if !result {
		log.Println("Authentication failed for user:", name)
		return result
	}
	return true
}

func (uc *UserUsecase) IsAccessTokenValid(token string) bool {
	result := uc.reepository.IsAccessTokenValid(token)
	if !result {
		log.Println("Invalid access token:", token)
		return result
	}
	return true
}

func (uc *UserUsecase) SelectAccessToken(username string, pass string) string {
	result := uc.reepository.SelectAccessToken(username, pass)
	if result == "" {
		log.Println("No access token found for user:", username)
		return result
	}
	return result
}
