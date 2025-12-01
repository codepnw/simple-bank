package consts

import "time"

// Token Duration
const (
	JWTAccessTokenDuration  = time.Hour * 24
	JWTRefreshTokenDuration = time.Hour * 24 * 7
)

type contextKey string

const (
	ContextTimeout = time.Second * 10

	// Context Key
	ContextUserClaimsKey contextKey = "user-claims"
	ContextUserIDKey     contextKey = "user-id"
)
