-- Drop RLS Policies
DROP POLICY IF EXISTS messages_own ON messages;
DROP POLICY IF EXISTS likes_delete_own ON likes;
DROP POLICY IF EXISTS likes_insert_own ON likes;
DROP POLICY IF EXISTS likes_select_all ON likes;
DROP POLICY IF EXISTS posts_delete_own ON posts;
DROP POLICY IF EXISTS posts_insert_own ON posts;
DROP POLICY IF EXISTS posts_select_all ON posts;
DROP POLICY IF EXISTS users_update_own ON users;
DROP POLICY IF EXISTS users_select_own ON users;

-- Disable RLS
ALTER TABLE messages DISABLE ROW LEVEL SECURITY;
ALTER TABLE likes DISABLE ROW LEVEL SECURITY;
ALTER TABLE posts DISABLE ROW LEVEL SECURITY;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;

-- Drop Indexes
DROP INDEX IF EXISTS idx_auth_logs_ip;
DROP INDEX IF EXISTS idx_auth_logs_user_id;
DROP INDEX IF EXISTS idx_admin_audit_logs_created_at;
DROP INDEX IF EXISTS idx_admin_audit_logs_target_user_id;
DROP INDEX IF EXISTS idx_admin_audit_logs_admin_id;
DROP INDEX IF EXISTS idx_messages_unread;
DROP INDEX IF EXISTS idx_messages_conversation;
DROP INDEX IF EXISTS idx_likes_post_id;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
