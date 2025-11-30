// Mocks Data for test business logic
package mocks

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/simple-bank/internal/features/user"
)

var ErrDatabase = errors.New("database error")

type MockDB struct{}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return nil
}
func (m *MockDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, nil
}
func (m *MockDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, nil
}

type MockTx struct{}

func (m *MockTx) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	return fn(nil)
}

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
