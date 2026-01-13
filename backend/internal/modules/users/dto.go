package users

import "time"

type CreateUserRequest struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	RoleAtMBL   string     `json:"role_at_mbl"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	State       *string    `json:"state,omitempty"`
	Department  string     `json:"department"`
}
