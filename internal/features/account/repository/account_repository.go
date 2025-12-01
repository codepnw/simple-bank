package accountrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/simple-bank/internal/features/account"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
)

//go:generate mockgen -source=account_repository.go -destination=mock_account_repository.go -package=accountrepository
type AccountRepository interface {
	Insert(ctx context.Context, input *account.Account) (*account.Account, error)
	FindByID(ctx context.Context, accountID int64) (*account.Account, error)
	List(ctx context.Context, ownerID int64, limit, offset int) ([]*account.Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Insert(ctx context.Context, input *account.Account) (*account.Account, error) {
	query := `
		INSERT INTO accounts (owner_id, balance, currency)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, input.OwnerID, input.Balance, input.Currency).Scan(
		&input.ID,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (r *accountRepository) FindByID(ctx context.Context, accountID int64) (*account.Account, error) {
	query := `
		SELECT id, owner_id, balance, currency, created_at, updated_at
		FROM accounts WHERE id = $1 LIMIT 1
	`
	acc := new(account.Account)
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&acc.ID,
		&acc.OwnerID,
		&acc.Balance,
		&acc.Currency,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrAccountNotFound
		}
		return nil, err
	}
	return acc, nil
}

func (r *accountRepository) List(ctx context.Context, ownerID int64, limit int, offset int) ([]*account.Account, error) {
	query := `
		SELECT id, owner_id, balance, currency, created_at, updated_at
		FROM accounts WHERE owner_id = $1 LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, ownerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accs := make([]*account.Account, 0)

	for rows.Next() {
		acc := new(account.Account)
		if err = rows.Scan(
			&acc.ID,
			&acc.OwnerID,
			&acc.Balance,
			&acc.Currency,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accs, nil
}
