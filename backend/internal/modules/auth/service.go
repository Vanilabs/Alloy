package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/notifications"
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	InviteUser(ctx context.Context, email string, role string, adminId string) error
	VerifyInvitation(ctx context.Context, token string, email string) error
	AcceptInvitation(ctx context.Context, token string, email string) error
}

type authService struct {
	repository     Repository
	logger         *zap.Logger
	notification   *notifications.Notification
	userRepository users.Repository
}

func NewService(repository Repository, logger *zap.Logger, notification *notifications.Notification, userRepository users.Repository) Service {
	return &authService{repository: repository, logger: logger, notification: notification, userRepository: userRepository}
}

func (s *authService) InviteUser(ctx context.Context, email string, role string, adminId string) error {
	_, err := s.userRepository.GetUserByEmail(ctx, email)
	if err == nil {
		return users.ErrEmailAlreadyExists
	}

	invitation := &models.Invitation{
		Email:     email,
		Role:      role,
		InvitedBy: uuid.MustParse(adminId),
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	err = s.repository.CreateInvitation(ctx, invitation)
	if err != nil {
		return err
	}

	html, err := s.notification.Email.TemplateParser.Parse("invite.html", map[string]interface{}{
		"inviteUrl": os.Getenv("FRONTEND_URL") + "/invite/" + invitation.Token,
		"token":     invitation.Token,
	})
	if err != nil {
		return err
	}

	notificationPayload := &notifications.NotificationPayload{
		To:      email,
		Subject: "You're Invited to Alloy",
		Body:    "You've been invited to join Alloy. Click the button below to accept the invitation and get started.",
		HTML:    html,
	}
	err = s.notification.Email.Send(notificationPayload)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) VerifyInvitation(ctx context.Context, token string, email string) error {
	invitation, err := s.repository.GetInvitationByTokenAndEmail(ctx, token, email)
	if err != nil {
		return err
	}

	if invitation.Status == "verified" {
		return ErrInvitationAlreadyVerified
	}

	return nil
}

func (s *authService) AcceptInvitation(ctx context.Context, token string, email string) error {
	invitation, err := s.repository.GetInvitationByTokenAndEmail(ctx, token, email)
	if err != nil {
		return err
	}

	if invitation.ExpiresAt.Before(time.Now()) {
		return ErrInvitationExpired
	}

	invitation.Status = "verified"
	now := time.Now()
	invitation.AcceptedAt = &now
	err = s.repository.UpdateInvitation(ctx, invitation)
	if err != nil {
		return err
	}
	return nil
}
