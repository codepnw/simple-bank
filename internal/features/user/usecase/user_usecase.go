package userusecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/user"
	userrepository "github.com/codepnw/simple-bank/internal/features/user/repository"
	"github.com/codepnw/simple-bank/pkg/jwt"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/password"
)

type UserUsecase interface {
	Register(ctx context.Context, input *user.User) (*TokenResponse, error)
	Login(ctx context.Context, email, pwd string) (*TokenResponse, error)
}

type userUsecase struct {
	repo  userrepository.UserRepository
	token *jwt.JWTToken
}

func NewUserUsecase(repo userrepository.UserRepository, token *jwt.JWTToken) UserUsecase {
	return &userUsecase{
		repo:  repo,
		token: token,
	}
}

func (u *userUsecase) Register(ctx context.Context, input *user.User) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	hashedPassword, err := password.HashedPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashedPassword

	userData, err := u.repo.Insert(ctx, input)
	if err != nil {
		return nil, err
	}

	// Generate Token
	response, err := u.generateToken(userData)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *userUsecase) Login(ctx context.Context, email, pwd string) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userData, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}

	if err = password.ComparePassword(userData.Password, pwd); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	// Generate Token
	response, err := u.generateToken(userData)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *userUsecase) generateToken(usr *user.User) (*TokenResponse, error) {
	accessToken, err := u.token.GenerateAccessToken(usr)
	if err != nil {
		return nil, fmt.Errorf("gen access token failed: %w", err)
	}
	refreshToken, err := u.token.GenerateRefreshToken(usr)
	if err != nil {
		return nil, fmt.Errorf("gen refresh token failed: %w", err)
	}

	response := &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}
