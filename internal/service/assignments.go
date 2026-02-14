package service

import (
	"context"
	"errors"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"
	"platform/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type AssignmentsService struct {
	assignmentRepo *repository.AssignmentRepository
	jobRepo        *repository.JobRepository
	contractorRepo *repository.ContractorRepository
}

func NewAssignmentsService(assignmentRepo *repository.AssignmentRepository, jobRepo *repository.JobRepository, contractorRepo *repository.ContractorRepository) *AssignmentsService {
	return &AssignmentsService{
		assignmentRepo: assignmentRepo,
		jobRepo:        jobRepo,
		contractorRepo: contractorRepo,
	}
}

func (s *AssignmentsService) Create(ctx context.Context, organizationID string, req dto.CreateAssignmentRequest) (*models.Assignment, error) {
	job, err := s.jobRepo.FindByID(ctx, req.JobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for assignment: %w", err)
	}
	if job == nil {
		return nil, ErrNotFound("Job not found")
	}

	contractor, err := s.contractorRepo.FindByID(ctx, req.ContractorID)
	if err != nil {
		return nil, fmt.Errorf("find contractor for assignment: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	assignment, err := s.assignmentRepo.Create(ctx, req.JobID, req.ContractorID)
	if err != nil {
		if isAssignmentConflict(err) {
			return nil, ErrConflict("Assignment already exists for this job and contractor")
		}
		return nil, fmt.Errorf("create assignment: %w", err)
	}

	return assignment, nil
}

func (s *AssignmentsService) ListByJob(ctx context.Context, organizationID, jobID string) ([]models.AssignmentWithContractor, error) {
	job, err := s.jobRepo.FindByID(ctx, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find job for assignment list: %w", err)
	}
	if job == nil {
		return nil, ErrNotFound("Job not found")
	}

	items, err := s.assignmentRepo.FindByJobInOrganization(ctx, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list assignments by job: %w", err)
	}
	return items, nil
}

func (s *AssignmentsService) UpdateStatus(ctx context.Context, organizationID, assignmentID string, req dto.UpdateAssignmentStatusRequest) (*models.Assignment, error) {
	assignment, err := s.assignmentRepo.FindByIDInOrganization(ctx, assignmentID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find assignment for status update: %w", err)
	}
	if assignment == nil {
		return nil, ErrNotFound("Assignment not found")
	}

	updated, err := s.assignmentRepo.UpdateStatusInOrganization(ctx, assignmentID, organizationID, models.AssignmentStatus(req.Status))
	if err != nil {
		return nil, fmt.Errorf("update assignment status: %w", err)
	}
	if updated == nil {
		return nil, ErrNotFound("Assignment not found")
	}

	return updated, nil
}

func isAssignmentConflict(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
