package main

import (
	"github.com/grayj/go-json-rest-middleware-tokenauth"
	"sync"
	"time"
)

type UserAuthService struct {
	sync.RWMutex
	LoginToUserAuth map[string]*UserAuth
	TokenToUserAuth map[string]*UserAuth
}

func NewUserAuthService() *UserAuthService {
	return &UserAuthService{
		LoginToUserAuth: map[string]*UserAuth{},
		TokenToUserAuth: map[string]*UserAuth{},
	}
}

func (uas *UserAuthService) Set(login string) (userAuth *UserAuth, err error) {
	uas.Lock()
	defer uas.Unlock()

	// remove old auth info, if exist
	uas.removeByLogin(login)

	// generate token and expiration datetime
	token, err := tokenauth.New()
	if err != nil {
		return
	}
	token = tokenauth.Hash(token)
	expTime := time.Now().Add(SessionDurationInMinutes * time.Minute)

	// set auth info
	userAuth = &UserAuth{Token: token, Login: login, ExpTime: expTime}
	uas.LoginToUserAuth[login] = userAuth
	uas.TokenToUserAuth[token] = userAuth

	return
}

func (uas *UserAuthService) Get(token string) (userAuth *UserAuth, ok bool) {
	userAuth, ok = nil, false
	if token == "" {
		return
	}

	uas.RLock()
	defer uas.RUnlock()

	// get userAuth data and check expiration
	userAuth, ok = uas.TokenToUserAuth[token]
	if ok && userAuth.ExpTime.Before(time.Now()) {
		// remove expired token
		uas.removeByToken(token)
		userAuth = nil
		ok = false
	}

	return
}

func (uas *UserAuthService) Remove(login string) {
	uas.Lock()
	defer uas.Unlock()

	uas.removeByLogin(login)
}

func (uas *UserAuthService) removeByLogin(login string) {
	userAuth, ok := uas.LoginToUserAuth[login]
	if ok {
		delete(uas.TokenToUserAuth, userAuth.Token)
		delete(uas.LoginToUserAuth, userAuth.Login)
	}
}

func (uas *UserAuthService) removeByToken(token string) {
	userAuth, ok := uas.TokenToUserAuth[token]
	if ok {
		delete(uas.TokenToUserAuth, userAuth.Token)
		delete(uas.LoginToUserAuth, userAuth.Login)
	}
}
