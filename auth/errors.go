package auth

import "errors"

var (
	ErrTokenNotSet = errors.New("a security token is not set")

	ErrAccessDenied = errors.New("access denied")
)
