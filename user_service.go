package main

import (
	"errors"
	"log"
)

// Service for User.
type UserService struct {

	// Data service for User.
	dataService *UserDataService
}

// Create new service for User.
func NewUserService(dataConnector DataConnector) *UserService {
	return &UserService{
		dataService: &UserDataService{dataConnector: dataConnector},
	}
}

// Save User in database.
func (us *UserService) AddUser(user *User) error {
	// check properties.
	if user.Login == "" {
		return errors.New("Login must be provided")
	}
	if user.Password == "" {
		return errors.New("Password must be provided")
	}

	// search existing in db.
	_, exists, err := us.dataService.Find(user.Login)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("User with login '" + user.Login + "' already exists")
	}

	// save in db.
	return us.dataService.Save(user)
}

// Try to find User with these login/password.
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
