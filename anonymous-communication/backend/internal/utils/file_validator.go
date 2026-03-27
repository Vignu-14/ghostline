package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"

	"anonymous-communication/backend/internal/models"

	_ "golang.org/x/image/webp"
)

const (
	maxUploadSizeBytes = 5 * 1024 * 1024
	maxImageDimension  = 5000
)

type ValidatedImage struct {
	Bytes     []byte
	MIMEType  string
	Extension string
}

func ValidateImageFile(file multipart.File, header *multipart.FileHeader) (*ValidatedImage, error) {
	if file == nil || header == nil {
		return nil, models.NewValidationError(map[string]string{
			"image": "image file is required",
		})
	}

	if header.Size <= 0 {
		return nil, models.NewValidationError(map[string]string{
			"image": "image file is empty",
		})
	}

	if header.Size > maxUploadSizeBytes {
		return nil, models.NewValidationError(map[string]string{
			"image": "image must be 5MB or smaller",
		})
	}

	content, err := io.ReadAll(io.LimitReader(file, maxUploadSizeBytes+1))
	if err != nil {
		return nil, err
	}

	if int64(len(content)) > maxUploadSizeBytes {
		return nil, models.NewValidationError(map[string]string{
			"image": "image must be 5MB or smaller",
		})
	}

	mimeType := http.DetectContentType(content)
	extension := ""

	switch mimeType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/gif":
		extension = ".gif"
	case "image/webp":
		extension = ".webp"
	default:
		return nil, models.NewValidationError(map[string]string{
			"image": "only JPEG, PNG, GIF, or WebP images are allowed",
		})
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		return nil, models.NewValidationError(map[string]string{
			"image": "image data is invalid",
		})
	}

	if config.Width > maxImageDimension || config.Height > maxImageDimension {
		return nil, models.NewValidationError(map[string]string{
			"image": "image width and height must be 5000px or smaller",
		})
	}

	return &ValidatedImage{
		Bytes:     content,
		MIMEType:  mimeType,
		Extension: extension,
	}, nil
}
