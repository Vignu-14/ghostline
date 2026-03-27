package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/repositories"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type postRepository interface {
	Create(ctx context.Context, params repositories.CreatePostParams) (*models.Post, error)
	Feed(ctx context.Context, limit, offset int) ([]models.PostFeedItem, error)
	FeedByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.PostFeedItem, error)
	FindByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	FindFeedByID(ctx context.Context, postID uuid.UUID) (*models.PostFeedItem, error)
	DeleteByID(ctx context.Context, postID uuid.UUID) error
}

type postUploadService interface {
	UploadPostImage(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteByPublicURL(ctx context.Context, publicURL string) error
	CreateSignedPostUpload(ctx context.Context, userID uuid.UUID, request models.CreatePostUploadRequest) (*models.PostUploadTarget, error)
	DeleteByObjectPath(ctx context.Context, objectPath string) error
	PublicURLForObject(objectPath string) string
	ObjectBelongsToUser(userID uuid.UUID, objectPath string) bool
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
	captionValue, err := sanitizeOptionalCaption(caption)
	if err != nil {
		return nil, err
	}

	var imageURL *string
	if file != nil && header != nil {
		if s.uploads == nil {
			return nil, models.ErrStorageNotConfigured
		}

		uploadedURL, err := s.uploads.UploadPostImage(ctx, userID, file, header)
		if err != nil {
			return nil, err
		}

		imageURL = &uploadedURL
	}

	if imageURL == nil && captionValue == nil {
		return nil, models.NewValidationError(map[string]string{
			"post": "add an image or write a thought before posting",
		})
	}

	post, err := s.posts.Create(ctx, repositories.CreatePostParams{
		UserID:   userID,
		ImageURL: imageURL,
		Caption:  captionValue,
	})
	if err != nil {
		if imageURL != nil && s.uploads != nil {
			_ = s.uploads.DeleteByPublicURL(ctx, *imageURL)
		}
		return nil, fmt.Errorf("create post: %w", err)
	}

	postFeedItem, err := s.posts.FindFeedByID(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created post: %w", err)
	}

	return postFeedItem, nil
}

func (s *PostService) CreateUploadURL(ctx context.Context, userID uuid.UUID, request models.CreatePostUploadRequest) (*models.PostUploadTarget, error) {
	if s.uploads == nil {
		return nil, models.ErrStorageNotConfigured
	}

	return s.uploads.CreateSignedPostUpload(ctx, userID, request)
}

func (s *PostService) CreateFromUploadedObject(ctx context.Context, userID uuid.UUID, request models.CreatePostFromUploadRequest) (*models.PostFeedItem, error) {
	objectPath := strings.TrimSpace(request.ObjectPath)
	captionValue, err := sanitizeOptionalCaption(request.Caption)
	if err != nil {
		return nil, err
	}

	var imageURL *string
	if objectPath != "" {
		if s.uploads == nil {
			return nil, models.ErrStorageNotConfigured
		}

		if !s.uploads.ObjectBelongsToUser(userID, objectPath) {
			return nil, models.NewValidationError(map[string]string{
				"object_path": "object_path is invalid",
			})
		}

		publicURL := s.uploads.PublicURLForObject(objectPath)
		imageURL = &publicURL
	}

	if imageURL == nil && captionValue == nil {
		return nil, models.NewValidationError(map[string]string{
			"post": "add an image or write a thought before posting",
		})
	}

	post, err := s.posts.Create(ctx, repositories.CreatePostParams{
		UserID:   userID,
		ImageURL: imageURL,
		Caption:  captionValue,
	})
	if err != nil {
		if objectPath != "" && s.uploads != nil {
			_ = s.uploads.DeleteByObjectPath(ctx, objectPath)
		}
		return nil, fmt.Errorf("create post from uploaded object: %w", err)
	}

	postFeedItem, err := s.posts.FindFeedByID(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created post: %w", err)
	}

	return postFeedItem, nil
}

func (s *PostService) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.PostFeedItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.posts.FeedByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list user posts: %w", err)
	}

	return posts, nil
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

	if s.uploads != nil && post.ImageURL != nil {
		_ = s.uploads.DeleteByPublicURL(ctx, *post.ImageURL)
	}

	return nil
}

func sanitizeOptionalCaption(caption string) (*string, error) {
	cleanCaption := utils.SanitizeText(caption)
	if cleanCaption == "" {
		return nil, nil
	}

	if len(cleanCaption) > 2200 {
		return nil, models.NewValidationError(map[string]string{
			"caption": "thought must be 2200 characters or fewer",
		})
	}

	return &cleanCaption, nil
}
