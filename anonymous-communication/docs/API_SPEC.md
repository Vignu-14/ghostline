# 📡 API Specification

Complete REST API and WebSocket specification for Ghostline Anonymous Messaging Platform.

---

## 📋 Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Codes](#error-codes)
- [Rate Limiting](#rate-limiting)
- [Authentication Endpoints](#authentication-endpoints)
- [User Endpoints](#user-endpoints)
- [Post Endpoints](#post-endpoints)
- [Message Endpoints](#message-endpoints)
- [Like Endpoints](#like-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Health Endpoints](#health-endpoints)
- [WebSocket](#websocket)

---

## Base URL

**Development:** `http://localhost:3000`

**Production:** `https://ghostline-backend-production-a17a.up.railway.app`

---

## Authentication

### JWT-based Authentication

All protected endpoints require a valid JWT token in the `auth_token` cookie.

**Token Format:**
```
Header: {
  "alg": "HS256",
  "typ": "JWT"
}

Payload: {
  "user_id": "uuid",
  "role": "user|admin",
  "exp": 1234567890
}
```

**Token Expiration:** 15 minutes

**Cookie Properties:**
- Name: `auth_token`
- HttpOnly: true (XSS protection)
- Secure: true (HTTPS only in production)
- SameSite: Strict (CSRF protection)
- Path: /
- Domain: (auto-set)

---

## Response Format

### Success Response (2xx)

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

### Error Response (4xx, 5xx)

```json
{
  "status": "error",
  "error": "Invalid credentials",
  "code": "INVALID_CREDENTIALS",
  "details": {
    "validation": [
      {
        "field": "password",
        "message": "Password must be at least 8 characters"
      }
    ]
  }
}
```

---

## Error Codes

| Status | Code | Message | Cause |
|--------|------|---------|-------|
| 400 | BAD_REQUEST | Bad request | Invalid parameters |
| 400 | VALIDATION_ERROR | Validation failed | Invalid input data |
| 401 | UNAUTHORIZED | Unauthorized | Missing/Invalid token |
| 401 | INVALID_CREDENTIALS | Invalid credentials | Wrong username/password |
| 403 | FORBIDDEN | Forbidden | Insufficient permissions |
| 404 | NOT_FOUND | Not found | Resource doesn't exist |
| 409 | CONFLICT | Conflict | Resource already exists |
| 429 | RATE_LIMITED | Too many requests | Rate limit exceeded |
| 500 | INTERNAL_ERROR | Internal server error | Server error |
| 503 | SERVICE_UNAVAILABLE | Service unavailable | Maintenance/down |

---

## Rate Limiting

### Limits per User

| Endpoint | Limit | Window |
|----------|-------|--------|
| POST /api/auth/login | 5 | 15 minutes |
| POST /api/auth/register | 3 | 1 hour |
| POST /api/messages | 10 | 1 second |
| POST /api/posts/:id/like | 100 | 1 hour |
| All other endpoints | 100 | 1 minute |

### Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1234567890
```

---

## Authentication Endpoints

### Register User

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user"
    },
    "message": "User created successfully"
  }
}
```

**Errors:**
- 400: Username already exists
- 400: Email already exists
- 400: Invalid email format
- 400: Password too weak

---

### Login User

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user"
    },
    "message": "Login successful"
  }
}
```

**Headers:**
```
Set-Cookie: auth_token=eyJhbG...; HttpOnly; Secure; SameSite=Strict; Path=/
```

**Errors:**
- 401: Invalid credentials
- 429: Too many login attempts

---

### Get Current User

```http
GET /api/auth/me
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 401: Token expired

---

### Logout

```http
POST /api/auth/logout
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "message": "Logged out successfully"
  }
}
```

---

## User Endpoints

### Get User Profile

```http
GET /api/users/:user_id
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-03-01T08:00:00Z"
    }
  }
}
```

**Errors:**
- 404: User not found

---

### Update User Profile

```http
PUT /api/users/:user_id
Content-Type: application/json
Cookie: auth_token=<jwt_token>

{
  "email": "newemail@example.com",
  "password": "NewSecurePass123!"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "newemail@example.com",
      "updated_at": "2026-03-28T10:30:00Z"
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 403: Cannot update another user's profile
- 400: Email already in use

---

### Get User Posts

```http
GET /api/users/:user_id/posts?page=1&limit=20
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "posts": [
      {
        "id": "post-uuid",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "caption": "My first post",
        "image_url": "https://storage.example.com/posts/image.jpg",
        "likes_count": 5,
        "created_at": "2026-03-28T09:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "total": 25,
      "limit": 20
    }
  }
}
```

---

### Get User Statistics

```http
GET /api/users/:user_id/stats
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "stats": {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "total_posts": 15,
      "total_likes_received": 42,
      "total_messages_sent": 128,
      "joined_date": "2026-03-01T08:00:00Z"
    }
  }
}
```

---

## Post Endpoints

### Get Feed

```http
GET /api/posts?page=1&limit=20
Cookie: auth_token=<jwt_token>
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Posts per page (default: 20, max: 100)
- `user_id`: Filter by user (optional)

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "posts": [
      {
        "id": "post-uuid",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "john_doe",
        "caption": "Amazing view!",
        "image_url": "https://storage.example.com/posts/image.jpg",
        "likes_count": 12,
        "user_liked": false,
        "created_at": "2026-03-28T09:15:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "total": 156,
      "limit": 20
    }
  }
}
```

---

### Create Post

```http
POST /api/posts
Content-Type: application/json
Cookie: auth_token=<jwt_token>

{
  "caption": "Beautiful sunset 🌅",
  "image_url": "https://storage.example.com/posts/sunset.jpg"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "post-uuid",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "caption": "Beautiful sunset 🌅",
      "image_url": "https://storage.example.com/posts/sunset.jpg",
      "likes_count": 0,
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 400: Image URL required

---

### Get Single Post

```http
GET /api/posts/:post_id
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "post-uuid",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "caption": "Beautiful sunset 🌅",
      "image_url": "https://storage.example.com/posts/sunset.jpg",
      "likes_count": 5,
      "user_liked": false,
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

**Errors:**
- 404: Post not found

---

### Delete Post

```http
DELETE /api/posts/:post_id
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "message": "Post deleted successfully"
  }
}
```

**Errors:**
- 401: Unauthorized
- 403: Cannot delete another user's post
- 404: Post not found

---

### Get Upload URL

```http
POST /api/posts/upload-url
Content-Type: application/json
Cookie: auth_token=<jwt_token>

{
  "file_name": "photo.jpg",
  "file_size": 1024000,
  "file_type": "image/jpeg"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "upload_url": "https://supabase.example.com/upload?token=xyz",
    "file_path": "posts/user-uuid/file-uuid.jpg",
    "expires_in": 3600
  }
}
```

**Errors:**
- 401: Unauthorized
- 400: File too large (max 5MB)
- 400: Invalid file type

---

### Finalize Post

```http
POST /api/posts/finalize
Content-Type: application/json
Cookie: auth_token=<jwt_token>

{
  "caption": "My photo",
  "image_url": "posts/user-uuid/file-uuid.jpg"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "post-uuid",
      "caption": "My photo",
      "image_url": "https://storage.example.com/posts/user-uuid/file-uuid.jpg",
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

---

## Message Endpoints

### Get Conversations

```http
GET /api/messages?limit=20
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "conversations": [
      {
        "conversation_id": "user-uuid-1",
        "username": "jane_doe",
        "last_message": "See you tomorrow!",
        "last_message_time": "2026-03-28T15:30:00Z",
        "unread_count": 2
      }
    ]
  }
}
```

---

### Get Conversation

```http
GET /api/messages/:user_id?limit=50
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "messages": [
      {
        "id": "msg-uuid",
        "sender_id": "550e8400-e29b-41d4-a716-446655440000",
        "receiver_id": "other-user-uuid",
        "content": "Hello!",
        "is_read": true,
        "created_at": "2026-03-28T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 125,
      "limit": 50
    }
  }
}
```

---

### Send Message

```http
POST /api/messages
Content-Type: application/json
Cookie: auth_token=<jwt_token>

{
  "receiver_id": "other-user-uuid",
  "content": "Hi, how are you?"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "message": {
      "id": "msg-uuid",
      "sender_id": "550e8400-e29b-41d4-a716-446655440000",
      "receiver_id": "other-user-uuid",
      "content": "Hi, how are you?",
      "is_read": false,
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 404: Receiver not found
- 429: Message rate limit exceeded

---

### Mark Message as Read

```http
PUT /api/messages/:message_id/read
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "message": {
      "id": "msg-uuid",
      "is_read": true,
      "updated_at": "2026-03-28T10:05:00Z"
    }
  }
}
```

---

## Like Endpoints

### Like Post

```http
POST /api/posts/:post_id/like
Cookie: auth_token=<jwt_token>
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "post-uuid",
      "likes_count": 13,
      "user_liked": true,
      "updated_at": "2026-03-28T10:15:00Z"
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 404: Post not found
- 409: Already liked

---

### Unlike Post

```http
DELETE /api/posts/:post_id/like
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "post-uuid",
      "likes_count": 12,
      "user_liked": false,
      "updated_at": "2026-03-28T10:20:00Z"
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 404: Post/Like not found

---

## Admin Endpoints

### List All Users

```http
GET /api/admin/users?page=1&limit=50
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "users": [
      {
        "id": "user-uuid",
        "username": "john_doe",
        "email": "john@example.com",
        "role": "user",
        "created_at": "2026-03-01T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "total": 150,
      "limit": 50
    }
  }
}
```

**Errors:**
- 401: Unauthorized
- 403: Admin access required

---

### Delete User

```http
DELETE /api/admin/users/:user_id
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "message": "User deleted successfully"
  }
}
```

**Errors:**
- 401: Unauthorized
- 403: Admin access required
- 404: User not found

---

### List All Posts

```http
GET /api/admin/posts?page=1&limit=50
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "posts": [
      {
        "id": "post-uuid",
        "user_id": "user-uuid",
        "caption": "Post content",
        "created_at": "2026-03-28T10:00:00Z"
      }
    ]
  }
}
```

---

### Delete Post

```http
DELETE /api/admin/posts/:post_id
Cookie: auth_token=<jwt_token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "message": "Post deleted successfully"
  }
}
```

---

## Health Endpoints

### Health Check

```http
GET /health
```

**Response (200):**
```json
{
  "status": "ok",
  "timestamp": "2026-03-28T10:00:00Z"
}
```

---

### Ready Check

```http
GET /health/ready
```

**Response (200):**
```json
{
  "status": "ready",
  "database": "connected",
  "timestamp": "2026-03-28T10:00:00Z"
}
```

---

## WebSocket

See [WEBSOCKET_PROTOCOL.md](./WEBSOCKET_PROTOCOL.md) for complete WebSocket documentation.

### Connection

```
ws://localhost:3000/ws/chat
wss://ghostline-backend-production-a17a.up.railway.app/ws/chat
```

### Authentication

WebSocket connection requires valid JWT token in auth token cookie.

### Message Format

**Send:**
```json
{
  "type": "message",
  "receiver_id": "other-user-uuid",
  "content": "Hello!"
}
```

**Receive:**
```json
{
  "type": "message",
  "id": "msg-uuid",
  "sender_id": "user-uuid",
  "receiver_id": "other-user-uuid",
  "content": "Hello!",
  "is_read": false,
  "created_at": "2026-03-28T10:00:00Z"
}
```

---

See also: [ARCHITECTURE.md](./ARCHITECTURE.md), [WEBSOCKET_PROTOCOL.md](./WEBSOCKET_PROTOCOL.md)
