package main

import (
	"github.com/gocql/gocql"
	"log"
)

// Connector to Cassandra DB
type DataConnector interface {
	GetSession() *gocql.Session
	GetKeyspace() string
}

// Simple data connector for demo. There is one predefined opened session.
// TODO: errors handling, implement session pool, auto reconnect, close opened sessions
type SimpleDataConnector struct {
	clusterAddress string
	cluster        *gocql.ClusterConfig
	session        *gocql.Session
	keyspace       string
}

// Creates SimpleDataConnector.
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

// Get database session.
func (dc *SimpleDataConnector) GetSession() *gocql.Session {
	return dc.session
}

// Get database keyspace.
func (dc *SimpleDataConnector) GetKeyspace() string {
	return dc.keyspace
}
