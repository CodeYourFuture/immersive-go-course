package util

import (
	"errors"
	"os"
)

func ReadPasswdFile() (string, error) {
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
