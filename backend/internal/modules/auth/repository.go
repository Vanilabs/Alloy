package auth

import (
	"alloy/internal/shared/database/models"
	"context"

	"github.com/google/uuid"
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
}

type authRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &authRepository{db: db}
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
