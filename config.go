package main

type Config struct {
	CassandraCluster             string
	Keyspace                     string
	ApiHost                      string
	AuthRealm                    string
	AuthSessionDurationInMinutes int
}
