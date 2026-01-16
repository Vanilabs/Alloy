package auth

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/router"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
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

	authGroup.Post("/invite", h.inviteUser)
	authGroup.Patch("/accept", h.acceptInvitation)
	return nil
}

func (h *Handler) inviteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var req InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.service.InviteUser(ctx, req.Email, req.Role, "54dfcc8f-a7ee-44c0-9758-fa17e28a90d0") // temp
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

	err := h.service.AcceptInvitation(ctx, req.Token, req.Email)
	if err != nil {
		if errors.Is(err, ErrInvitationAlreadyVerified) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "invitation already verified"})
		}
		if errors.Is(err, ErrInvitationNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invitation not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to verify invitation"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "invitation verified successfully"})
}
