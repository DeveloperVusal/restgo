package auth

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
