package repository

import (
	"context"
	"fmt"

	"platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContractorBlacklistRepository struct {
	pool *pgxpool.Pool
}

func NewContractorBlacklistRepository(pool *pgxpool.Pool) *ContractorBlacklistRepository {
	return &ContractorBlacklistRepository{pool: pool}
}

func (r *ContractorBlacklistRepository) Add(ctx context.Context, organizationID, contractorID string, reason *string) (*models.ContractorBlacklist, error) {
	var entry models.ContractorBlacklist
	err := r.pool.QueryRow(ctx,
		`INSERT INTO contractor_blacklist (organization_id, contractor_id, reason)
		 VALUES ($1, $2, $3)
		 RETURNING id, organization_id, contractor_id, reason, created_at`,
		organizationID, contractorID, reason,
	).Scan(&entry.ID, &entry.OrganizationID, &entry.ContractorID, &entry.Reason, &entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add contractor to blacklist: %w", err)
	}
	return &entry, nil
}

func (r *ContractorBlacklistRepository) Exists(ctx context.Context, organizationID, contractorID string) (bool, error) {
	var exists bool
	if err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM contractor_blacklist
			WHERE organization_id = $1 AND contractor_id = $2
		)`,
		organizationID, contractorID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("check blacklist existence: %w", err)
	}
	return exists, nil
}

func (r *ContractorBlacklistRepository) FindByID(ctx context.Context, id string) (*models.ContractorBlacklist, error) {
	var entry models.ContractorBlacklist
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, contractor_id, reason, created_at
		 FROM contractor_blacklist WHERE id = $1`,
		id,
	).Scan(&entry.ID, &entry.OrganizationID, &entry.ContractorID, &entry.Reason, &entry.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find blacklist entry by id: %w", err)
	}
	return &entry, nil
}
