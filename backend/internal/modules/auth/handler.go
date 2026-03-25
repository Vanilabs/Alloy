package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/middlewares"
	"alloy/internal/shared/router"
	"alloy/internal/shared/validations"
	"alloy/internal/shared/validations/schemas"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	logger  *zap.Logger
	env     *router.Environment
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Init(basePath string, env *router.Environment) error {
	h.logger = env.Logger
	h.env = env

	authGroup := env.Fiber.Group(basePath + "/auth")

	authGroup.Post("/magic-link", h.requestMagicLink)
	authGroup.Post("/magic-link/verify", h.verifyMagicLink)
	authGroup.Patch("/accept", h.acceptInvitation)

	// Protected routes
	authChecker := middlewares.JWTMiddleware(h.env)
	authGroup.Post("/invite", authChecker, h.inviteUser)

	return nil
}
func (h *Handler) inviteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var req schemas.InviteUserSchema
	if err := validations.ParseAndValidateBodyWithMessages(c, &req, req.Messages()); err != nil {
		// Handle validation errors
		if errList, ok := err.(*validations.ErrorList); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errList.Error()})
		}
		if valErr, ok := err.(*validations.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": valErr.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.service.InviteUser(ctx, &req, c.Locals("user_id").(string))
	if err != nil {
		if errors.Is(err, users.ErrEmailAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user invited successfully"})
}
func (h *Handler) acceptInvitation(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var req VerifyInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := h.service.AcceptInvitation(ctx, req.Token, req.Email)
	if err != nil {
		if errors.Is(err, ErrInvitationAlreadyVerified) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "invitation already verified"})
		}
		if errors.Is(err, ErrInvitationNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invitation not found"})
		}
		if errors.Is(err, ErrInvitationAlreadyAccepted) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "invitation already accepted"})
		}
		if errors.Is(err, ErrInvitationExpired) {
			return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "invitation expired"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to verify invitation"})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
func (h *Handler) requestMagicLink(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var req schemas.RequestMagicLinkSchema
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.service.RequestMagicLink(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to request magic link. err: %s", err.Error())})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "magic link sent successfully"})
}
func (h *Handler) verifyMagicLink(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var req schemas.MagicLinkVerifySchema
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sessionInfo := models.UserSessionInfo{
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
		TokenID:   uuid.New().String(),
	}

	loginResponse, err := h.service.VerifyMagicLink(ctx, req.Token, &sessionInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to verify magic link. err: %s", err.Error())})
	}

	h.logger.Info("user authenticated successfully", zap.String("user_id", loginResponse.User.ID.String()))

	return c.Status(fiber.StatusOK).JSON(loginResponse)
}
