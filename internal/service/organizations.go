package service

import (
	"context"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
)

type OrganizationsService struct {
	orgRepo *repository.OrganizationRepository
}

func NewOrganizationsService(orgRepo *repository.OrganizationRepository) *OrganizationsService {
	return &OrganizationsService{orgRepo: orgRepo}
}

func (s *OrganizationsService) GetByID(ctx context.Context, organizationID, requestedID string) (*models.Organization, error) {
	if _, err := uuid.Parse(requestedID); err != nil {
		return nil, ErrBadRequest("Invalid organization ID")
	}
	if requestedID != organizationID {
		return nil, ErrForbidden("Access denied for this organization")
	}

	org, err := s.orgRepo.FindByID(ctx, requestedID)
	if err != nil {
		return nil, fmt.Errorf("find organization by id: %w", err)
	}
	if org == nil {
		return nil, ErrNotFound("Organization not found")
	}

	return org, nil
}

func (s *OrganizationsService) Update(ctx context.Context, organizationID, requestedID string, req dto.UpdateOrganizationRequest) (*models.Organization, error) {
	if _, err := uuid.Parse(requestedID); err != nil {
		return nil, ErrBadRequest("Invalid organization ID")
	}
	if requestedID != organizationID {
		return nil, ErrForbidden("Access denied for this organization")
	}

	org, err := s.orgRepo.Update(ctx, requestedID, req.Name, req.INN, req.ContactEmail, req.ContactPhone)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}
	if org == nil {
		return nil, ErrNotFound("Organization not found")
	}

	return org, nil
}
