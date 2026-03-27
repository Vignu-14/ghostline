package services

import (
	"context"
	"fmt"
	"strings"

	"anonymous-communication/backend/internal/models"

	"github.com/google/uuid"
)

type userRepository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	SearchByUsername(ctx context.Context, query string, excludeUserID uuid.UUID, limit int) ([]models.UserSearchResult, error)
}

type userPostRepository interface {
	FeedByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.PostFeedItem, error)
}

type UserService struct {
	users userRepository
	posts userPostRepository
}

func NewUserService(users userRepository, posts userPostRepository) *UserService {
	return &UserService{
		users: users,
		posts: posts,
	}
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) SearchByUsername(ctx context.Context, requesterID uuid.UUID, query string, limit int) ([]models.UserSearchResult, error) {
	trimmedQuery := strings.TrimSpace(query)
	if len(trimmedQuery) < 1 {
		return make([]models.UserSearchResult, 0), nil
	}

	results, err := s.users.SearchByUsername(ctx, trimmedQuery, requesterID, limit)
	if err != nil {
		return nil, fmt.Errorf("search users by username: %w", err)
	}

	return results, nil
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string, limit, offset int) (*models.PublicUserProfile, []models.PostFeedItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("get profile by username: %w", err)
	}

	posts, err := s.posts.FeedByUserID(ctx, user.ID, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("get posts for profile: %w", err)
	}

	profile := user.ToPublicProfile()
	return &profile, posts, nil
}
