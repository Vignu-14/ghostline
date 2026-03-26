package repositories

import (
	"context"
	"fmt"

	"anonymous-communication/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthLogRepository struct {
	db *pgxpool.Pool
}

func NewAuthLogRepository(db *pgxpool.Pool) *AuthLogRepository {
	return &AuthLogRepository{db: db}
}

func (r *AuthLogRepository) Create(ctx context.Context, params models.CreateAuthLogParams) error {
	const query = `
		INSERT INTO auth_logs (user_id, status, ip_address, user_agent, failure_reason)
		VALUES ($1, $2, $3, $4, $5)
	`

	var userID any
	if params.UserID != nil {
		userID = *params.UserID
	}

	if _, err := r.db.Exec(
		ctx,
		query,
		userID,
		params.Status,
		params.IPAddress,
		nullableString(params.UserAgent),
		nullableString(params.FailureReason),
	); err != nil {
		return fmt.Errorf("insert auth log: %w", err)
	}

	return nil
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}

	return value
}
