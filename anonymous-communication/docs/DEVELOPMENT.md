# 💻 Development Guide

Complete guide for local development of Ghostline Anonymous Messaging Platform.

---

## 🚀 Quick Start (5 minutes)

### Backend Setup

```bash
cd anonymous-communication/backend

# 1. Create .env
cp .env.example .env

# 2. Edit .env
# Set DATABASE_URL to local PostgreSQL
# Set other required variables

# 3. Download dependencies  
go mod download

# 4. Run migrations and start
go run cmd/server/main.go
```

Server runs at `http://localhost:3000`

### Frontend Setup

```bash
cd anonymous-communication/frontend

# 1. Create .env
echo "VITE_API_BASE_URL=http://localhost:3000" > .env

# 2. Install dependencies
npm install

# 3. Start dev server
npm run dev
```

Frontend runs at `http://localhost:5173`

### Test It Works

1. Open `http://localhost:5173` in browser
2. Register new account
3. Create a post with image
4. Open another browser tab
5. Login as different user
6. Send message to first user
7. See message appear in real-time

---

## 📋 Prerequisites

### Backend
- **Go 1.25+** - https://golang.org/dl
- **PostgreSQL 15+** - https://postgres.org/download
  OR use **Supabase** (cloud PostgreSQL)
- **Git**

### Frontend
- **Node.js 18+** - https://nodejs.org
- **npm 9+** (comes with Node.js)
- **Git**

---

## Backend Development

### Environment Setup

#### Option A: Local PostgreSQL

```bash
# Install PostgreSQL (macOS)
brew install postgresql@15

# Start service
brew services start postgresql@15

# Create database
createdb ghostline

# Set connection string
DATABASE_URL=postgresql://localhost:5432/ghostline
```

#### Option B: Supabase (Recommended)

```bash
# Create free account at https://supabase.com
# Create project "ghostline"
# Get connection string from Settings → Database

DATABASE_URL=postgresql://postgres:[PASSWORD]@[HOST]:[PORT]/postgres
```

### File Structure

```
backend/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   ├── repositories/            # Database access
│   ├── models/                  # Data structures
│   ├── middleware/              # Request middleware
│   ├── database/                # DB setup & migrations
│   ├── config/                  # Configuration
│   ├── utils/                   # Utilities
│   ├── routes/                  # Route definitions
│   └── websocket/               # WebSocket logic
├── pkg/                         # Public packages
├── tests/                       # Test files
├── go.mod & go.sum             # Dependencies
├── Makefile                     # Build commands
├── Dockerfile                   # Docker image
└── .env.example                # Example env vars
```

### Configuration

Create `.env` in `backend/`:

```env
# Server
PORT=3000
ENVIRONMENT=development
ALLOWED_ORIGIN=http://localhost:5173

# Database
DATABASE_URL=postgresql://localhost:5432/ghostline
DB_MAX_CONNECTIONS=10
DB_MIN_CONNECTIONS=2

# JWT
JWT_SECRET=dev-secret-key-min-32-chars-needed
JWT_EXPIRATION_MINUTES=15
AUTH_COOKIE_NAME=auth_token
COOKIE_SECURE=false

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_KEY=your-key
STORAGE_BUCKET_NAME=user-uploads

# Logging
LOG_LEVEL=debug
```

### Running Server

#### Development (Hot Reload)

```bash
# Install air (Go hot reload)
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air
```

#### Production Build

```bash
# Build binary
go build -o server ./cmd/server

# Run
./server
```

#### Docker

```bash
# Build image
docker build -t ghostline-backend .

#Run container
docker run -p 3000:3000 \
  -e DATABASE_URL=postgresql://... \
  -e JWT_SECRET=... \
  ghostline-backend
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/services

# With verbose output
go test -v ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Logging

Logs appear in console:

```
2026-03-28T10:00:00 INFO server started listening on :3000
2026-03-28T10:00:05 DEBUG database connected
2026-03-28T10:00:10 INFO user created id=550e8400-e29b-41d4-a716-446655440000
```

Control log level in `.env`:
- `debug` - Everything (most verbose)
- `info` - Important events
- `warn` - Warnings
- `error` - Errors only

### Debugging

#### Using Printf

```go
fmt.Printf("user: %+v\n", user)
slog.Debug("created user", "id", user.ID)
```

#### Using Delve Debugger

```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug ./cmd/server

# In dlv prompt
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print variable
```

#### VS Code Debugging

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Backend",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${fileDirname}",
      "env": {},
      "args": []
    }
  ]
}
```

---

## Frontend Development

### Environment Setup

```bash
cd frontend

# Create .env
echo "VITE_API_BASE_URL=http://localhost:3000" > .env
echo "VITE_WS_BASE_URL=ws://localhost:3000" >> .env
```

### File Structure

```
frontend/
├── src/
│   ├── components/              # React components
│   ├── pages/                   # Page components
│   ├── context/                 # Context providers
│   ├── services/                # API clients
│   ├── hooks/                   # Custom hooks
│   ├── types/                   # TypeScript types
│   ├── utils/                   # Utilities
│   ├── styles/                  # CSS files
│   ├── App.tsx                  # Root component
│   └── main.tsx                 # Entry point
├── public/                      # Static assets
├── index.html                   # HTML template
├── package.json                 # Dependencies
├── tsconfig.json                # TypeScript config
├── vite.config.ts               # Vite config
├── tailwind.config.js           # Tailwind config
├── Dockerfile                   # Docker image
└── .env.example                # Example env vars
```

### Running Dev Server

```bash
# Start with hot reload
npm run dev

# Output:
# VITE v5.0.0 ready in 245 ms
# 
# ➜  Local:   http://localhost:5173/
# ➜  Network: use --host to access from network
```

Open http://localhost:5173

### Commands

```bash
# Development
npm run dev                 # Start dev server

# Build
npm run build              # Production build
npm run preview            # Preview build locally

# Type checking
npm run type-check         # Check TypeScript errors

# Linting & Formatting
npm run lint               # ESLint
npm run format             # Format code

# Testing
npm run test               # Run tests
npm run test:coverage      # Coverage report
```

### Component Development

#### Create New Component

```typescript
// src/components/MyComponent.tsx
import React from 'react';

interface MyComponentProps {
  title: string;
  onClose?: () => void;
}

export const MyComponent: React.FC<MyComponentProps> = ({ title, onClose }) => {
  return (
    <div className="p-4 bg-white rounded">
      <h2>{title}</h2>
      {onClose && <button onClick={onClose}>Close</button>}
    </div>
  );
};
```

#### Use Component

```typescript
import { MyComponent } from '../components/MyComponent';

export function MyPage() {
  const [isOpen, setIsOpen] = React.useState(true);
  
  return (
    <div>
      {isOpen && (
        <MyComponent 
          title="Hello" 
          onClose={() => setIsOpen(false)}
        />
      )}
    </div>
  );
}
```

### Debugging

#### Console Logging

```typescript
console.log('User:', user);
console.table(posts);
console.error('Error:', error);
```

#### React DevTools

1. Install Chrome extension: "React Developer Tools"
2. Open DevTools: F12
3. Go to "Components" tab
4. Inspect components, see props/state

#### Network Tab

1. Open DevTools: F12
2. Go to "Network" tab
3. Make API calls
4. See requests and responses
5. Check response status (200 ok, 401 auth error, etc.)

#### Breakpoints

```typescript
// Set breakpoint
debugger;

// Will pause execution
// Use DevTools to step through code
```

### Styling

#### Add Global Styles

Edit `src/index.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}
```

#### Use Tailwind Classes

```typescript
<div className="flex items-center justify-center min-h-screen bg-gray-100">
  <button className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
    Click me
  </button>
</div>
```

#### Component-specific Styles

```typescript
export function MyComponent() {
  const styles = {
    container: 'flex flex-col gap-4 p-6',
    title: 'text-2xl font-bold',
    button: 'px-4 py-2 bg-blue-500 text-white rounded'
  };
  
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Hello</h1>
      <button className={styles.button}>Click</button>
    </div>
  );
}
```

---

## Making API Calls

### HTTP Requests

```typescript
// In services/authService.ts
import api from './api';

export const authService = {
  login: async (username: string, password: string) => {
    const response = await api.post('/api/auth/login', {
      username,
      password
    });
    return response.data;
  },

  register: async (username: string, email: string, password: string) => {
    const response = await api.post('/api/auth/register', {
      username,
      email,
      password
    });
    return response.data;
  }
};
```

### WebSocket Connection

```typescript
// In services/websocketService.ts
let socket: WebSocket | null = null;

export const websocketService = {
  connect: (onMessage: (data: any) => void) => {
    const wsURL = import.meta.env.VITE_WS_BASE_URL;
    socket = new WebSocket(`${wsURL}/ws/chat`);

    socket.onopen = () => {
      console.log('Connected');
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      onMessage(data);
    };

    socket.onerror = (error) => {
      console.error('Error:', error);
    };
  },

  send: (message: any) => {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(message));
    }
  }
};
```

---

## Docker Development

### Build Image

```bash
cd backend
docker build -t ghostline-backend .

cd ../frontend
docker build -t ghostline-frontend .
```

### Run with Docker Compose

```bash
# In project root
docker-compose -f docker-compose.yml up

# Services start:
# - Frontend on localhost:3000
# - Backend on localhost:3001
# - Database on localhost:5432
```

### View Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs backend
docker-compose logs frontend

# Follow logs
docker-compose logs -f backend
```

---

## Troubleshooting

### "address already in use"

```bash
# Find process using port 3000
lsof -i :3000

# Kill process
kill -9 <PID>

# Or use different port
PORT=3001 go run cmd/server/main.go
```

### "module not found"

```bash
# Go
go mod download
go mod tidy

# Node
rm -rf node_modules package-lock.json
npm install
```

### "cannot connect to database"

```bash
# Check database is running
psql -U postgres -d ghostline -c "SELECT 1"

# Check connection string in .env
# Check PostgreSQL is accepting connections
```

### "WebSocket connection failed"

```bash
# Check backend is running
curl http://localhost:3000/health

# Check WebSocket URL in .env
# Make sure browser allows WebSocket
```

---

See also: [ARCHITECTURE.md](./ARCHITECTURE.md), [DEPLOYMENT.md](./DEPLOYMENT.md), [API_SPEC.md](./API_SPEC.md)
