package auth

import (
	"alloy/internal/shared/cache"
	"alloy/internal/shared/constants"
	"alloy/internal/shared/database/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	GetInvitationByTokenAndEmail(ctx context.Context, token string, email string) (*models.Invitation, error)
	GetInvitationByEmail(ctx context.Context, email string) (*models.Invitation, error)
	GetInvitationByID(ctx context.Context, id uuid.UUID) (*models.Invitation, error)
	GetInvitations(ctx context.Context) ([]models.Invitation, error)
	UpdateInvitation(ctx context.Context, invitation *models.Invitation) error
	DeleteInvitation(ctx context.Context, id uuid.UUID) error
	Cache_SetMagicLinkToken(ctx context.Context, token string, data *cache.MagicLinkCacheData, ttl time.Duration) error
	Cache_GetMagicLinkToken(ctx context.Context, token string) (*cache.MagicLinkCacheData, error)
	Cache_DeleteMagicLinkToken(ctx context.Context, token string) error
}

type authRepository struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewRepository(store *constants.DataStores) Repository {
	return &authRepository{db: store.PostGres, redis: store.Redis}
}

func (r *authRepository) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *authRepository) GetInvitationByTokenAndEmail(ctx context.Context, token string, email string) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := r.db.WithContext(ctx).Where("token = ? AND email = ?", token, email).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *authRepository) GetInvitationByEmail(ctx context.Context, email string) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *authRepository) GetInvitationByID(ctx context.Context, id uuid.UUID) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *authRepository) GetInvitations(ctx context.Context) ([]models.Invitation, error) {
	var invitations []models.Invitation
	if err := r.db.WithContext(ctx).Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *authRepository) UpdateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *authRepository) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Invitation{}, "id = ?", id).Error
}

func (r *authRepository) Cache_SetMagicLinkToken(ctx context.Context, token string, data *cache.MagicLinkCacheData, ttl time.Duration) error {
	key := fmt.Sprintf("magic_link:%s", token)
	json, err := json.Marshal(data)
	if err != nil {
		r.logger.Error("Failed to marshal magic link data", zap.Error(err))
		return err
	}

	err = r.redis.Set(ctx, key, json, ttl).Err()
	if err != nil {
		r.logger.Error("Failed to set magic link token", zap.Error(err))
		return err
	}
	r.logger.Debug("User magic link token stored",
		zap.String("user_id", data.UserID),
		zap.String("email", data.Email),
		zap.Duration("ttl", ttl))

	return nil
}

func (r *authRepository) Cache_GetMagicLinkToken(ctx context.Context, token string) (*cache.MagicLinkCacheData, error) {
	key := fmt.Sprintf("magic_link:%s", token)
	jsonData, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var data cache.MagicLinkCacheData
	err = json.Unmarshal([]byte(jsonData), &data)
	return &data, nil
}

func (r *authRepository) Cache_DeleteMagicLinkToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("magic_link:%s", token)
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Failed to delete magic link token", zap.Error(err))
		return err
	}
	r.logger.Debug("Magic link token deleted", zap.String("token", token))
	return nil
}
