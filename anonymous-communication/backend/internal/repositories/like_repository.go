package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(ctx context.Context, userID, postID uuid.UUID) error {
	const query = `
		INSERT INTO likes (user_id, post_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, post_id) DO NOTHING
	`

	if _, err := r.db.Exec(ctx, query, userID, postID); err != nil {
		return fmt.Errorf("create like: %w", err)
	}

	return nil
}

func (r *LikeRepository) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	const query = `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`

	if _, err := r.db.Exec(ctx, query, userID, postID); err != nil {
		return fmt.Errorf("delete like: %w", err)
	}

	return nil
}
