package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/cache"
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/notifications"
	"alloy/internal/shared/utils"
	"alloy/internal/shared/validations/schemas"
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	InviteUser(ctx context.Context, payload *schemas.InviteUserSchema, adminId string) error
	VerifyInvitation(ctx context.Context, token string, email string) error
	AcceptInvitation(ctx context.Context, token string, email string) (*models.LoginResponse, error)
	RequestMagicLink(ctx context.Context, email string) error
	VerifyMagicLink(ctx context.Context, token string, sessionInfo *models.UserSessionInfo) (*models.LoginResponse, error)
}

type authService struct {
	repository     Repository
	logger         *zap.Logger
	notification   *notifications.Notification
	userRepository users.Repository
	jwtManager     *utils.JWTManager
}

func NewService(repository Repository, logger *zap.Logger, notification *notifications.Notification, userRepository users.Repository, jwtManager *utils.JWTManager) Service {
	return &authService{
		repository:     repository,
		logger:         logger,
		notification:   notification,
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

func (s *authService) InviteUser(ctx context.Context, payload *schemas.InviteUserSchema, adminId string) error {
	email := payload.Email
	role := payload.Role
	roleAtOrg := payload.RoleAtOrg
	department := payload.Department
	phone := payload.Phone

	_, err := s.userRepository.GetUserByEmail(ctx, email)
	if err == nil {
		return users.ErrEmailAlreadyExists
	}

	invitation := &models.Invitation{
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Email:      email,
		Role:       role,
		RoleAtOrg:  roleAtOrg,
		Department: department,
		Phone:      phone,
		InvitedBy:  uuid.MustParse(adminId),
		Token:      uuid.New().String(),
		ExpiresAt:  time.Now().Add(time.Hour * 24 * 7),
	}

	err = s.repository.CreateInvitation(ctx, invitation)
	if err != nil {
		return err
	}

	html, err := s.notification.Email.TemplateParser.Parse("invite.html", map[string]interface{}{
		"inviteUrl": os.Getenv("FRONTEND_URL") + "/invite?token=" + invitation.Token,
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

func (s *authService) AcceptInvitation(ctx context.Context, token string, email string) (*models.LoginResponse, error) {
	invitation, err := s.repository.GetInvitationByTokenAndEmail(ctx, token, email)
	if err != nil {
		s.logger.Info("invitation not found", zap.Error(err), zap.String("token", token), zap.String("email", email))
		return nil, err
	}
	if invitation.Status == models.InvitationStatusAccepted {
		s.logger.Info("invitation already accepted", zap.String("token", token), zap.String("email", email))
		return nil, ErrInvitationAlreadyAccepted
	}

	if invitation.ExpiresAt.Before(time.Now()) {
		s.logger.Info("invitation expired", zap.Error(err), zap.String("token", token), zap.String("email", email))
		return nil, ErrInvitationExpired
	}

	// Generate a unique employee number with retry logic
	var employeeNumber string
	maxRetries := 10
	foundAvailable := false

	for i := 1; i <= maxRetries; i++ {
		var err error
		employeeNumber, err = utils.GenerateEmployeeNumber()
		if err != nil {
			s.logger.Error("failed to generate employee number", zap.Error(err), zap.String("token", token), zap.String("email", email), zap.Int("attempt", i))
			return nil, err
		}

		// Check if employee number already exists
		_, err = s.userRepository.GetUserByEmployeeNumber(ctx, employeeNumber)
		if err != nil {
			// If record not found, the number is available - we can use it
			if errors.Is(err, gorm.ErrRecordNotFound) {
				foundAvailable = true
				break
			}
			// For other database errors, log and return
			s.logger.Error("database error while checking employee number", zap.Error(err), zap.String("token", token), zap.String("email", email), zap.String("employee_number", employeeNumber))
			return nil, err
		}

		// If no error, user exists - collision detected, try again
		s.logger.Info("employee number collision detected, generating new one", zap.String("employee_number", employeeNumber), zap.Int("attempt", i))
	}

	// If we exhausted all retries without finding an available number, return an error
	if !foundAvailable {
		s.logger.Error("failed to generate unique employee number after max retries", zap.String("token", token), zap.String("email", email), zap.Int("max_retries", maxRetries))
		return nil, errors.New("failed to generate unique employee number after maximum retries")
	}

	invitation.Status = models.InvitationStatusAccepted
	now := time.Now()
	invitation.AcceptedAt = &now
	err = s.repository.UpdateInvitation(ctx, invitation)
	if err != nil {
		s.logger.Info("failed to update invitation", zap.Error(err), zap.String("token", token), zap.String("email", email))
		return nil, err
	}

	user := &models.User{
		FirstName:      invitation.FirstName,
		LastName:       invitation.LastName,
		Email:          email,
		Role:           invitation.Role,
		RoleAtOrg:      invitation.RoleAtOrg,
		Department:     invitation.Department,
		Phone:          invitation.Phone,
		EmployeeNumber: employeeNumber,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Info("failed to create user", zap.Error(err), zap.String("token", token), zap.String("email", email))
		return nil, err
	}

	s.logger.Info("user created", zap.String("user_id", user.ID.String()), zap.String("email", email), zap.String("employee_number", employeeNumber))

	// handle email notificaion
	html, err := s.notification.Email.TemplateParser.Parse("welcome.html", map[string]interface{}{
		"firstName":      user.FirstName,
		"lastName":       user.LastName,
		"employeeNumber": user.EmployeeNumber,
	})
	if err != nil {
		s.logger.Error("failed to parse welcome template", zap.Error(err))
		return nil, err
	}

	notificationPayload := &notifications.NotificationPayload{
		To:      email,
		Subject: "Welcome to Alloy",
		Body:    "Welcome to Alloy. You are now setup and ready to go.",
		HTML:    html,
	}

	err = s.notification.Email.Send(notificationPayload)
	if err != nil {
		s.logger.Error("failed to send welcome email", zap.Error(err))
		return nil, err
	}

	s.logger.Info("welcome email sent", zap.String("email", email))

	// generate jwt token
	tokenID := uuid.New().String()
	jwtData := &models.JWTData{
		UserID:  user.ID.String(),
		Email:   user.Email,
		TokenID: tokenID,
		Role:    user.Role,
	}

	accessToken, err := s.jwtManager.GenerateJWT(jwtData)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), tokenID)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	user, err = s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}

	return &models.LoginResponse{
		Auth: models.JwtAuthData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(utils.SessionExpiry.Seconds()),
			TokenType:    "Bearer",
		},
		User: *user,
	}, nil
}

func (s *authService) RequestMagicLink(ctx context.Context, email string) error {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("Magic link requested for non-existent user", zap.String("email", email))
		return nil
	}

	token, err := utils.GenerateSecureToken(models.MagicLinkTokenLength)
	if err != nil {
		s.logger.Error("Failed to generate magic link token", zap.Error(err))
		return ErrMagicLinkTokenGeneration
	}

	magicLinkData := &cache.MagicLinkCacheData{
		UserID:    user.ID.String(),
		Email:     email,
		ExpiresAt: time.Now().Add(models.MagicLinkTokenExpiry),
	}

	err = s.repository.Cache_SetMagicLinkToken(ctx, token, magicLinkData, models.MagicLinkTokenExpiry)
	if err != nil {
		s.logger.Error("Failed to set magic link token", zap.Error(err))
		return err
	}

	s.logger.Info("Magic link token set", zap.String("token", token))
	s.logger.Debug("User magic link token stored",
		zap.String("user_id", magicLinkData.UserID),
		zap.String("email", magicLinkData.Email),
		zap.Duration("ttl", models.MagicLinkTokenExpiry))

	// handle email notificaion
	html, err := s.notification.Email.TemplateParser.Parse("magic_link.html", map[string]interface{}{
		"magicLinkURL": os.Getenv("FRONTEND_URL") + "/auth/magic-link?token=" + token,
	})
	if err != nil {
		s.logger.Error("Failed to parse magic link template", zap.Error(err))
		return err
	}

	notificationPayload := &notifications.NotificationPayload{
		To:      email,
		Subject: "Alloy Login Link",
		Body:    "You've requested to login to Alloy. Click the button below to securely sign in.",
		HTML:    html,
	}

	err = s.notification.Email.Send(notificationPayload)
	if err != nil {
		s.logger.Error("Failed to send magic link email", zap.Error(err))
		return err
	}

	s.logger.Info("Magic link email sent", zap.String("email", email))
	s.logger.Debug("Magic link email sent",
		zap.String("user_id", magicLinkData.UserID),
		zap.String("email", magicLinkData.Email),
		zap.Duration("ttl", models.MagicLinkTokenExpiry))

	return nil
}

func (s *authService) VerifyMagicLink(ctx context.Context, token string, sessionInfo *models.UserSessionInfo) (*models.LoginResponse, error) {
	magicLinkData, err := s.repository.Cache_GetMagicLinkToken(ctx, token)
	if err != nil {
		s.logger.Error("invalid or expired magic link token", zap.Error(err))
		return nil, err
	}

	if time.Now().After(magicLinkData.ExpiresAt) {
		s.logger.Error("magic link token expired", zap.String("token", token))
		err = s.repository.Cache_DeleteMagicLinkToken(ctx, token)
		if err != nil {
			s.logger.Error("failed to delete expired magic link token", zap.Error(err))
		}
		return nil, ErrMagicLinkTokenExpired
	}

	user, err := s.userRepository.GetUserByID(ctx, uuid.MustParse(magicLinkData.UserID))
	if err != nil {
		s.logger.Error("user not found", zap.String("user_id", magicLinkData.UserID))
		return nil, err
	}

	// delete used magic link token
	err = s.repository.Cache_DeleteMagicLinkToken(ctx, token)
	if err != nil {
		s.logger.Error("failed to delete used magic link token", zap.Error(err))
	}

	sessionInfo.UserID = user.ID.String()
	sessionInfo.LoginTime = time.Now()
	sessionInfo.LastSeen = time.Now()
	sessionInfo.IsActive = true
	if sessionInfo.TokenID == "" {
		sessionInfo.TokenID = uuid.New().String()
	}

	jwtData := &models.JWTData{
		UserID:  user.ID.String(),
		Email:   user.Email,
		TokenID: sessionInfo.TokenID,
		Role:    user.Role,
	}

	accessToken, err := s.jwtManager.GenerateJWT(jwtData)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), sessionInfo.TokenID)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	return &models.LoginResponse{
		Auth: models.JwtAuthData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(utils.SessionExpiry.Seconds()),
			TokenType:    "Bearer",
		},
		User: models.User{
			ID:             user.ID,
			Email:          user.Email,
			Role:           user.Role,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			RoleAtOrg:      user.RoleAtOrg,
			EmployeeNumber: user.EmployeeNumber,
		},
	}, nil
}
