package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type ComplianceHandler struct {
	complianceService *service.ComplianceService
}

func NewComplianceHandler(complianceService *service.ComplianceService) *ComplianceHandler {
	return &ComplianceHandler{complianceService: complianceService}
}

func (h *ComplianceHandler) ListBlacklist(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	items, err := h.complianceService.ListBlacklist(r.Context(), claims.OrganizationID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, items)
}

func (h *ComplianceHandler) DeleteBlacklistEntry(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.complianceService.DeleteBlacklistEntry(r.Context(), claims.OrganizationID, id); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
