package users

import (
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/notifications"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type userService struct {
	repository   Repository
	logger       *zap.Logger
	notification *notifications.Notification
}

func NewService(repository Repository, logger *zap.Logger, notification *notifications.Notification) Service {
	return &userService{
		repository:   repository,
		logger:       logger,
		notification: notification,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.repository.GetAllUsers(ctx)
	if err != nil {
		s.logger.Error("failed to get all users", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		s.logger.Error("failed to get user by id", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existing, err := s.repository.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to check existing email", zap.Error(err))
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}

	// Check if phone already exists
	existing, err = s.repository.GetUserByPhone(ctx, user.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to check existing phone", zap.Error(err))
		return err
	}
	if existing != nil {
		return ErrPhoneAlreadyExists
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return err
	}

	html, err := s.notification.Email.TemplateParser.Parse("signup.html", map[string]interface{}{
		"name": user.FirstName,
	})
	if err != nil {
		s.logger.Error("failed to parse template", zap.Error(err))
		return err
	}

	notificationPayload := &notifications.NotificationPayload{
		To:      user.Email,
		Subject: "Welcome to Alloy",
		Body:    "Welcome to Alloy. You are now setup and ready to go.",
		HTML:    html,
	}
	err = s.notification.Email.Send(notificationPayload)
	if err != nil {
		s.logger.Error("failed to send email", zap.Error(err))
		return err
	}

	return nil
}
