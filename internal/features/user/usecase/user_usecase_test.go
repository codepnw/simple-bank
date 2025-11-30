package userusecase_test

import (
	"context"
	"testing"

	"github.com/codepnw/simple-bank/internal/features/user"
	userrepository "github.com/codepnw/simple-bank/internal/features/user/repository"
	userusecase "github.com/codepnw/simple-bank/internal/features/user/usecase"
	"github.com/codepnw/simple-bank/internal/mocks"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/jwt"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/password"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	type testCase struct {
		name        string
		input       *user.User
		mockFn      func(mockRepo *userrepository.MockUserRepository, input *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			input: &user.User{
				Username:  "testmock001",
				FirstName: "john",
				LastName:  "cena",
				Email:     "mock@example.com",
				Password:  "password",
			},
			mockFn: func(mockRepo *userrepository.MockUserRepository, input *user.User) {
				u := mocks.MockUserData()
				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), input).Return(u, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail insert user",
			input: &user.User{
				Username:  "testmock001",
				FirstName: "john",
				LastName:  "cena",
				Email:     "mock@example.com",
			},
			mockFn: func(mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), input).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, repo, _ := setup(t)

			tc.mockFn(repo, tc.input)

			result, err := uc.Register(context.Background(), tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name        string
		input       *user.User
		mockFn      func(mockRepo *userrepository.MockUserRepository, input *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			input: &user.User{
				Email:    "mock@example.com",
				Password: "password",
			},
			mockFn: func(mockRepo *userrepository.MockUserRepository, input *user.User) {
				hashedPassword, _ := password.HashedPassword(input.Password)
				u := mocks.MockUserData()
				u.Password = hashedPassword

				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(u, nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail user credentials",
			input: &user.User{
				Email:    "mock@example.com",
				Password: "wrong-password",
			},
			mockFn: func(mockRepo *userrepository.MockUserRepository, input *user.User) {
				u := mocks.MockUserData()
				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(u, nil).Times(1)
			},
			expectedErr: errs.ErrInvalidCredentials,
		},
		{
			name: "fail db error",
			input: &user.User{
				Email:    "mock@example.com",
				Password: "password",
			},
			mockFn: func(mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockRepo.EXPECT().FindByEmail(gomock.Any(), input.Email).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, repo, _ := setup(t)

			tc.mockFn(repo, tc.input)

			result, err := uc.Login(context.Background(), tc.input.Email, tc.input.Password)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type testCase struct {
		name        string
		token       string
		mockFn      func(mockRepo *userrepository.MockUserRepository, token string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				u := mocks.MockUserData()
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), token).Return(u.ID, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), u.ID).Return(u, nil).Times(1)

				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail not found",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), token).Return(int64(0), errs.ErrTokenNotFound).Times(1)
			},
			expectedErr: errs.ErrTokenNotFound,
		},
		{
			name: "fail db error",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				u := mocks.MockUserData()
				mockRepo.EXPECT().ValidateRefreshToken(gomock.Any(), token).Return(u.ID, nil).Times(1)

				mockRepo.EXPECT().FindByID(gomock.Any(), u.ID).Return(u, nil).Times(1)

				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(nil).Times(1)

				mockRepo.EXPECT().SaveRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, repo, _ := setup(t)

			tc.mockFn(repo, tc.token)

			result, err := uc.RefreshToken(context.Background(), tc.token)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	type testCase struct {
		name        string
		token       string
		mockFn      func(mockRepo *userrepository.MockUserRepository, token string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "fail db error",
			token: "mock_refresh_token",
			mockFn: func(mockRepo *userrepository.MockUserRepository, token string) {
				mockRepo.EXPECT().RevokedRefreshToken(gomock.Any(), gomock.Any(), token).Return(mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, repo, _ := setup(t)

			tc.mockFn(repo, tc.token)

			err := uc.Logout(context.Background(), tc.token)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func setup(t *testing.T) (userusecase.UserUsecase, *userrepository.MockUserRepository, mocks.MockDB) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := userrepository.NewMockUserRepository(ctrl)
	mockToken, err := jwt.InitJWT(&config.JWTConfig{
		SecretKey:  "mock-secret-key",
		RefreshKey: "mock-refresh-key",
	})
	if err != nil {
		t.Fatalf("init jwt failed: %v", err)
	}
	mockDB := mocks.MockDB{}
	mockTx := mocks.MockTx{}

	uc := userusecase.NewUserUsecase(mockRepo, mockToken, &mockTx)
	return uc, mockRepo, mockDB
}
