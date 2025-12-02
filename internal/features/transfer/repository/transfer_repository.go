package transferrepository

import (
	"context"
	"database/sql"

	"github.com/codepnw/simple-bank/internal/features/transfer"
)

//go:generate mockgen -source=transfer_repository.go -destination=mock_transfer_repository.go -package=transferrepository
type TransferRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, input *transfer.Transfer) (*transfer.Transfer, error)
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Insert(ctx context.Context, tx *sql.Tx, input *transfer.Transfer) (*transfer.Transfer, error) {
	query := `
		INSERT INTO transfers (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := tx.QueryRowContext(ctx, query, input.FromAccountID, input.ToAccountID, input.Amount).Scan(
		&input.ID,
		&input.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return input, nil
}
