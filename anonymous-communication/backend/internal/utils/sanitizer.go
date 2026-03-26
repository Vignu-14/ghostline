package utils

import "strings"

func SanitizeText(input string) string {
	return strings.TrimSpace(input)
}
