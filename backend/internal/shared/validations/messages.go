package validations

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldMessageOverride lets you specify custom messages per field and tag.
//
//	Example: {
//	  "Filename": {"required": "Filename is required"},
//	}
type FieldMessageOverride map[string]map[string]string

func ValidateWithMessages(payload interface{}, overrides FieldMessageOverride) error {
	if err := ValidateStruct(payload); err != nil {
		if verifications, ok := err.(validator.ValidationErrors); ok {
			messages := make([]string, 0, len(verifications))
			for _, fe := range verifications {

				if overrides != nil {
					if fieldOverrides, ok := overrides[fe.Field()]; ok {
						if msg, ok := fieldOverrides[fe.Tag()]; ok && msg != "" {
							messages = append(messages, msg)
							continue
						}
						if msg, ok := fieldOverrides["*"]; ok && msg != "" { // wildcard per field
							messages = append(messages, msg)
							continue
						}
					}
				}
				messages = append(messages, messageForTag(toReadableField(fe.Field()), fe))
			}
			return &ErrorList{Messages: messages}
		}
		return err
	}
	return nil
}

// ErrorList aggregates multiple validation messages into one error.
type ErrorList struct {
	Messages []string
}

func (e *ErrorList) Error() string {
	return strings.Join(e.Messages, "; ")
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// toReadableField converts struct field name to a spaced lower-case label
func toReadableField(field string) string {
	var out []rune
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, ' ')
		}
		out = append(out, r)
	}
	s := strings.ToLower(string(out))
	s = strings.ReplaceAll(s, " id", " ID")
	return s
}

// messageForTag returns a default human-friendly message for a field error
func messageForTag(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "len":
		return field + " must be exactly " + fe.Param() + " characters"
	case "oneof":
		return field + " must be one of: " + fe.Param()
	case "alphanum":
		return field + " must be alphanumeric"
	case "username":
		return field + " can only contain letters, numbers, underscores, and hyphens. It must start with a letter or number"
	case "numeric":
		return field + " must be numeric"
	case "decimal":
		return field + " must be a valid decimal number"
	case "decimal_positive":
		return field + " must be a valid positive decimal number greater than zero"
	case "eqfield":
		return field + " must match " + toReadableField(fe.Param())
	case "strong_password":
		return field + " must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	default:
		return field + " is invalid"
	}
}
