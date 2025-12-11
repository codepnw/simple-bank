// Mocks Data for test business logic
package mocks

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/simple-bank/internal/features/account"
	"github.com/codepnw/simple-bank/internal/features/transfer"
	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
	"github.com/codepnw/simple-bank/internal/features/user"
	"github.com/codepnw/simple-bank/pkg/token"
)

var ErrDatabase = errors.New("database error")

type MockToken struct{}

func (m MockToken) GenerateAccessToken(u *user.User) (string, error) {
	return "", nil
}

func (m MockToken) GenerateRefreshToken(u *user.User) (string, error) {
	return "", nil
}

func (m MockToken) VerifyAccessToken(tokenStr string) (*token.Payload, error) {
	return nil, nil
}

func (m MockToken) VerifyRefreshToken(tokenStr string) (*token.Payload, error) {
	return nil, nil
}

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

func MockAccountData() *account.Account {
	return &account.Account{
		ID:       10,
		OwnerID:  10,
		Balance:  1000,
		Currency: "THB",
	}
}

func MockTransferData(input *transferusecase.TransferParams) *transfer.Transfer {
	return &transfer.Transfer{
		ID:            1,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
	}
}
