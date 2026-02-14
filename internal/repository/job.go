package repository

import (
	"context"
	"fmt"
	"strings"

	"platform/internal/dto"
	"platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	pool *pgxpool.Pool
}

func NewJobRepository(pool *pgxpool.Pool) *JobRepository {
	return &JobRepository{pool: pool}
}

func (r *JobRepository) Create(ctx context.Context, organizationID, createdByUserID string, req dto.CreateJobRequest) (*models.Job, error) {
	var j models.Job
	err := r.pool.QueryRow(ctx,
		`INSERT INTO jobs (
			title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6::timestamptz, $7, $8, 'DRAFT', $9, $10
		)
		RETURNING id, title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id, created_at, updated_at`,
		req.Title, req.Description, req.Region, req.EquipmentType, req.Volume, req.Deadline, req.Price, req.Conditions, organizationID, createdByUserID,
	).Scan(
		&j.ID, &j.Title, &j.Description, &j.Region, &j.EquipmentType, &j.Volume, &j.Deadline, &j.Price, &j.Conditions, &j.Status, &j.OrganizationID, &j.CreatedByUserID, &j.CreatedAt, &j.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}
	return &j, nil
}

func (r *JobRepository) FindByID(ctx context.Context, id, organizationID string) (*models.Job, error) {
	var j models.Job
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id, created_at, updated_at
		 FROM jobs WHERE id = $1 AND organization_id = $2`,
		id, organizationID,
	).Scan(
		&j.ID, &j.Title, &j.Description, &j.Region, &j.EquipmentType, &j.Volume, &j.Deadline, &j.Price, &j.Conditions, &j.Status, &j.OrganizationID, &j.CreatedByUserID, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find job by id: %w", err)
	}
	return &j, nil
}

func (r *JobRepository) FindAllByOrganization(ctx context.Context, organizationID string, req dto.QueryJobsRequest) ([]models.JobWithCounts, int, error) {
	where, args := buildJobFilters(organizationID, req)

	countQuery := "SELECT COUNT(*) FROM jobs j WHERE " + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count jobs: %w", err)
	}

	argsWithPage := append(args, req.Limit, (req.Page-1)*req.Limit)
	query := `SELECT
		j.id, j.title, j.description, j.region, j.equipment_type, j.volume, j.deadline, j.price, j.conditions, j.status, j.organization_id, j.created_by_user_id, j.created_at, j.updated_at,
		COALESCE(d.dispatch_count, 0) AS dispatch_count,
		COALESCE(o.offer_count, 0) AS offer_count,
		COALESCE(a.assignment_count, 0) AS assignment_count
	FROM jobs j
	LEFT JOIN (
		SELECT job_id, COUNT(*) AS dispatch_count FROM job_dispatches GROUP BY job_id
	) d ON d.job_id = j.id
	LEFT JOIN (
		SELECT job_id, COUNT(*) AS offer_count FROM job_offers GROUP BY job_id
	) o ON o.job_id = j.id
	LEFT JOIN (
		SELECT job_id, COUNT(*) AS assignment_count FROM assignments GROUP BY job_id
	) a ON a.job_id = j.id
	WHERE ` + where + `
	ORDER BY j.created_at DESC
	LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.pool.Query(ctx, query, argsWithPage...)
	if err != nil {
		return nil, 0, fmt.Errorf("query jobs: %w", err)
	}
	defer rows.Close()

	jobs := make([]models.JobWithCounts, 0, req.Limit)
	for rows.Next() {
		var j models.JobWithCounts
		if err := rows.Scan(
			&j.ID, &j.Title, &j.Description, &j.Region, &j.EquipmentType, &j.Volume, &j.Deadline, &j.Price, &j.Conditions, &j.Status, &j.OrganizationID, &j.CreatedByUserID, &j.CreatedAt, &j.UpdatedAt,
			&j.DispatchCount, &j.OfferCount, &j.AssignmentCount,
		); err != nil {
			return nil, 0, fmt.Errorf("scan job with counts: %w", err)
		}
		jobs = append(jobs, j)
	}

	return jobs, total, nil
}

func (r *JobRepository) UpdateDraft(ctx context.Context, id, organizationID string, req dto.UpdateJobRequest) (*models.Job, error) {
	var j models.Job
	err := r.pool.QueryRow(ctx,
		`UPDATE jobs SET
			title = COALESCE($3, title),
			description = COALESCE($4, description),
			region = COALESCE($5, region),
			equipment_type = COALESCE($6, equipment_type),
			volume = COALESCE($7, volume),
			deadline = COALESCE($8::timestamptz, deadline),
			price = COALESCE($9, price),
			conditions = COALESCE($10, conditions)
		WHERE id = $1 AND organization_id = $2 AND status = 'DRAFT'
		RETURNING id, title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id, created_at, updated_at`,
		id, organizationID, req.Title, req.Description, req.Region, req.EquipmentType, req.Volume, req.Deadline, req.Price, req.Conditions,
	).Scan(
		&j.ID, &j.Title, &j.Description, &j.Region, &j.EquipmentType, &j.Volume, &j.Deadline, &j.Price, &j.Conditions, &j.Status, &j.OrganizationID, &j.CreatedByUserID, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update draft job: %w", err)
	}
	return &j, nil
}

func (r *JobRepository) SetStatus(ctx context.Context, id, organizationID string, from []models.JobStatus, to models.JobStatus) (*models.Job, error) {
	allowed := make([]string, 0, len(from))
	for _, s := range from {
		allowed = append(allowed, string(s))
	}

	var j models.Job
	err := r.pool.QueryRow(ctx,
		`UPDATE jobs
		 SET status = $3::job_status
		 WHERE id = $1 AND organization_id = $2 AND status = ANY($4::job_status[])
		 RETURNING id, title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id, created_at, updated_at`,
		id, organizationID, to, allowed,
	).Scan(
		&j.ID, &j.Title, &j.Description, &j.Region, &j.EquipmentType, &j.Volume, &j.Deadline, &j.Price, &j.Conditions, &j.Status, &j.OrganizationID, &j.CreatedByUserID, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("set job status: %w", err)
	}
	return &j, nil
}

func (r *JobRepository) DuplicateAsDraft(ctx context.Context, id, organizationID, createdByUserID string) (*models.Job, error) {
	var j models.Job
	err := r.pool.QueryRow(ctx, `
		INSERT INTO jobs (
			title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id
		)
		SELECT
			title, description, region, equipment_type, volume, deadline, price, conditions, 'DRAFT', organization_id, $3
		FROM jobs
		WHERE id = $1 AND organization_id = $2
		RETURNING id, title, description, region, equipment_type, volume, deadline, price, conditions, status, organization_id, created_by_user_id, created_at, updated_at
	`, id, organizationID, createdByUserID).Scan(
		&j.ID, &j.Title, &j.Description, &j.Region, &j.EquipmentType, &j.Volume, &j.Deadline, &j.Price, &j.Conditions, &j.Status, &j.OrganizationID, &j.CreatedByUserID, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("duplicate job as draft: %w", err)
	}
	return &j, nil
}

func buildJobFilters(organizationID string, req dto.QueryJobsRequest) (string, []interface{}) {
	parts := []string{"j.organization_id = $1"}
	args := []interface{}{organizationID}

	if req.Status != nil && *req.Status != "" {
		args = append(args, *req.Status)
		parts = append(parts, fmt.Sprintf("j.status = $%d::job_status", len(args)))
	}
	if req.Region != nil && *req.Region != "" {
		args = append(args, *req.Region)
		parts = append(parts, fmt.Sprintf("j.region = $%d", len(args)))
	}

	return strings.Join(parts, " AND "), args
}
