package auth

import "errors"

var (
	ErrInvitationNotFound        = errors.New("invitation not found")
	ErrInvitationAlreadyVerified = errors.New("invitation already verified")
	ErrInvitationExpired         = errors.New("invitation expired")
)
