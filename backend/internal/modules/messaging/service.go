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
	"sort"
	"sync"

	"github.com/gocql/gocql"
	"go.uber.org/zap"

	"github.com/google/uuid"

	"time"
)

type MemberInfo struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

type LastMessagePreview struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	SenderID  uuid.UUID `json:"sender_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ActiveConversation struct {
	ConversationID uuid.UUID           `json:"conversation_id"`
	ChannelID      uuid.UUID           `json:"channel_id"`
	ChannelType    constants.MessageChannelType `json:"channel_type"`
	ChannelName    *string             `json:"channel_name,omitempty"`
	Members        []MemberInfo        `json:"members"`
	LastMessage    *LastMessagePreview `json:"last_message,omitempty"`
	UnreadCount    int                 `json:"unread_count"`
	LastActivityAt time.Time           `json:"last_activity_at"`
}

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
	GetActiveConversations(ctx context.Context, userID uuid.UUID) ([]ActiveConversation, error)
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

func (ms *messagingService) GetActiveConversations(ctx context.Context, userID uuid.UUID) ([]ActiveConversation, error) {
	convsWithChannels, err := ms.messageRepository.GetConversationsWithChannels(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(convsWithChannels) == 0 {
		return []ActiveConversation{}, nil
	}

	channelIDs := make([]uuid.UUID, len(convsWithChannels))
	conversationIDs := make([]uuid.UUID, len(convsWithChannels))
	for i, c := range convsWithChannels {
		channelIDs[i] = c.ChannelID
		conversationIDs[i] = c.ConversationID
	}

	membersWithUsers, err := ms.messageRepository.GetChannelMembersWithUsers(ctx, channelIDs)
	if err != nil {
		return nil, err
	}
	membersByChannel := make(map[uuid.UUID][]MemberInfo)
	for _, m := range membersWithUsers {
		membersByChannel[m.ChannelID] = append(membersByChannel[m.ChannelID], MemberInfo{
			UserID:    m.UserID,
			FirstName: m.FirstName,
			LastName:  m.LastName,
			Email:     m.Email,
			Role:      m.Role,
		})
	}

	readMap, err := ms.messageRepository.GetConversationReads(ctx, userID, conversationIDs)
	if err != nil {
		return nil, err
	}

	type convData struct {
		conversationID uuid.UUID
		lastMsg        *models.Message
		unreadCount    int
		err            error
	}

	resultsCh := make(chan convData, len(convsWithChannels))
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for _, conv := range convsWithChannels {
		wg.Add(1)
		go func(convID uuid.UUID) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			lastMsg, err := ms.messageRepository.GetLastMessageInConversation(ctx, convID)
			if err != nil {
				resultsCh <- convData{conversationID: convID, err: err}
				return
			}

			lastReadAt := readMap[convID]
			unread, err := ms.messageRepository.CountMessagesAfter(ctx, convID, lastReadAt)
			if err != nil {
				resultsCh <- convData{conversationID: convID, err: err}
				return
			}

			resultsCh <- convData{conversationID: convID, lastMsg: lastMsg, unreadCount: unread}
		}(conv.ConversationID)
	}

	wg.Wait()
	close(resultsCh)

	convDataMap := make(map[uuid.UUID]convData)
	for r := range resultsCh {
		if r.err != nil {
			return nil, r.err
		}
		convDataMap[r.conversationID] = r
	}

	activeConvos := make([]ActiveConversation, 0, len(convsWithChannels))
	for _, conv := range convsWithChannels {
		data := convDataMap[conv.ConversationID]
		item := ActiveConversation{
			ConversationID: conv.ConversationID,
			ChannelID:      conv.ChannelID,
			ChannelType:    constants.MessageChannelType(conv.ChannelType),
			ChannelName:    conv.ChannelName,
			Members:        membersByChannel[conv.ChannelID],
			UnreadCount:    data.unreadCount,
			LastActivityAt: conv.CreatedAt,
		}
		if data.lastMsg != nil {
			item.LastMessage = &LastMessagePreview{
				ID:        GocqlToUUID(data.lastMsg.ID),
				Text:      data.lastMsg.Text,
				SenderID:  GocqlToUUID(data.lastMsg.SenderID),
				Timestamp: data.lastMsg.Timestamp,
			}
			item.LastActivityAt = data.lastMsg.Timestamp
		}
		activeConvos = append(activeConvos, item)
	}

	sort.Slice(activeConvos, func(i, j int) bool {
		return activeConvos[i].LastActivityAt.After(activeConvos[j].LastActivityAt)
	})

	return activeConvos, nil
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

