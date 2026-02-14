package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type JobsHandler struct {
	jobsService *service.JobsService
}

func NewJobsHandler(jobsService *service.JobsService) *JobsHandler {
	return &JobsHandler{jobsService: jobsService}
}

func (h *JobsHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateJobRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	job, err := h.jobsService.Create(r.Context(), claims.OrganizationID, claims.Sub, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusCreated, job)
}

func (h *JobsHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := parseQueryJobsRequest(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := dto.ValidateStruct(req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	resp, err := h.jobsService.List(r.Context(), claims.OrganizationID, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, resp)
}

func (h *JobsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	job, err := h.jobsService.GetByID(r.Context(), claims.OrganizationID, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, job)
}

func (h *JobsHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	var req dto.UpdateJobRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	job, err := h.jobsService.Update(r.Context(), claims.OrganizationID, id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, job)
}

func (h *JobsHandler) Publish(w http.ResponseWriter, r *http.Request) {
	h.handleStatusTransition(w, r, func(orgID, id string) (interface{}, error) {
		return h.jobsService.Publish(r.Context(), orgID, id)
	})
}

func (h *JobsHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	h.handleStatusTransition(w, r, func(orgID, id string) (interface{}, error) {
		return h.jobsService.Cancel(r.Context(), orgID, id)
	})
}

func (h *JobsHandler) Complete(w http.ResponseWriter, r *http.Request) {
	h.handleStatusTransition(w, r, func(orgID, id string) (interface{}, error) {
		return h.jobsService.Complete(r.Context(), orgID, id)
	})
}

func (h *JobsHandler) Duplicate(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	job, err := h.jobsService.Duplicate(r.Context(), claims.OrganizationID, claims.Sub, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusCreated, job)
}

func (h *JobsHandler) handleStatusTransition(w http.ResponseWriter, r *http.Request, run func(organizationID, id string) (interface{}, error)) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	res, err := run(claims.OrganizationID, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, res)
}

func parseQueryJobsRequest(r *http.Request) (dto.QueryJobsRequest, error) {
	query := r.URL.Query()
	page, limit, err := parsePagination(r)
	if err != nil {
		return dto.QueryJobsRequest{}, err
	}

	req := dto.QueryJobsRequest{
		Page:  page,
		Limit: limit,
	}

	if v := query.Get("status"); v != "" {
		req.Status = &v
	}
	if v := query.Get("region"); v != "" {
		req.Region = &v
	}

	return req, nil
}
