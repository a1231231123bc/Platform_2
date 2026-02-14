package service

import (
	"context"
	"fmt"

	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ComplianceService struct {
	blacklistRepo *repository.ContractorBlacklistRepository
}

func NewComplianceService(blacklistRepo *repository.ContractorBlacklistRepository) *ComplianceService {
	return &ComplianceService{blacklistRepo: blacklistRepo}
}

func (s *ComplianceService) ListBlacklist(ctx context.Context, organizationID string) ([]models.ContractorBlacklistWithContractor, error) {
	items, err := s.blacklistRepo.FindAllByOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list blacklist: %w", err)
	}
	return items, nil
}

func (s *ComplianceService) DeleteBlacklistEntry(ctx context.Context, organizationID, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrBadRequest("Invalid blacklist entry ID")
	}

	err := s.blacklistRepo.DeleteInOrganization(ctx, id, organizationID)
	if err == pgx.ErrNoRows {
		return ErrNotFound("Blacklist entry not found")
	}
	if err != nil {
		return fmt.Errorf("delete blacklist entry: %w", err)
	}
	return nil
}
