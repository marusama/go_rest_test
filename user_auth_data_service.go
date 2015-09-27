package main

import (
	"strconv"
	"time"
)

// Data service for user authentication.
type UserAuthDataService struct {

	// Connector to database.
	dataConnector DataConnector
}

// Save UserAuth in database.
func (ds *UserAuthDataService) Save(userAuth *UserAuth) error {
	query := "INSERT INTO " +
		ds.dataConnector.GetKeyspace() + ".user_auths (ftoken, login, exp_time) " +
		"VALUES (" +
		"'" + userAuth.Token + "', " +
		"'" + userAuth.Login + "', " +
		strconv.FormatInt(userAuth.ExpTime.Unix(), 10) +
		")"
	return ds.dataConnector.GetSession().Query(query).Exec()
}

// Find UserAuth by login.
func (ds *UserAuthDataService) FindByLogin(login string) (userAuth *UserAuth, ok bool, err error) {
	query := "SELECT ftoken, login, exp_time FROM " + ds.dataConnector.GetKeyspace() + ".user_auths " +
		"WHERE login = '" + login + "'"

	var dbToken, dbLogin string
	var dbExpTime int64

	dbErr := ds.dataConnector.GetSession().Query(query).Scan(&dbToken, &dbLogin, &dbExpTime)
	if dbErr != nil {
		if dbErr.Error() == "not found" {
			return nil, false, nil
		} else {
			return nil, false, dbErr
		}
	}

	return &UserAuth{Token: dbToken, Login: dbLogin, ExpTime: time.Unix(dbExpTime, 0)}, true, nil
}

// Find stored UserAuth by token.
func (ds *UserAuthDataService) FindByToken(token string) (userAuth *UserAuth, ok bool, err error) {
	query := "SELECT ftoken, login, exp_time FROM " + ds.dataConnector.GetKeyspace() + ".user_auths " +
		"WHERE ftoken = '" + token + "'"

	var dbToken, dbLogin string
	var dbExpTime int64

	dbErr := ds.dataConnector.GetSession().Query(query).Scan(&dbToken, &dbLogin, &dbExpTime)
	if dbErr != nil {
		if dbErr.Error() == "not found" {
			return nil, false, nil
		} else {
			return nil, false, dbErr
		}
	}

	return &UserAuth{Token: dbToken, Login: dbLogin, ExpTime: time.Unix(dbExpTime, 0)}, true, nil
}

// Remove UserAuth by login.
func (ds *UserAuthDataService) Remove(login string) error {
	query := "DELETE FROM " + ds.dataConnector.GetKeyspace() + ".user_auths " +
		"WHERE login = '" + login + "'"

	return ds.dataConnector.GetSession().Query(query).Exec()
}
