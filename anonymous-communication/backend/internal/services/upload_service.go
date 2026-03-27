package services

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type UploadService struct {
	storage config.StorageConfig
	client  *http.Client
}

func NewUploadService(storage config.StorageConfig) *UploadService {
	return &UploadService{
		storage: storage,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *UploadService) UploadPostImage(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (string, error) {
	if !s.storage.Enabled() {
		return "", models.ErrStorageNotConfigured
	}

	validatedImage, err := utils.ValidateImageFile(file, header)
	if err != nil {
		return "", err
	}

	objectPath := fmt.Sprintf("posts/%s/%s%s", userID.String(), utils.NewUUID().String(), validatedImage.Extension)
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		url.PathEscape(s.storage.BucketName),
		escapeObjectPath(objectPath),
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(validatedImage.Bytes))
	if err != nil {
		return "", fmt.Errorf("create upload request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+s.storage.SupabaseServiceKey)
	request.Header.Set("apikey", s.storage.SupabaseServiceKey)
	request.Header.Set("Content-Type", validatedImage.MIMEType)
	request.Header.Set("x-upsert", "false")

	response, err := s.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("upload image to storage: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("upload image to storage: unexpected status %d", response.StatusCode)
	}

	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		s.storage.BucketName,
		objectPath,
	), nil
}

func (s *UploadService) DeleteByPublicURL(ctx context.Context, publicURL string) error {
	if !s.storage.Enabled() {
		return models.ErrStorageNotConfigured
	}

	objectPath, err := s.objectPathFromPublicURL(publicURL)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		url.PathEscape(s.storage.BucketName),
		escapeObjectPath(objectPath),
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create delete request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+s.storage.SupabaseServiceKey)
	request.Header.Set("apikey", s.storage.SupabaseServiceKey)

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("delete image from storage: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil
	}

	if response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("delete image from storage: unexpected status %d", response.StatusCode)
	}

	return nil
}

func (s *UploadService) objectPathFromPublicURL(publicURL string) (string, error) {
	publicPrefix := fmt.Sprintf("%s/storage/v1/object/public/%s/",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		s.storage.BucketName,
	)

	if !strings.HasPrefix(publicURL, publicPrefix) {
		return "", fmt.Errorf("public url does not belong to configured storage bucket")
	}

	return strings.TrimPrefix(publicURL, publicPrefix), nil
}

func escapeObjectPath(path string) string {
	segments := strings.Split(path, "/")
	for index, segment := range segments {
		segments[index] = url.PathEscape(segment)
	}

	return strings.Join(segments, "/")
}
