package service

import (
	"context"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"
	"platform/internal/repository"
)

type RatingsService struct {
	ratingRepo     *repository.RatingRepository
	jobRepo        *repository.JobRepository
	contractorRepo *repository.ContractorRepository
}

func NewRatingsService(ratingRepo *repository.RatingRepository, jobRepo *repository.JobRepository, contractorRepo *repository.ContractorRepository) *RatingsService {
	return &RatingsService{
		ratingRepo:     ratingRepo,
		jobRepo:        jobRepo,
		contractorRepo: contractorRepo,
	}
}

func (s *RatingsService) Create(ctx context.Context, organizationID string, req dto.CreateRatingRequest) (*models.Rating, error) {
	job, err := s.jobRepo.FindByID(ctx, req.JobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for rating: %w", err)
	}
	if job == nil {
		return nil, ErrNotFound("Job not found")
	}

	contractor, err := s.contractorRepo.FindByID(ctx, req.ContractorID)
	if err != nil {
		return nil, fmt.Errorf("find contractor for rating: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	rating, err := s.ratingRepo.Create(ctx, organizationID, req)
	if err != nil {
		return nil, fmt.Errorf("create rating: %w", err)
	}
	return rating, nil
}

func (s *RatingsService) ListByContractor(ctx context.Context, organizationID, contractorID string) ([]models.Rating, error) {
	contractor, err := s.contractorRepo.FindByID(ctx, contractorID)
	if err != nil {
		return nil, fmt.Errorf("find contractor for ratings list: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	items, err := s.ratingRepo.FindByContractorInOrganization(ctx, contractorID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list ratings by contractor: %w", err)
	}
	return items, nil
}

func (s *RatingsService) ListByJob(ctx context.Context, organizationID, jobID string) ([]models.Rating, error) {
	job, err := s.jobRepo.FindByID(ctx, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for ratings list: %w", err)
	}
	if job == nil {
		return nil, ErrNotFound("Job not found")
	}

	items, err := s.ratingRepo.FindByJobInOrganization(ctx, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list ratings by job: %w", err)
	}
	return items, nil
}
