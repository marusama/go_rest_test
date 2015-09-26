package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

// some options for API
const (
	ApiHost                  = ":8080"
	AuthRealm                = "test_realm"
	SessionDurationInMinutes = 60
)

func main() {
	runApiServer()
}

func runApiServer() {
	api := rest.NewApi()
	//api.Use(rest.DefaultDevStack...)

	// set middleware for authentication
	api.Use(GetAuthMiddleware())

	// router
	router, err := rest.MakeRouter(GetRoutes()...)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	fmt.Println("Starting to listen " + ApiHost + "...")
	log.Fatal(http.ListenAndServe(ApiHost, api.MakeHandler()))
}
