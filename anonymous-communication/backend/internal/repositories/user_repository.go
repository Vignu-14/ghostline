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

type UserRepository struct {
	db *pgxpool.Pool
}

type CreateUserParams struct {
	Username     string
	Email        string
	PasswordHash string
	Role         string
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, params CreateUserParams) (*models.User, error) {
	const query = `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, password_hash, role, impersonation_password_hash, profile_picture_url, created_at
	`

	user, err := scanUser(r.db.QueryRow(ctx, query, params.Username, params.Email, params.PasswordHash, params.Role))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	const query = `
		SELECT id, username, email, password_hash, role, impersonation_password_hash, profile_picture_url, created_at
		FROM users
		WHERE LOWER(username) = LOWER($1)
		LIMIT 1
	`

	user, err := scanUser(r.db.QueryRow(ctx, query, username))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, username, email, password_hash, role, impersonation_password_hash, profile_picture_url, created_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	user, err := scanUser(r.db.QueryRow(ctx, query, userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`

	var exists bool
	if err := r.db.QueryRow(ctx, query, username).Scan(&exists); err != nil {
		return false, fmt.Errorf("check username exists: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) = LOWER($1))`

	var exists bool
	if err := r.db.QueryRow(ctx, query, email).Scan(&exists); err != nil {
		return false, fmt.Errorf("check email exists: %w", err)
	}

	return exists, nil
}

func scanUser(row pgx.Row) (*models.User, error) {
	var user models.User
	var impersonationHash sql.NullString
	var profilePictureURL sql.NullString

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&impersonationHash,
		&profilePictureURL,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, fmt.Errorf("scan user: %w", err)
	}

	if impersonationHash.Valid {
		user.ImpersonationPasswordHash = &impersonationHash.String
	}

	if profilePictureURL.Valid {
		user.ProfilePictureURL = &profilePictureURL.String
	}

	return &user, nil
}
