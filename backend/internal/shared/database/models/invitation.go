package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusExpired  = "expired"
)

type Invitation struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName  string     `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName   string     `gorm:"type:varchar(255);not null" json:"last_name"`
	Email      string     `gorm:"type:varchar(255);not null" json:"email"`
	Phone      string     `gorm:"type:varchar(255);not null" json:"phone"`
	Department string     `gorm:"type:varchar(255);not null" json:"department"`
	Role       string     `gorm:"type:varchar(255);not null" json:"role"`
	RoleAtOrg  string     `gorm:"type:varchar(255);not null" json:"role_at_org"`
	Token      string     `gorm:"type:varchar(255);not null" json:"token"`
	InvitedBy  uuid.UUID  `gorm:"type:uuid;not null" json:"invited_by"`
	ExpiresAt  time.Time  `gorm:"type:timestamp;not null" json:"expires_at"`
	AcceptedAt *time.Time `gorm:"type:timestamp" json:"accepted_at,omitempty"`
	Status     string     `gorm:"type:varchar(255);not null;default:'pending'" json:"status"`
	CreatedAt  time.Time  `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}
