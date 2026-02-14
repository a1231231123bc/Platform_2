package repository

import (
	"context"
	"fmt"

	"platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, organization_id, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.OrganizationID, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, organization_id, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.OrganizationID, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, tx pgx.Tx, email, passwordHash, name string, role models.UserRole, organizationID string) (*models.User, error) {
	var u models.User
	err := tx.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name, role, organization_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, email, password_hash, name, role, organization_id, created_at, updated_at`,
		email, passwordHash, name, role, organizationID,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.OrganizationID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindAllByOrganization(ctx context.Context, organizationID string) ([]models.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, email, password_hash, name, role, organization_id, created_at, updated_at
		 FROM users WHERE organization_id = $1 ORDER BY created_at DESC`, organizationID,
	)
	if err != nil {
		return nil, fmt.Errorf("find users by org: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.OrganizationID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, name *string, role *string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET
			name = COALESCE($2, name),
			role = COALESCE($3::user_role, role)
		 WHERE id = $1
		 RETURNING id, email, password_hash, name, role, organization_id, created_at, updated_at`,
		id, name, role,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.OrganizationID, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Pool exposes the pool for transaction usage.
func (r *UserRepository) Pool() *pgxpool.Pool {
	return r.pool
}
