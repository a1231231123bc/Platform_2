package service

import (
	"context"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"
	"platform/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UsersService struct {
	userRepo *repository.UserRepository
}

func NewUsersService(userRepo *repository.UserRepository) *UsersService {
	return &UsersService{userRepo: userRepo}
}

func (s *UsersService) ListByOrganization(ctx context.Context, organizationID string) ([]models.User, error) {
	users, err := s.userRepo.FindAllByOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list users by organization: %w", err)
	}
	return users, nil
}

func (s *UsersService) GetByID(ctx context.Context, organizationID, userID string) (*models.User, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, ErrBadRequest("Invalid user ID")
	}

	user, err := s.userRepo.FindByIDAndOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, ErrNotFound("User not found")
	}
	return user, nil
}

func (s *UsersService) Update(ctx context.Context, organizationID, userID string, req dto.UpdateUserRequest) (*models.User, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, ErrBadRequest("Invalid user ID")
	}

	user, err := s.userRepo.UpdateInOrganization(ctx, userID, organizationID, req.Name, req.Role)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if user == nil {
		return nil, ErrNotFound("User not found")
	}
	return user, nil
}

func (s *UsersService) Delete(ctx context.Context, organizationID, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return ErrBadRequest("Invalid user ID")
	}

	err := s.userRepo.DeleteInOrganization(ctx, userID, organizationID)
	if err == pgx.ErrNoRows {
		return ErrNotFound("User not found")
	}
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
