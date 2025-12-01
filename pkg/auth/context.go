package auth

import (
	"context"
	"errors"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/user"
	"github.com/codepnw/simple-bank/pkg/jwt"
)

func GetCurrentUser(ctx context.Context) (*user.User, error) {
	claims, ok := ctx.Value(consts.ContextUserClaimsKey).(*jwt.UserClaims)
	if !ok {
		return nil, errors.New("invalid user context")
	}

	usr := &user.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}
	return usr, nil
}

func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(consts.ContextUserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, consts.ContextUserIDKey, userID)
}
