package main

import (
	"errors"
	"sync"
)

type UserService struct {
	sync.RWMutex
	Store map[string]*User
}

func NewUserService() *UserService {
	return &UserService{
		Store: map[string]*User{},
	}
}

func (us *UserService) AddUser(user *User) error {
	if user.Login == "" {
		return errors.New("Login must be provided")
	}
	if user.Password == "" {
		return errors.New("Password must be provided")
	}

	us.Lock()
	defer us.Unlock()

	_, exists := us.Store[user.Login]
	if exists {
		return errors.New("User with login '" + user.Login + "' already exists")
	}

	us.Store[user.Login] = user
	return nil
}

func (us *UserService) CheckUser(login string, password string) bool {
	if login == "" || password == "" {
		return false
	}

	us.RLock()
	defer us.RUnlock()

	v, ok := us.Store[login]
	if ok {
		return v.Password == password
	}

	return false
}
