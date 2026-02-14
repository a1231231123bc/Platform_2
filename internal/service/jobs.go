package service

import (
	"context"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
)

type JobsService struct {
	jobRepo *repository.JobRepository
}

func NewJobsService(jobRepo *repository.JobRepository) *JobsService {
	return &JobsService{jobRepo: jobRepo}
}

func (s *JobsService) Create(ctx context.Context, organizationID, createdByUserID string, req dto.CreateJobRequest) (*models.Job, error) {
	job, err := s.jobRepo.Create(ctx, organizationID, createdByUserID, req)
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}
	return job, nil
}

func (s *JobsService) List(ctx context.Context, organizationID string, req dto.QueryJobsRequest) (*dto.PaginatedResponse, error) {
	jobs, total, err := s.jobRepo.FindAllByOrganization(ctx, organizationID, req)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}

	totalPages := 0
	if req.Limit > 0 {
		totalPages = (total + req.Limit - 1) / req.Limit
	}

	return &dto.PaginatedResponse{
		Data:       jobs,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *JobsService) GetByID(ctx context.Context, organizationID, id string) (*models.Job, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid job ID")
	}

	job, err := s.jobRepo.FindByID(ctx, id, organizationID)
	if err != nil {
		return nil, fmt.Errorf("get job by id: %w", err)
	}
	if job == nil {
		return nil, ErrNotFound("Job not found")
	}

	return job, nil
}

func (s *JobsService) Update(ctx context.Context, organizationID, id string, req dto.UpdateJobRequest) (*models.Job, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid job ID")
	}

	existing, err := s.jobRepo.FindByID(ctx, id, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for update: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound("Job not found")
	}
	if existing.Status != models.JobStatusDraft {
		return nil, ErrBadRequest("Only DRAFT jobs can be updated")
	}

	job, err := s.jobRepo.UpdateDraft(ctx, id, organizationID, req)
	if err != nil {
		return nil, fmt.Errorf("update job: %w", err)
	}
	if job == nil {
		return nil, ErrBadRequest("Only DRAFT jobs can be updated")
	}
	return job, nil
}

func (s *JobsService) Publish(ctx context.Context, organizationID, id string) (*models.Job, error) {
	return s.changeStatus(ctx, organizationID, id, []models.JobStatus{models.JobStatusDraft}, models.JobStatusActive, "Only DRAFT jobs can be published")
}

func (s *JobsService) Cancel(ctx context.Context, organizationID, id string) (*models.Job, error) {
	return s.changeStatus(ctx, organizationID, id, []models.JobStatus{models.JobStatusDraft, models.JobStatusActive, models.JobStatusInProgress}, models.JobStatusCancelled, "Job cannot be cancelled from current status")
}

func (s *JobsService) Complete(ctx context.Context, organizationID, id string) (*models.Job, error) {
	return s.changeStatus(ctx, organizationID, id, []models.JobStatus{models.JobStatusActive, models.JobStatusInProgress}, models.JobStatusCompleted, "Only ACTIVE or IN_PROGRESS jobs can be completed")
}

func (s *JobsService) Duplicate(ctx context.Context, organizationID, createdByUserID, id string) (*models.Job, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid job ID")
	}

	duplicate, err := s.jobRepo.DuplicateAsDraft(ctx, id, organizationID, createdByUserID)
	if err != nil {
		return nil, fmt.Errorf("duplicate job: %w", err)
	}
	if duplicate == nil {
		return nil, ErrNotFound("Job not found")
	}
	return duplicate, nil
}

func (s *JobsService) changeStatus(ctx context.Context, organizationID, id string, from []models.JobStatus, to models.JobStatus, invalidTransitionMsg string) (*models.Job, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid job ID")
	}

	existing, err := s.jobRepo.FindByID(ctx, id, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for status change: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound("Job not found")
	}

	updated, err := s.jobRepo.SetStatus(ctx, id, organizationID, from, to)
	if err != nil {
		return nil, fmt.Errorf("change job status: %w", err)
	}
	if updated == nil {
		return nil, ErrBadRequest(invalidTransitionMsg)
	}
	return updated, nil
}
