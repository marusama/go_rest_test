package main

// Simple User with login and password.
// Warning! Password is stored as plain text. DO NOT USE IN PRODUCTION.
type User struct {

	// User login.
	Login string `json:"login"`

	// User password. TODO: Replace with hash.
	Password string `json:"password"`
}
