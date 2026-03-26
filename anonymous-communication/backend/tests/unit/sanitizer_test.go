package unit

import (
	"testing"

	"anonymous-communication/backend/internal/utils"
)

func TestSanitizeText(t *testing.T) {
	input := "  hello world  "

	got := utils.SanitizeText(input)
	if got != "hello world" {
		t.Fatalf("expected trimmed text, got %q", got)
	}
}
