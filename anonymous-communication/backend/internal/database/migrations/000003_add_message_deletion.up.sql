ALTER TABLE messages
    ADD COLUMN deleted_for_sender_at TIMESTAMP,
    ADD COLUMN deleted_for_receiver_at TIMESTAMP,
    ADD COLUMN deleted_for_everyone_at TIMESTAMP;

CREATE INDEX idx_messages_sender_visibility
    ON messages(sender_id, receiver_id, created_at DESC)
    WHERE deleted_for_sender_at IS NULL;

CREATE INDEX idx_messages_receiver_visibility
    ON messages(receiver_id, sender_id, created_at DESC)
    WHERE deleted_for_receiver_at IS NULL;

CREATE INDEX idx_messages_unread_visible
    ON messages(receiver_id, sender_id, created_at DESC)
    WHERE is_read = false AND deleted_for_receiver_at IS NULL;
