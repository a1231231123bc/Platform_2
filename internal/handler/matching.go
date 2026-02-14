package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type MatchingHandler struct {
	matchingService *service.MatchingService
}

func NewMatchingHandler(matchingService *service.MatchingService) *MatchingHandler {
	return &MatchingHandler{matchingService: matchingService}
}

func (h *MatchingHandler) MatchByJobID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobId")
	contractors, err := h.matchingService.MatchByJobID(r.Context(), claims.OrganizationID, jobID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, contractors)
}
