package entryrepository

import (
	"context"
	"database/sql"

	"github.com/codepnw/simple-bank/internal/features/entry"
)

//go:generate mockgen -source=entry_repository.go -destination=mock_entry_repository.go -package=entryrepository
type EntryRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, input *entry.Entry) (*entry.Entry, error)
}

type entryRepository struct {
	db *sql.DB
}

func NewEntryRepository(db *sql.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (r *entryRepository) Insert(ctx context.Context, tx *sql.Tx, input *entry.Entry) (*entry.Entry, error) {
	query := `
		INSERT INTO entries (account_id, amount)
		VALUES ($1, $2) RETURNING id, created_at
	`
	err := tx.QueryRowContext(ctx, query, input.AccountID, input.Amount).Scan(
		&input.ID,
		&input.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return input, nil
}
