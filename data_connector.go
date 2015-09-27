package main

import (
	"github.com/gocql/gocql"
	"log"
)

type DataConnector interface {
	GetSession() *gocql.Session
	GetKeyspace() string
}

// TODO: errors handling, implement session pool, auto reconnect, close opened sessions
// for simple demo there is one predefined opened session

type SimpleDataConnector struct {
	clusterAddress string
	cluster        *gocql.ClusterConfig
	session        *gocql.Session
	keyspace       string
}

func NewDataConnector(config *Config) (dataConnector DataConnector, err error) {
	simpleDataConnector := &SimpleDataConnector{}
	simpleDataConnector.clusterAddress = config.CassandraCluster
	simpleDataConnector.keyspace = config.Keyspace

	log.Println("Connecting to Cassandra cluster...")
	simpleDataConnector.cluster = gocql.NewCluster(simpleDataConnector.clusterAddress)
	simpleDataConnector.session, err = simpleDataConnector.cluster.CreateSession()

	dataConnector = simpleDataConnector
	return
}

func (dc *SimpleDataConnector) GetSession() *gocql.Session {
	return dc.session
}

func (dc *SimpleDataConnector) GetKeyspace() string {
	return dc.keyspace
}
