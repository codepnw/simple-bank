package errs

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrTokenNotFound         = errors.New("token not found")
	ErrTokenRevoked          = errors.New("token revoked")
	ErrTokenExpires          = errors.New("token is expires")
	ErrUnauthorized          = errors.New("unauthorized")
)
