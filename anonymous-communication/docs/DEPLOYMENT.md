# 🚀 Deployment Guide

Complete guide for deploying Ghostline to production.

---

## 📋 Overview

The Ghostline application consists of:
- **Backend:** Go + Fiber (deployed on Railway)
- **Frontend:** React + Vite (deployed on Vercel)  
- **Database:** PostgreSQL (Supabase)
- **Storage:** S3-compatible (Supabase Storage)

---

## Prerequisites

### GitHub

- GitHub account
- Repository access
- SSH keys configured

### Railway (Backend)

- Railway account (https://railway.app)
- Connected GitHub account
- Railway CLI installed (optional)

### Vercel (Frontend)

- Vercel account (https://vercel.com)
- Connected GitHub account
- Vercel CLI installed (optional)

### Supabase (Database & Storage)

- Supabase account (https://supabase.com)
- Project created
- Service key generated

---

## Step 1: Prepare GitHub Repositories

### Create Two Repositories

1. **Backend Repository**
   - Name: `ghostline-backend`
   - URL: https://github.com/YOUR_USERNAME/ghostline-backend
   - Make it public

2. **Frontend Repository**
   - Name: `ghostline-frontend`
   - URL: https://github.com/YOUR_USERNAME/ghostline-frontend
   - Make it public

### Push Code

```bash
# Backend
cd ghostline-backend
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/ghostline-backend.git
git push -u origin main

# Frontend
cd ghostline-frontend
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/ghostline-frontend.git
git push -u origin main
```

---

## Step 2: Set Up Supabase Database

### Create Supabase Project

1. Go to https://supabase.com
2. Click "New Project"
3. Enter project details:
   - Organization: Create one
   - Project Name: `ghostline`
   - Region: Select closest region
   - Database Password: Strong password (save it!)
4. Click "Create new project"

### Wait for Project Creation

Takes 1-2 minutes. You'll receive email when ready.

### Get Connection Details

```
Project Settings → Database → URI

Format: postgresql://postgres:[PASSWORD]@[HOST]:[PORT]/postgres

Save these values:
- DATABASE_URL (full URI)
- SUPABASE_URL (see Settings → API)
- SUPABASE_SERVICE_KEY (see Settings → API)
```

### Create Storage Bucket

1. Go to Supabase → Storage → Buckets
2. Create new bucket: `user-uploads`
3. Set to Public (allow public reads)
4. Create additional buckets if needed:
   - `posts` - for post content
   - `profiles` - for profile pictures

---

## Step 3: Deploy Backend to Railway

### Create Railway Project

1. Go to https://railway.app
2. Click "New Project"
3. Select "GitHub Repo"
4. Authorize GitHub if needed
5. Select `ghostline-backend` repository
6. Railway auto-detects Go project
7. Click "Deploy"

### Configure Environment Variables

Railway Dashboard → Variables:

```
PORT=3000
ENVIRONMENT=production
ALLOWED_ORIGIN=https://ghostline-frontend-five.vercel.app

DATABASE_URL=postgresql://...@... (from Supabase)
DB_MAX_CONNECTIONS=25
DB_MIN_CONNECTIONS=5
DB_MAX_CONN_LIFETIME_MINUTES=60
DB_MAX_CONN_IDLE_MINUTES=15
DB_HEALTH_CHECK_SECONDS=30
DB_CONNECT_TIMEOUT_SECONDS=5

JWT_SECRET=<generate-with-openssl-rand-base64-32> (min 32 chars)
JWT_EXPIRATION_MINUTES=15
AUTH_COOKIE_NAME=auth_token
COOKIE_SECURE=true

SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_KEY=<your-service-key>
STORAGE_BUCKET_NAME=user-uploads

LOG_LEVEL=info
```

### Generate JWT Secret

```bash
# On your machine
openssl rand -base64 32
# Output: AbCdEfGhIjKlMnOpQrStUvWxYz1234567890+/=
```

### Monitor Deployment

1. Railway Dashboard → Logs
2. Wait for "Build complete"
3. Wait for "Listening on port 3000"
4. Check Health: `https://ghostline-backend-production-xxxx.up.railway.app/health`

### Get Backend URL

```
Railway Dashboard → Deployments
Copy domain: ghostline-backend-production-xxxx.up.railway.app
```

---

## Step 4: Deploy Frontend to Vercel

### Create Vercel Project

1. Go to https://vercel.com
2. Click "Add New" → "Project"
3. Select "Import Git Repository"
4. Authorize GitHub if needed
5. Select `ghostline-frontend` repository
6. Click "Import"

### Configure Build Settings

- **Framework:** Vite
- **Build Command:** `npm run build`
- **Output Directory:** `dist`
- **Install Command:** `npm install`

### Set Environment Variables

Before deploying, add environment variables:

Project Settings → Environment Variables:

```
VITE_API_BASE_URL=https://ghostline-backend-production-xxxx.up.railway.app
VITE_WS_BASE_URL=wss://ghostline-backend-production-xxxx.up.railway.app
```

Replace `xxxx` with your actual Railway project ID.

### Deploy

Click "Deploy" - Vercel auto-builds and deploys.

### Get Frontend URL

```
Vercel Dashboard → Domains
Default: ghostline-frontend-<hash>.vercel.app
```

---

## Step 5: Update Backend ALLOWED_ORIGIN

Your frontend URL is now live. Update backend:

Railway Dashboard → Variables:

```
ALLOWED_ORIGIN=https://ghostline-frontend-<hash>.vercel.app
```

Redeploy backend:
```
Railway Dashboard → Deployments → Latest → Redeploy
```

Wait 30 seconds for new deployment.

---

## Step 6: Test Deployment

### Test Health Endpoints

```bash
# Backend health
curl https://ghostline-backend-production-xxxx.up.railway.app/health

# Response should be:
{"status":"ok","timestamp":"2026-03-28T10:00:00Z"}
```

### Test Frontend

```
Open: https://ghostline-frontend-<hash>.vercel.app
Should see login page
```

### Test Login

1. Register new account
2. Login with credentials
3. Should redirect to home page
4. Should see "Loading posts..." or feed

### Test API Connection

1. Go to http home page
2. Open DevTools (F12)
3. Network tab
4. Try creating a post
5. Should see POST /api/posts request
6. Status should be 200 or 201

### Test WebSocket

1. Open Chat page
2. Find another user
3. Send message
4. Message should appear instantly

---

## Monitoring & Maintenance

### Check Logs

#### Backend (Railway)

```
Railway Dashboard → Project → Logs
```

Look for:
- Deployment errors
- Database connection issues
- Application crashes

#### Frontend (Vercel)

```
Vercel Dashboard → Project → Deployments
```

Look for:
- Build failures
- Runtime errors

### Check Metrics

#### Railway

```
Railway Dashboard → Metrics
- CPU usage
- Memory usage
- Network I/O
```

#### Vercel

```
Vercel Dashboard → Analytics
- Page load times
- Core Web Vitals
- Error rate
```

### Database Monitoring

```
Supabase Dashboard → Monitoring
- Connection count
- Query performance
- Storage usage
```

---

## Troubleshooting

### Backend Won't Deploy

**Error:** "failed to build: Dockerfile cannot be empty"

**Solution:**
```bash
# Check Dockerfile exists and has content
ls -la Dockerfile
cat Dockerfile

# If empty, create it:
# See DEVELOPMENT.md for Dockerfile content
```

**Error:** "go.mod requires go >= 1.25"

**Solution:**
- Update go.mod: `go 1.25`
- Update Dockerfile: `FROM golang:1.25-alpine`

### API Connection Error

**Error:** "No 'Access-Control-Allow-Origin' header"

**Solution:**
1. Check frontend URL is correct
2. Update ALLOWED_ORIGIN in Railway
3. Redeploy backend
4. Hard refresh frontend (Ctrl+Shift+R)

### WebSocket Not Connecting

**Error:** "Failed to establish WebSocket connection"

**Solution:**
1. Check JWT token is valid (login first)
2. Verify backend is running
3. Check WebSocket URL uses `wss://` (not `ws://`)
4. Check browser allows WebSocket connections

### Database Connection Timeout

**Error:** "database connection failed"

**Solution:**
1. Check DATABASE_URL is correct
2. Verify Supabase database is running
3. Check connection limits (DB_MAX_CONNECTIONS)
4. Try using connection pooler:
   ```
   db.xxx.supabase.co:6543
   ```

---

## Custom Domain (Optional)

### Add Custom Domain to Frontend

1. Go to Vercel → Project → Settings → Domains
2. Add your domain (e.g., ghostline.com)
3. Follow DNS setup instructions
4. Deployment auto-creates SSL certificate

### Add Custom Domain to Backend

1. Go to Railway → Project → Settings
2. Add custom domain
3. Update DNS records
4. Wait for SSL certificate (5-10 minutes)
5. Update ALLOWED_ORIGIN in backend env vars

---

## CI/CD Pipeline

### Automatic Deployment

Both platforms automatically deploy when you push to `main`:

```bash
# Make changes locally
git add .
git commit -m "Feature: add new feature"
git push origin main

# Railway automatically:
# 1. Detects push
# 2. Runs build
# 3. Runs tests
# 4. Deploys new version

# Vercel automatically:
# 1. Detects push
# 2. Installs dependencies
# 3. Builds optimized bundle
# 4. Deploys to CDN
```

### Rollback

If deployment fails:

#### Railway
```
Deployments → Click previous version → Rollback
```

#### Vercel
```
Deployments → Click previous deployment → Redeploy
```

---

## Security Checklist

- ✅ All secrets in environment variables (not in code)
- ✅ JWT_SECRET is strong (32+ characters)
- ✅ COOKIE_SECURE=true in production
- ✅ ALLOWED_ORIGIN matches frontend URL exactly
- ✅ Database password is strong
- ✅ No hardcoded API keys in frontend
- ✅ HTTPS enabled on both services
- ✅ Database backups enabled (Supabase default)

---

## Performance Optimization

### Frontend

- Enable Vercel Analytics
- Use image optimization
- Enable gzip compression (default)
- Configure CDN caching

### Backend

- Monitor Railway metrics
- Scale if needed (more containers)
- Optimize database queries
- Enable connection pooling

### Database

- Check slow query logs
- Add indexes if needed
- Monitor connection count
- Optimize migration files

---

See also: [DEVELOPMENT.md](./DEVELOPMENT.md), [ARCHITECTURE.md](./ARCHITECTURE.md)
