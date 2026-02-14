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

func (r *ContractorBlacklistRepository) FindAllByOrganization(ctx context.Context, organizationID string) ([]models.ContractorBlacklistWithContractor, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			cb.id, cb.organization_id, cb.contractor_id, cb.reason, cb.created_at,
			c.id, c.name, c.phone, c.email, c.type, c.regions, c.equipment_types, c.price_expectations, c.is_available, c.consent_given, c.consent_date, c.created_at, c.updated_at
		FROM contractor_blacklist cb
		JOIN contractors c ON c.id = cb.contractor_id
		WHERE cb.organization_id = $1
		ORDER BY cb.created_at DESC
	`, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find blacklist by organization: %w", err)
	}
	defer rows.Close()

	items := make([]models.ContractorBlacklistWithContractor, 0, 16)
	for rows.Next() {
		var item models.ContractorBlacklistWithContractor
		if err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.ContractorID, &item.Reason, &item.CreatedAt,
			&item.Contractor.ID, &item.Contractor.Name, &item.Contractor.Phone, &item.Contractor.Email, &item.Contractor.Type, &item.Contractor.Regions, &item.Contractor.EquipmentTypes, &item.Contractor.PriceExpectations, &item.Contractor.IsAvailable, &item.Contractor.ConsentGiven, &item.Contractor.ConsentDate, &item.Contractor.CreatedAt, &item.Contractor.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan blacklist with contractor: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ContractorBlacklistRepository) DeleteInOrganization(ctx context.Context, id, organizationID string) error {
	ct, err := r.pool.Exec(ctx, `
		DELETE FROM contractor_blacklist
		WHERE id = $1 AND organization_id = $2
	`, id, organizationID)
	if err != nil {
		return fmt.Errorf("delete blacklist entry in organization: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
