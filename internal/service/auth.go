package service

import (
	"context"
	"fmt"
	"time"

	"platform/internal/config"
	"platform/internal/dto"
	"platform/internal/middleware"
	"platform/internal/models"
	"platform/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	orgRepo  *repository.OrganizationRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, orgRepo *repository.OrganizationRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
		cfg:      cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(ctx, req.AdminEmail)
	if err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrConflict("User with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	tx, err := s.userRepo.Pool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	org, err := s.orgRepo.Create(ctx, tx, req.OrganizationName, req.INN, req.ContactEmail, req.ContactPhone)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	user, err := s.userRepo.Create(ctx, tx, req.AdminEmail, string(hash), req.AdminName, models.UserRoleAdmin, org.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User: dto.AuthUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
		Organization: &dto.AuthOrg{
			ID:   org.ID,
			Name: org.Name,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrUnauthorized("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrUnauthorized("Invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User: dto.AuthUser{
			ID:             user.ID,
			Email:          user.Email,
			Name:           user.Name,
			Role:           string(user.Role),
			OrganizationID: user.OrganizationID,
		},
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := middleware.JWTClaims{
		Sub:            user.ID,
		Email:          user.Email,
		Role:           string(user.Role),
		OrganizationID: user.OrganizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return middleware.GenerateToken(s.cfg.JWTSecret, claims)
}
