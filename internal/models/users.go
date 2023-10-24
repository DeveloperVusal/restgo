package models

import "time"

type Users struct {
	Id             uint
	Email          string
	Activation     bool
	Name           string
	Surname        string
	TokenSecretKey string
	ConfirmCode    string
	ConfirmedAt    time.Time
	ConfirmStatus  string
	UpdatedAt      time.Time
	CreatedAt      time.Time
}
