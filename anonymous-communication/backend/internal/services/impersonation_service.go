package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type impersonationUserRepository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type adminAuditRepository interface {
	CreateAuditLog(ctx context.Context, params models.CreateAdminAuditLogParams) error
}

type ImpersonationService struct {
	users  impersonationUserRepository
	audits adminAuditRepository
	jwt    config.JWTConfig
}

func NewImpersonationService(users impersonationUserRepository, audits adminAuditRepository, jwt config.JWTConfig) *ImpersonationService {
	return &ImpersonationService{
		users:  users,
		audits: audits,
		jwt:    jwt,
	}
}

func (s *ImpersonationService) Start(ctx context.Context, adminID uuid.UUID, input models.ImpersonateRequest, ipAddress string) (*models.AuthSession, error) {
	if err := utils.ValidateImpersonateRequest(input); err != nil {
		return nil, err
	}

	admin, err := s.users.FindByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("find admin: %w", err)
	}

	if admin.Role != config.RoleAdmin {
		return nil, models.ErrAdminOnly
	}

	if admin.ImpersonationPasswordHash == nil || strings.TrimSpace(*admin.ImpersonationPasswordHash) == "" {
		s.logFailure(ctx, adminID, nil, ipAddress, "impersonation_password_not_configured")
		return nil, models.ErrImpersonationNotConfigured
	}

	targetUserID, err := uuid.Parse(strings.TrimSpace(input.TargetUserID))
	if err != nil {
		return nil, models.NewValidationError(map[string]string{
			"target_user_id": "target_user_id must be a valid UUID",
		})
	}

	if targetUserID == adminID {
		s.logFailure(ctx, adminID, &targetUserID, ipAddress, "cannot_impersonate_self")
		return nil, models.ErrCannotImpersonateSelf
	}

	targetUser, err := s.users.FindByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.logFailure(ctx, adminID, &targetUserID, ipAddress, "target_user_not_found")
			return nil, err
		}

		return nil, fmt.Errorf("find target user: %w", err)
	}

	if err := utils.ComparePassword(*admin.ImpersonationPasswordHash, input.ImpersonationPassword); err != nil {
		s.logFailure(ctx, adminID, &targetUser.ID, ipAddress, "invalid_impersonation_password")
		return nil, models.ErrInvalidImpersonationPassword
	}

	token, err := utils.GenerateToken(s.jwt.Secret, s.jwt.Expiration, targetUser.ID, targetUser.Role, &admin.ID)
	if err != nil {
		return nil, fmt.Errorf("generate impersonation jwt: %w", err)
	}

	if err := s.audits.CreateAuditLog(ctx, models.CreateAdminAuditLogParams{
		AdminID:      admin.ID,
		TargetUserID: &targetUser.ID,
		Action:       "impersonate",
		IPAddress:    ipAddress,
		Metadata: map[string]any{
			"target_role": targetUser.Role,
		},
	}); err != nil {
		return nil, fmt.Errorf("create admin audit log: %w", err)
	}

	return &models.AuthSession{
		Token: token,
		User:  targetUser.ToResponse(),
	}, nil
}

func (s *ImpersonationService) Stop(ctx context.Context, targetUserID, impersonatorID uuid.UUID, ipAddress string) (*models.AuthSession, error) {
	admin, err := s.users.FindByID(ctx, impersonatorID)
	if err != nil {
		return nil, fmt.Errorf("find impersonator: %w", err)
	}

	if admin.Role != config.RoleAdmin {
		return nil, models.ErrAdminOnly
	}

	token, err := utils.GenerateToken(s.jwt.Secret, s.jwt.Expiration, admin.ID, admin.Role, nil)
	if err != nil {
		return nil, fmt.Errorf("generate admin jwt: %w", err)
	}

	if err := s.audits.CreateAuditLog(ctx, models.CreateAdminAuditLogParams{
		AdminID:      admin.ID,
		TargetUserID: &targetUserID,
		Action:       "impersonate_end",
		IPAddress:    ipAddress,
	}); err != nil {
		return nil, fmt.Errorf("create admin audit log: %w", err)
	}

	return &models.AuthSession{
		Token: token,
		User:  admin.ToResponse(),
	}, nil
}

func (s *ImpersonationService) logFailure(ctx context.Context, adminID uuid.UUID, targetUserID *uuid.UUID, ipAddress, reason string) {
	if s.audits == nil {
		return
	}

	_ = s.audits.CreateAuditLog(ctx, models.CreateAdminAuditLogParams{
		AdminID:      adminID,
		TargetUserID: targetUserID,
		Action:       "impersonate_failed",
		IPAddress:    ipAddress,
		Metadata: map[string]any{
			"reason": reason,
		},
	})
}
