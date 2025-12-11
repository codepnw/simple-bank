package userusecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/user"
	userrepository "github.com/codepnw/simple-bank/internal/features/user/repository"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/codepnw/simple-bank/pkg/token"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/helper"
)

type UserUsecase interface {
	Register(ctx context.Context, input *user.User) (*TokenResponse, error)
	Login(ctx context.Context, email, pwd string) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type userUsecase struct {
	repo  userrepository.UserRepository
	token token.TokenMaker
	tx    database.TxManager
}

func NewUserUsecase(repo userrepository.UserRepository, token token.TokenMaker, tx database.TxManager) UserUsecase {
	return &userUsecase{
		repo:  repo,
		token: token,
		tx:    tx,
	}
}

func (u *userUsecase) Register(ctx context.Context, input *user.User) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	hashedPassword, err := helper.HashedPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashedPassword

	var response *TokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Insert User
		userData, err := u.repo.Insert(ctx, tx, input)
		if err != nil {
			return err
		}

		// Generate Token
		resp, err := u.generateToken(userData)
		if err != nil {
			return err
		}

		// Save Refresh Token
		err = u.repo.SaveRefreshToken(ctx, tx, &user.Auth{
			UserID:    userData.ID,
			Token:     resp.RefreshToken,
			ExpiresAt: time.Now().Add(consts.JWTRefreshTokenDuration),
		})
		if err != nil {
			return err
		}

		response = resp
		return nil
	})
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

	if err = helper.ComparePassword(userData.Password, pwd); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	var response *TokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Generate Token
		resp, err := u.generateToken(userData)
		if err != nil {
			return err
		}

		// Save Refresh Token
		err = u.repo.SaveRefreshToken(ctx, tx, &user.Auth{
			UserID:    userData.ID,
			Token:     resp.RefreshToken,
			ExpiresAt: time.Now().Add(consts.JWTRefreshTokenDuration),
		})
		if err != nil {
			return err
		}

		response = resp
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *userUsecase) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	userID, err := u.repo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	userData, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response *TokenResponse
	err = u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Revoked Refresh Token
		if err = u.repo.RevokedRefreshToken(ctx, tx, refreshToken); err != nil {
			return nil
		}

		// Generate Token
		resp, err := u.generateToken(userData)
		if err != nil {
			return nil
		}

		// Save Refresh Token
		err = u.repo.SaveRefreshToken(ctx, tx, &user.Auth{
			UserID:    userData.ID,
			Token:     resp.RefreshToken,
			ExpiresAt: time.Now().Add(consts.JWTRefreshTokenDuration),
		})
		if err != nil {
			return err
		}

		response = resp
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *userUsecase) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	err := u.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := u.repo.RevokedRefreshToken(ctx, tx, refreshToken); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
