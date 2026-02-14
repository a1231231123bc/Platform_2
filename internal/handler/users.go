package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	usersService *service.UsersService
}

func NewUsersHandler(usersService *service.UsersService) *UsersHandler {
	return &UsersHandler{usersService: usersService}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	users, err := h.usersService.ListByOrganization(r.Context(), claims.OrganizationID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, users)
}

func (h *UsersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	user, err := h.usersService.GetByID(r.Context(), claims.OrganizationID, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	user, err := h.usersService.Update(r.Context(), claims.OrganizationID, id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.usersService.Delete(r.Context(), claims.OrganizationID, id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
