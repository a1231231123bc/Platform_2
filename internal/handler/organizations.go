package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type OrganizationsHandler struct {
	orgService *service.OrganizationsService
}

func NewOrganizationsHandler(orgService *service.OrganizationsService) *OrganizationsHandler {
	return &OrganizationsHandler{orgService: orgService}
}

func (h *OrganizationsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	org, err := h.orgService.GetByID(r.Context(), claims.OrganizationID, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, org)
}

func (h *OrganizationsHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	var req dto.UpdateOrganizationRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	org, err := h.orgService.Update(r.Context(), claims.OrganizationID, id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, org)
}
