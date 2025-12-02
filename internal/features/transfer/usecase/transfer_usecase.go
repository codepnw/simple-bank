package transferusecase

import (
	"context"
	"database/sql"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/account"
	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	"github.com/codepnw/simple-bank/internal/features/entry"
	entryrepository "github.com/codepnw/simple-bank/internal/features/entry/repository"
	"github.com/codepnw/simple-bank/internal/features/transfer"
	transferrepository "github.com/codepnw/simple-bank/internal/features/transfer/repository"
	"github.com/codepnw/simple-bank/pkg/auth"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
)

type TransferUsecase interface {
	Transfer(ctx context.Context, input *TransferParams) (*TransferResult, error)
}

type transferUsecase struct {
	tranRepo transferrepository.TransferRepository
	accRepo  accountrepository.AccountRepository
	entRepo  entryrepository.EntryRepository
	tx       database.TxManager
}

func NewTransferUsecase(
	tranRepo transferrepository.TransferRepository,
	accRepo accountrepository.AccountRepository,
	entRepo entryrepository.EntryRepository,
	tx database.TxManager,
) TransferUsecase {
	return &transferUsecase{
		tranRepo: tranRepo,
		accRepo:  accRepo,
		entRepo:  entRepo,
		tx:       tx,
	}
}

func (u *transferUsecase) Transfer(ctx context.Context, input *TransferParams) (*TransferResult, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userID := auth.GetUserID(ctx)

	if input.FromAccountID == input.ToAccountID {
		return nil, errs.ErrTransferToSelf
	}

	fromAcc, err := u.accRepo.FindByID(ctx, input.FromAccountID)
	if err != nil {
		return nil, err
	}
	// Check Owner
	if fromAcc.OwnerID != userID {
		return nil, errs.ErrAccountNotFound
	}
	// Check Currency
	if fromAcc.Currency != account.AccountCurrency(input.Currency) {
		return nil, errs.ErrCurrencyMismatch
	}

	toAcc, err := u.accRepo.FindByID(ctx, input.ToAccountID)
	if err != nil {
		return nil, err
	}
	// Check Currency
	if fromAcc.Currency != toAcc.Currency {
		return nil, errs.ErrCurrencyMismatch
	}
	// Check Balance
	if fromAcc.Balance < input.Amount {
		return nil, errs.ErrMoneyNotEnough
	}

	result := new(TransferResult)
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		var err error

		// Create Transfer
		result.Transfer, err = u.tranRepo.Insert(ctx, tx, &transfer.Transfer{
			FromAccountID: input.FromAccountID,
			ToAccountID:   input.ToAccountID,
			Amount:        input.Amount,
		})
		if err != nil {
			return err
		}

		// Create Entry From Account (minus)
		result.FromEntry, err = u.entRepo.Insert(ctx, tx, &entry.Entry{
			AccountID: input.FromAccountID,
			Amount:    -input.Amount, // minus
		})
		if err != nil {
			return err
		}

		// Create Entry To Account (plus)
		result.ToEntry, err = u.entRepo.Insert(ctx, tx, &entry.Entry{
			AccountID: input.ToAccountID,
			Amount:    input.Amount, // plus
		})
		if err != nil {
			return err
		}

		// Update Balance
		// NOTE: Prevent "Deadlock" sort by ID
		if input.FromAccountID < input.ToAccountID {
			// Lock 1 (From) -> Lock 2 (To)
			result.FromAccount, result.ToAccount, err = u.addMoney(ctx, tx, input.FromAccountID, -input.Amount, input.ToAccountID, input.Amount)
		} else {
			// Lock 1 (To) -> Lock 2 (From)
			result.FromAccount, result.ToAccount, err = u.addMoney(ctx, tx, input.ToAccountID, input.Amount, input.FromAccountID, -input.Amount)
		}

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *transferUsecase) addMoney(ctx context.Context, tx *sql.Tx, accID1, amount1, accID2, amount2 int64) (acc1 *account.Account, acc2 *account.Account, err error) {
	acc1, err = u.accRepo.AddAccountBalance(ctx, tx, accID1, amount1)
	if err != nil {
		return nil, nil, err
	}

	acc2, err = u.accRepo.AddAccountBalance(ctx, tx, accID2, amount2)
	if err != nil {
		return nil, nil, err
	}

	return acc1, acc2, nil
}
