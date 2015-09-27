package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/grayj/go-json-rest-middleware-tokenauth"
	"log"
	"net/http"
)

// Run demo API host.
func RunApiServer(services *Services, config *Config) error {
	api := rest.NewApi()

	//api.Use(rest.DefaultDevStack...)

	// set middleware for authentications
	api.Use(getAuthMiddleware(services, config.AuthRealm))

	// router
	router, err := rest.MakeRouter(getRoutes(services)...)
	if err != nil {
		return err
	}
	api.SetApp(router)

	log.Println("Starting to listen " + config.ApiHost + "...")
	return http.ListenAndServe(config.ApiHost, api.MakeHandler())
}

// Get authentication middleware for REST API.
func getAuthMiddleware(services *Services, authRealm string) rest.Middleware {

	// token-based authentication
	var tokenAuthMiddleware = &tokenauth.AuthTokenMiddleware{
		Realm: authRealm,
		Authenticator: func(token string) string {
			userAuth, ok, err := services.UserAuthService.Get(token)
			if err != nil {
				log.Fatal(err)
				return ""
			}
			if !ok {
				log.Println("Expired or invalid token: " + token)
				return ""
			}
			return userAuth.Login
		},
	}

	// basic access authentication
	var basicAuthMiddleware = &rest.AuthBasicMiddleware{
		Realm: authRealm,
		Authenticator: func(user string, password string) bool {
			ok := services.UserService.CheckUser(user, password)
			if !ok {
				log.Println("Failed login for user: '" + user + "', password: '" + password + "'")
			}
			return ok
		},
	}

	// return middleware func that selects appropriate auth middleware for different route paths.
	return rest.MiddlewareSimple(func(handler rest.HandlerFunc) rest.HandlerFunc {
		return func(w rest.ResponseWriter, request *rest.Request) {
			path := request.URL.Path

			// switch by path:
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

// Get routes for REST API.
func getRoutes(services *Services) []*rest.Route {
	routes := []*rest.Route{}
	routes = append(routes,
		rest.Post("/register", func(w rest.ResponseWriter, r *rest.Request) { register(services, w, r) }),
		rest.Post("/login", func(w rest.ResponseWriter, r *rest.Request) { login(services, w, r) }),
		rest.Get("/auth_test", func(w rest.ResponseWriter, r *rest.Request) { auth_test(services, w, r) }),
		rest.Post("/logout", func(w rest.ResponseWriter, r *rest.Request) { logout(services, w, r) }),
	)
	return routes
}

// Handler for /register. Register new User.
func register(s *Services, w rest.ResponseWriter, r *rest.Request) {
	// decoding new user
	user := User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// trying to add user
	log.Println("Registering user: " + user.Login)
	err = s.UserService.AddUser(&user)
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

// Handler for /login. Log in the client and return generated token.
func login(s *Services, w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["REMOTE_USER"].(string)

	// login. generates token and stores in UserAuth.
	userAuth, err := s.UserAuthService.Set(user)

	if err != nil {
		log.Println("Login error: " + err.Error())
		rest.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Logged in: " + user)
	w.WriteJson(userAuth)
}

// Handler for /auth_test. Test function that requires authorization. Return authed login.
func auth_test(s *Services, w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["REMOTE_USER"].(string)
	log.Println("Testing auth: " + user)

	w.WriteJson(map[string]string{"authed": user})
}

// Handler for /logout. Log out the client.
func logout(s *Services, w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["REMOTE_USER"].(string)

	s.UserAuthService.Remove(user)

	log.Println("Logged out: " + user)
	w.WriteJson(map[string]string{"status": "OK"})
}
