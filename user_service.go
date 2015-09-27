package main

import (
	"errors"
	"log"
)

type UserService struct {
	dataService *UserDataService
}

func NewUserService(dataConnector DataConnector) *UserService {
	return &UserService{
		dataService: &UserDataService{dataConnector: dataConnector},
	}
}

func (us *UserService) AddUser(user *User) error {
	if user.Login == "" {
		return errors.New("Login must be provided")
	}
	if user.Password == "" {
		return errors.New("Password must be provided")
	}

	_, exists, err := us.dataService.Find(user.Login)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("User with login '" + user.Login + "' already exists")
	}

	return us.dataService.Save(user)
}

func (us *UserService) CheckUser(login string, password string) bool {
	if login == "" || password == "" {
		return false
	}

	v, ok, err := us.dataService.Find(login)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if ok {
		return v.Password == password
	}

	return false
}
