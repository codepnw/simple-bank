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
				mockRepo.EXPECT().Insert(gomock.Any(), input).Return(u, nil).Times(1)
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
				mockRepo.EXPECT().Insert(gomock.Any(), input).Return(nil, mocks.ErrDatabase).Times(1)
			},
			expectedErr: mocks.ErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc, repo := setup(t)

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
			uc, repo := setup(t)

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

func setup(t *testing.T) (userusecase.UserUsecase, *userrepository.MockUserRepository) {
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

	uc := userusecase.NewUserUsecase(mockRepo, mockToken)
	return uc, mockRepo
}
