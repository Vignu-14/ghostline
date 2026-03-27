package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"anonymous-communication/backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db *pgxpool.Pool
}

type CreatePostParams struct {
	UserID   uuid.UUID
	ImageURL *string
	Caption  *string
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, params CreatePostParams) (*models.Post, error) {
	const query = `
		INSERT INTO posts (user_id, image_url, caption)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, image_url, caption, created_at
	`

	post, err := scanPost(r.db.QueryRow(ctx, query, params.UserID, params.ImageURL, params.Caption))
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return post, nil
}

func (r *PostRepository) Feed(ctx context.Context, limit, offset int) ([]models.PostFeedItem, error) {
	const query = `
		SELECT
			p.id,
			p.user_id,
			p.image_url,
			p.caption,
			p.created_at,
			u.id,
			u.username,
			u.profile_picture_url,
			COUNT(l.user_id)::bigint AS like_count
		FROM posts p
		INNER JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts feed: %w", err)
	}
	defer rows.Close()

	posts := make([]models.PostFeedItem, 0)
	for rows.Next() {
		post, err := scanPostFeedRow(rows)
		if err != nil {
			return nil, err
		}

		posts = append(posts, *post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate posts feed: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) FeedByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.PostFeedItem, error) {
	const query = `
		SELECT
			p.id,
			p.user_id,
			p.image_url,
			p.caption,
			p.created_at,
			u.id,
			u.username,
			u.profile_picture_url,
			COUNT(l.user_id)::bigint AS like_count
		FROM posts p
		INNER JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id
		WHERE p.user_id = $1
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts by user: %w", err)
	}
	defer rows.Close()

	posts := make([]models.PostFeedItem, 0)
	for rows.Next() {
		post, err := scanPostFeedRow(rows)
		if err != nil {
			return nil, err
		}

		posts = append(posts, *post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user posts: %w", err)
	}

	return posts, nil
}

func (r *PostRepository) FindByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	const query = `
		SELECT id, user_id, image_url, caption, created_at
		FROM posts
		WHERE id = $1
		LIMIT 1
	`

	post, err := scanPost(r.db.QueryRow(ctx, query, postID))
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) FindFeedByID(ctx context.Context, postID uuid.UUID) (*models.PostFeedItem, error) {
	const query = `
		SELECT
			p.id,
			p.user_id,
			p.image_url,
			p.caption,
			p.created_at,
			u.id,
			u.username,
			u.profile_picture_url,
			COUNT(l.user_id)::bigint AS like_count
		FROM posts p
		INNER JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.post_id = p.id
		WHERE p.id = $1
		GROUP BY p.id, u.id
	`

	post, err := scanPostFeedRow(r.db.QueryRow(ctx, query, postID))
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) DeleteByID(ctx context.Context, postID uuid.UUID) error {
	const query = `DELETE FROM posts WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return models.ErrPostNotFound
	}

	return nil
}

func (r *PostRepository) ExistsByID(ctx context.Context, postID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, postID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check post exists: %w", err)
	}

	return exists, nil
}

func scanPost(row pgx.Row) (*models.Post, error) {
	var post models.Post
	var imageURL sql.NullString
	var caption sql.NullString

	err := row.Scan(&post.ID, &post.UserID, &imageURL, &caption, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrPostNotFound
		}

		return nil, fmt.Errorf("scan post: %w", err)
	}

	if imageURL.Valid {
		post.ImageURL = &imageURL.String
	}

	if caption.Valid {
		post.Caption = &caption.String
	}

	return &post, nil
}

func scanPostFeedRow(row pgx.Row) (*models.PostFeedItem, error) {
	var (
		post         models.PostFeedItem
		postID       uuid.UUID
		userID       uuid.UUID
		authorID     uuid.UUID
		imageURL     sql.NullString
		caption      sql.NullString
		profilePhoto sql.NullString
	)

	err := row.Scan(
		&postID,
		&userID,
		&imageURL,
		&caption,
		&post.CreatedAt,
		&authorID,
		&post.User.Username,
		&profilePhoto,
		&post.LikeCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrPostNotFound
		}

		return nil, fmt.Errorf("scan post feed item: %w", err)
	}

	post.ID = postID.String()
	post.UserID = userID.String()
	post.User.ID = authorID.String()

	if imageURL.Valid {
		post.ImageURL = &imageURL.String
	}

	if caption.Valid {
		post.Caption = &caption.String
	}

	if profilePhoto.Valid {
		post.User.ProfilePictureURL = &profilePhoto.String
	}

	return &post, nil
}
