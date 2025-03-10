package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) < 5 {
		return "", errors.New("password is too short (at least 5 characters)")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	if len(password) < 5 {
		return errors.New("password is too short")
	}
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
}
