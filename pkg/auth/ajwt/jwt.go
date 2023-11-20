package ajwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserId uint `json:"user_id"`
}

type JWT struct {
	Secret           string
	UserId           int
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func (j *JWT) NewPairTokens() (string, string) {
	var signedKey interface{} = []byte(j.Secret)

	accessExpiresAt := time.Now().Add(time.Minute * 15)
	refreshExpiresAt := time.Now().Add(time.Minute * 43830)

	if !j.AccessExpiresAt.IsZero() {
		accessExpiresAt = j.AccessExpiresAt
	}

	if !j.RefreshExpiresAt.IsZero() {
		refreshExpiresAt = j.RefreshExpiresAt
	}

	initToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserId: uint(j.UserId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			Issuer:    strconv.FormatInt(time.Now().Unix(), 10),
		},
	})

	access, _ := initToken.SignedString(signedKey)

	initToken = jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserId: uint(j.UserId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			Issuer:    strconv.FormatInt(time.Now().Unix(), 10),
		},
	})
	refresh, _ := initToken.SignedString(signedKey)

	return access, refresh
}
