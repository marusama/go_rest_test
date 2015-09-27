package main

import (
	"github.com/grayj/go-json-rest-middleware-tokenauth"
	"time"
)

type UserAuthService struct {
	sessionDurationInMinutes int
	dataService              *UserAuthDataService
}

func NewUserAuthService(dataConnector DataConnector, sessionDurationInMinutes int) *UserAuthService {
	return &UserAuthService{
		sessionDurationInMinutes: sessionDurationInMinutes,
		dataService:              &UserAuthDataService{dataConnector: dataConnector},
	}
}

func (uas *UserAuthService) Set(login string) (userAuth *UserAuth, err error) {
	// remove old auth info, if exist
	err = uas.dataService.Remove(login)
	if err != nil {
		return
	}

	// generate token and expiration datetime
	token, err := tokenauth.New()
	if err != nil {
		return
	}
	token = tokenauth.Hash(token)
	expTime := time.Now().Add(time.Duration(uas.sessionDurationInMinutes) * time.Minute)

	// set auth info
	userAuth = &UserAuth{Token: token, Login: login, ExpTime: expTime}

	err = uas.dataService.Save(userAuth)

	return userAuth, err
}

func (uas *UserAuthService) Get(token string) (userAuth *UserAuth, ok bool, err error) {
	userAuth, ok, err = nil, false, nil
	if token == "" {
		return
	}

	// get userAuth data and check expiration
	userAuth, ok, err = uas.dataService.FindByToken(token)
	if err != nil {
		return
	}

	if ok && userAuth.ExpTime.Before(time.Now()) {
		// remove expired token
		uas.dataService.Remove(userAuth.Login)
		userAuth = nil
		ok = false
	}

	return
}

func (uas *UserAuthService) Remove(login string) error {
	return uas.dataService.Remove(login)
}
