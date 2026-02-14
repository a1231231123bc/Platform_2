package repository

import (
	"context"
	"fmt"

	"platform/internal/dto"
	"platform/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RatingRepository struct {
	pool *pgxpool.Pool
}

func NewRatingRepository(pool *pgxpool.Pool) *RatingRepository {
	return &RatingRepository{pool: pool}
}

func (r *RatingRepository) Create(ctx context.Context, organizationID string, req dto.CreateRatingRequest) (*models.Rating, error) {
	var rating models.Rating
	err := r.pool.QueryRow(ctx, `
		INSERT INTO ratings (job_id, contractor_id, organization_id, score, comment, author_type)
		VALUES ($1, $2, $3::uuid, $4, $5, 'CUSTOMER')
		RETURNING id, job_id, contractor_id, author_contractor_id, organization_id, score, comment, author_type, created_at
	`, req.JobID, req.ContractorID, organizationID, req.Score, req.Comment).Scan(
		&rating.ID, &rating.JobID, &rating.ContractorID, &rating.AuthorContractorID, &rating.OrganizationID, &rating.Score, &rating.Comment, &rating.AuthorType, &rating.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create rating: %w", err)
	}
	return &rating, nil
}

func (r *RatingRepository) FindByContractorInOrganization(ctx context.Context, contractorID, organizationID string) ([]models.Rating, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, job_id, contractor_id, author_contractor_id, organization_id, score, comment, author_type, created_at
		FROM ratings
		WHERE contractor_id = $1 AND organization_id = $2::uuid
		ORDER BY created_at DESC
	`, contractorID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find ratings by contractor in organization: %w", err)
	}
	defer rows.Close()

	ratings := make([]models.Rating, 0, 16)
	for rows.Next() {
		var rating models.Rating
		if err := rows.Scan(
			&rating.ID, &rating.JobID, &rating.ContractorID, &rating.AuthorContractorID, &rating.OrganizationID, &rating.Score, &rating.Comment, &rating.AuthorType, &rating.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rating by contractor: %w", err)
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

func (r *RatingRepository) FindByJobInOrganization(ctx context.Context, jobID, organizationID string) ([]models.Rating, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, job_id, contractor_id, author_contractor_id, organization_id, score, comment, author_type, created_at
		FROM ratings
		WHERE job_id = $1 AND organization_id = $2::uuid
		ORDER BY created_at DESC
	`, jobID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("find ratings by job in organization: %w", err)
	}
	defer rows.Close()

	ratings := make([]models.Rating, 0, 16)
	for rows.Next() {
		var rating models.Rating
		if err := rows.Scan(
			&rating.ID, &rating.JobID, &rating.ContractorID, &rating.AuthorContractorID, &rating.OrganizationID, &rating.Score, &rating.Comment, &rating.AuthorType, &rating.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rating by job: %w", err)
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}
