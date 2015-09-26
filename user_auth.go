package main

import (
	"time"
)

type UserAuth struct {
	Token   string    `json:"access_token"`
	Login   string    `json:"-"`
	ExpTime time.Time `json:"exp_time"`
}
