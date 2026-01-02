package jwtmaker

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/user"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/token"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	secretKey  string
	refreshKey string
}

func NewJWTMaker(cfg *config.JWTConfig) (token.TokenMaker, error) {
	if cfg.SecretKey == "" || cfg.RefreshKey == "" {
		return nil, errors.New("key is empty")
	}
	return &JWTToken{
		secretKey:  cfg.SecretKey,
		refreshKey: cfg.RefreshKey,
	}, nil
}

type userClaims struct {
	UserID int64
	Email  string
	*jwt.RegisteredClaims
}

func (j *JWTToken) GenerateAccessToken(u *user.User) (string, error) {
	return j.generateToken(j.secretKey, u, consts.TokenAccessDuration)
}

func (j *JWTToken) GenerateRefreshToken(u *user.User) (string, error) {
	return j.generateToken(j.refreshKey, u, consts.TokenRefreshDuration)
}

func (j *JWTToken) generateToken(key string, u *user.User, duration time.Duration) (string, error) {
	claims := &userClaims{
		UserID: u.ID,
		Email:  u.Email,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "simple-bank-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("signed token failed: %w", err)
	}
	return ss, nil
}

func (j *JWTToken) VerifyAccessToken(tokenStr string) (*token.Payload, error) {
	return j.verifyToken(j.secretKey, tokenStr)
}

func (j *JWTToken) VerifyRefreshToken(tokenStr string) (*token.Payload, error) {
	return j.verifyToken(j.refreshKey, tokenStr)
}

func (j *JWTToken) verifyToken(key, tokenStr string) (*token.Payload, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*userClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	payload := &token.Payload{
		UserID:    claims.UserID,
		Email:     claims.Email,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
