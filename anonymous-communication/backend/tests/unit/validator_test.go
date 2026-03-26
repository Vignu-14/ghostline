package unit

import (
	"errors"
	"testing"

	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"
)

func TestValidateRegisterRequest(t *testing.T) {
	valid := models.RegisterRequest{
		Username: "demo_user",
		Email:    "demo@example.com",
		Password: "StrongPass1!",
	}

	if err := utils.ValidateRegisterRequest(valid); err != nil {
		t.Fatalf("expected valid register request, got error: %v", err)
	}
}

func TestValidateRegisterRequestReturnsFieldErrors(t *testing.T) {
	invalid := models.RegisterRequest{
		Username: "x",
		Email:    "not-an-email",
		Password: "weak",
	}

	err := utils.ValidateRegisterRequest(invalid)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *models.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected validation error type, got %T", err)
	}

	if len(validationErr.Fields) != 3 {
		t.Fatalf("expected 3 field errors, got %d", len(validationErr.Fields))
	}
}
