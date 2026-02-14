package dto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// PaginatedResponse is a generic paginated response.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// DecodeAndValidate reads JSON body, decodes into dst, and validates.
func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return ValidateStruct(dst)
}

// ValidateStruct validates a struct using validator tags.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidationErrors formats validator errors into a map.
func ValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field()[:1]) + fe.Field()[1:]
			switch fe.Tag() {
			case "required":
				errs[field] = "This field is required"
			case "email":
				errs[field] = "Must be a valid email"
			case "min":
				errs[field] = fmt.Sprintf("Must be at least %s characters", fe.Param())
			case "max":
				errs[field] = fmt.Sprintf("Must be at most %s characters", fe.Param())
			case "oneof":
				errs[field] = fmt.Sprintf("Must be one of: %s", fe.Param())
			case "uuid":
				errs[field] = "Must be a valid UUID"
			default:
				errs[field] = fmt.Sprintf("Failed on '%s' validation", fe.Tag())
			}
		}
	}
	return errs
}

// WriteJSON writes a JSON response.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes a JSON error response.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

// WriteValidationError writes a 400 response with validation details.
func WriteValidationError(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusBadRequest, ErrorResponse{
		Error:   "Validation failed",
		Details: ValidationErrors(err),
	})
}
