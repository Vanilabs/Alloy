package validations

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParseAndValidateBodyWithMessages parses JSON and validates with custom messages per field/tag.
// Pass nil for overrides to use default humanized messages.
func ParseAndValidateBodyWithMessages(c *fiber.Ctx, dest interface{}, overrides FieldMessageOverride) error {
	// Use Fiber's BodyParser which is optimized and secure
	if err := c.BodyParser(dest); err != nil {
		// Check if it's a JSON syntax error
		errMsg := err.Error()
		if strings.Contains(errMsg, "Unprocessable Entity") || strings.Contains(errMsg, "cannot parse") {
			// Provide a more helpful message for parsing errors
			return &ValidationError{
				Field:   "body",
				Message: "Invalid JSON format or missing required fields in request body",
			}
		}
		// Return the original error for other cases
		return err
	}

	// Now validate the parsed struct with custom messages
	return ValidateWithMessages(dest, overrides)
}
