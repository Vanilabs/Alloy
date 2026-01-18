package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName      string     `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName       string     `gorm:"type:varchar(255);not null" json:"last_name"`
	RoleAtOrg      string     `gorm:"type:varchar(255);not null" json:"role_at_org"`
	Role           string     `gorm:"type:varchar(50);not null" json:"role"`
	Email          string     `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone          string     `gorm:"type:varchar(255);not null;unique" json:"phone"`
	DateOfBirth    *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	State          *string    `gorm:"type:varchar(100)" json:"state,omitempty"`
	Department     string     `gorm:"type:varchar(255);not null" json:"department"`
	EmployeeNumber string     `gorm:"type:varchar(50);unique" json:"employee_number"`
	CreatedAt      time.Time  `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}
