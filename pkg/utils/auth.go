package utils

import (
	"net/http"
	"strings"
)

func RealIp(r *http.Request) string {
	addr := strings.Split(r.RemoteAddr, ":")

	return addr[0]
}
