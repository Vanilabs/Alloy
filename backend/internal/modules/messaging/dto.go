package messaging

import (

	"alloy/internal/shared/constants"
	"github.com/google/uuid"

	"time"

)

type ChatDTO struct {
	Text        string `json:"text"`
	Timestamp      time.Time  `json:"timestamp"`
	ConversationID uuid.UUID `json:"conversation_id"`
}

type UpdateLastReadDTO struct {
	UserID         uuid.UUID `json:"user_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	MessageID      uuid.UUID `json:"message_id"`
	ReadTime       time.Time `json:"read_time"`
}

type channelMember struct {
	Email string `json:"email"`
	Role constants.MessageChannelRole `json:"role"`
}

type CreateChannelDTO struct {
	Type constants.MessageChannelType `json:"channel_type"`
	Name *string		`json:"channel_name,omitempty"`
	Description *string		`json:"channel_description,omitempty"`
	Members []channelMember	`json:"channel_members"`
} 