package users

import (
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/router"
	"context"
	"errors"
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
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init(basePath string, env *router.Environment) error {
	h.logger = env.Logger
	h.env = env

	// Get service from environment if not already set
	if h.service == nil {
		if svc := env.Services.Get("users"); svc != nil {
			h.service = svc.(Service)
		}
	}

	userGroup := env.Fiber.Group(basePath + "/users")
	userGroup.Get("/", h.getAllUsers)
	userGroup.Get("/:id", h.getUserByID)
	userGroup.Post("/", h.createUser)

	return nil
}

func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	usersList, err := h.service.GetAllUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch users"})
	}

	return c.Status(fiber.StatusOK).JSON(usersList)
}

func (h *Handler) getUserByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id format"})
	}

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch user"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handler) createUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Minute)
	defer cancel()

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := &models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		RoleAtMBL:   req.RoleAtMBL,
		Email:       req.Email,
		Phone:       req.Phone,
		DateOfBirth: req.DateOfBirth,
		State:       req.State,
		Department:  req.Department,
		Role:        "user", // Default role for new users
	}

	if err := h.service.CreateUser(ctx, user); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exists"})
		}
		if errors.Is(err, ErrPhoneAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "phone already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
