package transferusecase_test

import (
	"context"
	"testing"

	"github.com/codepnw/simple-bank/internal/features/account"
	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	"github.com/codepnw/simple-bank/internal/features/entry"
	entryrepository "github.com/codepnw/simple-bank/internal/features/entry/repository"
	"github.com/codepnw/simple-bank/internal/features/transfer"
	transferrepository "github.com/codepnw/simple-bank/internal/features/transfer/repository"
	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
	"github.com/codepnw/simple-bank/internal/mocks"
	"github.com/codepnw/simple-bank/pkg/auth"
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
			name: "success",
			input: &transferusecase.TransferParams{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        10,
				Currency:      "THB",
			},
			mockFn: func(tranRepo *transferrepository.MockTransferRepository, accRepo *accountrepository.MockAccountRepository, entRepo *entryrepository.MockEntryRepository, input *transferusecase.TransferParams) {
				mockFromAcc := &account.Account{
					ID:       1,
					OwnerID:  10,
					Balance:  100,
					Currency: "THB",
				}
				accRepo.EXPECT().FindByID(gomock.Any(), input.FromAccountID).Return(mockFromAcc, nil).Times(1)

				mockToAcc := &account.Account{
					ID:       2,
					OwnerID:  11,
					Balance:  0,
					Currency: "THB",
				}
				accRepo.EXPECT().FindByID(gomock.Any(), input.ToAccountID).Return(mockToAcc, nil).Times(1)

				mockTransInput := &transfer.Transfer{
					FromAccountID: input.FromAccountID,
					ToAccountID:   input.ToAccountID,
					Amount:        input.Amount,
				}
				mockTrans := &transfer.Transfer{
					ID:            1,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        10,
				}
				tranRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), mockTransInput).Return(mockTrans, nil).Times(1)

				mockEntFrom := &entry.Entry{
					AccountID: 1,
					Amount:    -10,
				}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockEntFrom, nil).Times(1)

				mockEntTo := &entry.Entry{
					AccountID: 2,
					Amount:    10,
				}
				entRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockEntTo, nil).Times(1)

				mockAcc1 := &account.Account{
					ID:      1,
					Balance: 90,
				}
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.FromAccountID, -input.Amount).Return(mockAcc1, nil).Times(1)

				mockAcc2 := &account.Account{
					ID:      2,
					Balance: 10,
				}
				accRepo.EXPECT().AddAccountBalance(gomock.Any(), gomock.Any(), input.ToAccountID, input.Amount).Return(mockAcc2, nil).Times(1)
			},
			expectedErr: nil,
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
