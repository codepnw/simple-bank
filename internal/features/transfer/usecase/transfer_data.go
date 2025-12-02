package transferusecase

import (
	"github.com/codepnw/simple-bank/internal/features/account"
	"github.com/codepnw/simple-bank/internal/features/entry"
	"github.com/codepnw/simple-bank/internal/features/transfer"
)

type TransferParams struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

type TransferResult struct {
	Transfer    *transfer.Transfer `json:"transfer"`
	FromAccount *account.Account   `json:"from_account"`
	ToAccount   *account.Account   `json:"to_account"`
	FromEntry   *entry.Entry       `json:"from_entry"`
	ToEntry     *entry.Entry       `json:"to_entry"`
}
