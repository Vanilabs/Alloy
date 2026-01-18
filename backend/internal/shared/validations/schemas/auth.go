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

type InviteUserSchema struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Role       string `json:"role" validate:"required,oneof=admin user"`
	RoleAtOrg  string `json:"role_at_org" validate:"required"`
	Department string `json:"department" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
}

func (InviteUserSchema) Messages() validations.FieldMessageOverride {
	return validations.FieldMessageOverride{
		"FirstName": {
			"required": "First name is required",
		},
		"LastName": {
			"required": "Last name is required",
		},
		"Email": {
			"required": "Email is required",
			"email":    "Email must be a valid email address",
		},
		"Role": {
			"required": "Role is required",
		},
		"RoleAtOrg": {
			"required": "Role at organization is required",
		},
		"Department": {
			"required": "Department is required",
		},
		"Phone": {
			"required": "Phone is required",
		},
	}
}
