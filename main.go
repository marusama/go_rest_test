package main

import (
	"log"
)

func main() {
	config := &Config{
		CassandraCluster:             "127.0.0.1",
		Keyspace:                     "test_rest",
		ApiHost:                      ":8080",
		AuthRealm:                    "test_realm",
		AuthSessionDurationInMinutes: 60,
	}

	log.Println("Init...")

	// connecting to cassandra db
	dataConnector, err := NewDataConnector(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	// creating services
	services, err := NewServices(dataConnector, config)
	if err != nil {
		log.Fatal(err)
		return
	}

	// run API
	log.Fatal(RunApiServer(services, config))
}
