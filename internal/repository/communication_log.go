package repository

import (
	"context"
	"fmt"

	"platform/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunicationLogRepository struct {
	pool *pgxpool.Pool
}

func NewCommunicationLogRepository(pool *pgxpool.Pool) *CommunicationLogRepository {
	return &CommunicationLogRepository{pool: pool}
}

func (r *CommunicationLogRepository) Create(
	ctx context.Context,
	jobID *string,
	contractorID, channel, message string,
	direction models.CommunicationDirection,
) (*models.CommunicationLog, error) {
	var log models.CommunicationLog
	err := r.pool.QueryRow(ctx, `
		INSERT INTO communication_logs (job_id, contractor_id, channel, message, direction)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5::communication_direction)
		RETURNING id, job_id, contractor_id, channel, message, direction, created_at
	`, jobID, contractorID, channel, message, direction).Scan(
		&log.ID, &log.JobID, &log.ContractorID, &log.Channel, &log.Message, &log.Direction, &log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create communication log: %w", err)
	}
	return &log, nil
}

func (r *CommunicationLogRepository) FindByContractor(ctx context.Context, contractorID string, page, limit int) ([]models.CommunicationLog, error) {
	offset := (page - 1) * limit
	rows, err := r.pool.Query(ctx, `
		SELECT id, job_id, contractor_id, channel, message, direction, created_at
		FROM communication_logs
		WHERE contractor_id = $1::uuid
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, contractorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find communication logs by contractor: %w", err)
	}
	defer rows.Close()

	items := make([]models.CommunicationLog, 0, limit)
	for rows.Next() {
		var item models.CommunicationLog
		if err := rows.Scan(&item.ID, &item.JobID, &item.ContractorID, &item.Channel, &item.Message, &item.Direction, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan communication log: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
