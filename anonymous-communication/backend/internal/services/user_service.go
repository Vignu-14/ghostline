package services

import (
	"context"
	"fmt"

	"anonymous-communication/backend/internal/models"

	"github.com/google/uuid"
)

type userRepository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type UserService struct {
	users userRepository
}

func NewUserService(users userRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}
