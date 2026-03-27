package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"anonymous-communication/backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*models.Message, error) {
	const query = `
		INSERT INTO messages (sender_id, receiver_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, sender_id, receiver_id, content, is_read, created_at,
		          deleted_for_sender_at, deleted_for_receiver_at, deleted_for_everyone_at
	`

	message, err := scanMessage(r.db.QueryRow(ctx, query, senderID, receiverID, content))
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return message, nil
}

func (r *MessageRepository) Conversation(ctx context.Context, userID, otherUserID uuid.UUID, limit, offset int) ([]models.MessageResponse, error) {
	const query = `
		SELECT id, sender_id, receiver_id, content, is_read, created_at,
		       deleted_for_sender_at, deleted_for_receiver_at, deleted_for_everyone_at
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2 AND deleted_for_sender_at IS NULL)
		   OR (sender_id = $2 AND receiver_id = $1 AND deleted_for_receiver_at IS NULL)
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, userID, otherUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get conversation: %w", err)
	}
	defer rows.Close()

	messages := make([]models.MessageResponse, 0)
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message.ToResponse())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversation messages: %w", err)
	}

	return messages, nil
}

func (r *MessageRepository) ListConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.ConversationSummary, error) {
	const query = `
		WITH visible_messages AS (
			SELECT
				m.id,
				m.sender_id,
				m.receiver_id,
				m.content,
				m.is_read,
				m.created_at,
				m.deleted_for_everyone_at,
				CASE
					WHEN m.sender_id = $1 THEN m.receiver_id
					ELSE m.sender_id
				END AS other_user_id
			FROM messages m
			WHERE (m.sender_id = $1 AND m.deleted_for_sender_at IS NULL)
			   OR (m.receiver_id = $1 AND m.deleted_for_receiver_at IS NULL)
		),
		ranked_messages AS (
			SELECT
				vm.id,
				vm.sender_id,
				vm.receiver_id,
				vm.content,
				vm.is_read,
				vm.created_at,
				vm.deleted_for_everyone_at,
				vm.other_user_id,
				ROW_NUMBER() OVER (
					PARTITION BY vm.other_user_id
					ORDER BY vm.created_at DESC
				) AS rank
			FROM visible_messages vm
		),
		unread_counts AS (
			SELECT sender_id AS other_user_id, COUNT(*)::bigint AS unread_count
			FROM messages
			WHERE receiver_id = $1
			  AND is_read = false
			  AND deleted_for_receiver_at IS NULL
			GROUP BY sender_id
		)
		SELECT
			rm.id,
			rm.sender_id,
			rm.receiver_id,
			rm.content,
			rm.is_read,
			rm.created_at,
			rm.deleted_for_everyone_at,
			u.id,
			u.username,
			u.profile_picture_url,
			COALESCE(uc.unread_count, 0)
		FROM ranked_messages rm
		INNER JOIN users u ON u.id = rm.other_user_id
		LEFT JOIN unread_counts uc ON uc.other_user_id = rm.other_user_id
		WHERE rm.rank = 1
		ORDER BY rm.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list conversations: %w", err)
	}
	defer rows.Close()

	conversations := make([]models.ConversationSummary, 0)
	for rows.Next() {
		var (
			message      models.Message
			otherUserID  uuid.UUID
			profilePhoto sql.NullString
			summary      models.ConversationSummary
		)

		if err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Content,
			&message.IsRead,
			&message.CreatedAt,
			&message.DeletedForEveryoneAt,
			&otherUserID,
			&summary.User.Username,
			&profilePhoto,
			&summary.UnreadCount,
		); err != nil {
			return nil, fmt.Errorf("scan conversation summary: %w", err)
		}

		summary.User.ID = otherUserID.String()
		if profilePhoto.Valid {
			summary.User.ProfilePictureURL = &profilePhoto.String
		}
		summary.LastMessage = message.ToResponse()
		conversations = append(conversations, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversations: %w", err)
	}

	return conversations, nil
}

func (r *MessageRepository) MarkConversationAsRead(ctx context.Context, userID, otherUserID uuid.UUID) error {
	const query = `
		UPDATE messages
		SET is_read = true
		WHERE sender_id = $1
		  AND receiver_id = $2
		  AND is_read = false
		  AND deleted_for_receiver_at IS NULL
	`

	if _, err := r.db.Exec(ctx, query, otherUserID, userID); err != nil {
		return fmt.Errorf("mark conversation as read: %w", err)
	}

	return nil
}

func (r *MessageRepository) FindByIDsForUser(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) ([]models.Message, error) {
	const query = `
		SELECT id, sender_id, receiver_id, content, is_read, created_at,
		       deleted_for_sender_at, deleted_for_receiver_at, deleted_for_everyone_at
		FROM messages
		WHERE id = ANY($1::uuid[])
		  AND ($2 = sender_id OR $2 = receiver_id)
	`

	rows, err := r.db.Query(ctx, query, messageIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("find messages for user: %w", err)
	}
	defer rows.Close()

	messages := make([]models.Message, 0, len(messageIDs))
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}

		messages = append(messages, *message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user messages: %w", err)
	}

	return messages, nil
}

func (r *MessageRepository) DeleteForUser(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) (int64, error) {
	const query = `
		UPDATE messages
		SET deleted_for_sender_at = CASE
				WHEN sender_id = $1 AND deleted_for_sender_at IS NULL THEN NOW()
				ELSE deleted_for_sender_at
			END,
			deleted_for_receiver_at = CASE
				WHEN receiver_id = $1 AND deleted_for_receiver_at IS NULL THEN NOW()
				ELSE deleted_for_receiver_at
			END
		WHERE id = ANY($2::uuid[])
		  AND ($1 = sender_id OR $1 = receiver_id)
	`

	result, err := r.db.Exec(ctx, query, userID, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("delete messages for user: %w", err)
	}

	return result.RowsAffected(), nil
}

func (r *MessageRepository) DeleteForEveryone(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) (int64, error) {
	const query = `
		UPDATE messages
		SET deleted_for_everyone_at = COALESCE(deleted_for_everyone_at, NOW())
		WHERE id = ANY($2::uuid[])
		  AND sender_id = $1
	`

	result, err := r.db.Exec(ctx, query, userID, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("delete messages for everyone: %w", err)
	}

	return result.RowsAffected(), nil
}

func (r *MessageRepository) ClearConversationForUser(ctx context.Context, userID, otherUserID uuid.UUID) (int64, error) {
	const query = `
		UPDATE messages
		SET deleted_for_sender_at = CASE
				WHEN sender_id = $1 AND deleted_for_sender_at IS NULL THEN NOW()
				ELSE deleted_for_sender_at
			END,
			deleted_for_receiver_at = CASE
				WHEN receiver_id = $1 AND deleted_for_receiver_at IS NULL THEN NOW()
				ELSE deleted_for_receiver_at
			END
		WHERE (sender_id = $1 AND receiver_id = $2 AND deleted_for_sender_at IS NULL)
		   OR (sender_id = $2 AND receiver_id = $1 AND deleted_for_receiver_at IS NULL)
	`

	result, err := r.db.Exec(ctx, query, userID, otherUserID)
	if err != nil {
		return 0, fmt.Errorf("clear conversation for user: %w", err)
	}

	return result.RowsAffected(), nil
}

func (r *MessageRepository) DeleteConversationForEveryone(ctx context.Context, userID, otherUserID uuid.UUID) (int64, error) {
	const query = `
		UPDATE messages
		SET deleted_for_everyone_at = COALESCE(deleted_for_everyone_at, NOW())
		WHERE sender_id = $1
		  AND receiver_id = $2
	`

	result, err := r.db.Exec(ctx, query, userID, otherUserID)
	if err != nil {
		return 0, fmt.Errorf("delete sent conversation messages for everyone: %w", err)
	}

	return result.RowsAffected(), nil
}

func scanMessage(row pgx.Row) (*models.Message, error) {
	var message models.Message

	err := row.Scan(
		&message.ID,
		&message.SenderID,
		&message.ReceiverID,
		&message.Content,
		&message.IsRead,
		&message.CreatedAt,
		&message.DeletedForSenderAt,
		&message.DeletedForReceiverAt,
		&message.DeletedForEveryoneAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrMessageNotFound
		}

		return nil, fmt.Errorf("scan message: %w", err)
	}

	return &message, nil
}
