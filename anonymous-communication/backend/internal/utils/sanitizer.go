package utils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var strictSanitizer = bluemonday.StrictPolicy()

func SanitizeText(input string) string {
	return strings.TrimSpace(strictSanitizer.Sanitize(input))
}
