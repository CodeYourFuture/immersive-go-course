package util

import (
	"errors"
	"os"
)

// Get the Postgres password from the environment, either via $POSTGRES_PASSWORD
// or $POSTGRES_PASSWORD_FILE
func ReadPasswd() (string, error) {
	if os.Getenv("POSTGRES_PASSWORD") != "" {
		return os.Getenv("POSTGRES_PASSWORD"), nil
	}

	passwordFile := os.Getenv("POSTGRES_PASSWORD_FILE")
	if passwordFile == "" {
		return "", errors.New("please set POSTGRES_PASSWORD_FILE environment variable")
	}

	pwdFile, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", err
	}
	return string(pwdFile), nil
}
