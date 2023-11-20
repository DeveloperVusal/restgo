package pswd

import "golang.org/x/crypto/bcrypt"

func HashPassword(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
