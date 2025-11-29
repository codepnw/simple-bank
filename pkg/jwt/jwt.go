package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/user"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	secretKey  string
	refreshKey string
}

func InitJWT(cfg *config.JWTConfig) (*JWTToken, error) {
	if cfg.SecretKey == "" || cfg.RefreshKey == "" {
		return nil, errors.New("key is empty")
	}
	return &JWTToken{
		secretKey:  cfg.SecretKey,
		refreshKey: cfg.RefreshKey,
	}, nil
}

type UserClaims struct {
	UserID int64
	Email  string
	*jwt.RegisteredClaims
}

func (j *JWTToken) GenerateAccessToken(u *user.User) (string, error) {
	return j.generateToken(j.secretKey, u, consts.JWTAccessTokenDuration)
}

func (j *JWTToken) GenerateRefreshToken(u *user.User) (string, error) {
	return j.generateToken(j.refreshKey, u, consts.JWTRefreshTokenDuration)
}

func (j *JWTToken) generateToken(key string, u *user.User, duration time.Duration) (string, error) {
	claims := &UserClaims{
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
