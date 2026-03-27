-- Indexes for performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_likes_post_id ON likes(post_id);
CREATE INDEX idx_messages_conversation ON messages(sender_id, receiver_id, created_at DESC);
CREATE INDEX idx_messages_unread ON messages(receiver_id) WHERE is_read = false;
CREATE INDEX idx_auth_logs_user_id ON auth_logs(user_id);
CREATE INDEX idx_auth_logs_ip ON auth_logs(ip_address);
CREATE INDEX idx_admin_audit_logs_admin_id ON admin_audit_logs(admin_id);
CREATE INDEX idx_admin_audit_logs_target_user_id ON admin_audit_logs(target_user_id);
CREATE INDEX idx_admin_audit_logs_created_at ON admin_audit_logs(created_at DESC);

-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE posts ENABLE ROW LEVEL SECURITY;
ALTER TABLE likes ENABLE ROW LEVEL SECURITY;
ALTER TABLE messages ENABLE ROW LEVEL SECURITY;

-- RLS Policies
CREATE POLICY users_select_own ON users FOR SELECT USING (true);
CREATE POLICY users_update_own ON users FOR UPDATE USING (id = current_setting('app.user_id', true)::uuid);

CREATE POLICY posts_select_all ON posts FOR SELECT USING (true);
CREATE POLICY posts_insert_own ON posts FOR INSERT WITH CHECK (user_id = current_setting('app.user_id', true)::uuid);
CREATE POLICY posts_delete_own ON posts FOR DELETE USING (user_id = current_setting('app.user_id', true)::uuid);

CREATE POLICY likes_select_all ON likes FOR SELECT USING (true);
CREATE POLICY likes_insert_own ON likes FOR INSERT WITH CHECK (user_id = current_setting('app.user_id', true)::uuid);
CREATE POLICY likes_delete_own ON likes FOR DELETE USING (user_id = current_setting('app.user_id', true)::uuid);

CREATE POLICY messages_own ON messages FOR ALL USING (
    sender_id = current_setting('app.user_id', true)::uuid OR 
    receiver_id = current_setting('app.user_id', true)::uuid
);
