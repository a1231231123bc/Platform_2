package repository

import (
	"context"
	"fmt"

	"platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

func (r *OrganizationRepository) Create(ctx context.Context, tx pgx.Tx, name string, inn, contactEmail, contactPhone *string) (*models.Organization, error) {
	var o models.Organization
	err := tx.QueryRow(ctx,
		`INSERT INTO organizations (name, inn, contact_email, contact_phone)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, inn, contact_email, contact_phone, created_at, updated_at`,
		name, inn, contactEmail, contactPhone,
	).Scan(&o.ID, &o.Name, &o.INN, &o.ContactEmail, &o.ContactPhone, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}
	return &o, nil
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*models.Organization, error) {
	var o models.Organization
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, inn, contact_email, contact_phone, created_at, updated_at
		 FROM organizations WHERE id = $1`, id,
	).Scan(&o.ID, &o.Name, &o.INN, &o.ContactEmail, &o.ContactPhone, &o.CreatedAt, &o.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find organization by id: %w", err)
	}
	return &o, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, id string, name, inn, contactEmail, contactPhone *string) (*models.Organization, error) {
	var o models.Organization
	err := r.pool.QueryRow(ctx,
		`UPDATE organizations SET
			name = COALESCE($2, name),
			inn = COALESCE($3, inn),
			contact_email = COALESCE($4, contact_email),
			contact_phone = COALESCE($5, contact_phone)
		 WHERE id = $1
		 RETURNING id, name, inn, contact_email, contact_phone, created_at, updated_at`,
		id, name, inn, contactEmail, contactPhone,
	).Scan(&o.ID, &o.Name, &o.INN, &o.ContactEmail, &o.ContactPhone, &o.CreatedAt, &o.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}
	return &o, nil
}
