package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/grayj/go-json-rest-middleware-tokenauth"
	"log"
	"net/http"
)

var (
	userService     = NewUserService()
	userAuthService = NewUserAuthService()
)

func GetAuthMiddleware() rest.Middleware {
	var tokenAuthMiddleware = &tokenauth.AuthTokenMiddleware{
		Realm: AuthRealm,
		Authenticator: func(token string) string {
			userAuth, ok := userAuthService.Get(token)
			if !ok {
				log.Println("Expired or invalid token: " + token)
				return ""
			}
			return userAuth.Login
		},
	}

	var basicAuthMiddleware = &rest.AuthBasicMiddleware{
		Realm: AuthRealm,
		Authenticator: func(user string, password string) bool {
			ok := userService.CheckUser(user, password)
			if !ok {
				log.Println("Failed login for user: '" + user + "', password: '" + password + "'")
			}
			return ok
		},
	}

	return rest.MiddlewareSimple(func(handler rest.HandlerFunc) rest.HandlerFunc {
		return func(w rest.ResponseWriter, request *rest.Request) {
			path := request.URL.Path
			switch {

			// no auth for registering new user
			case path == "/register":
				handler(w, request)

			// basic auth for log in
			case path == "/login":
				basicAuthHandler := basicAuthMiddleware.MiddlewareFunc(handler)
				basicAuthHandler(w, request)

			// token auth for all other API resources
			default:
				tokenAuthHandler := tokenAuthMiddleware.MiddlewareFunc(handler)
				tokenAuthHandler(w, request)
			}
		}
	})
}

func GetRoutes() []*rest.Route {
	routes := []*rest.Route{}
	routes = append(routes,
		rest.Post("/register", register),
		rest.Post("/login", login),
		rest.Get("/auth_test", auth_test),
		rest.Post("/logout", logout),
	)
	return routes
}

func register(w rest.ResponseWriter, r *rest.Request) {
	// decoding new user
	user := User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// trying to add user
	log.Println("Registering user: " + user.Login)
	err = userService.AddUser(&user)
	if err != nil {
		log.Println("Register error: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("User registered: " + user.Login)

	// response
	w.WriteJson(map[string]string{
		"status": "OK",
	})
}

func login(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["REMOTE_USER"].(string)

	userAuth, err := userAuthService.Set(user)

	if err != nil {
		log.Println("Login error: " + err.Error())
		rest.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Logged in: " + user)
	w.WriteJson(userAuth)
}

func auth_test(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["REMOTE_USER"].(string)
	log.Println("Testing auth: " + user)

	w.WriteJson(map[string]string{"authed": user})
}

func logout(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["REMOTE_USER"].(string)

	userAuthService.Remove(user)

	log.Println("Logged out: " + user)
	w.WriteJson(map[string]string{"status": "OK"})
}
