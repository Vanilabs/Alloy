package models

import (
	"github.com/gocql/gocql"

	"time"

	"github.com/google/uuid"

	"alloy/internal/shared/constants"

)

type Message struct {
	ID gocql.UUID
	ConversationID gocql.UUID
	SenderID gocql.UUID
	Text string
	Attachments *string
	Timestamp time.Time
}

type Channel struct {
	ID uuid.UUID	`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Type 	constants.MessageChannelType `gorm:"type:varchar(10);not null"`
	Name *string `gorm:"type:varchar(100)"`
	Description *string `gorm:"type:varchar(400)"`
	CreatedAt       time.Time	`gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	Members []User `gorm:"many2many:channel_members;"`
}

type ChannelMember struct {
	ChannelID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`

	Role      constants.MessageChannelRole    `gorm:"not null;default:member"`
	JoinedAt  time.Time `gorm:"autoCreateTime"`
}

type Conversation  struct {
	ID  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ChannelID		uuid.UUID   `gorm:"type:uuid;not null" json:"channel_id"`
	CreatedAt       time.Time	`gorm:"type:timestamp;not null;default:now()" json:"created_at"`
}

type ConversationRead struct {
	ConversationID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;primaryKey"`

	LastReadMessageID uuid.UUID `gorm:"type:uuid;not null"`
	LastReadAt        time.Time `gorm:"not null"`

	LastDeliveredAt *time.Time `gorm:""`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}