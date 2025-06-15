package services

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	password_bytes := []byte(password)
	hashed_password_bytes, err := bcrypt.GenerateFromPassword(password_bytes,bcrypt.MinCost)
	if err != nil {
		log.Println("failed to hash password:",err.Error())
		return "", errors.New("failed to hash password")
	}
	return string(hashed_password_bytes), nil
}