package generate

import (
	crand "crypto/rand"
	"encoding/base64"
	mrand "math/rand"
	"strconv"
)

func RandomNumbers(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		result += strconv.Itoa(mrand.Intn(9))
	}

	return result
}

func RandomStringBytes(n int) (string, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}
