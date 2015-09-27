package main

import (
	"log"
)

// Backend services.
type Services struct {

	// User service.
	UserService *UserService

	// UserAuth service.
	UserAuthService *UserAuthService
}

// Creates new services.
func NewServices(dataConnector DataConnector, config *Config) (services *Services, err error) {

	// check and create Cassandra database schema if needed.
	err = ensureInitDbSchema(dataConnector, config)
	if err != nil {
		return
	}

	// create new services.
	services = &Services{}
	services.UserService = NewUserService(dataConnector)
	services.UserAuthService = NewUserAuthService(dataConnector, config.AuthSessionDurationInMinutes)

	return
}

// Check database schema.
// TODO: use migration tools
func ensureInitDbSchema(dataConnector DataConnector, config *Config) error {
	keyspace := config.Keyspace

	queryCountTables := "SELECT COUNT(*) FROM system.schema_columnfamilies WHERE keyspace_name='" + keyspace + "';"
	var tableCount int

	// Check that our tables exist
	err := dataConnector.GetSession().Query(queryCountTables).Scan(&tableCount)
	if err != nil {
		return err
	}

	if tableCount == 0 {
		log.Println("Initializing DB schema...")

		// keyspace
		queryCreateKeyspace := "CREATE KEYSPACE " + keyspace + " WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"
		dataConnector.GetSession().Query(queryCreateKeyspace).Exec()
		if err != nil {
			return err
		}

		// users
		queryCreateTableUsers := "CREATE TABLE " + keyspace + ".users (login text, password text, PRIMARY KEY (login));"
		dataConnector.GetSession().Query(queryCreateTableUsers).Exec()
		if err != nil {
			return err
		}

		// user_auths
		queryCreateTableUserAuths := "CREATE TABLE " + keyspace + ".user_auths (ftoken text, login text, exp_time int, PRIMARY KEY (login));"
		dataConnector.GetSession().Query(queryCreateTableUserAuths).Exec()
		if err != nil {
			return err
		}

		// secondary index for search by token
		querySecondaryIndexUserAuths := "CREATE INDEX user_auths_token_idx ON " + keyspace + ".user_auths (ftoken)"
		dataConnector.GetSession().Query(querySecondaryIndexUserAuths).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
