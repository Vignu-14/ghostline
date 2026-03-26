package unit

import (
	"testing"

	"anonymous-communication/backend/internal/utils"
)

func TestHashPasswordAndComparePassword(t *testing.T) {
	hash, err := utils.HashPassword("StrongPass1!")
	if err != nil {
		t.Fatalf("expected password hash, got error: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty password hash")
	}

	if err := utils.ComparePassword(hash, "StrongPass1!"); err != nil {
		t.Fatalf("expected matching password, got error: %v", err)
	}

	if err := utils.ComparePassword(hash, "WrongPass1!"); err == nil {
		t.Fatal("expected wrong password comparison to fail")
	}
}
