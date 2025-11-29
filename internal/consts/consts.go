package consts

import "time"

const ContextTimeout = time.Second * 10

const (
	JWTAccessTokenDuration  = time.Hour * 24
	JWTRefreshTokenDuration = time.Hour * 24 * 7
)
