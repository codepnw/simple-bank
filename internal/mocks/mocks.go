// Mocks Data for test business logic
package mocks

import (
	"errors"

	"github.com/codepnw/simple-bank/internal/features/user"
)

var ErrDatabase = errors.New("database error")

func MockUserData() *user.User {
	return &user.User{
		ID:        10,
		Username:  "testmock001",
		FirstName: "john",
		LastName:  "cena",
		Email:     "mock@example.com",
		Password:  "password",
	}
}
