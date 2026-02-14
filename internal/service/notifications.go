package service

import (
	"context"
	"fmt"

	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
)

type NotificationsService struct {
	logRepo        *repository.CommunicationLogRepository
	contractorRepo *repository.ContractorRepository
}

func NewNotificationsService(logRepo *repository.CommunicationLogRepository, contractorRepo *repository.ContractorRepository) *NotificationsService {
	return &NotificationsService{
		logRepo:        logRepo,
		contractorRepo: contractorRepo,
	}
}

func (s *NotificationsService) CreateLog(
	ctx context.Context,
	jobID *string,
	contractorID, channel, message, direction string,
) (*models.CommunicationLog, error) {
	if _, err := uuid.Parse(contractorID); err != nil {
		return nil, ErrBadRequest("Invalid contractor ID")
	}
	if jobID != nil {
		if _, err := uuid.Parse(*jobID); err != nil {
			return nil, ErrBadRequest("Invalid job ID")
		}
	}
	if channel == "" {
		return nil, ErrBadRequest("Channel is required")
	}
	if message == "" {
		return nil, ErrBadRequest("Message is required")
	}

	parsedDirection, err := parseCommunicationDirection(direction)
	if err != nil {
		return nil, err
	}

	contractor, err := s.contractorRepo.FindByID(ctx, contractorID)
	if err != nil {
		return nil, fmt.Errorf("find contractor for communication log: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	log, err := s.logRepo.Create(ctx, jobID, contractorID, channel, message, parsedDirection)
	if err != nil {
		return nil, fmt.Errorf("create communication log entry: %w", err)
	}
	return log, nil
}

func (s *NotificationsService) ListByContractor(ctx context.Context, contractorID string, page, limit int) ([]models.CommunicationLog, error) {
	if _, err := uuid.Parse(contractorID); err != nil {
		return nil, ErrBadRequest("Invalid contractor ID")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	contractor, err := s.contractorRepo.FindByID(ctx, contractorID)
	if err != nil {
		return nil, fmt.Errorf("find contractor for communication logs: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	items, err := s.logRepo.FindByContractor(ctx, contractorID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("list communication logs by contractor: %w", err)
	}
	return items, nil
}

func parseCommunicationDirection(direction string) (models.CommunicationDirection, error) {
	switch direction {
	case string(models.CommunicationDirectionOutgoing):
		return models.CommunicationDirectionOutgoing, nil
	case string(models.CommunicationDirectionIncoming):
		return models.CommunicationDirectionIncoming, nil
	default:
		return "", ErrBadRequest("Invalid communication direction")
	}
}
