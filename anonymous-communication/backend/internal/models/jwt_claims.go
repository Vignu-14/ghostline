package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID         string  `json:"user_id"`
	Role           string  `json:"role"`
	ImpersonatorID *string `json:"impersonator_id,omitempty"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) UserUUID() (uuid.UUID, error) {
	return uuid.Parse(c.UserID)
}
