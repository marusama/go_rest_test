package main

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/grayj/go-json-rest-middleware-tokenauth"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	userService = UserService{
		Store: map[string]*User{},
	}
	userAuthService = UserAuthService{
		LoginToUserAuth: map[string]*UserAuth{},
		TokenToUserAuth: map[string]*UserAuth{},
	}
)

func main() {

	var authRealm = "test"

	var tokenAuthMiddleware = &tokenauth.AuthTokenMiddleware{
		Realm: authRealm,
		Authenticator: func(token string) string {
			userAuth, ok := userAuthService.Get(token)
			if !ok {
				return ""
			}
			return userAuth.Login
		},
	}

	var basicAuthMiddleware = &rest.AuthBasicMiddleware{
		Realm: authRealm,
		Authenticator: func(user string, password string) bool {
			return userService.CheckUser(user, password)
		},
	}
	api := rest.NewApi()
	//api.Use(rest.DefaultDevStack...)

	// auth types middleware
	api.Use(rest.MiddlewareSimple(func(h rest.HandlerFunc) rest.HandlerFunc {
		return func(w rest.ResponseWriter, request *rest.Request) {
			path := request.URL.Path
			fmt.Println("switcher: " + path)
			switch {

			// no auth for registering new user
			case path == "/register":
				fmt.Println("/register")
				h(w, request)

			// basic auth for log in
			case path == "/login":
				fmt.Println("/login")
				basicAuthHandler := basicAuthMiddleware.MiddlewareFunc(h)
				basicAuthHandler(w, request)

			// token auth for all other API resources
			default:
				fmt.Println("default")
				tokenAuthHandler := tokenAuthMiddleware.MiddlewareFunc(h)
				tokenAuthHandler(w, request)
			}
		}
	}))

	// router
	router, err := rest.MakeRouter(
		rest.Post("/register", register),
		rest.Post("/login", login),
		rest.Get("/auth_test", auth_test),
		rest.Post("/logout", logout),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	address := "localhost:8080"
	fmt.Println("Listening " + address + "...")
	log.Fatal(http.ListenAndServe(address, api.MakeHandler()))
}

func register(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println("register func")

	// decoding new user
	user := User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// trying to add user
	err = userService.AddUser(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// response
	w.WriteJson(map[string]string{
		"status": "OK",
	})
}

func login(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println("login func")

	userAuth, err := userAuthService.Set(r.Env["REMOTE_USER"].(string))

	if err != nil {
		rest.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(userAuth)
}

func auth_test(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println("auth_test func")

	w.WriteJson(map[string]string{"authed": r.Env["REMOTE_USER"].(string)})
}

func logout(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println("logout func")

	userAuthService.Remove(r.Env["REMOTE_USER"].(string))
	w.WriteJson(map[string]string{"status": "OK"})
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserService struct {
	sync.RWMutex
	Store map[string]*User
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

type UserAuth struct {
	Token   string    `json:"access_token"`
	Login   string    `json:"login"`
	ExpTime time.Time `json:"exp_time"`
}

type UserAuthService struct {
	sync.RWMutex
	LoginToUserAuth map[string]*UserAuth
	TokenToUserAuth map[string]*UserAuth
}

func (uas *UserAuthService) Set(login string) (userAuth *UserAuth, err error) {
	uas.Lock()
	defer uas.Unlock()

	// remove old auth info, if exist
	uas.removeByLogin(login)

	token, err := tokenauth.New()
	if err != nil {
		return
	}
	tokenHash := tokenauth.Hash(token)
	expTime := time.Now().Add(60 * time.Second)

	// set auth info
	userAuth = &UserAuth{Token: tokenHash, Login: login, ExpTime: expTime}
	uas.LoginToUserAuth[login] = userAuth
	uas.TokenToUserAuth[token] = userAuth

	return
}

func (uas *UserAuthService) Get(token string) (userAuth *UserAuth, ok bool) {
	uas.RLock()
	defer uas.RUnlock()

	tokenHash := tokenauth.Hash(token)

	userAuth, ok = uas.TokenToUserAuth[tokenHash]
	if ok && userAuth.ExpTime.Before(time.Now()) {
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
