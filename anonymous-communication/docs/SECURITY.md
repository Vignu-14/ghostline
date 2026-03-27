# 🔐 Security Implementation Guide

Complete security analysis and implementation details for Ghostline.

---

## 🎯 Security Objectives

- **Confidentiality:** Data encrypted in transit and at rest
- **Integrity:** Tamper detection and prevention
- **Availability:** Protection against abuse and attacks
- **Authentication:** Verify user identity
- **Authorization:** Control what users can access

---

## Authentication & JWT

### JWT Implementation

**Token Format:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJpYXQiOjE2ODMwMDAwMDAsImV4cCI6MTY4MzAwOTAwMH0.signature
```

**Payload:**
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",  // user ID
  "iat": 1683000000,                              // issued at
  "exp": 1683009000                               // expires in 15 mins
}
```

### JWT Generation

```go
// Backend: internal/services/authService.go
import (
    "github.com/golang-jwt/jwt/v5"
    "os"
    "time"
)

func GenerateJWT(userID string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    expirationMinutes := 15
    
    claims := jwt.MapClaims{
        "sub": userID,
        "iat": time.Now().Unix(),
        "exp": time.Now().Add(time.Minute * time.Duration(expirationMinutes)).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### JWT Validation

```go
// Middleware: internal/middleware/jwt_middleware.go
func ValidateJWT(tokenString string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    
    token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err != nil {
        return "", err
    }
    
    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok || !token.Valid {
        return "", errors.New("invalid token")
    }
    
    return (*claims)["sub"].(string), nil
}
```

### Token Expiration

JWT tokens expire after **15 minutes**:
- Forces token refresh
- Limits damage if token is stolen
- User stays logged in via cookie refresh

**Token Refresh Flow:**
```
1. User logs in → JWT issued
2. After 10 mins → Token still valid
3. After 15 mins → Token expired
4. Frontend detects expiration
5. Frontend calls /auth/me endpoint
6. Backend issues new token
7. User continues session
```

---

## Password Security

### Password Requirements

```
Minimum 8 characters
At least 1 uppercase letter (A-Z)
At least 1 lowercase letter (a-z)
At least 1 digit (0-9)
At least 1 special character (!@#$%^&*)

Example: MyP@ssw0rd123
```

### Password Storage

Passwords are **never stored in plaintext**. Instead, Bcrypt hash is stored:

```go
// Backend: internal/services/authService.go
import "golang.org/x/crypto/bcrypt"

const BcryptCost = 12

// Hash password
passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)

// Verify password
err := bcrypt.CompareHashAndPassword(
    []byte(storedHash), 
    []byte(providedPassword),
)
```

### Why Bcrypt?

- **Slow:** Takes ~100ms per hash (makes brute force slow)
- **Salted:** Each hash different, even for same password
- **Adaptive:** Can increase cost factor as computers get faster
- **Industry standard:** Used by major platforms

### Password Reset

**Reset Flow:**
```
1. User clicks "Forgot Password"
2. System sends email with link
3. Link contains time-limited token (30 mins)
4. User clicks link, enters new password
5. New password hashed and stored
6. Old sessions invalidated
```

---

## Input Validation & Sanitization

### Frontend Validation

```typescript
// Validate inputs before sending
function validateUsername(username: string): boolean {
  // Length validation
  if (username.length < 3 || username.length > 50) {
    return false;
  }
  
  // Pattern validation (alphanumeric + underscore/hyphen)
  if (!/^[a-zA-Z0-9_-]+$/.test(username)) {
    return false;
  }
  
  return true;
}

function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

function validatePassword(password: string): boolean {
  const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
  return regex.test(password);
}
```

### Backend Validation

```go
// Backend validation happens regardless of frontend
type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// Validate struct
validator := validator.New()
if err := validator.Struct(req); err != nil {
    return c.Status(400).JSON(fiber.Map{"error": "validation failed"})
}
```

### HTML Content Sanitization (XSS Prevention)

Post captions and messages are sanitized to prevent XSS:

```go
// Backend: internal/utils/sanitizer.go
import "github.com/microcosm-cc/bluemonday"

func SanitizeHTML(content string) string {
    // Create policy that strips all HTML tags
    p := bluemonday.StrictPolicy()
    return p.Sanitize(content)
}

// Usage in post creation
caption := SanitizeHTML(req.Caption)
```

**What happens:**
```
Input:  <script>alert('xss')</script>
Output: (empty - script removed)

Input:  <b>Bold text</b>
Output: Bold text (tags removed)

Input:  Normal text
Output: Normal text (unchanged)
```

---

## Database Security

### SQL Injection Prevention

**Vulnerable (WRONG):**
```go
// NEVER DO THIS
query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
```

**Safe (CORRECT):**
```go
// Use parameterized queries
query := "SELECT * FROM users WHERE username = $1"
user := db.QueryRowContext(ctx, query, username)
```

All queries in backend use parameterized queries with placeholders (`$1`, `$2`, etc).

### Connection Pooling

Database connections are pooled to prevent exhaustion:

```go
// config/database.go
sqldb.SetMaxOpenConns(25)      // Max 25 concurrent connections
sqldb.SetMaxIdleConns(5)       // Keep 5 idle connections
sqldb.SetConnMaxLifetime(...)  // Recycle old connections
```

---

## Authentication Cookies

### Cookie Configuration

```javascript
// Set in backend response
Set-Cookie: auth_token=<JWT>; 
  Path=/;
  HttpOnly;
  Secure;
  SameSite=Strict;
  Max-Age=900
```

### Cookie Security

| Setting | Value | Why |
|---------|-------|-----|
| HttpOnly | true | Prevents JavaScript access (XSS protection) |
| Secure | true | Only sent over HTTPS |
| SameSite | Strict | Prevents CSRF attacks |
| Path | / | Available on all routes |
| Max-Age | 900 | Expires in 15 minutes |

### Cookie Flow

```
1. User logs in
2. Backend creates JWT
3. Backend sets cookie
4. Browser stores cookie
5. Browser auto-includes in requests
6. Backend validates JWT from cookie
7. User session valid
```

---

## CORS (Cross-Origin Resource Sharing)

### CORS Configuration

```go
// Backend: internal/middleware/cors_middleware.go
corsConfig := cors.Config{
    AllowOrigins: os.Getenv("ALLOWED_ORIGIN"), // "https://example.com"
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge: 3600,
}
app.Use(cors.New(corsConfig))
```

### What It Prevents

**Prevents:**
- Frontend on `example.com` accessing API on `evil.com`
- Cross-site attacks stealing data

**Allows:**
- Legitimate frontend accessing legitimate API
- Specified origins only

### ALLOWED_ORIGIN Values

```
Development:  http://localhost:5173
Production:   https://ghostline-frontend-xxxxx.vercel.app
Custom domain: https://ghostline.example.com
```

---

## CSRF (Cross-Site Request Forgery) Prevention

### Double-Submit Cookie Pattern

Modern browsers with `SameSite=Strict` cookies automatically prevent CSRF.

**How it works:**
```
1. Attacker tries to submit form from attacker.com
2. Browser doesn't include auth_token cookie (different origin)
3. Request lacks authentication
4. Server rejects request
```

### Backend Protection

```go
// Middleware verifies origin matches
csrf.New(csrf.Config{
    KeyLookup: "header:X-CSRF-Token",
    CookieName: "XSRF-TOKEN",
    Cookie: &http.Cookie{
        HttpOnly: true,
        Secure: true,
        SameSite: http.SameSiteLaxMode,
    },
})
```

---

## Rate Limiting

### Implementation

```go
// Backend: internal/middleware/rate_limiter.go
import "github.com/gofiber/fiber/v2/middleware/limiter"

limiterConfig := limiter.Config{
    Max: 100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP() // Rate limit by IP
    },
    Storage: store, // Redis or memory store
}

app.Use(limiter.New(limiterConfig))
```

### Rate Limits per Endpoint

| Endpoint | Limit | Duration |
|----------|-------|----------|
| /api/auth/login | 5 requests | 15 minutes |
| /api/auth/register | 3 requests | 1 hour |
| /api/posts | 10 requests | 1 minute |
| /api/messages | 20 requests | 1 minute |
| General | 100 requests | 1 minute |

### Prevents

- Brute force password attacks
- DoS attacks
- Account enumeration
- Spam/abuse

---

## Data Encryption

### In Transit (TLS)

**All** communication is encrypted:

```
HTTPS:  https://ghostline.com ✅ Encrypted
HTTP:   http://ghostline.com ❌ Not allowed
WSS:    wss://ghostline.com ✅ Encrypted
WS:     ws://ghostline.com ❌ Not allowed
```

**Certificate**: Auto-renewed by:
- Vercel (frontend)
- Railway (backend)
- Supabase (database)

### At Rest

Sensitive data in database:
- **Passwords:** Bcrypt hashed (salted)
- **JWTs:** Not stored (stateless)
- **Messages:** Stored in encrypted database (Supabase)
- **Images:** Stored in encrypted storage (Supabase)

---

## Row-Level Security (RLS)

Supabase enforces RLS on database:

```sql
-- Only users can see their own messages
CREATE POLICY "users_can_send_messages"
ON messages
FOR INSERT
WITH CHECK (auth.uid() = sender_id);

CREATE POLICY "users_can_view_messages"
ON messages
FOR SELECT
USING (auth.uid() = sender_id OR auth.uid() = receiver_id);
```

**Prevents:**
- Users accessing other users' private messages
- Direct database queries bypassing API

---

## API Security Headers

All responses include security headers:

```
X-Frame-Options: DENY                  # No iframes
X-Content-Type-Options: nosniff         # Prevent MIME sniffing
Content-Security-Policy: default-src 'self'  # Only trusted sources
X-XSS-Protection: 1; mode=block        # XSS protection
Strict-Transport-Security: max-age=...  # Force HTTPS
```

---

## Audit Logging

Login/logout events logged:

```go
// Track all authentication events
type AuthLog struct {
    ID           UUID
    UserID       UUID
    Action       string    // 'login', 'logout', 'register', 'failed_login'
    Status       string    // 'success', 'failed'
    IPAddress    string
    UserAgent    string
    ErrorMessage string
    CreatedAt    time.Time
}
```

**Queries:**
```sql
-- Find suspicious activity
SELECT * FROM auth_logs 
WHERE action = 'failed_login' 
AND created_at > NOW() - INTERVAL '1 hour'
GROUP BY user_id, ip_address
HAVING COUNT(*) > 5;
```

---

## Security Best Practices

### For Developers

- ✅ Never log passwords
- ✅ Always validate input
- ✅ Use parameterized queries
- ✅ Enable HTTPS everywhere
- ✅ Rotate secrets regularly
- ✅ Use environment variables
- ✅ Keep dependencies updated
- ✅ Run security audits
- ❌ Don't hardcode secrets
- ❌ Don't trust user input
- ❌ Don't use MD5/SHA1 for passwords
- ❌ Don't disable CORS validation

### For Administrators

- ✅ Monitor auth logs
- ✅ Check rate limiting
- ✅ Review database backups
- ✅ Update dependencies
- ✅ Rotate JWT_SECRET annually
- ✅ Monitor for abuse
- ✅ Check SSL certificates
- ❌ Don't disable authentication
- ❌ Don't expose database credentials
- ❌ Don't run unpatched versions

---

## Vulnerability Scanning

### Dependencies

```bash
# Check for vulnerable packages
go list -json -m all | nancy sleuth  # Go
npm audit                            # Node.js
```

### Code Analysis

```bash
# Static analysis
go vet ./...                    # Go linting
npm run lint                    # JavaScript
```

### Regular Security Audits

Perform quarterly security reviews:
- [ ] Check authentication implementation
- [ ] Verify rate limiting
- [ ] Review audit logs
- [ ] Test SQL injection
- [ ] Test XSS prevention
- [ ] Verify CORS configuration

---

## Incident Response

### If Password Compromised

```
1. Change password immediately
2. Check login history in settings
3. Logout all sessions
4. Enable 2FA (when implemented)
5. Contact admin if continued issues
```

### If Private Messages Exposed

```
1. Contact admin immediately
2. Check message deletion
3. Report account compromise
4. Change password
5. Enable 2FA
```

### If Database Breached

```
1. Immediate notification to all users
2. Force password reset for all accounts
3. Invalidate all JWT tokens
4. Rotate JWT_SECRET
5. Rotate database password
6. Review backups
7. Investigation and public disclosure
```

---

## Compliance & Standards

### OWASP Top 10

Ghostline protects against:
- ✅ A1: SQL Injection (parameterized queries)
- ✅ A2: Broken Authentication (JWT, secure cookies)
- ✅ A3: Broken Access Control (RLS, role-based)
- ✅ A4: Sensitive Data Exposure (HTTPS, hashing)
- ✅ A5: XML/XXE (no XML used)
- ✅ A6: Broken Access Control (RLS)
- ✅ A7: XSS (HTML sanitization)
- ✅ A8: CSRF (SameSite cookies)
- ✅ A9: Using Components with Known Vulnerabilities (dependency scans)
- ✅ A10: Insufficient Logging (audit logs)

---

See also: [API_SPEC.md](./API_SPEC.md), [ARCHITECTURE.md](./ARCHITECTURE.md), [DEPLOYMENT.md](./DEPLOYMENT.md)
