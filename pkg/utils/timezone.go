package utils

import (
	"os"
	"time"
)

var defaultZone string = "Asia/Baku"

func GetTimezone() (*time.Location, error) {
	tz := os.Getenv("APP_TIMEZONE")

	if tz == "" {
		return time.LoadLocation(tz)
	} else {
		return time.LoadLocation(defaultZone)
	}
}
