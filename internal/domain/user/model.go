package user

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type User struct {
	Id             uint              `db:"id"`
	Email          string            `db:"email"`
	Password       string            `db:"password"`
	Activation     bool              `db:"activation"`
	Name           sql.NullString    `db:"name"`
	Surname        sql.NullString    `db:"surname"`
	TokenSecretKey string            `db:"token_secret_key,omitempty"`
	ConfirmCode    sql.NullString    `db:"confirm_code,omitempty"`
	ConfirmAction  sql.NullString    `db:"confirm_action,omitempty"`
	ConfirmedAt    sql.NullTime      `db:"confirmed_at,omitempty"`
	ConfirmStatus  ConfirmStatusEnum `db:"confirm_status,omitempty"`
	UpdatedAt      sql.NullTime      `db:"updated_at,omitempty"`
	CreatedAt      time.Time         `db:"created_at"`
}

func (a *User) TableName() string {
	return "users"
}

type ConfirmStatusEnum string

const (
	ConfirmStatus_QUEST   ConfirmStatusEnum = "quest"
	ConfirmStatus_WAIT                      = "waiting"
	ConfirmStatus_ERROR                     = "error"
	ConfirmStatus_OK                        = "success"
	ConfirmStatus_UNKNOWN                   = "unknown"
)

func (ge *ConfirmStatusEnum) Scan(value interface{}) error {
	if value == nil {
		*ge = ConfirmStatus_UNKNOWN
		return nil
	}

	switch v := value.(type) {
	case string:
		switch string(v) {
		case "quest":
			*ge = ConfirmStatus_QUEST
		case "waiting":
			*ge = ConfirmStatus_WAIT
		case "error":
			*ge = ConfirmStatus_ERROR
		case "success":
			*ge = ConfirmStatus_OK
		default:
			*ge = ConfirmStatus_UNKNOWN
		}
	default:
		return fmt.Errorf("unexpected type for ConfirmStatusEnum: %T", value)
	}
	return nil
}

func (ge ConfirmStatusEnum) Value() (driver.Value, error) {
	return string(ge), nil
}
