package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

const maxDirectUploadSizeBytes = 5 * 1024 * 1024

type UploadService struct {
	storage config.StorageConfig
	client  *http.Client
}

type signedUploadResponse struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func NewUploadService(storage config.StorageConfig) *UploadService {
	return &UploadService{
		storage: storage,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *UploadService) CreateSignedPostUpload(ctx context.Context, userID uuid.UUID, request models.CreatePostUploadRequest) (*models.PostUploadTarget, error) {
	if !s.storage.Enabled() {
		return nil, models.ErrStorageNotConfigured
	}

	extension, err := validateUploadMetadata(request)
	if err != nil {
		return nil, err
	}

	objectPath := fmt.Sprintf("posts/%s/%s%s", userID.String(), utils.NewUUID().String(), extension)
	endpoint := fmt.Sprintf("%s/storage/v1/object/upload/sign/%s/%s",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		url.PathEscape(s.storage.BucketName),
		escapeObjectPath(objectPath),
	)

	requestBody := bytes.NewBufferString(`{"upsert":false}`)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, requestBody)
	if err != nil {
		return nil, fmt.Errorf("create signed upload request: %w", err)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+s.storage.SupabaseServiceKey)
	httpRequest.Header.Set("apikey", s.storage.SupabaseServiceKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("sign upload request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("sign upload request: unexpected status %d", response.StatusCode)
	}

	var payload signedUploadResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode signed upload response: %w", err)
	}

	uploadURL := fmt.Sprintf("%s/storage/v1%s",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		payload.URL,
	)

	return &models.PostUploadTarget{
		ObjectPath: objectPath,
		UploadURL:  uploadURL,
		PublicURL:  s.PublicURLForObject(objectPath),
		Method:     http.MethodPut,
	}, nil
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

func (s *UploadService) DeleteByObjectPath(ctx context.Context, objectPath string) error {
	if !s.storage.Enabled() {
		return models.ErrStorageNotConfigured
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

func (s *UploadService) PublicURLForObject(objectPath string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		strings.TrimSuffix(s.storage.SupabaseURL, "/"),
		s.storage.BucketName,
		objectPath,
	)
}

func (s *UploadService) ObjectBelongsToUser(userID uuid.UUID, objectPath string) bool {
	expectedPrefix := fmt.Sprintf("posts/%s/", userID.String())
	return strings.HasPrefix(strings.TrimSpace(objectPath), expectedPrefix)
}

func escapeObjectPath(path string) string {
	segments := strings.Split(path, "/")
	for index, segment := range segments {
		segments[index] = url.PathEscape(segment)
	}

	return strings.Join(segments, "/")
}

func validateUploadMetadata(request models.CreatePostUploadRequest) (string, error) {
	if request.FileSize <= 0 {
		return "", models.NewValidationError(map[string]string{
			"file_size": "file size must be greater than zero",
		})
	}

	if request.FileSize > maxDirectUploadSizeBytes {
		return "", models.NewValidationError(map[string]string{
			"file_size": "image must be 5MB or smaller",
		})
	}

	switch strings.TrimSpace(request.ContentType) {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	case "image/webp":
		return ".webp", nil
	}

	extension := strings.ToLower(strings.TrimSpace(filepath.Ext(request.FileName)))
	switch extension {
	case ".jpg", ".jpeg":
		return ".jpg", nil
	case ".png":
		return ".png", nil
	case ".gif":
		return ".gif", nil
	case ".webp":
		return ".webp", nil
	default:
		return "", models.NewValidationError(map[string]string{
			"content_type": "only JPEG, PNG, GIF, or WebP images are allowed",
		})
	}
}
