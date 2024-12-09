package middleware

import (
	"os"
	"strconv"
)

func GetCookieSettings() (string, bool, bool, error) {
	env := os.Getenv("ENV")

	var domain string
	var secure, httpOnly bool
	var err error

	if env == "development" {
		domain = os.Getenv("DEV_ALLOWED_ORIGIN")
		secure, err = strconv.ParseBool(os.Getenv("DEV_SECURE_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
		httpOnly, err = strconv.ParseBool(os.Getenv("DEV_HTTP_ONLY_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
	} else {
		domain = os.Getenv("PROD_ALLOWED_ORIGIN")
		secure, err = strconv.ParseBool(os.Getenv("PROD_SECURE_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
		httpOnly, err = strconv.ParseBool(os.Getenv("PROD_HTTP_ONLY_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
	}

	return domain, secure, httpOnly, nil
}
