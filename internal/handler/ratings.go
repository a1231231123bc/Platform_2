package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type RatingsHandler struct {
	ratingsService *service.RatingsService
}

func NewRatingsHandler(ratingsService *service.RatingsService) *RatingsHandler {
	return &RatingsHandler{ratingsService: ratingsService}
}

func (h *RatingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateRatingRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	rating, err := h.ratingsService.Create(r.Context(), claims.OrganizationID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusCreated, rating)
}

func (h *RatingsHandler) ListByContractor(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	contractorID := chi.URLParam(r, "contractorId")
	items, err := h.ratingsService.ListByContractor(r.Context(), claims.OrganizationID, contractorID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, items)
}

func (h *RatingsHandler) ListByJob(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobId")
	items, err := h.ratingsService.ListByJob(r.Context(), claims.OrganizationID, jobID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, items)
}
