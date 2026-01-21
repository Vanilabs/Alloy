package models

import "errors"

var (
	ErrMissingAuthorization                = errors.New("authorization token missing")
	ErrInvalidToken                        = errors.New("unable to authenticate token")
	ErrMissingToken                        = errors.New("no token found")
	ErrInvalidAuthorizationHeaderFormat    = errors.New("invalid authorization header format")
	ErrMissingOrInvalidAuthorizationHeader = errors.New("missing or invalid authorization header")

	// JWT Token
	ErrInvalidTokenFormat          = errors.New("invalid token format")
	ErrorInvalidORExpiredToken     = errors.New("invalid or expired token")
	ErrSessionExpiredOrInvalidated = errors.New("session expired or invalidated")
	ErrInvalidOrExpiredToken       = errors.New("invalid or expired token")
	ErrTokenRevoked                = errors.New("token revoked")
)
