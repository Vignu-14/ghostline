package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/repositories"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type postRepository interface {
	Create(ctx context.Context, params repositories.CreatePostParams) (*models.Post, error)
	Feed(ctx context.Context, limit, offset int) ([]models.PostFeedItem, error)
	FindByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	FindFeedByID(ctx context.Context, postID uuid.UUID) (*models.PostFeedItem, error)
	DeleteByID(ctx context.Context, postID uuid.UUID) error
}

type postUploadService interface {
	UploadPostImage(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteByPublicURL(ctx context.Context, publicURL string) error
}

type PostService struct {
	posts   postRepository
	uploads postUploadService
}

func NewPostService(posts postRepository, uploads postUploadService) *PostService {
	return &PostService{
		posts:   posts,
		uploads: uploads,
	}
}

func (s *PostService) ListFeed(ctx context.Context, limit, offset int) ([]models.PostFeedItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.posts.Feed(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts feed: %w", err)
	}

	return posts, nil
}

func (s *PostService) Create(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader, caption string) (*models.PostFeedItem, error) {
	if s.uploads == nil {
		return nil, models.ErrStorageNotConfigured
	}

	imageURL, err := s.uploads.UploadPostImage(ctx, userID, file, header)
	if err != nil {
		return nil, err
	}

	var captionValue *string
	cleanCaption := utils.SanitizeText(caption)
	if cleanCaption != "" {
		captionValue = &cleanCaption
	}

	post, err := s.posts.Create(ctx, repositories.CreatePostParams{
		UserID:   userID,
		ImageURL: imageURL,
		Caption:  captionValue,
	})
	if err != nil {
		_ = s.uploads.DeleteByPublicURL(ctx, imageURL)
		return nil, fmt.Errorf("create post: %w", err)
	}

	postFeedItem, err := s.posts.FindFeedByID(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created post: %w", err)
	}

	return postFeedItem, nil
}

func (s *PostService) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	post, err := s.posts.FindByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return models.ErrForbidden
	}

	if err := s.posts.DeleteByID(ctx, postID); err != nil {
		return err
	}

	if s.uploads != nil {
		_ = s.uploads.DeleteByPublicURL(ctx, post.ImageURL)
	}

	return nil
}
