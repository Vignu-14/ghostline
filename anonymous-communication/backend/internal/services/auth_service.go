package services

import (
	"context"
	"errors"
	"fmt"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/repositories"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type authUserRepository interface {
	Create(ctx context.Context, params repositories.CreateUserParams) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

type authLogRepository interface {
	Create(ctx context.Context, params models.CreateAuthLogParams) error
}

type AuthService struct {
	users    authUserRepository
	authLogs authLogRepository
	jwt      config.JWTConfig
}

func NewAuthService(users authUserRepository, authLogs authLogRepository, jwt config.JWTConfig) *AuthService {
	return &AuthService{
		users:    users,
		authLogs: authLogs,
		jwt:      jwt,
	}
}

func (s *AuthService) Register(ctx context.Context, input models.RegisterRequest) (*models.AuthSession, error) {
	if err := utils.ValidateRegisterRequest(input); err != nil {
		return nil, err
	}

	username := utils.NormalizeUsername(input.Username)
	email := utils.NormalizeEmail(input.Email)

	usernameExists, err := s.users.UsernameExists(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("check username availability: %w", err)
	}
	if usernameExists {
		return nil, models.ErrUsernameTaken
	}

	emailExists, err := s.users.EmailExists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email availability: %w", err)
	}
	if emailExists {
		return nil, models.ErrEmailTaken
	}

	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, repositories.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         config.RoleUser,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := utils.GenerateToken(s.jwt.Secret, s.jwt.Expiration, user.ID, user.Role, nil)
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
	}

	return &models.AuthSession{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginRequest, ipAddress, userAgent string) (*models.AuthSession, error) {
	if err := utils.ValidateLoginRequest(input); err != nil {
		return nil, err
	}

	username := utils.NormalizeUsername(input.Username)

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.logFailure(ctx, nil, ipAddress, userAgent, "user not found")
			return nil, models.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("find user by username: %w", err)
	}

	if err := utils.ComparePassword(user.PasswordHash, input.Password); err != nil {
		s.logFailure(ctx, &user.ID, ipAddress, userAgent, "wrong password")
		return nil, models.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(s.jwt.Secret, s.jwt.Expiration, user.ID, user.Role, nil)
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
	}

	s.logSuccess(ctx, &user.ID, ipAddress, userAgent)

	return &models.AuthSession{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *AuthService) logSuccess(ctx context.Context, userID *uuid.UUID, ipAddress, userAgent string) {
	if s.authLogs == nil {
		return
	}

	_ = s.authLogs.Create(ctx, models.CreateAuthLogParams{
		UserID:    userID,
		Status:    "success",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

func (s *AuthService) logFailure(ctx context.Context, userID *uuid.UUID, ipAddress, userAgent, reason string) {
	if s.authLogs == nil {
		return
	}

	_ = s.authLogs.Create(ctx, models.CreateAuthLogParams{
		UserID:        userID,
		Status:        "failed",
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		FailureReason: reason,
	})
}
