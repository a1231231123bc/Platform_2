package repository

import (
	"context"
	"fmt"

	"platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobOfferRepository struct {
	pool *pgxpool.Pool
}

func NewJobOfferRepository(pool *pgxpool.Pool) *JobOfferRepository {
	return &JobOfferRepository{pool: pool}
}

func (r *JobOfferRepository) FindByToken(ctx context.Context, token string) (*models.JobOfferWithJob, error) {
	var item models.JobOfferWithJob

	err := r.pool.QueryRow(ctx, `
		SELECT
			jo.id, jo.job_id, jo.dispatch_id, jo.contractor_id, jo.status, jo.token, jo.sent_at, jo.responded_at,
			j.id, j.title, j.description, j.region, j.equipment_type, j.volume, j.deadline, j.price, j.conditions, j.status, j.organization_id, j.created_by_user_id, j.created_at, j.updated_at
		FROM job_offers jo
		JOIN jobs j ON j.id = jo.job_id
		WHERE jo.token = $1::uuid
	`, token).Scan(
		&item.ID, &item.JobID, &item.DispatchID, &item.ContractorID, &item.Status, &item.Token, &item.SentAt, &item.RespondedAt,
		&item.Job.ID, &item.Job.Title, &item.Job.Description, &item.Job.Region, &item.Job.EquipmentType, &item.Job.Volume, &item.Job.Deadline, &item.Job.Price, &item.Job.Conditions, &item.Job.Status, &item.Job.OrganizationID, &item.Job.CreatedByUserID, &item.Job.CreatedAt, &item.Job.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find job offer by token: %w", err)
	}

	return &item, nil
}

func (r *JobOfferRepository) UpdateStatusByToken(ctx context.Context, token string, from []models.OfferStatus, to models.OfferStatus) (*models.JobOfferWithJob, error) {
	allowed := make([]string, 0, len(from))
	for _, status := range from {
		allowed = append(allowed, string(status))
	}

	var item models.JobOfferWithJob
	err := r.pool.QueryRow(ctx, `
		WITH updated AS (
			UPDATE job_offers
			SET
				status = $2::offer_status,
				responded_at = now()
			WHERE token = $1::uuid
				AND status = ANY($3::offer_status[])
			RETURNING id, job_id, dispatch_id, contractor_id, status, token, sent_at, responded_at
		)
		SELECT
			u.id, u.job_id, u.dispatch_id, u.contractor_id, u.status, u.token, u.sent_at, u.responded_at,
			j.id, j.title, j.description, j.region, j.equipment_type, j.volume, j.deadline, j.price, j.conditions, j.status, j.organization_id, j.created_by_user_id, j.created_at, j.updated_at
		FROM updated u
		JOIN jobs j ON j.id = u.job_id
	`, token, to, allowed).Scan(
		&item.ID, &item.JobID, &item.DispatchID, &item.ContractorID, &item.Status, &item.Token, &item.SentAt, &item.RespondedAt,
		&item.Job.ID, &item.Job.Title, &item.Job.Description, &item.Job.Region, &item.Job.EquipmentType, &item.Job.Volume, &item.Job.Deadline, &item.Job.Price, &item.Job.Conditions, &item.Job.Status, &item.Job.OrganizationID, &item.Job.CreatedByUserID, &item.Job.CreatedAt, &item.Job.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update job offer status by token: %w", err)
	}

	return &item, nil
}
