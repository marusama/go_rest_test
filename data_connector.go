package main

import (
	"github.com/gocql/gocql"
	"log"
)

// TODO: errors handling, implement session pool, auto reconnect, close opened sessions
// for simple demo there is one predefined opened session

type DataConnector struct {
	clusterAddress string
	cluster        *gocql.ClusterConfig
	session        *gocql.Session
}

func NewDataConnector(config *Config) (dataConnector *DataConnector, err error) {
	dataConnector = &DataConnector{}
	dataConnector.clusterAddress = config.CassandraCluster

	log.Println("Connecting to Cassandra cluster...")
	dataConnector.cluster = gocql.NewCluster(dataConnector.clusterAddress)
	dataConnector.session, err = dataConnector.cluster.CreateSession()

	return
}

func (dc *DataConnector) GetSession() *gocql.Session {
	return dc.session
}
