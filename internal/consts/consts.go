package consts

import "time"

// Token Duration
const (
	TokenAccessDuration  = time.Hour * 1
	TokenRefreshDuration = time.Hour * 24 * 7
)

type contextKey string

// Context Config
const (
	ContextTimeout = time.Second * 10

	// Context Key
	ContextUserClaimsKey contextKey = "user-claims"
	ContextUserIDKey     contextKey = "user-id"
)

const ParamAccountID = "account_id"
