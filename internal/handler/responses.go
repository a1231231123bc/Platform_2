package handler

import (
	"net/http"

	"platform/internal/dto"
	"platform/internal/service"

	"github.com/go-chi/chi/v5"
)

type ResponsesHandler struct {
	responsesService *service.ResponsesService
}

func NewResponsesHandler(responsesService *service.ResponsesService) *ResponsesHandler {
	return &ResponsesHandler{responsesService: responsesService}
}

func (h *ResponsesHandler) GetByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	offer, err := h.responsesService.GetByToken(r.Context(), token)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, offer)
}

func (h *ResponsesHandler) Accept(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	offer, err := h.responsesService.Accept(r.Context(), token)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, offer)
}

func (h *ResponsesHandler) Decline(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	offer, err := h.responsesService.Decline(r.Context(), token)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	dto.WriteJSON(w, http.StatusOK, offer)
}
