package main

// Simple config for REST API and Cassandra DB
type Config struct {
	CassandraCluster             string
	Keyspace                     string
	ApiHost                      string
	AuthRealm                    string
	AuthSessionDurationInMinutes int
}
