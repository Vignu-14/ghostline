package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(strings.TrimSpace(value))
}
