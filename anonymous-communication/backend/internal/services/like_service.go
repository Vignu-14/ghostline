package services

import (
	"context"
	"fmt"

	"anonymous-communication/backend/internal/models"

	"github.com/google/uuid"
)

type likePostRepository interface {
	FindByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
}

type likeRepository interface {
	Create(ctx context.Context, userID, postID uuid.UUID) error
	Delete(ctx context.Context, userID, postID uuid.UUID) error
}

type LikeService struct {
	posts likePostRepository
	likes likeRepository
}

func NewLikeService(posts likePostRepository, likes likeRepository) *LikeService {
	return &LikeService{
		posts: posts,
		likes: likes,
	}
}

func (s *LikeService) Like(ctx context.Context, userID, postID uuid.UUID) error {
	post, err := s.posts.FindByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.UserID == userID {
		return models.ErrCannotLikeOwnPost
	}

	if err := s.likes.Create(ctx, userID, postID); err != nil {
		return fmt.Errorf("like post: %w", err)
	}

	return nil
}

func (s *LikeService) Unlike(ctx context.Context, userID, postID uuid.UUID) error {
	if _, err := s.posts.FindByID(ctx, postID); err != nil {
		return err
	}

	if err := s.likes.Delete(ctx, userID, postID); err != nil {
		return fmt.Errorf("unlike post: %w", err)
	}

	return nil
}
