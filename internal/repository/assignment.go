package repository

import (
	"context"
	"fmt"

	"platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssignmentRepository struct {
	pool *pgxpool.Pool
}

func NewAssignmentRepository(pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{pool: pool}
}

func (r *AssignmentRepository) Create(ctx context.Context, jobID, contractorID string) (*models.Assignment, error) {
	var a models.Assignment
	err := r.pool.QueryRow(ctx, `
		INSERT INTO assignments (job_id, contractor_id, status)
		VALUES ($1, $2, 'ASSIGNED')
		RETURNING id, job_id, contractor_id, status, assigned_at, updated_at
	`, jobID, contractorID).Scan(
		&a.ID, &a.JobID, &a.ContractorID, &a.Status, &a.AssignedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create assignment: %w", err)
	}
	return &a, nil
}

func (r *AssignmentRepository) FindByIDInOrganization(ctx context.Context, assignmentID, organizationID string) (*models.Assignment, error) {
	var a models.Assignment
	err := r.pool.QueryRow(ctx, `
		SELECT a.id, a.job_id, a.contractor_id, a.status, a.assigned_at, a.updated_at
		FROM assignments a
		JOIN jobs j ON j.id = a.job_id
		WHERE a.id = $1 AND j.organization_id = $2
	`, assignmentID, organizationID).Scan(
		&a.ID, &a.JobID, &a.ContractorID, &a.Status, &a.AssignedAt, &a.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find assignment by id in organization: %w", err)
	}
	return &a, nil
}

func (r *AssignmentRepository) FindByJobInOrganization(ctx context.Context, jobID, organizationID string) ([]models.AssignmentWithContractor, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			a.id, a.job_id, a.contractor_id, a.status, a.assigned_at, a.updated_at,
			c.id, c.name, c.phone, c.email, c.type, c.regions, c.equipment_types, c.price_expectations, c.is_available, c.consent_given, c.consent_date, c.created_at, c.updated_at
		FROM assignments a
		JOIN jobs j ON j.id = a.job_id
		JOIN contractors c ON c.id = a.contractor_id
		WHERE a.job_id = $1 AND j.organization_id = $2
		ORDER BY a.assigned_at DESC
	`, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find assignments by job in organization: %w", err)
	}
	defer rows.Close()

	result := make([]models.AssignmentWithContractor, 0, 8)
	for rows.Next() {
		var item models.AssignmentWithContractor
		if err := rows.Scan(
			&item.ID, &item.JobID, &item.ContractorID, &item.Status, &item.AssignedAt, &item.UpdatedAt,
			&item.Contractor.ID, &item.Contractor.Name, &item.Contractor.Phone, &item.Contractor.Email, &item.Contractor.Type, &item.Contractor.Regions, &item.Contractor.EquipmentTypes, &item.Contractor.PriceExpectations, &item.Contractor.IsAvailable, &item.Contractor.ConsentGiven, &item.Contractor.ConsentDate, &item.Contractor.CreatedAt, &item.Contractor.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan assignment with contractor: %w", err)
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *AssignmentRepository) UpdateStatusInOrganization(ctx context.Context, assignmentID, organizationID string, status models.AssignmentStatus) (*models.Assignment, error) {
	var a models.Assignment
	err := r.pool.QueryRow(ctx, `
		UPDATE assignments a
		SET status = $3::assignment_status
		FROM jobs j
		WHERE a.id = $1
			AND a.job_id = j.id
			AND j.organization_id = $2
		RETURNING a.id, a.job_id, a.contractor_id, a.status, a.assigned_at, a.updated_at
	`, assignmentID, organizationID, status).Scan(
		&a.ID, &a.JobID, &a.ContractorID, &a.Status, &a.AssignedAt, &a.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update assignment status in organization: %w", err)
	}
	return &a, nil
}
