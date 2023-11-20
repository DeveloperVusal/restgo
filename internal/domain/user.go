package domain

import "time"

type User struct {
	Id             uint
	Email          string
	Password       string
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

func (a *User) TableName() string {
	return "users"
}
