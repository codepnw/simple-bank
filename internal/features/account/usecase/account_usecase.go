package accountusecase

import (
	"context"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/account"
	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	"github.com/codepnw/simple-bank/pkg/auth"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, currency account.AccountCurrency) (*account.Account, error)
	GetAccount(ctx context.Context, id int64) (*account.Account, error)
	ListAccounts(ctx context.Context, pageID, pageSize int) ([]*account.Account, error)
}

type accountUsecase struct {
	repo accountrepository.AccountRepository
}

func NewAccountUsecase(repo accountrepository.AccountRepository) AccountUsecase {
	return &accountUsecase{repo: repo}
}

func (u *accountUsecase) CreateAccount(ctx context.Context, currency account.AccountCurrency) (*account.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, errs.ErrNoUserID
	}

	curr, err := u.validateCurrency(currency)
	if err != nil {
		return nil, err
	}

	accountData, err := u.repo.Insert(ctx, &account.Account{
		OwnerID:  userID,
		Balance:  0,
		Currency: curr,
	})
	if err != nil {
		return nil, err
	}
	return accountData, nil
}

func (u *accountUsecase) validateCurrency(currency account.AccountCurrency) (account.AccountCurrency, error) {
	switch currency {
	case "THB", "thb":
		return account.CurrencyTHB, nil
	case "USD", "usd":
		return account.CurrencyUSD, nil
	default:
		return "", errs.ErrInvalidCurrency
	}
}

func (u *accountUsecase) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, errs.ErrNoUserID
	}

	accountData, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if accountData.OwnerID != userID {
		return nil, errs.ErrAccountNotFound
	}
	return accountData, nil
}

func (u *accountUsecase) ListAccounts(ctx context.Context, pageID, pageSize int) ([]*account.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, errs.ErrNoUserID
	}

	if pageID < 1 {
		pageID = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}
	offset := (pageID - 1) * pageSize

	accounts, err := u.repo.List(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
