package validations

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func init() {
	// Register custom validation for signup schema
	validate.RegisterValidation("signup_contact", validateSignupContact)
	// Register custom validation for strong password
	validate.RegisterValidation("strong_password", validateStrongPassword)
	// Register custom validation for country code
	validate.RegisterValidation("country_code", validateCountryCode)
	// Register custom validation for username
	validate.RegisterValidation("username", validateUsername)
	// Register custom validation for decimal/price values
	validate.RegisterValidation("decimal", validateDecimal)
	// Register custom validation for positive decimal values
	validate.RegisterValidation("decimal_positive", validateDecimalPositive)
}

// validateSignupContact ensures either email or phone_number is provided
func validateSignupContact(fl validator.FieldLevel) bool {
	// This will be called on the struct level, not individual fields
	// We'll handle this in the SignupSchema validation method
	return true
}

// validateStrongPassword validates password strength:
// - At least 8 characters
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one number
// - Contains at least one special character
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum length check (already handled by min=8, but we check here too)
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsername validates username format:
// - Allows letters (a-z, A-Z), numbers (0-9), underscores (_), and hyphens (-)
// - Must start with a letter or number (not underscore or hyphen)
// - Must not end with underscore or hyphen
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) == 0 {
		return true // Empty is handled by omitempty
	}

	// Check if starts with letter or number
	if len(username) > 0 {
		firstChar := rune(username[0])
		if !unicode.IsLetter(firstChar) && !unicode.IsDigit(firstChar) {
			return false
		}
	}

	// Check if ends with underscore or hyphen
	if len(username) > 0 {
		lastChar := rune(username[len(username)-1])
		if lastChar == '_' || lastChar == '-' {
			return false
		}
	}

	// Check all characters are valid (letters, numbers, underscores, hyphens)
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			return false
		}
	}

	return true
}

// validateDecimal validates decimal/currency format:
// - Allows digits (0-9)
// - Allows optional decimal point (.)
// - Allows optional negative sign (-) at the start
// Examples: "100", "99.99", "-50", "0.01"
func validateDecimal(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) == 0 {
		return true // Empty is handled by omitempty
	}

	// Check for negative sign at the start
	startIdx := 0
	if len(value) > 0 && value[0] == '-' {
		startIdx = 1
		if len(value) == 1 {
			return false // Just a minus sign is invalid
		}
	}

	hasDecimalPoint := false
	hasDigit := false

	// Validate each character
	for i := startIdx; i < len(value); i++ {
		char := value[i]
		switch {
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '.':
			if hasDecimalPoint {
				return false // Multiple decimal points
			}
			hasDecimalPoint = true
		default:
			return false // Invalid character
		}
	}

	// Must have at least one digit
	if !hasDigit {
		return false
	}

	// Decimal point cannot be at the start or end
	if hasDecimalPoint {
		if startIdx == len(value)-1 || value[len(value)-1] == '.' {
			return false
		}
	}

	return true
}

// validateDecimalPositive validates that a string is a valid positive decimal number:
// - Must be a valid decimal format (digits and optional decimal point)
// - Must be positive (greater than zero)
// - Must be parseable as a decimal
// Examples: "100", "99.99", "0.01" (valid), "-50", "0", "-0.5", "0.0" (invalid)
func validateDecimalPositive(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) == 0 {
		return true // Empty is handled by omitempty
	}

	// First check if it's a valid decimal format
	if !validateDecimal(fl) {
		return false
	}

	// Check for negative sign
	if len(value) > 0 && value[0] == '-' {
		return false // Negative values are not allowed
	}

	// Check if value is exactly "0" or starts with "0" followed by only zeros and optional decimal point
	// This catches "0", "0.0", "0.00", etc.
	trimmed := strings.TrimLeft(value, "0")
	if trimmed == "" || trimmed == "." {
		return false // Zero is not positive
	}

	// Check if it starts with multiple zeros before decimal (like "00.5" which is invalid format)
	if len(value) > 1 && value[0] == '0' && value[1] != '.' {
		return false
	}

	return true
}

// ValidateStruct validates a struct based on `validate` tags
func ValidateStruct(payload interface{}) error { return validate.Struct(payload) }

// ValidateSignupContact validates that either email or phone_number is provided
func ValidateSignupContact(email, phoneNumber, countryCode string) error {
	email = strings.TrimSpace(email)
	phoneNumber = strings.TrimSpace(phoneNumber)
	countryCode = strings.TrimSpace(countryCode)

	// Check if both are empty
	if email == "" && phoneNumber == "" {
		return &ValidationError{
			Field:   "contact",
			Message: "Either email or phone number must be provided",
		}
	}

	// Check if both are provided
	if email != "" && phoneNumber != "" {
		return &ValidationError{
			Field:   "contact",
			Message: "Please provide either email or phone number, not both",
		}
	}

	// If phone number is provided, validate it
	if phoneNumber != "" {
		if err := ValidatePhoneNumber(phoneNumber, countryCode); err != nil {
			return err
		}
	}

	return nil
}
