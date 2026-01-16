package schemas

import "alloy/internal/shared/validations"

type RequestMagicLinkSchema struct {
	Email string `json:"email" validate:"required,email"`
}

func (RequestMagicLinkSchema) Messages() validations.FieldMessageOverride {
	return validations.FieldMessageOverride{
		"Email": {
			"required": "Email is required",
			"email":    "Email must be a valid email address",
		},
	}
}

type MagicLinkVerifySchema struct {
	Token string `json:"token" validate:"required,min=32"`
}

func (MagicLinkVerifySchema) Messages() validations.FieldMessageOverride {
	return validations.FieldMessageOverride{
		"Token": {
			"required": "Token is required",
			"min":      "Token must be at least 32 characters",
		},
	}
}
