# 🏗️ System Architecture

Complete technical architecture for the Ghostline Anonymous Messaging Platform.

---

## 📊 System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Browser / Client                          │
│                  (React 19 + TypeScript)                     │
│                                                              │
│  URL: https://ghostline-frontend-five.vercel.app            │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
      HTTP          WebSocket      API Calls
     (REST)         (Real-time)    (Axios)
        │              │              │
┌───────▼──────────────▼──────────────▼──────────────────────┐
│                                                              │
│              Backend (Go + Fiber Framework)                │
│                                                              │
│  URL: https://ghostline-backend-production-a17a.up.railway.app
│  Port: 3000 (dev), 443 (prod)                              │
│                                                              │
│  ├── HTTP Handlers (REST API)                              │
│  ├── WebSocket Hub (Real-time messaging)                   │
│  ├── Middleware (Auth, CORS, Logging)                      │
│  ├── Services (Business Logic)                             │
│  └── Repositories (Database Layer)                         │
│                                                              │
└───────┬──────────────┬──────────────┬──────────────────────┘
        │              │              │
        │          Database        File Storage
        │          Connection       (Upload/Download)
        │                           │
┌───────▼──────────────────┐  ┌─────▼──────────────┐
│   PostgreSQL Database     │  │ Supabase Storage   │
│   (Supabase Hosted)       │  │ (S3-compatible)    │
│                           │  │                    │
│  • Users table            │  │ Buckets:           │
│  • Posts table            │  │ - user-uploads/   │
│  • Messages table         │  │ - posts/          │
│  • Likes table            │  │ - profiles/       │
│  • Auth logs table        │  │                    │
└───────────────────────────┘  └────────────────────┘
```

---

## 🔄 Request/Response Flow

### Authentication Flow

```
User Input (LoginForm)
    ↓
POST /api/auth/login
    ↓
Backend: authService.login()
    ↓
Validate credentials
    ↓
Hash password with bcrypt
    ↓
Password matches?
    ├─ Yes → Generate JWT token
    │        Set secure cookie
    │        Return user data
    └─ No  → Return 401 Unauthorized
    ↓
Frontend: Store user in AuthContext
    ↓
Protected routes now accessible
    ↓
All subsequent requests include auth_token cookie
```

### Post Creation Flow

```
User selects image + types caption
    ↓
POST /api/posts/upload-url (get signed URL)
    ↓
Frontend: Upload to Supabase Storage
    ↓
PUT <signed-url> with image data
    ↓
Supabase: Store file
    ↓
Frontend: POST /api/posts/finalize
    ↓
Backend: Create post record in database
    ↓
Post appears in feed
```

### Real-time Chat Flow

```
User clicks "Message"
    ↓
Upgrade HTTP to WebSocket
    ↓
ws://backend/ws/chat
    ↓
Backend: Validate JWT from cookie
    ↓
Add connection to Hub
    ↓
Send "connected" event
    ↓
User A types message
    ↓
Client sends: {type: "message", receiver_id, content}
    ↓
Backend: Validate & rate limit
    ↓
Save to database
    ↓
Broadcast to User B (if online)
    ├─ Online  → Message delivered instantly
    └─ Offline → Message stored, retrieved on next login
    ↓
User B receives via WebSocket
    ↓
UI updates in real-time
```

### Like Flow

```
User clicks Like button
    ↓
POST /api/posts/:id/like
    ↓
Backend: Check if already liked
    ├─ Yes → Return 409 Conflict
    └─ No  → Create like record
    ↓
Increment likes_count
    ↓
Return updated post
    ↓
Frontend: Update UI immediately
    ↓
Like button changes state
    ↓
Count increments
```

---

## 🗄️ Database Schema

### Users Table

```
users
├── id (UUID) - Primary Key
├── username (VARCHAR 50) - UNIQUE, indexed
├── email (VARCHAR 100) - UNIQUE, indexed
├── password_hash (VARCHAR 60) - Bcrypt hash
├── role (VARCHAR 20) - 'user' or 'admin'
├── impersonation_password_hash (VARCHAR 60) - For admin impersonation
├── profile_picture_url (TEXT) - nullable
├── created_at (TIMESTAMP) - indexed
└── updated_at (TIMESTAMP)

Indexes:
  - username_idx on (username)
  - email_idx on (email)
  - created_at_idx on (created_at)
```

### Posts Table

```
posts
├── id (UUID) - Primary Key
├── user_id (UUID) - Foreign Key → users.id
├── caption (TEXT) - nullable
├── image_url (TEXT) - Supabase storage URL
├── likes_count (INTEGER) - Cached value
├── created_at (TIMESTAMP) - indexed
└── updated_at (TIMESTAMP)

Indexes:
  - user_id_idx on (user_id)
  - created_at_idx on (created_at)
```

### Messages Table

```
messages
├── id (UUID) - Primary Key
├── sender_id (UUID) - Foreign Key → users.id
├── receiver_id (UUID) - Foreign Key → users.id
├── content (TEXT) - Sanitized
├── is_read (BOOLEAN)
├── deleted_by_sender (BOOLEAN)
├── deleted_by_receiver (BOOLEAN)
├── created_at (TIMESTAMP) - indexed
└── updated_at (TIMESTAMP)

Indexes:
  - sender_id_idx on (sender_id)
  - receiver_id_idx on (receiver_id)
  - created_at_idx on (created_at)
  - composite_idx on (sender_id, receiver_id)
```

### Likes Table

```
likes
├── user_id (UUID) - Primary Key (part 1)
├── post_id (UUID) - Primary Key (part 2)
├── created_at (TIMESTAMP)

Indexes:
  - user_id_idx on (user_id)
  - post_id_idx on (post_id)
```

### Auth Logs Table

```
auth_logs
├── id (UUID) - Primary Key
├── user_id (UUID) - Foreign Key → users.id
├── action (VARCHAR) - 'login', 'logout', 'register'
├── status (VARCHAR) - 'success', 'failed'
├── ip_address (INET)
├── user_agent (TEXT)
├── created_at (TIMESTAMP) - indexed

Indexes:
  - user_id_idx on (user_id)
  - created_at_idx on (created_at)
```

---

## 🔐 Security Architecture

### Authentication Layer

```
Request comes in
    ↓
Extract auth_token cookie
    ↓
Decode JWT using JWT_SECRET
    ↓
Validate signature
    ↓
Check expiration (15 minutes)
    ↓
Extract user_id, role from token
    ↓
Add to request context
    ↓
Proceed to handler
```

### Password Security

```
User enters password
    ↓
Validate format:
  • 8+ characters
  • At least 1 uppercase
  • At least 1 number
  • At least 1 special character
    ↓
Hash with bcrypt (cost 12)
    ↓
Store hash in database
    ↓
Never store plaintext
```

### Input Sanitization

```
User input (post caption, message)
    ↓
Check length limits
    ↓
Sanitize with Bluemonday:
  • Remove HTML tags
  • Remove scripts
  • Remove dangerous attributes
    ↓
Store cleaned content
```

### SQL Injection Prevention

```
All database queries use parameterized statements:
    
// SAFE ✅
db.QueryRow("SELECT * FROM users WHERE id = $1", userID)

// UNSAFE ❌
db.QueryRow("SELECT * FROM users WHERE id = " + userID)
```

### XSS Prevention

```
React automatically escapes JSX content
    ↓
Never use dangerouslySetInnerHTML
    ↓
Use Bluemonday for user-generated HTML
    ↓
Content Safe to display
```

### CORS Protection

```
Request from frontend
    ↓
Check Origin header
    ↓
Compare against ALLOWED_ORIGIN
    ↓
If match:
  • Set Access-Control-Allow-Origin
  • Set appropriate headers
  ↓
If no match:
  • Block request (browser CORS policy)
```

### CSRF Protection

```
HttpOnly Cookie:
  ✅ Not accessible from JavaScript
  ✅ Automatically sent with requests
  
SameSite=Strict:
  ✅ Only sent for same-site requests
  ✅ Not sent for cross-site requests
  
Result:
  ✅ CSRF attacks prevented
```

---

## 🔌 WebSocket Architecture

### Connection Management

```
Client connects
    ↓
Backend creates WebSocket connection
    ↓
Validate JWT from cookie
    ↓
Create client struct
    ↓
Add to Hub
    ↓
Send "connected" event
    ↓
Client ready to send/receive
```

### Message Broadcasting

```
Hub structure:
  ├── Active Connections Map
  │   ├── User A: WebSocket conn 1
  │   ├── User B: WebSocket conn 1
  │   └── User C: WebSocket conn 2 (multiple connections)
  │
  └── Broadcast Channel
      ├── Send channel
      ├── Register channel
      └── Unregister channel
```

### Message Delivery

```
User A sends message
    ↓
WebSocket handler validates
    ↓
Save to database
    ↓
Broadcast to all connections of User B
    ↓
User B online?
    ├─ Yes → Deliver immediately via WebSocket
    └─ No  → Stored in database
    ↓
User B next login/connection
    ↓
Fetch from database
    ↓
Display in UI
```

---

## 📱 Frontend Architecture

### Component Hierarchy

```
App
├── AuthProvider
│   ├── ChatProvider
│   │   ├── NotificationProvider
│   │   │   └── Routes
│   │   │       ├── LoginPage
│   │   │       ├── RegisterPage
│   │   │       ├── HomePage
│   │   │       │   ├── FeedSection
│   │   │       │   │   └── PostCard[]
│   │   │       │   └── Sidebar
│   │   │       ├── ChatPage
│   │   │       │   ├── ChatList
│   │   │       │   │   └── ChatListItem[]
│   │   │       │   ├── MessagePanel
│   │   │       │   │   ├── MessageList
│   │   │       │   │   │   └── MessageBubble[]
│   │   │       │   │   └── MessageInput
│   │   │       │   └── UserInfo
│   │   │       ├── ProfilePage
│   │   │       │   ├── ProfileHeader
│   │   │       │   ├── UserStats
│   │   │       │   └── UserPosts[]
│   │   │       ├── AdminPage (if admin)
│   │   │       │   ├── UserManagement
│   │   │       │   ├── PostsManagement
│   │   │       │   └── SystemStats
│   │   │       └── NotFoundPage
```

### State Management

```
AuthContext
├── currentUser: User | null
├── isAuthenticated: boolean
├── loading: boolean
└── Methods: login, logout, register, updateUser

ChatContext
├── conversations: Chat[]
├── selectedChat: Chat | null
├── messages: Message[]
├── isConnected: boolean
└── Methods: connect, disconnect, sendMessage

NotificationContext
├── notifications: Notification[]
└── Methods: add, remove

ThemeContext
├── isDarkMode: boolean
└── Methods: toggle
```

### Custom Hooks

```
useAuth()
├── Access AuthContext
└── Returns: {user, isAuthenticated, login, logout}

useChat()
├── Access ChatContext
├── Manage WebSocket
└── Returns: {messages, send, connect}

useWebSocket(url)
├── Create WebSocket connection
├── Handle reconnection
├── Return connection state

useInfiniteScroll()
├── Detect scroll near bottom
├── Load more data
└── Returns: isLoading, hasMore

useDebounce(value, delay)
├── Delay value updates
└── Used for search optimization
```

---

## 🚀 API Request Flow

### Example: Get Feed

```
Frontend:
  GET /api/posts?page=1&limit=20
  Headers: Cookie: auth_token=<jwt>

Backend:
  1. Middleware: Extract & validate JWT
  2. Handler: postHandler.getFeed()
  3. Service: postService.getFeed(userID, page, limit)
  4. Repository: postRepository.GetFeed()
  5. Query: SELECT * FROM posts ORDER BY created_at DESC LIMIT 20
  6. Format response
  7. Return to client

Response:
  {
    "status": "success",
    "data": {
      "posts": [...],
      "pagination": {...}
    }
  }

Frontend:
  Receive response
  Update state
  Re-render PostList
  Display posts
```

---

## 📊 Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      CLIENT (React)                          │
│                                                              │
│  HomePage                  ChatPage              AdminPage   │
│     │                          │                    │        │
│     └──────────┬───────────────┴────────────────────┘        │
│                │                                              │
│        AuthContext state                                      │
│        ChatContext state                                      │
│        NotificationContext state                              │
└────────┬──────────────────────────────────────────────────────┘
         │
         │ HTTP (REST) + WebSocket
         │
┌────────▼──────────────────────────────────────────────────────┐
│                     SERVER (Go + Fiber)                        │
│                                                              │
│  Routes → Middleware → Handlers → Services → Repositories  │
│                          ↓                      ↓             │
│                    Business Logic         Database Queries   │
└────────┬──────────────────────────────────────────────────────┘
         │
         │ SQL Queries + Storage API
         │
┌────────▴──────────────────────────────────────────────────────┐
│            DATABASE & STORAGE (Supabase)                      │
│                                                              │
│  PostgreSQL        +        S3 Storage                        │
│  (Structured)              (Files)                            │
└───────────────────────────────────────────────────────────────┘
```

---

## ⚡ Performance Optimizations

### Database

- **Connection Pooling:** Min 5, Max 25 connections
- **Indexes:** On frequently queried columns
- **Queries:** Optimized with EXPLAIN ANALYZE
- **Caching:** In-memory cache for static data

### Backend

- **Gzip Compression:** For responses
- **Response Caching:** ETag for feed
- **Request Validation:** Early exit on invalid input
- **Rate Limiting:** Prevent abuse

### Frontend

- **Code Splitting:** Lazy load routes
- **Image Optimization:** Lazy load images
- **Component Memoization:** Prevent re-renders
- **WebSocket:** Real-time instead of polling

### Network

- **CDN:** Vercel global CDN for frontend
- **WebSocket:** Persistent connection for chat
- **API:** Minimal payload sizes
- **Compression:** Gzip for all responses

---

## 🔄 Deployment Architecture

```
GitHub Repositories
├── ghostline-backend
│   └── Push → Railway
│       ├── Docker build
│       ├── Deploy
│       └── Set env variables
│
└── ghostline-frontend
    └── Push → Vercel
        ├── Build (npm run build)
        ├── Deploy
        └── Set env variables
```

---

See also: [API_SPEC.md](./API_SPEC.md), [DEPLOYMENT.md](./DEPLOYMENT.md), [WEBSOCKET_PROTOCOL.md](./WEBSOCKET_PROTOCOL.md)
