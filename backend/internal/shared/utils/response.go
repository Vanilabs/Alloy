package utils

import (
	"alloy/internal/shared/database/models"

	"github.com/gofiber/fiber/v2"
)

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	response := models.ErrorResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}

	return c.Status(statusCode).JSON(response)
}
