DROP INDEX IF EXISTS idx_messages_unread_visible;
DROP INDEX IF EXISTS idx_messages_receiver_visibility;
DROP INDEX IF EXISTS idx_messages_sender_visibility;

ALTER TABLE messages
    DROP COLUMN IF EXISTS deleted_for_everyone_at,
    DROP COLUMN IF EXISTS deleted_for_receiver_at,
    DROP COLUMN IF EXISTS deleted_for_sender_at;
