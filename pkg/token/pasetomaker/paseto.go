package pasetomaker

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/user"
	"github.com/codepnw/simple-bank/pkg/token"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (token.TokenMaker, error) {
	if len(symmetricKey) != 32 {
		return nil, fmt.Errorf("invalid key size: must be exactly 32 characters")
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (p *PasetoMaker) GenerateAccessToken(u *user.User) (string, error) {
	return p.generateToken(u, consts.JWTAccessTokenDuration)
}

func (p *PasetoMaker) GenerateRefreshToken(u *user.User) (string, error) {
	return p.generateToken(u, consts.JWTRefreshTokenDuration)
}

func (p *PasetoMaker) VerifyAccessToken(tokenStr string) (*token.Payload, error) {
	return p.verifyToken(tokenStr)
}

func (p *PasetoMaker) VerifyRefreshToken(tokenStr string) (*token.Payload, error) {
	return p.verifyToken(tokenStr)
}

func (p *PasetoMaker) generateToken(u *user.User, duration time.Duration) (string, error) {
	payload, err := token.NewTokenPayload(u, duration)
	if err != nil {
		return "", err
	}

	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	return token, err
}

func (p *PasetoMaker) verifyToken(tokenStr string) (*token.Payload, error) {
	payload := new(token.Payload)

	err := p.paseto.Decrypt(tokenStr, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
