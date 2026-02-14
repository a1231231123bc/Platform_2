package service

import (
	"context"
	"fmt"

	"platform/internal/repository"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{analyticsRepo: analyticsRepo}
}

func (s *AnalyticsService) Dashboard(ctx context.Context, organizationID string) (*repository.DashboardMetrics, error) {
	metrics, err := s.analyticsRepo.Dashboard(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("load dashboard analytics: %w", err)
	}
	return metrics, nil
}
