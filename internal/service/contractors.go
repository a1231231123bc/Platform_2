package service

import (
	"context"
	"errors"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type ContractorsService struct {
	contractorRepo *repository.ContractorRepository
	blacklistRepo  *repository.ContractorBlacklistRepository
}

func NewContractorsService(contractorRepo *repository.ContractorRepository, blacklistRepo *repository.ContractorBlacklistRepository) *ContractorsService {
	return &ContractorsService{
		contractorRepo: contractorRepo,
		blacklistRepo:  blacklistRepo,
	}
}

func (s *ContractorsService) Create(ctx context.Context, req dto.CreateContractorRequest) (*models.Contractor, error) {
	contractor, err := s.contractorRepo.Create(ctx, req)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict("Contractor with this phone already exists")
		}
		return nil, fmt.Errorf("create contractor: %w", err)
	}
	return contractor, nil
}

func (s *ContractorsService) List(ctx context.Context, req dto.QueryContractorsRequest) (*dto.PaginatedResponse, error) {
	contractors, total, err := s.contractorRepo.FindAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list contractors: %w", err)
	}

	totalPages := 0
	if req.Limit > 0 {
		totalPages = (total + req.Limit - 1) / req.Limit
	}

	return &dto.PaginatedResponse{
		Data:       contractors,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *ContractorsService) GetByID(ctx context.Context, id string) (*models.Contractor, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid contractor ID")
	}

	contractor, err := s.contractorRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get contractor by id: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	return contractor, nil
}

func (s *ContractorsService) Update(ctx context.Context, id string, req dto.UpdateContractorRequest) (*models.Contractor, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid contractor ID")
	}

	contractor, err := s.contractorRepo.Update(ctx, id, req)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict("Contractor with this phone already exists")
		}
		return nil, fmt.Errorf("update contractor: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	return contractor, nil
}

func (s *ContractorsService) AddToBlacklist(ctx context.Context, organizationID, contractorID string, req dto.AddToBlacklistRequest) (*models.ContractorBlacklist, error) {
	if _, err := uuid.Parse(contractorID); err != nil {
		return nil, ErrBadRequest("Invalid contractor ID")
	}

	contractor, err := s.contractorRepo.FindByID(ctx, contractorID)
	if err != nil {
		return nil, fmt.Errorf("find contractor for blacklist: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	exists, err := s.blacklistRepo.Exists(ctx, organizationID, contractorID)
	if err != nil {
		return nil, fmt.Errorf("check blacklist: %w", err)
	}
	if exists {
		return nil, ErrConflict("Contractor is already blacklisted")
	}

	entry, err := s.blacklistRepo.Add(ctx, organizationID, contractorID, req.Reason)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict("Contractor is already blacklisted")
		}
		return nil, fmt.Errorf("add to blacklist: %w", err)
	}
	return entry, nil
}

func (s *ContractorsService) History(ctx context.Context, id string, page, limit int) ([]models.ContractorHistoryItem, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrBadRequest("Invalid contractor ID")
	}

	contractor, err := s.contractorRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find contractor for history: %w", err)
	}
	if contractor == nil {
		return nil, ErrNotFound("Contractor not found")
	}

	items, err := s.contractorRepo.FindHistory(ctx, id, page, limit)
	if err != nil {
		return nil, fmt.Errorf("get contractor history: %w", err)
	}

	return items, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
