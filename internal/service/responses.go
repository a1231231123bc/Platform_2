package service

import (
	"context"
	"fmt"

	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
)

type ResponsesService struct {
	jobOfferRepo *repository.JobOfferRepository
}

func NewResponsesService(jobOfferRepo *repository.JobOfferRepository) *ResponsesService {
	return &ResponsesService{jobOfferRepo: jobOfferRepo}
}

func (s *ResponsesService) GetByToken(ctx context.Context, token string) (*models.JobOfferWithJob, error) {
	if _, err := uuid.Parse(token); err != nil {
		return nil, ErrBadRequest("Invalid offer token")
	}

	offer, err := s.jobOfferRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("get offer by token: %w", err)
	}
	if offer == nil {
		return nil, ErrNotFound("Offer not found")
	}

	return offer, nil
}

func (s *ResponsesService) Accept(ctx context.Context, token string) (*models.JobOfferWithJob, error) {
	return s.respond(ctx, token, models.OfferStatusAccepted)
}

func (s *ResponsesService) Decline(ctx context.Context, token string) (*models.JobOfferWithJob, error) {
	return s.respond(ctx, token, models.OfferStatusDeclined)
}

func (s *ResponsesService) respond(ctx context.Context, token string, target models.OfferStatus) (*models.JobOfferWithJob, error) {
	if _, err := uuid.Parse(token); err != nil {
		return nil, ErrBadRequest("Invalid offer token")
	}

	current, err := s.jobOfferRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("find offer before response: %w", err)
	}
	if current == nil {
		return nil, ErrNotFound("Offer not found")
	}
	if !isOfferRespondable(current.Status) {
		return nil, ErrBadRequest("Offer is no longer available for response")
	}

	updated, err := s.jobOfferRepo.UpdateStatusByToken(ctx, token, []models.OfferStatus{models.OfferStatusSent}, target)
	if err != nil {
		return nil, fmt.Errorf("update offer response status: %w", err)
	}
	if updated == nil {
		return nil, ErrBadRequest("Offer is no longer available for response")
	}

	return updated, nil
}

func isOfferRespondable(status models.OfferStatus) bool {
	return status == models.OfferStatusSent
}
