package service

import (
	"context"
	"fmt"

	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
)

type MatchingService struct {
	jobRepo        *repository.JobRepository
	contractorRepo *repository.ContractorRepository
}

func NewMatchingService(jobRepo *repository.JobRepository, contractorRepo *repository.ContractorRepository) *MatchingService {
	return &MatchingService{
		jobRepo:        jobRepo,
		contractorRepo: contractorRepo,
	}
}

func (s *MatchingService) MatchByJobID(ctx context.Context, organizationID, jobID string) ([]models.Contractor, error) {
	if _, err := uuid.Parse(jobID); err != nil {
		return nil, ErrBadRequest("Invalid job ID")
	}

	job, err := s.jobRepo.FindByID(ctx, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for matching: %w", err)
	}
	if job == nil {
		return nil, ErrNotFound("Job not found")
	}

	contractors, err := s.contractorRepo.FindMatchingByJob(ctx, job.Region, job.EquipmentType)
	if err != nil {
		return nil, fmt.Errorf("find matching contractors: %w", err)
	}

	return contractors, nil
}
