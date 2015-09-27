package main

import (
	"time"
)

// Simple user authentication storage.
type UserAuth struct {

	// Token string.
	Token string `json:"access_token"`

	// User login.
	Login string `json:"-"`

	// Token expriration datetime.
	ExpTime time.Time `json:"exp_time"`
}
