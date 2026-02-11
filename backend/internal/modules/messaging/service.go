package messaging

import (
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/socket"

	"alloy/internal/modules/users"

	"alloy/internal/shared/constants"

	"alloy/internal/shared/router"

	"context"

	"errors"

	"fmt"
	"encoding/json"
	
	"github.com/gocql/gocql"
	"go.uber.org/zap"

	"github.com/google/uuid"

	"time"
)

type Service interface {
	SaveMessage (text string, conversation_id uuid.UUID, sender_id uuid.UUID, timestamp time.Time) (*gocql.UUID, error)
	ListConversationMessages(conversation_id string, limit int) ([]models.Message, error)
	CheckUserOnline(userID string) (bool, error)
	GetOfflineMessages(ctx context.Context, userID uuid.UUID, limit int,) ([]models.Message, error)
	SendMessage(ctx context.Context, senderSocketID string, msg MessageMetadata) error 
	InitiateChat(ctx context.Context, cfg CreateChannelDTO) (*ConversationMetadata, error) 
	ValidateConversation(ctx context.Context, conversation_id uuid.UUID) (*ConversationMetadata, error)
	GetChannelMembers(ctx context.Context, channel_id uuid.UUID) ([]models.ChannelMember, error)
	UpdateLastReadMessage(ctx context.Context, userID, conversationID, messageID uuid.UUID, readTime time.Time) error
}

type ConversationMetadata struct {
	ID uuid.UUID
	channelId uuid.UUID
}

type ChannelMetadata struct {
	Type constants.MessageChannelType
	Name *string
	Members []channelMember

}

type MessageMetadata struct {
	ID uuid.UUID
	SenderID uuid.UUID
	RecipientID uuid.UUID
	ConversationID uuid.UUID
	Text       string 
	Timestamp      time.Time
}

type messagingService struct {
	messageRepository Repository
	userRepository users.Repository
	socketManager *socket.ConnectionManager
	env     *router.Environment
}


func ToGocqlUUID(id uuid.UUID) gocql.UUID {
	return gocql.UUID(id)
}

func GocqlToUUID(gid gocql.UUID) uuid.UUID {
	return uuid.UUID(gid)
}

func NewService(ms_repository Repository, user_repository users.Repository,
	socketManager *socket.ConnectionManager, env *router.Environment) Service {
	return &messagingService{
		messageRepository: ms_repository,
		userRepository: user_repository,
		socketManager: socketManager,
		env: env,
	}
}

func (ms *messagingService) SaveMessage(text string, conversation_id uuid.UUID, 
				sender_id uuid.UUID, timestamp time.Time) (*gocql.UUID, error) {
	msg := &models.Message{
		ID: gocql.TimeUUID(),
		ConversationID: ToGocqlUUID(conversation_id),		
		SenderID: ToGocqlUUID(sender_id),
		Text: text,
		Timestamp: timestamp,
	}
	if err:= ms.messageRepository.Save(msg); err != nil{
		return nil, err
	}

	return &msg.ID, nil
}

func (ms *messagingService) ListConversationMessages(conversation_id string, limit int) ([]models.Message, error) {
	return  ms.messageRepository.ListByConversation(conversation_id, limit)
}

func (ms *messagingService) CheckUserOnline(userID string) (bool, error) {
	active_sockets, err := ms.socketManager.SocketTracker.GetUserSockets(userID)
	if err != nil {
		return false, err
	}
	if len(active_sockets) == 0{
		return false, nil
	}
	ms.env.Logger.Info("User online ==> ", zap.String("user id", userID))
	ms.env.Logger.Info("User Active connections ==> ", zap.Int("sockets", len(active_sockets)))
	return true, nil
}

func (ms *messagingService) ValidateConversation(ctx context.Context, id uuid.UUID) (*ConversationMetadata, error) {
	conversation, err := ms.messageRepository.GetConversationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, nil
	}
	meta := &ConversationMetadata{
		ID: id,
		channelId: conversation.ChannelID,
	}
	return meta, nil


}

func (ms *messagingService) UpdateLastReadMessage(ctx context.Context, userID, 
			conversationID, messageID uuid.UUID, readTime time.Time) error {
	read := &models.ConversationRead{
		UserID:            userID,
		ConversationID:    conversationID,
		LastReadMessageID: messageID,
		LastReadAt:        readTime,
		UpdatedAt:         time.Now(),
	}
	return ms.messageRepository.UpdateLastRead(ctx, read)
			}


func (ms *messagingService)  GetOfflineMessages(
	ctx context.Context,
    userID uuid.UUID,
	limit int,
) ([]models.Message, error) {


   return ms.messageRepository.FetchAllOfflineMessagesForUser(ctx, userID, limit, 10)

}

func (ms *messagingService) SendMessage(ctx context.Context, senderSocketID string, msg MessageMetadata) error {
	payload, err := json.Marshal(msg)
    if err != nil {
        return err
    }
	if msg.SenderID == msg.RecipientID {
		return ms.socketManager.BroadcastToSelf(senderSocketID, msg.SenderID.String(), payload)
	}
	targetOnline, err := ms.CheckUserOnline(msg.RecipientID.String()) 
	if err != nil {
		return err
	}
	if !targetOnline {
		ms.env.Logger.Info("Recipient Not Online")
		return nil
		}
	
	ms.env.Logger.Info("Recipient Online")

	return ms.socketManager.BroadcastToAnotherUser(msg.RecipientID.String(), payload)
	
}

func (ms *messagingService) GetChannelMembers(ctx context.Context, channel_id uuid.UUID) ([]models.ChannelMember, error) {
	return ms.messageRepository.GetChannelMembers(ctx, channel_id)
}

func (ms *messagingService) InitiateChat(ctx context.Context, channel_config CreateChannelDTO) (*ConversationMetadata, error) {

	valid_members := make([]models.ChannelMember, 0, len(channel_config.Members))
	for _, prospect := range channel_config.Members {
		user, err := ms.userRepository.GetUserByEmail(ctx, prospect.Email)
		if err != nil {
			return nil, fmt.Errorf(
			"%w: %v",
			errors.New("email not valid"),
			prospect.Email,
		)
		}
		valid_members = append(valid_members, models.ChannelMember{
			UserID:    user.ID,
			Role:      prospect.Role,
		})
	}

	channel := &models.Channel{
			Type: channel_config.Type,
			Name: channel_config.Name,
			Description: channel_config.Description,
	}
	channel_id, err := ms.messageRepository.CreateChannel(ctx, channel)
	if err != nil {
		return nil, err
	}
	conversation := &models.Conversation{
		ChannelID: channel_id,
	}
	conversation_id, err := ms.messageRepository.CreateConversation(ctx, conversation)
	if err != nil {
		return nil, err
	}

	for i := range valid_members {
    	valid_members[i].ChannelID = channel_id
		}

	err = ms.messageRepository.AddMembersToChannel(ctx, valid_members)
	if err != nil {
		return nil, err
	}
	return &ConversationMetadata{
		ID: conversation_id,
		channelId: channel_id,
	}, nil
}

