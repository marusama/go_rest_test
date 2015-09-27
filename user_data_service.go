package main

import ()

type UserDataService struct {
	dataConnector DataConnector
}

func (ds *UserDataService) Save(user *User) error {
	query := "INSERT INTO " +
		ds.dataConnector.GetKeyspace() + ".users (login, password) " +
		"VALUES (" +
		"'" + user.Login + "', " +
		"'" + user.Password + "'" +
		")"
	return ds.dataConnector.GetSession().Query(query).Exec()
}

func (ds *UserDataService) Find(login string) (user *User, ok bool, err error) {
	query := "SELECT login, password FROM " + ds.dataConnector.GetKeyspace() + ".users " +
		"WHERE login = '" + login + "'"

	var dbLogin, dbPassword string

	dbErr := ds.dataConnector.GetSession().Query(query).Scan(&dbLogin, &dbPassword)
	if dbErr != nil {
		if dbErr.Error() == "not found" {
			return nil, false, nil
		} else {
			return nil, false, dbErr
		}
	}

	return &User{Login: dbLogin, Password: dbPassword}, true, nil
}
