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

type ContractorRepository struct {
	pool *pgxpool.Pool
}

func NewContractorRepository(pool *pgxpool.Pool) *ContractorRepository {
	return &ContractorRepository{pool: pool}
}

func (r *ContractorRepository) Create(ctx context.Context, req dto.CreateContractorRequest) (*models.Contractor, error) {
	var c models.Contractor
	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	err := r.pool.QueryRow(ctx,
		`INSERT INTO contractors (
			name, phone, email, type, regions, equipment_types,
			price_expectations, is_available, consent_given, consent_date
		) VALUES ($1, $2, $3, $4::contractor_type, $5, $6, $7, $8, $9, CASE WHEN $9 THEN now() ELSE NULL END)
		RETURNING id, name, phone, email, type, regions, equipment_types, price_expectations, is_available, consent_given, consent_date, created_at, updated_at`,
		req.Name, req.Phone, req.Email, req.Type, req.Regions, req.EquipmentTypes,
		req.PriceExpectations, isAvailable, req.ConsentGiven,
	).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.Type, &c.Regions, &c.EquipmentTypes,
		&c.PriceExpectations, &c.IsAvailable, &c.ConsentGiven, &c.ConsentDate, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create contractor: %w", err)
	}

	return &c, nil
}

func (r *ContractorRepository) FindByID(ctx context.Context, id string) (*models.Contractor, error) {
	var c models.Contractor
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, phone, email, type, regions, equipment_types, price_expectations, is_available, consent_given, consent_date, created_at, updated_at
		 FROM contractors
		 WHERE id = $1`,
		id,
	).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.Type, &c.Regions, &c.EquipmentTypes,
		&c.PriceExpectations, &c.IsAvailable, &c.ConsentGiven, &c.ConsentDate, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find contractor by id: %w", err)
	}
	return &c, nil
}

func (r *ContractorRepository) FindAll(ctx context.Context, req dto.QueryContractorsRequest) ([]models.Contractor, int, error) {
	where, args := buildContractorFilters(req)

	countQuery := "SELECT COUNT(*) FROM contractors c WHERE " + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count contractors: %w", err)
	}

	limitArg := append(args, req.Limit, (req.Page-1)*req.Limit)
	query := `SELECT id, name, phone, email, type, regions, equipment_types, price_expectations, is_available, consent_given, consent_date, created_at, updated_at
		FROM contractors c
		WHERE ` + where + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.pool.Query(ctx, query, limitArg...)
	if err != nil {
		return nil, 0, fmt.Errorf("query contractors: %w", err)
	}
	defer rows.Close()

	contractors := make([]models.Contractor, 0, req.Limit)
	for rows.Next() {
		var c models.Contractor
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.Email, &c.Type, &c.Regions, &c.EquipmentTypes,
			&c.PriceExpectations, &c.IsAvailable, &c.ConsentGiven, &c.ConsentDate, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan contractor: %w", err)
		}
		contractors = append(contractors, c)
	}

	return contractors, total, nil
}

func (r *ContractorRepository) Update(ctx context.Context, id string, req dto.UpdateContractorRequest) (*models.Contractor, error) {
	var c models.Contractor
	err := r.pool.QueryRow(ctx,
		`UPDATE contractors SET
			name = COALESCE($2, name),
			phone = COALESCE($3, phone),
			email = COALESCE($4, email),
			type = COALESCE($5::contractor_type, type),
			regions = CASE WHEN $6::text[] IS NULL THEN regions ELSE $6 END,
			equipment_types = CASE WHEN $7::text[] IS NULL THEN equipment_types ELSE $7 END,
			price_expectations = COALESCE($8, price_expectations),
			is_available = COALESCE($9, is_available),
			consent_given = COALESCE($10, consent_given),
			consent_date = CASE
				WHEN COALESCE($10, consent_given) = true AND consent_date IS NULL THEN now()
				WHEN COALESCE($10, consent_given) = false THEN NULL
				ELSE consent_date
			END
		WHERE id = $1
		RETURNING id, name, phone, email, type, regions, equipment_types, price_expectations, is_available, consent_given, consent_date, created_at, updated_at`,
		id, req.Name, req.Phone, req.Email, req.Type, req.Regions, req.EquipmentTypes, req.PriceExpectations, req.IsAvailable, req.ConsentGiven,
	).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.Type, &c.Regions, &c.EquipmentTypes,
		&c.PriceExpectations, &c.IsAvailable, &c.ConsentGiven, &c.ConsentDate, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update contractor: %w", err)
	}
	return &c, nil
}

func (r *ContractorRepository) FindHistory(ctx context.Context, contractorID string, page, limit int) ([]models.ContractorHistoryItem, error) {
	offset := (page - 1) * limit
	rows, err := r.pool.Query(ctx, `
		SELECT source, job_id, job_title, status, happened_at
		FROM (
			SELECT
				'offer' AS source,
				jo.job_id,
				j.title AS job_title,
				jo.status::text AS status,
				jo.sent_at AS happened_at
			FROM job_offers jo
			JOIN jobs j ON j.id = jo.job_id
			WHERE jo.contractor_id = $1

			UNION ALL

			SELECT
				'assignment' AS source,
				a.job_id,
				j.title AS job_title,
				a.status::text AS status,
				a.assigned_at AS happened_at
			FROM assignments a
			JOIN jobs j ON j.id = a.job_id
			WHERE a.contractor_id = $1
		) history
		ORDER BY happened_at DESC
		LIMIT $2 OFFSET $3
	`, contractorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find contractor history: %w", err)
	}
	defer rows.Close()

	items := make([]models.ContractorHistoryItem, 0, limit)
	for rows.Next() {
		var item models.ContractorHistoryItem
		if err := rows.Scan(&item.Source, &item.JobID, &item.JobTitle, &item.Status, &item.HappenedAt); err != nil {
			return nil, fmt.Errorf("scan history item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func buildContractorFilters(req dto.QueryContractorsRequest) (string, []interface{}) {
	parts := []string{"1=1"}
	args := make([]interface{}, 0, 4)

	if req.Region != nil && *req.Region != "" {
		args = append(args, *req.Region)
		parts = append(parts, fmt.Sprintf("$%d = ANY(c.regions)", len(args)))
	}
	if req.Type != nil && *req.Type != "" {
		args = append(args, *req.Type)
		parts = append(parts, fmt.Sprintf("c.type = $%d::contractor_type", len(args)))
	}
	if req.IsAvailable != nil {
		args = append(args, *req.IsAvailable)
		parts = append(parts, fmt.Sprintf("c.is_available = $%d", len(args)))
	}
	if req.Search != nil && *req.Search != "" {
		args = append(args, "%"+*req.Search+"%")
		parts = append(parts, fmt.Sprintf("(c.name ILIKE $%d OR c.phone ILIKE $%d)", len(args), len(args)))
	}

	return strings.Join(parts, " AND "), args
}
