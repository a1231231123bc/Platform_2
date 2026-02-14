package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type AssignmentsHandler struct {
	assignmentsService *service.AssignmentsService
}

func NewAssignmentsHandler(assignmentsService *service.AssignmentsService) *AssignmentsHandler {
	return &AssignmentsHandler{assignmentsService: assignmentsService}
}

func (h *AssignmentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateAssignmentRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	assignment, err := h.assignmentsService.Create(r.Context(), claims.OrganizationID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusCreated, assignment)
}

func (h *AssignmentsHandler) ListByJob(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jobID := chi.URLParam(r, "jobId")
	items, err := h.assignmentsService.ListByJob(r.Context(), claims.OrganizationID, jobID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, items)
}

func (h *AssignmentsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assignmentID := chi.URLParam(r, "id")
	var req dto.UpdateAssignmentStatusRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	updated, err := h.assignmentsService.UpdateStatus(r.Context(), claims.OrganizationID, assignmentID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, updated)
}
