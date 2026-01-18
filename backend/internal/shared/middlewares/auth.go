package middlewares

import (
	"alloy/internal/modules/users"
	"alloy/internal/shared/constants"
	"alloy/internal/shared/database/models"
	"alloy/internal/shared/router"
	"alloy/internal/shared/utils"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func JWTMiddleware(env *router.Environment) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if authHeader == "" {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				models.ErrMissingOrInvalidAuthorizationHeader.Error(),
			)
		}

		tokenStr := env.JWTManager.ExtractTokenFromHeader(authHeader)
		if tokenStr == "" {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				models.ErrMissingOrInvalidAuthorizationHeader.Error(),
			)
		}

		claims, err := env.JWTManager.ParseJWT(tokenStr)
		if err != nil {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				models.ErrInvalidOrExpiredToken.Error(),
			)
		}

		if claims == nil {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				models.ErrInvalidOrExpiredToken.Error(),
			)
		}

		if claims.UserID == "" {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				models.ErrInvalidOrExpiredToken.Error(),
			)
		}

		// Retrieve userService from the environment
		userService := env.Services.Get(constants.USERS_MODULE_NAME)
		if userService != nil {
			ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
			defer cancel()

			_, err = userService.(users.Service).GetUserByID(ctx, uuid.MustParse(claims.UserID))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) || err == users.ErrUserNotFound {
					return utils.SendErrorResponse(
						c, fiber.StatusUnauthorized,
						models.ErrInvalidOrExpiredToken.Error(), // Don't reveal user existence
					)
				}
				// For other errors (timeout, DB errors), return generic error
				return utils.SendErrorResponse(
					c, fiber.StatusUnauthorized,
					models.ErrInvalidOrExpiredToken.Error(),
				)
			}
		}

		// Set validated claims in context for downstream handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("token_id", claims.ID)
		c.Locals("claims", claims)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
