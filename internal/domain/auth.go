package domain

import "time"

type Auth struct {
	Id           uint      `db:"id"`
	UserId       uint      `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	Ip           string    `db:"ip"`
	Device       string    `db:"device"`
	UserAgent    string    `db:"user_agent"`
	CreatedAt    time.Time `db:"created_at"`
}

func (a *Auth) TableName() string {
	return "auths"
}

type AuthDto struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Device    string `json:"device" validate:"omitempty"`
	Ip        string
	UserAgent string
}
