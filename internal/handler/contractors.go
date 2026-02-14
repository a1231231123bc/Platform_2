package handler

import (
	"net/http"
	"strconv"

	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type ContractorsHandler struct {
	contractorsService *service.ContractorsService
}

func NewContractorsHandler(contractorsService *service.ContractorsService) *ContractorsHandler {
	return &ContractorsHandler{contractorsService: contractorsService}
}

func (h *ContractorsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateContractorRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	contractor, err := h.contractorsService.Create(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusCreated, contractor)
}

func (h *ContractorsHandler) List(w http.ResponseWriter, r *http.Request) {
	req, err := parseQueryContractorsRequest(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := dto.ValidateStruct(req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	resp, err := h.contractorsService.List(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, resp)
}

func (h *ContractorsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	contractor, err := h.contractorsService.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, contractor)
}

func (h *ContractorsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateContractorRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	contractor, err := h.contractorsService.Update(r.Context(), id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusOK, contractor)
}

func (h *ContractorsHandler) AddToBlacklist(w http.ResponseWriter, r *http.Request) {
	claims := middleware.CurrentUser(r.Context())
	if claims == nil {
		dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	var req dto.AddToBlacklistRequest
	if err := dto.DecodeAndValidate(r, &req); err != nil {
		dto.WriteValidationError(w, err)
		return
	}

	entry, err := h.contractorsService.AddToBlacklist(r.Context(), claims.OrganizationID, id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	dto.WriteJSON(w, http.StatusCreated, entry)
}

func (h *ContractorsHandler) History(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	page, limit, err := parsePagination(r)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.contractorsService.History(r.Context(), id, page, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, dto.PaginatedResponse{
		Data:       items,
		Total:      len(items),
		Page:       page,
		Limit:      limit,
		TotalPages: 1,
	})
}

func parseQueryContractorsRequest(r *http.Request) (dto.QueryContractorsRequest, error) {
	query := r.URL.Query()

	page, limit, err := parsePagination(r)
	if err != nil {
		return dto.QueryContractorsRequest{}, err
	}

	req := dto.QueryContractorsRequest{
		Page:  page,
		Limit: limit,
	}

	if v := query.Get("region"); v != "" {
		req.Region = &v
	}
	if v := query.Get("type"); v != "" {
		req.Type = &v
	}
	if v := query.Get("search"); v != "" {
		req.Search = &v
	}
	if v := query.Get("isAvailable"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return dto.QueryContractorsRequest{}, err
		}
		req.IsAvailable = &parsed
	}

	return req, nil
}

func parsePagination(r *http.Request) (int, int, error) {
	query := r.URL.Query()
	page := 1
	limit := 20

	if v := query.Get("page"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, err
		}
		page = parsed
	}
	if v := query.Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, err
		}
		limit = parsed
	}

	return page, limit, nil
}
