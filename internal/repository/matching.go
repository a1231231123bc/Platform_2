package repository

import (
	"context"
	"fmt"
	"strings"

	"platform/internal/models"
)

// FindMatchingByJob returns available contractors matched by job region and equipment.
func (r *ContractorRepository) FindMatchingByJob(ctx context.Context, region string, equipmentType *string) ([]models.Contractor, error) {
	where, args := buildMatchingFilters(region, equipmentType)
	query := `SELECT id, name, phone, email, type, regions, equipment_types, price_expectations, is_available, consent_given, consent_date, created_at, updated_at
		FROM contractors c
		WHERE ` + where + `
		ORDER BY c.created_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query matching contractors: %w", err)
	}
	defer rows.Close()

	contractors := make([]models.Contractor, 0, 16)
	for rows.Next() {
		var c models.Contractor
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Email, &c.Type, &c.Regions, &c.EquipmentTypes,
			&c.PriceExpectations, &c.IsAvailable, &c.ConsentGiven, &c.ConsentDate, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan matching contractor: %w", err)
		}
		contractors = append(contractors, c)
	}

	return contractors, nil
}

func buildMatchingFilters(region string, equipmentType *string) (string, []interface{}) {
	parts := []string{"c.is_available = true"}
	args := []interface{}{region}
	parts = append(parts, "$1 = ANY(c.regions)")

	if equipmentType != nil && *equipmentType != "" {
		args = append(args, *equipmentType)
		parts = append(parts, fmt.Sprintf("c.equipment_types @> ARRAY[$%d]::text[]", len(args)))
	}

	return strings.Join(parts, " AND "), args
}
