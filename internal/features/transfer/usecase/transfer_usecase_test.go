package transferusecase_test

import (
	"context"
	"testing"

	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	"github.com/codepnw/simple-bank/internal/features/entry"
	entryrepository "github.com/codepnw/simple-bank/internal/features/entry/repository"
	transferrepository "github.com/codepnw/simple-bank/internal/features/transfer/repository"
	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
	"github.com/codepnw/simple-bank/internal/mocks"
	"github.com/codepnw/simple-bank/pkg/auth"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	type testCase struct {
		name        string
		input       *transferusecase.TransferParams
		mockFn      func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success FromID < ToID",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.OwnerID = 100
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				mockTrans := mocks.MockTransferData(input)
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrans, nil).Times(1)

				mockFromEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: -input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFromEnt, nil).Times(1)

				mockToEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockToEnt, nil).Times(1)

				addFromAcc := mocks.MockAccountData()
				addFromAcc.Balance += -input.Amount
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.FromAccountID, -input.Amount).Return(addFromAcc, nil).Times(1)

				addToAcc := mocks.MockAccountData()
				addToAcc.Balance += input.Amount
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.ToAccountID, input.Amount).Return(addToAcc, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "success FromID > ToID",
			input: &transferusecase.TransferParams{
				FromAccountID: 3,
				ToAccountID:   2,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.OwnerID = 100
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				mockTrans := mocks.MockTransferData(input)
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrans, nil).Times(1)

				mockFromEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: -input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFromEnt, nil).Times(1)

				mockToEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockToEnt, nil).Times(1)

				addToAcc := mocks.MockAccountData()
				addToAcc.Balance += input.Amount
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.ToAccountID, input.Amount).Return(addToAcc, nil).Times(1)

				addFromAcc := mocks.MockAccountData()
				addFromAcc.Balance += -input.Amount
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.FromAccountID, -input.Amount).Return(addFromAcc, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail transfer to self",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   1,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
			},
			expectedErr: errs.ErrTransferToSelf,
		},
		{
			name: "fail user no account",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				a := mocks.MockAccountData()
				a.OwnerID = 100
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(a, nil).Times(1)
			},
			expectedErr: errs.ErrAccountNotFound,
		},
		{
			name: "fail get from account",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
		{
			name: "fail get to account",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				fromAcc.Balance = 20
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
		{
			name: "fail currency not equal",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.Currency = "USD"
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)
			},
			expectedErr: errs.ErrCurrencyMismatch,
		},
		{
			name: "fail currency mismatch",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				a := mocks.MockAccountData()
				a.Currency = "USD"
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(a, nil).Times(1)
			},
			expectedErr: errs.ErrCurrencyMismatch,
		},
		{
			name: "fail money not enough",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				fromAcc.Balance = 20
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.ID = 2
				toAcc.OwnerID = 100
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)
			},
			expectedErr: errs.ErrMoneyNotEnough,
		},
		{
			name: "fail insert transfer",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.ID = 2
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
		{
			name: "fail insert from account entry",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.ID = 2
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				mockTrans := mocks.MockTransferData(input)
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrans, nil).Times(1)

				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
		{
			name: "fail insert to account entry",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.ID = 2
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				mockTrans := mocks.MockTransferData(input)
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrans, nil).Times(1)

				mockEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: -input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockEnt, nil).Times(1)

				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
		{
			name: "fail add acc1 balance",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.ID = 2
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				mockTrans := mocks.MockTransferData(input)
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrans, nil).Times(1)

				mockFromEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: -input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFromEnt, nil).Times(1)

				mockToEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockToEnt, nil).Times(1)

				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.FromAccountID, -input.Amount).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
		{
			name: "fail add acc2 balance",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				fromAcc := mocks.MockAccountData()
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(fromAcc, nil).Times(1)

				toAcc := mocks.MockAccountData()
				toAcc.ID = 2
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(toAcc, nil).Times(1)

				mockTrans := mocks.MockTransferData(input)
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrans, nil).Times(1)

				mockFromEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: -input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockFromEnt, nil).Times(1)

				mockToEnt := &entry.Entry{AccountID: input.FromAccountID, Amount: input.Amount}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockToEnt, nil).Times(1)

				addFromAcc := mocks.MockAccountData()
				addFromAcc.Balance += -input.Amount
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.FromAccountID, -input.Amount).Return(addFromAcc, nil).Times(1)

				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.ToAccountID, input.Amount).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, tranRepo, accRepo, entRepo := setup(t)

			tc.mockFn(tranRepo, accRepo, entRepo, tc.input)

			ctx := auth.SetUserID(context.Background(), int64(10))

			result, err := uc.Transfer(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func setup(t *testing.T) (transferusecase.TransferUsecase, *transferrepository.MockTransferRepository, *accountrepository.MockAccountRepository, *entryrepository.MockEntryRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tranRepo := transferrepository.NewMockTransferRepository(ctrl)
	accRepo := accountrepository.NewMockAccountRepository(ctrl)
	entRepo := entryrepository.NewMockEntryRepository(ctrl)
	mockTx := &mocks.MockTx{}

	uc := transferusecase.NewTransferUsecase(tranRepo, accRepo, entRepo, mockTx)
	return uc, tranRepo, accRepo, entRepo
}
