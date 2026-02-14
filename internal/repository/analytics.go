package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardMetrics struct {
	TotalJobs        int `json:"totalJobs"`
	ActiveJobs       int `json:"activeJobs"`
	TotalContractors int `json:"totalContractors"`
	TotalAssignments int `json:"totalAssignments"`
}

type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

func (r *AnalyticsRepository) Dashboard(ctx context.Context, organizationID string) (*DashboardMetrics, error) {
	var metrics DashboardMetrics

	err := r.pool.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM jobs j WHERE j.organization_id = $1) AS total_jobs,
			(SELECT COUNT(*) FROM jobs j WHERE j.organization_id = $1 AND j.status IN ('ACTIVE', 'IN_PROGRESS')) AS active_jobs,
			(
				SELECT COUNT(DISTINCT a.contractor_id)
				FROM assignments a
				JOIN jobs j ON j.id = a.job_id
				WHERE j.organization_id = $1
			) AS total_contractors,
			(
				SELECT COUNT(*)
				FROM assignments a
				JOIN jobs j ON j.id = a.job_id
				WHERE j.organization_id = $1
			) AS total_assignments
	`, organizationID).Scan(
		&metrics.TotalJobs,
		&metrics.ActiveJobs,
		&metrics.TotalContractors,
		&metrics.TotalAssignments,
	)
	if err != nil {
		return nil, fmt.Errorf("get dashboard metrics: %w", err)
	}

	return &metrics, nil
}
