package accountusecase_test

import (
	"context"
	"testing"

	"github.com/codepnw/simple-bank/internal/features/account"
	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	accountusecase "github.com/codepnw/simple-bank/internal/features/account/usecase"
	"github.com/codepnw/simple-bank/internal/mocks"
	"github.com/codepnw/simple-bank/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		currency    account.AccountCurrency
		mockFn      func(mockRepo *accountrepository.MockAccountRepository, currency account.AccountCurrency)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:     "success",
			userID:   10,
			currency: account.AccountCurrency("THB"),
			mockFn: func(mockRepo *accountrepository.MockAccountRepository, currency account.AccountCurrency) {
				a := mocks.MockAccountData()
				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(a, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:     "fail db error",
			userID:   10,
			currency: account.AccountCurrency("THB"),
			mockFn: func(mockRepo *accountrepository.MockAccountRepository, currency account.AccountCurrency) {
				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := setup(t)

			tc.mockFn(mockRepo, tc.currency)

			ctx := context.Background()
			if tc.userID != 0 {
				ctx = auth.SetUserID(ctx, tc.userID)
			}

			result, err := uc.CreateAccount(ctx, tc.currency)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		accountID   int64
		mockFn      func(mockRepo *accountrepository.MockAccountRepository, accountID int64)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:      "success",
			userID:    10,
			accountID: 10,
			mockFn: func(mockRepo *accountrepository.MockAccountRepository, accountID int64) {
				a := mocks.MockAccountData()
				mockRepo.EXPECT().FindByID(gomock.Any(), accountID).Return(a, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:      "fail db error",
			userID:    10,
			accountID: 10,
			mockFn: func(mockRepo *accountrepository.MockAccountRepository, accountID int64) {
				mockRepo.EXPECT().FindByID(gomock.Any(), accountID).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := setup(t)

			tc.mockFn(mockRepo, tc.accountID)

			ctx := context.Background()
			if tc.userID != 0 {
				ctx = auth.SetUserID(ctx, tc.userID)
			}

			result, err := uc.GetAccount(ctx, tc.accountID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestListAccounts(t *testing.T) {
	type testCase struct {
		name        string
		userID      int64
		mockFn      func(mockRepo *accountrepository.MockAccountRepository, userID int64)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:   "success",
			userID: 10,
			mockFn: func(mockRepo *accountrepository.MockAccountRepository, userID int64) {
				accs := []*account.Account{
					{ID: 10, OwnerID: 10},
					{ID: 11, OwnerID: 10},
				}
				mockRepo.EXPECT().List(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(accs, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:   "fail db error",
			userID: 10,
			mockFn: func(mockRepo *accountrepository.MockAccountRepository, userID int64) {
				mockRepo.EXPECT().List(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, mockRepo := setup(t)

			tc.mockFn(mockRepo, tc.userID)

			ctx := context.Background()
			if tc.userID != 0 {
				ctx = auth.SetUserID(ctx, tc.userID)
			}

			result, err := uc.ListAccounts(ctx, 1, 5)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func setup(t *testing.T) (accountusecase.AccountUsecase, *accountrepository.MockAccountRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := accountrepository.NewMockAccountRepository(ctrl)
	uc := accountusecase.NewAccountUsecase(mockRepo)

	return uc, mockRepo
}
