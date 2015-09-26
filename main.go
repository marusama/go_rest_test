package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gocql/gocql"
	"log"
	"net/http"
)

// some options for API
const (
	CassandraCluster = "127.0.0.1"
	Keyspace         = "test_rest"

	ApiHost                  = ":8080"
	AuthRealm                = "test_realm"
	SessionDurationInMinutes = 60
)

func main() {
	cluster := gocql.NewCluster(CassandraCluster)
	cluster.Keyspace = "demo"
	session, _ := cluster.CreateSession()
	defer session.Close()

	ensureDbSchema(session)

	//runApiServer()
}

func ensureDbSchema(session *gocql.Session) {
	queryCountTables := "SELECT COUNT(*) FROM system.schema_columnfamilies WHERE keyspace_name='" + Keyspace + "';"
	var tableCount int
	session.Query(queryCountTables).Scan(&tableCount)
	if tableCount == 0 {
		queryCreateKeyspace := "CREATE KEYSPACE " + Keyspace + ";"
		queryCreateTableUsers := "CREATE TABLE " + Keyspace + ".users (login text, password text, PRIMARY KEY (login));"
		queryCreateTableUserAuths := "CREATE TABLE " + Keyspace + ".user_auths (token text, login text, exp_time timestamp, PRIMARY KEY (login));"
		session.Query(queryCreateKeyspace).Exec()
		session.Query(queryCreateTableUsers).Exec()
		session.Query(queryCreateTableUserAuths).Exec()
	}
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
