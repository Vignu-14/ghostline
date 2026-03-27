# 🗄️ Database Schema

Complete database schema for Ghostline Anonymous Messaging Platform.

---

## 📋 Overview

The Ghostline database uses PostgreSQL with the following tables:

1. **users** - User accounts and authentication
2. **posts** - User posts with images
3. **messages** - Private messages between users
4. **likes** - Post likes
5. **auth_logs** - Login/logout audit trail

---

## Users Table

Stores user account information.

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  impersonation_password_hash VARCHAR(255),
  role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
  profile_picture_url TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);
```

### Fields

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique user identifier |
| username | VARCHAR(50) | UNIQUE, NOT NULL | Login name (3-50 chars) |
| email | VARCHAR(100) | UNIQUE, NOT NULL | Email address |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| impersonation_password_hash | VARCHAR(255) | nullable | Hash for admin impersonation |
| role | VARCHAR(20) | 'user' \| 'admin' | User role for permissions |
| profile_picture_url | TEXT | nullable | Supabase storage URL |
| created_at | TIMESTAMP | DEFAULT NOW() | Account creation date |
| updated_at | TIMESTAMP | DEFAULT NOW() | Last update date |

---

## Posts Table

Stores user posts with image URLs.

```sql
CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  caption TEXT,
  image_url TEXT NOT NULL,
  likes_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

### Fields

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique post identifier |
| user_id | UUID | FOREIGN KEY | Reference to users.id |
| caption | TEXT | nullable | Post description |
| image_url | TEXT | NOT NULL | Supabase storage URL |
| likes_count | INTEGER | DEFAULT 0 | Cached like count |
| created_at | TIMESTAMP | DEFAULT NOW() | Post creation date |
| updated_at | TIMESTAMP | DEFAULT NOW() | Last update date |

### Cascade Delete

When a user is deleted, all their posts are automatically deleted:
```sql
REFERENCES users(id) ON DELETE CASCADE
```

---

## Messages Table

Stores private messages between users.

```sql
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_by_sender BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_by_receiver BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX idx_messages_conversation ON messages(sender_id, receiver_id);
```

### Fields

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique message identifier |
| sender_id | UUID | FOREIGN KEY | Sender user ID |
| receiver_id | UUID | FOREIGN KEY | Receiver user ID |
| content | TEXT | NOT NULL | Message content (sanitized) |
| is_read | BOOLEAN | DEFAULT FALSE | Read status |
| deleted_by_sender | BOOLEAN | DEFAULT FALSE | Soft delete flag |
| deleted_by_receiver | BOOLEAN | DEFAULT FALSE | Soft delete flag |
| created_at | TIMESTAMP | DEFAULT NOW() | Send time |
| updated_at | TIMESTAMP | DEFAULT NOW() | Last update time |

### Soft Deletes

Messages are soft-deleted (not removed from DB):
- If sender deletes: `deleted_by_sender = TRUE`
- If receiver deletes: `deleted_by_receiver = TRUE`
- Only both true completely hides message

---

## Likes Table

Stores post likes (many-to-many relationship).

```sql
CREATE TABLE likes (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_post_id ON likes(post_id);
```

### Fields

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| user_id | UUID | FOREIGN KEY | User who liked |
| post_id | UUID | FOREIGN KEY | Post being liked |
| created_at | TIMESTAMP | DEFAULT NOW() | Like timestamp |
| PRIMARY KEY | (user_id, post_id) | Composite | One like per user per post |

---

## Auth Logs Table

Stores authentication events for audit trail.

```sql
CREATE TABLE auth_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  action VARCHAR(50) NOT NULL CHECK (action IN ('login', 'logout', 'register', 'failed_login')),
  status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failed')),
  ip_address INET,
  user_agent TEXT,
  error_message TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auth_logs_user_id ON auth_logs(user_id);
CREATE INDEX idx_auth_logs_created_at ON auth_logs(created_at DESC);
CREATE INDEX idx_auth_logs_action ON auth_logs(action);
```

### Fields

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique log entry ID |
| user_id | UUID | FOREIGN KEY | User ID (nullable for failed login) |
| action | VARCHAR(50) | 'login' \| 'logout' \| 'register' \| 'failed_login' | Action performed |
| status | VARCHAR(20) | 'success' \| 'failed' | Result status |
| ip_address | INET | nullable | Client IP address |
| user_agent | TEXT | nullable | Browser/client info |
| error_message | TEXT | nullable | Error description if failed |
| created_at | TIMESTAMP | DEFAULT NOW() | Event timestamp |

---

## Views (Optional)

### User Statistics View

```sql
CREATE VIEW user_stats AS
SELECT
  u.id,
  u.username,
  COUNT(DISTINCT p.id) as total_posts,
  COALESCE(SUM(p.likes_count), 0) as total_likes_received,
  COUNT(DISTINCT m.id) as total_messages_sent,
  u.created_at as joined_date
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
LEFT JOIN messages m ON u.id = m.sender_id
GROUP BY u.id, u.username, u.created_at;
```

---

## Relationships

### Entity-Relationship Diagram

```
users
  ├── POST (1 to Many)
  │   ├── id
  │   ├── user_id → users.id
  │   └── image_url
  │
  ├── MESSAGE (1 to Many as sender)
  │   ├── id
  │   ├── sender_id → users.id
  │   └── receiver_id → users.id
  │
  ├── MESSAGE (1 to Many as receiver)
  │   └── similar structure
  │
  └── LIKE (Many to Many)
      └── (user_id, post_id)
```

---

## Indexing Strategy

### Primary Indexes

1. **users(username)** - Fast login lookups
2. **users(email)** - Email verification
3. **posts(user_id)** - Get user's posts
4. **posts(created_at DESC)** - Feed ordering
5. **messages(sender_id, receiver_id)** - Conversation queries
6. **messages(created_at DESC)** - Recent messages
7. **likes(user_id)** - User's likes
8. **likes(post_id)** - Post's likes

### Composite Indexes

```sql
-- Fast conversation lookups
CREATE INDEX idx_messages_conversation ON messages(sender_id, receiver_id);

-- Fast stats queries
CREATE INDEX idx_posts_user_created ON posts(user_id, created_at DESC);
```

---

## Query Examples

### Get User Profile

```sql
SELECT id, username, email, role, created_at
FROM users
WHERE username = 'john_doe';
```

### Get User's Posts

```sql
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 20;
```

### Get Conversation

```sql
SELECT * FROM messages
WHERE (sender_id = $1 AND receiver_id = $2)
   OR (sender_id = $2 AND receiver_id = $1)
ORDER BY created_at DESC
LIMIT 50;
```

### Get Post with Like Count

```sql
SELECT
  p.id,
  p.caption,
  p.image_url,
  p.likes_count,
  EXISTS(SELECT 1 FROM likes WHERE post_id = p.id AND user_id = $1) as user_liked
FROM posts p
WHERE p.id = $2;
```

### Get Feed (First 20 Posts)

```sql
SELECT
  p.id,
  p.user_id,
  u.username,
  p.caption,
  p.image_url,
  p.likes_count,
  EXISTS(SELECT 1 FROM likes WHERE post_id = p.id AND user_id = $1) as user_liked,
  p.created_at
FROM posts p
JOIN users u ON p.user_id = u.id
ORDER BY p.created_at DESC
LIMIT 20
OFFSET 0;
```

### Get Unread Messages Count

```sql
SELECT COUNT(*) as unread_count
FROM messages
WHERE receiver_id = $1
  AND is_read = FALSE
  AND deleted_by_receiver = FALSE;
```

---

## Data Constraints

### Username

- **Length:** 3-50 characters
- **Pattern:** Alphanumeric, underscore, hyphen
- **Unique:** Across all users
- **Case-sensitive**

### Email

- **Format:** Valid email format
- **Unique:** Across all users
- **Case-insensitive**

### Password

- **Length:** 8+ characters
- **Requirements:**
  - At least 1 uppercase letter
  - At least 1 number
  - At least 1 special character
- **Not stored:** Only bcrypt hash stored

### Post Caption

- **Max length:** 2000 characters
- **Content:** Sanitized (no HTML tags)
- **Optional:** Can be empty if image present

### Message

- **Max length:** 5000 characters
- **Content:** Sanitized (no HTML tags)
- **Required:** Cannot be empty

---

## Migrations

### Initial Schema Migration

```sql
-- Create users table
CREATE TABLE users (...);

-- Create posts table
CREATE TABLE posts (...);

-- Create messages table
CREATE TABLE messages (...);

-- Create likes table
CREATE TABLE likes (...);

-- Create auth_logs table
CREATE TABLE auth_logs (...);
```

Automatically runs on server startup via Go migrations.

---

## Backup & Recovery

### Daily Backups

Supabase automatically backs up database daily:
- ✅ Point-in-time recovery (30 days)
- ✅ Geographic redundancy
- ✅ Automatic failover

### Manual Backup

```bash
pg_dump postgresql://user:pass@host:5432/ghostline > backup.sql
```

### Restore from Backup

```bash
psql postgresql://user:pass@host:5432/ghostline < backup.sql
```

---

See also: [ARCHITECTURE.md](./ARCHITECTURE.md), [DEPLOYMENT.md](./DEPLOYMENT.md)
