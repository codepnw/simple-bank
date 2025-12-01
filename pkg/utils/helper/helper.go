package helper

import (
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash password failed: %v", err)
		return "", err
	}
	return string(hashed), nil
}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("compare password failed: %v", err)
		return err
	}
	return nil
}

// Validate : validate handler request
func Validate(input any) error {
	v := validator.New()
	return v.Struct(input)
}

func ParseInt64(key string) (int64, error) {
	val, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func ParseInt(key string) int {
	val, err := strconv.Atoi(key)
	if err != nil {
		return 0
	}
	return val
}
