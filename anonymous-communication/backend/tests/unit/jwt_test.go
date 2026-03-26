package unit

import (
	"testing"
	"time"

	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

func TestGenerateAndParseToken(t *testing.T) {
	userID := uuid.New()

	token, err := utils.GenerateToken("test-secret", time.Minute, userID, "user", nil)
	if err != nil {
		t.Fatalf("expected token generation to succeed, got error: %v", err)
	}

	claims, err := utils.ParseToken("test-secret", token)
	if err != nil {
		t.Fatalf("expected token parsing to succeed, got error: %v", err)
	}

	if claims.UserID != userID.String() {
		t.Fatalf("expected user id %s, got %s", userID, claims.UserID)
	}

	if claims.Role != "user" {
		t.Fatalf("expected role user, got %s", claims.Role)
	}
}
