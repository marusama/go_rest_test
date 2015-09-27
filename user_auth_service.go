package main

import (
	"github.com/grayj/go-json-rest-middleware-tokenauth"
	"time"
)

// Service for UserAuth.
type UserAuthService struct {

	// Auth session duration in minutes.
	sessionDurationInMinutes int

	// Data service.
	dataService *UserAuthDataService
}

// Create new UserAuth service.
func NewUserAuthService(dataConnector DataConnector, sessionDurationInMinutes int) *UserAuthService {
	return &UserAuthService{
		sessionDurationInMinutes: sessionDurationInMinutes,
		dataService:              &UserAuthDataService{dataConnector: dataConnector},
	}
}

// Create (or replace existing) and save in database UserAuth for User login.
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

	// save in db.
	err = uas.dataService.Save(userAuth)

	return userAuth, err
}

// Find and return UserAuth data by token.
// If the token is expired, UserAuth will be deleted and will not be returned.
func (uas *UserAuthService) Get(token string) (userAuth *UserAuth, ok bool, err error) {
	userAuth, ok, err = nil, false, nil
	if token == "" {
		return
	}

	// get userAuth data
	userAuth, ok, err = uas.dataService.FindByToken(token)
	if err != nil {
		return
	}

	// and check expiration
	if ok && userAuth.ExpTime.Before(time.Now()) {
		// remove expired token
		uas.dataService.Remove(userAuth.Login)
		userAuth = nil
		ok = false
	}

	return
}

// Remove stored UserAuth by login.
func (uas *UserAuthService) Remove(login string) error {
	return uas.dataService.Remove(login)
}
