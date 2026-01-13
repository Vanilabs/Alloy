package messaging

import (
	"alloy/internal/shared/database/models"
	"time"

	"github.com/gocql/gocql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"context"

	"sort"

	"sync"

	"github.com/google/uuid"
)



type Repository interface {
	Save(msg *models.Message) error
	ListByConversation(conversationID string, limit int) ([]models.Message, error)
	CreateConversation(ctx context.Context, convo *models.Conversation)(uuid.UUID, error)
	CreateChannel(ctx context.Context, channel *models.Channel) (uuid.UUID, error)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error) 
	GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) 
	AddMembersToChannel(ctx context.Context, members []models.ChannelMember,) error
	GetChannelMembers(ctx context.Context, channel_id uuid.UUID) ([]models.ChannelMember, error)
	GetUserConversations(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	FetchAllOfflineMessagesForUser(ctx context.Context,userID uuid.UUID, limit int,concurrency int,) ([]models.Message, error)
	UpdateLastRead(ctx context.Context, convo_read *models.ConversationRead,) error
}


type messageRepository struct {
	session *gocql.Session
	db *gorm.DB
}

func NewRepository(
	session *gocql.Session,
	db *gorm.DB,
) Repository {
	return &messageRepository{session: session, db: db}
}

func (r *messageRepository) Save(msg *models.Message) error {
	return r.session.Query(
		`INSERT INTO messages
		 (id, conversation_id, sender_id, text, timestamp)
		 VALUES (?, ?, ?, ?, ?)`,
		msg.ID,
		msg.ConversationID,
		msg.SenderID,
		msg.Text,
		msg.Timestamp,
	).Exec()
}

func (r *messageRepository) ListByConversation(
	conversationID string,
	limit int,
) ([]models.Message, error) {

	iter := r.session.Query(
		`SELECT timestamp, id, sender_id, text, conversation_id
		 FROM messages
		 WHERE conversation_id = ?
		 LIMIT ?`,
		conversationID,
		limit,
	).Iter()

	var messages []models.Message
	var m models.Message

	for iter.Scan(
		&m.Timestamp,
		&m.ID,
		&m.SenderID,
		&m.Text,
		&m.ConversationID,
	) {
		messages = append(messages, m)
	}

	return messages, iter.Close()
}

func (r *messageRepository) GetChannelByID(ctx context.Context, id uuid.UUID) (*models.Channel, error) {
	var channel models.Channel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&channel).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *messageRepository) GetChannelMembers(ctx context.Context, channel_id uuid.UUID) ([]models.ChannelMember, error) {
	var members []models.ChannelMember
	if err := r.db.WithContext(ctx).Where("channel_id = ?", channel_id).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}


func (r *messageRepository) GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	var convo models.Conversation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&convo).Error; err != nil {
		return nil, err
	}
	return &convo, nil
}

func (r *messageRepository) CreateChannel(ctx context.Context, channel *models.Channel) (uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(channel).Error
	if err != nil {
		return uuid.Nil, err
	}

	return channel.ID, nil
}

func (r *messageRepository) FetchMessagesAfter(ctx context.Context, conversationID uuid.UUID,
				after time.Time, limit int) ([]models.Message, error) {

	query := `
		SELECT id, conversation_id, sender_id, text, timestamp
		FROM messages
		WHERE conversation_id = ?
		  AND timestamp > ?
		LIMIT ?
	`
	parsedConversationID := ToGocqlUUID(conversationID)
	iter := r.session.Query(
		query,
		parsedConversationID,
		after,
		limit,
	).WithContext(ctx).Iter()

	var messages []models.Message
	var msg models.Message

	for iter.Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.Text,
		&msg.Timestamp,
	) {
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil	
				}


func (r *messageRepository) GetUserConversations(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var conversationIDs []uuid.UUID

	err := r.db.WithContext(ctx).
		Table("conversations").
		Select("conversations.id").
		Joins(`
			JOIN channel_members 
			  ON channel_members.channel_id = conversations.channel_id
		`).
		Where("channel_members.user_id = ?", userID).
		Scan(&conversationIDs).
		Error

	if err != nil {
		return nil, err
	}

	return conversationIDs, nil
}


func (r *messageRepository) UpdateLastRead(ctx context.Context, convo_read *models.ConversationRead,) error {

	var existing models.ConversationRead
	err := r.db.WithContext(ctx).
    	Where("user_id = ? AND conversation_id = ?", convo_read.UserID, convo_read.ConversationID).
    	First(&existing).Error

	if err == nil && existing.LastReadAt.After(convo_read.LastReadAt) {
    return nil
}

	return r.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "conversation_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"last_read_message_id": convo_read.LastReadMessageID,
				"last_read_at": convo_read.LastReadAt,
				"updated_at": time.Now(),
			}),
		},
	).Create(&convo_read).Error
}


func (r *messageRepository) getConversationReads(
	ctx context.Context,
	userID uuid.UUID,
	conversationIDs []uuid.UUID,
) (map[uuid.UUID]time.Time, error) {

	var reads []models.ConversationRead
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND conversation_id IN ?", userID, conversationIDs).
		Find(&reads).Error
	if err != nil {
		return nil, err
	}

	readMap := make(map[uuid.UUID]time.Time)
	for _, read := range reads {
		readMap[read.ConversationID] = read.LastReadAt
	}
	return readMap, nil
}



func (r *messageRepository) FetchAllOfflineMessagesForUser(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	concurrency int,
) ([]models.Message, error) {

	conversationIDs, err := r.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, err
	}


	readMap, err := r.getConversationReads(ctx, userID, conversationIDs)
	if err != nil {
		return nil, err
	}

	type convResult struct {
		messages []models.Message
		err      error
	}

	sem := make(chan struct{}, concurrency)
	resultsCh := make(chan convResult, len(conversationIDs))
	var wg sync.WaitGroup

	for _, convID := range conversationIDs {
		wg.Add(1)
		go func(convID uuid.UUID) {
			defer wg.Done()
			sem <- struct{}{}           // acquire
			defer func() { <-sem }()    // release

			lastReadAt, ok := readMap[convID]
			if !ok {
				lastReadAt = time.Time{} // never read
			}

			msgs, err := r.FetchMessagesAfter(ctx, convID, lastReadAt, limit)
			resultsCh <- convResult{
				messages: msgs,
				err:      err,
			}
		}(convID)
	}

	wg.Wait()
	close(resultsCh)

	var allMessages []models.Message
	for res := range resultsCh {
		if res.err != nil {
			return nil, res.err
		}
		allMessages = append(allMessages, res.messages...)
	}

	// 5. Sort by timestamp ascending
	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].Timestamp.Before(allMessages[j].Timestamp)
	})

	return allMessages, nil
}


func (r *messageRepository) CreateConversation(ctx context.Context, convo *models.Conversation)(uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(convo).Error
	if err != nil {
		return uuid.Nil, err
	}

	return convo.ID, nil
}

func (r *messageRepository) AddMembersToChannel(
	ctx context.Context,
	members []models.ChannelMember,
) error {
	return r.db.WithContext(ctx).Create(&members).Error
}
