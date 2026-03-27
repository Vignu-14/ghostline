# 📡 WebSocket Protocol Specification

Complete WebSocket implementation for real-time messaging.

---

## 🎯 Overview

WebSocket enables real-time bidirectional communication for instant messaging and live updates.

**Endpoint:** `wss://api.example.com/ws/chat` (WSS = WebSocket Secure)

**Protocol:** Custom JSON-based message protocol over WebSocket

---

## Connection

### Establishing Connection

```javascript
// Frontend
const token = localStorage.getItem('auth_token');
const wsURL = `${import.meta.env.VITE_WS_BASE_URL}/ws/chat`;
const socket = new WebSocket(wsURL);

socket.addEventListener('open', () => {
    console.log('Connected');
    // Send authentication
    socket.send(JSON.stringify({
        type: 'authenticate',
        token: token
    }));
});
```

### Backend Connection Handler

```go
// Backend: internal/websocket/handler.go
func HandleWebSocket(c *fiber.Ctx) error {
    return websocket.New(func(ws *websocket.Conn) {
        // Connection established
        clientID := uuid.New().String()
        
        for {
            var message map[string]interface{}
            err := ws.ReadJSON(&message)
            if err != nil {
                // Connection closed
                break
            }
            
            // Process message
            handleMessage(clientID, ws, message)
        }
    })(c)
}
```

### Connection States

```
1. CONNECTING        → Client: attempting to connect
2. OPEN             → Connected: ready to send/receive
3. CLOSING          → Client: closing connection
4. CLOSED           → Connection terminated
```

---

## Authentication

### Handshake

```javascript
// Step 1: Connect
const socket = new WebSocket('wss://api.example.com/ws/chat');

// Step 2: On open, send authentication
socket.addEventListener('open', () => {
    socket.send(JSON.stringify({
        type: 'authenticate',
        token: jwtToken
    }));
});

// Step 3: Wait for auth response
socket.addEventListener('message', (event) => {
    const msg = JSON.parse(event.data);
    if (msg.type === 'auth_success') {
        // Authenticated, can now send messages
    } else if (msg.type === 'auth_error') {
        // Authentication failed
        socket.close();
    }
});
```

### Authenticate Message

```json
{
  "type": "authenticate",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Auth Response (Success)

```json
{
  "type": "auth_success",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "timestamp": "2026-03-28T10:00:00Z"
}
```

### Auth Response (Failure)

```json
{
  "type": "auth_error",
  "error": "invalid_token",
  "message": "Token expired"
}
```

---

## Message Types

### 1. Send Message

**Client → Server:**

```json
{
  "type": "message",
  "receiver_id": "550e8400-e29b-41d4-a716-446655440001",
  "content": "Hello, how are you?",
  "timestamp": "2026-03-28T10:00:30Z"
}
```

**Server → Receiver:**

```json
{
  "type": "message",
  "sender_id": "550e8400-e29b-41d4-a716-446655440000",
  "sender_username": "john_doe",
  "content": "Hello, how are you?",
  "message_id": "a7c2f8b1-9e3d-4a6f-b2c8-1e7d9a3f5c2e",
  "timestamp": "2026-03-28T10:00:30Z"
}
```

### 2. Message Acknowledgment

**Server → Sender:**

```json
{
  "type": "message_ack",
  "message_id": "a7c2f8b1-9e3d-4a6f-b2c8-1e7d9a3f5c2e",
  "status": "delivered",
  "timestamp": "2026-03-28T10:00:30Z"
}
```

Statuses:
- `pending` - Not yet delivered
- `delivered` - Received by server
- `seen` - Recipient has seen
- `failed` - Send failed

### 3. Typing Indicator

**Client → Server:**

```json
{
  "type": "typing",
  "receiver_id": "550e8400-e29b-41d4-a716-446655440001",
  "is_typing": true
}
```

**Server → Receiver:**

```json
{
  "type": "typing",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "is_typing": true
}
```

### 4. User Online Status

**Server → Connected Clients:**

```json
{
  "type": "user_online",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "status": "online"
}
```

Statuses:
- `online` - User connected
- `away` - Idle for 5+ minutes
- `offline` - Disconnected

### 5. Mark as Read

**Client → Server:**

```json
{
  "type": "mark_read",
  "conversation_id": "550e8400-e29b-41d4-a716-446655440001",
  "last_message_id": "a7c2f8b1-9e3d-4a6f-b2c8-1e7d9a3f5c2e"
}
```

**Server → Sender:**

```json
{
  "type": "message_read",
  "reader_id": "550e8400-e29b-41d4-a716-446655440001",
  "reader_username": "jane_doe",
  "last_message_id": "a7c2f8b1-9e3d-4a6f-b2c8-1e7d9a3f5c2e",
  "timestamp": "2026-03-28T10:05:00Z"
}
```

### 6. Ping (Keep-Alive)

**Client → Server (every 30 seconds):**

```json
{
  "type": "ping"
}
```

**Server → Client:**

```json
{
  "type": "pong",
  "timestamp": "2026-03-28T10:00:30Z"
}
```

Prevents connection timeout.

### 7. Error Message

**Server → Client:**

```json
{
  "type": "error",
  "code": "INVALID_MESSAGE",
  "message": "Message content cannot be empty",
  "timestamp": "2026-03-28T10:00:30Z"
}
```

Error Codes:
- `INVALID_MESSAGE` - Message format invalid
- `MESSAGE_TOO_LONG` - Content exceeds limit
- `RECIPIENT_NOT_FOUND` - User not found
- `PERMISSION_DENIED` - Not authorized
- `RATE_LIMITED` - Too many messages
- `SERVER_ERROR` - Internal error

---

## Message Flow Diagram

```
Client (Sender)                Server                  Client (Receiver)
    |                            |                            |
    |-- authenticate -->         |                            |
    |                    <-- auth_success --                  |
    |                            |                  <-- auth_success --
    |                            |                            |
    |-- message -->              |-- message -->              |
    |                            |                            |
    |<-- message_ack --          |-- typing_indicator         |
    |                            |                            |
    |                            |                            |
    |                            |-- message -->              |
    |                            |                            |
    |                            |                  <-- message_ack --
    |                            |                            |
    |<-- message_read --         |<-- message_read --         |
    |                            |                            |
```

---

## Connection Management

### Automatic Reconnection

```javascript
class WebSocketClient {
    constructor(url, maxRetries = 5) {
        this.url = url;
        this.maxRetries = maxRetries;
        this.retryCount = 0;
        this.connect();
    }

    connect() {
        this.socket = new WebSocket(this.url);

        this.socket.addEventListener('open', () => {
            this.retryCount = 0;
            this.authenticate();
        });

        this.socket.addEventListener('close', () => {
            this.attemptReconnect();
        });

        this.socket.addEventListener('error', (error) => {
            console.error('WebSocket error:', error);
        });

        this.socket.addEventListener('message', (event) => {
            this.handleMessage(JSON.parse(event.data));
        });
    }

    attemptReconnect() {
        if (this.retryCount < this.maxRetries) {
            const delay = Math.pow(2, this.retryCount) * 1000; // Exponential backoff
            setTimeout(() => {
                this.retryCount++;
                this.connect();
            }, delay);
        }
    }

    send(message) {
        if (this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message));
        }
    }
}
```

### Connection Timeout

```go
// Backend: Close idle connections after 2 minutes
const readDeadline = 2 * time.Minute
ws.SetReadDeadline(time.Now().Add(readDeadline))

// Client must send ping every 30 seconds
for {
    _, message, err := ws.ReadMessage()
    if err != nil {
        // Connection closed
        break
    }
    // Reset deadline on every message
    ws.SetReadDeadline(time.Now().Add(readDeadline))
}
```

---

## Rate Limiting

### Per-Connection Limits

| Action | Limit | Duration |
|--------|-------|----------|
| Messages | 10 | 1 minute |
| Typing | 1 | per second |
| Disconnects | 5 | 1 minute |

### Enforcement

```go
// Backend: Rate limiting per connection
type ConnectionLimiter struct {
    messageCount int
    lastReset    time.Time
    maxMessages  int
    resetPeriod  time.Duration
}

func (cl *ConnectionLimiter) CanSendMessage() bool {
    now := time.Now()
    
    if now.Sub(cl.lastReset) > cl.resetPeriod {
        cl.messageCount = 0
        cl.lastReset = now
    }
    
    if cl.messageCount >= cl.maxMessages {
        return false
    }
    
    cl.messageCount++
    return true
}
```

### Rate Limited Response

```json
{
  "type": "error",
  "code": "RATE_LIMITED",
  "message": "Too many messages. Try again in 60 seconds",
  "retry_after": 60
}
```

---

## Offline Message Delivery

### When User Offline

If sender sends message to offline recipient:

1. Message stored in database
2. Sender gets `delivered` ack
3. When recipient connects, receives stored messages
4. Recipient gets messages in chronological order

### Database Storage

```sql
SELECT * FROM messages
WHERE 
  receiver_id = $1 
  AND deleted_by_receiver = FALSE
  AND created_at > (SELECT last_seen FROM users WHERE id = $1)
ORDER BY created_at DESC
LIMIT 50
```

---

## Broadcast Events

### New Post Created

```json
{
  "type": "post_created",
  "post_id": "a7c2f8b1-9e3d-4a6f-b2c8-1e7d9a3f5c2e",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "caption": "Beautiful sunset!",
  "image_url": "https://...",
  "likes_count": 0,
  "created_at": "2026-03-28T10:00:00Z"
}
```

### Post Liked

```json
{
  "type": "post_liked",
  "post_id": "a7c2f8b1-9e3d-4a6f-b2c8-1e7d9a3f5c2e",
  "liker_id": "550e8400-e29b-41d4-a716-446655440001",
  "liker_username": "jane_doe",
  "likes_count": 42,
  "timestamp": "2026-03-28T10:05:00Z"
}
```

---

## Security Considerations

### Token Expiration

JWT tokens expire every 15 minutes. When expired:

```javascript
socket.addEventListener('message', (event) => {
    const msg = JSON.parse(event.data);
    if (msg.type === 'auth_error' && msg.error === 'token_expired') {
        // Refresh token from HTTP API
        const newToken = await refreshToken();
        // Reconnect WebSocket
        socket.close();
        connectWebSocket(newToken);
    }
});
```

### Input Validation

All messages validated server-side:

```go
// Backend validation
if len(message.Content) == 0 || len(message.Content) > 5000 {
    return fmt.Errorf("invalid message length")
}

// Sanitize content
message.Content = sanitizeHTML(message.Content)
```

### Rate Limiting

Prevents abuse and DoS attacks (see Rate Limiting section above).

---

## Testing WebSocket

### Using WebSocket Browser API

```javascript
const socket = new WebSocket('wss://api.example.com/ws/chat');

socket.onopen = () => {
    // Send test message
    socket.send(JSON.stringify({
        type: 'authenticate',
        token: 'YOUR_JWT_TOKEN'
    }));
};

socket.onmessage = (event) => {
    console.log('Received:', event.data);
};

socket.onerror = (error) => {
    console.error('Error:', error);
};

socket.onclose = () => {
    console.log('Disconnected');
};
```

### Using websocat CLI

```bash
# Install websocat
cargo install websocat

# Connect to WebSocket
websocat 'wss://api.example.com/ws/chat'

# In the prompt, send JSON messages:
{"type":"authenticate","token":"YOUR_JWT_TOKEN"}
{"type":"message","receiver_id":"...","content":"Hello"}
```

### Using curl for Testing

```bash
# Note: curl 7.70.0+ supports WebSocket
curl --include \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
     --header "Sec-WebSocket-Version: 13" \
     http://localhost:3000/ws/chat
```

---

## Performance Optimization

### Message Batching

```javascript
// Don't send every keystroke
class MessageBatcher {
    constructor(flushInterval = 100) {
        this.batch = [];
        this.flushInterval = flushInterval;
        this.timer = null;
    }

    add(message) {
        this.batch.push(message);
        
        if (!this.timer) {
            this.timer = setTimeout(() => this.flush(), this.flushInterval);
        }
    }

    flush() {
        if (this.batch.length > 0) {
            socket.send(JSON.stringify({
                type: 'batch',
                messages: this.batch
            }));
            this.batch = [];
        }
        this.timer = null;
    }
}
```

### Connection Pooling

Backend maintains connection pool:

```go
type Hub struct {
    clients   map[string]*Client
    broadcast chan *Message
    register  chan *Client
    unregister chan *Client
    mutex     sync.RWMutex
}

// Limit max connections per user
const maxConnectionsPerUser = 3
```

---

## Debugging

### Browser DevTools

```javascript
// Intercept all WebSocket messages
const originalSend = WebSocket.prototype.send;
WebSocket.prototype.send = function(data) {
    console.log('==> SENT:', JSON.parse(data));
    originalSend.call(this, data);
};

const originalOnMessage = 'onmessage';
Object.defineProperty(WebSocket.prototype, 'onmessage', {
    set(handler) {
        return this.addEventListener('message', (event) => {
            console.log('<== RECEIVED:', JSON.parse(event.data));
            handler.call(this, event);
        });
    }
});
```

### Server-Side Logging

```go
// Log all WebSocket messages
func handleMessage(ws *websocket.Conn, msg map[string]interface{}) {
    slog.Debug("websocket message", 
        "type", msg["type"], 
        "user_id", userID,
        "timestamp", time.Now())
    
    // Process message
}
```

---

See also: [API_SPEC.md](./API_SPEC.md), [ARCHITECTURE.md](./ARCHITECTURE.md), [SECURITY.md](./SECURITY.md)
