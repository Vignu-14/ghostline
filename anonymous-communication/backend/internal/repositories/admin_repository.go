package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"anonymous-communication/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateAuditLog(ctx context.Context, params models.CreateAdminAuditLogParams) error {
	const query = `
		INSERT INTO admin_audit_logs (admin_id, target_user_id, action, ip_address, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`

	var targetUserID any
	if params.TargetUserID != nil {
		targetUserID = *params.TargetUserID
	}

	var metadata any
	if len(params.Metadata) > 0 {
		payload, err := json.Marshal(params.Metadata)
		if err != nil {
			return fmt.Errorf("marshal admin audit metadata: %w", err)
		}

		metadata = payload
	}

	if _, err := r.db.Exec(ctx, query, params.AdminID, targetUserID, params.Action, params.IPAddress, metadata); err != nil {
		return fmt.Errorf("insert admin audit log: %w", err)
	}

	return nil
}
