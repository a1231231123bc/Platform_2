package dto

// RegisterRequest is the request body for POST /auth/register.
type RegisterRequest struct {
	OrganizationName string  `json:"organizationName" validate:"required,min=2,max=255"`
	INN              *string `json:"inn"`
	ContactEmail     *string `json:"contactEmail" validate:"omitempty,email"`
	ContactPhone     *string `json:"contactPhone"`
	AdminName        string  `json:"adminName" validate:"required,min=2"`
	AdminEmail       string  `json:"adminEmail" validate:"required,email"`
	AdminPassword    string  `json:"adminPassword" validate:"required,min=8"`
}

// LoginRequest is the request body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// AuthResponse is the response for register/login.
type AuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	User         AuthUser     `json:"user"`
	Organization *AuthOrg     `json:"organization,omitempty"`
}

type AuthUser struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	OrganizationID string `json:"organizationId,omitempty"`
}

type AuthOrg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
