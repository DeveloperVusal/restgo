package domain

import (
	"database/sql/driver"
	"time"
)

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
	ConfirmStatus  ConfirmStatusEnum
	UpdatedAt      time.Time
	CreatedAt      time.Time
}

func (a *User) TableName() string {
	return "users"
}

type ConfirmStatusEnum string

const (
	ConfirmStatus_QUEST ConfirmStatusEnum = "quest"
	ConfirmStatus_WAIT                    = "waiting"
	ConfirmStatus_ERROR                   = "error"
	ConfirmStatus_OK                      = "success"
)

func (ge *ConfirmStatusEnum) Scan(value interface{}) error {
	*ge = ConfirmStatusEnum(value.([]byte))
	return nil
}

func (ge ConfirmStatusEnum) Value() (driver.Value, error) {
	return string(ge), nil
}
