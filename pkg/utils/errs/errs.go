package errs

import "errors"

// User
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
)

// Account
var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrCurrencyAlreadyExists = errors.New("account with this currency already exists")
	ErrInvalidCurrency       = errors.New("invalid account currency ['THB', 'USD']")
)

// Auth
var (
	ErrNoUserID      = errors.New("no user id")
	ErrNoPermission  = errors.New("no permission")
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenRevoked  = errors.New("token revoked")
	ErrTokenExpires  = errors.New("token is expires")
	ErrUnauthorized  = errors.New("unauthorized")
)

// Transfer
var (
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrMoneyNotEnough   = errors.New("money not enough")
	ErrTransferToSelf   = errors.New("transfer to self")
)
