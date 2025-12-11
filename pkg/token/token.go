package token

import (
	"errors"
	"time"

	"github.com/codepnw/simple-bank/internal/features/user"
)

type TokenMaker interface {
	GenerateAccessToken(u *user.User) (string, error)
	GenerateRefreshToken(u *user.User) (string, error)
	VerifyAccessToken(tokenStr string) (*Payload, error)
	VerifyRefreshToken(tokenStr string) (*Payload, error)
}

type Payload struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expires_at"`
}

func NewTokenPayload(u *user.User, duration time.Duration) (*Payload, error) {
	return &Payload{
		UserID:    u.ID,
		Email:     u.Email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}
