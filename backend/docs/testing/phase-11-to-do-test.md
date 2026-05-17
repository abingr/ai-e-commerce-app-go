# Phase 11 To-Do Test: Free Hosting Deployment Preparation

This checklist verifies the app is ready for Vercel frontend hosting, hosted Go API deployment, and Neon PostgreSQL.

## 1. Run local automated checks

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm ci
npm run build
npm audit --audit-level=high
```

Expected:

- Go tests pass
- frontend build passes
- audit reports no high or critical vulnerabilities

## 2. Verify backend CORS config locally

Open:

```text
backend/.env.example
```

Verify this value exists:

```text
ECOMMERCE_CORS_ALLOWED_ORIGINS=http://localhost:5173
```

Code being tested:

- `Load` in `backend/internal/config/config.go`
- `cors` in `backend/internal/routes/middleware.go`

## 3. Verify frontend API config locally

Open:

```text
frontend/.env.example
```

Verify:

```text
VITE_API_BASE_URL=http://localhost:8080
```

Code being used:

- `API_BASE_URL` in `frontend/src/App.jsx`

## 4. Start local full stack

Start PostgreSQL:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go
docker compose up -d postgres
```

Run migrations:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/migrate up
```

Run backend:

```powershell
go run ./cmd/api
```

Run frontend:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm run dev
```

Open:

```text
http://localhost:5173
```

Verify:

- product list loads
- login/register works
- cart works
- checkout creates an order

## 5. Create Neon database

In Neon:

1. Create a project.
2. Copy the pooled PostgreSQL connection string.
3. Keep `sslmode=require`.

Set locally for migration:

```powershell
$env:ECOMMERCE_DATABASE_URL="<neon-pooled-connection-string>"
```

Run:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/migrate up
```

Expected:

```text
Migrations completed.
```

## 6. Deploy backend to Render or Koyeb

Use:

```text
backend/Dockerfile
```

For Render, this repo includes:

```text
render.yaml
```

Set backend environment variables:

```text
ECOMMERCE_APP_ENV=production
ECOMMERCE_APP_NAME=ai-e-commerce-api
ECOMMERCE_HTTP_PORT=8080
ECOMMERCE_DATABASE_URL=<neon-pooled-connection-string>
ECOMMERCE_JWT_SECRET=<long-random-secret>
ECOMMERCE_JWT_ISSUER=ai-e-commerce-api
ECOMMERCE_CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

Verify backend:

```text
https://your-backend-url/health
https://your-backend-url/ready
```

## 7. Deploy frontend to Vercel

Set Vercel root directory to:

```text
frontend
```

Set frontend environment variable:

```text
VITE_API_BASE_URL=https://your-backend-url
```

Deploy.

After Vercel gives you a URL, update the backend:

```text
ECOMMERCE_CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

Redeploy or restart the backend.

## 8. Production browser test

Open your Vercel URL.

Verify:

- API status is `ok`
- products load from deployed backend
- register/login works
- cart works
- checkout creates an order
- order history appears

## 9. Common deployment errors

If the frontend says API is offline:

- check `VITE_API_BASE_URL`
- check backend `/health`
- check backend CORS allowed origin
- check browser DevTools Network tab

If backend `/ready` fails:

- check `ECOMMERCE_DATABASE_URL`
- check Neon database status
- confirm migrations ran

If auth works locally but not in production:

- check `ECOMMERCE_JWT_SECRET`
- check `ECOMMERCE_JWT_ISSUER`
- clear browser `localStorage.auth_token` and login again
