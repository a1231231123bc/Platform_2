package handler

import (
	"errors"
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	resp, err := h.authService.Register(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	resp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":             claims.Sub,
		"email":          claims.Email,
		"role":           claims.Role,
		"organizationId": claims.OrganizationID,
	})
}

// handleServiceError writes the appropriate HTTP error based on ServiceError type.
func handleServiceError(w http.ResponseWriter, err error) {
	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {
		dto.WriteError(w, svcErr.Status, svcErr.Message)
		return
	}
	dto.WriteError(w, http.StatusInternalServerError, "Internal server error")
}
