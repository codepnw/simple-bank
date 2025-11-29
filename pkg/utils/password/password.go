package password

import (
	"log"

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
