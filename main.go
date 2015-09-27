package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"strconv"
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
	session, _ := cluster.CreateSession()
	defer session.Close()

	ensureInitDbSchema(session)

	//runApiServer()
}

func ensureInitDbSchema(session *gocql.Session) {
	queryCountTables := "SELECT COUNT(*) FROM system.schema_columnfamilies WHERE keyspace_name='" + Keyspace + "';"
	var tableCount int
	session.Query(queryCountTables).Scan(&tableCount)
	fmt.Println("table count " + strconv.FormatInt(int64(tableCount), 10))
	if tableCount == 0 {
		fmt.Println("Init schema")
		queryCreateKeyspace := "CREATE KEYSPACE " + Keyspace + " WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"
		queryCreateTableUsers := "CREATE TABLE " + Keyspace + ".users (login text, password text, PRIMARY KEY (login));"
		queryCreateTableUserAuths := "CREATE TABLE " + Keyspace + ".user_auths (ftoken text, login text, exp_time int, PRIMARY KEY (login));"
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
