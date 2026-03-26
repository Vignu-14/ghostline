package utils

import (
	"errors"
	"fmt"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(secret string, expiration time.Duration, userID uuid.UUID, role string, impersonatorID *uuid.UUID) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is required")
	}

	issuedAt := time.Now().UTC()

	var impersonator *string
	if impersonatorID != nil {
		value := impersonatorID.String()
		impersonator = &value
	}

	claims := models.JWTClaims{
		UserID:         userID.String(),
		Role:           role,
		ImpersonatorID: impersonator,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.AppName,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signedToken, nil
}

func ParseToken(secret, tokenString string) (*models.JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}

		return []byte(secret), nil
	}, jwt.WithIssuer(config.AppName))
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*models.JWTClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
