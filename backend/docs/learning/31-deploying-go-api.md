# Deploying the Go API

The backend is a long-running Go/Gin HTTP server. That makes container-style hosting a better first deployment target than serverless functions.

Recommended free-hosting options:

- Render Free Web Service
- Koyeb Free Web Service

## Why not deploy the Go API to Vercel first?

Vercel can run Go functions, but its Go runtime expects files inside an `/api` directory that export an `http.HandlerFunc`.

This project currently has:

```text
backend/cmd/api/main.go
```

That file starts a normal HTTP server. Converting it to Vercel Go Functions would be possible, but it would teach serverless adaptation before we have a stable cloud deployment.

For Phase 11, keep the backend as a normal Go service.

## Render deployment shape

This repo includes:

```text
render.yaml
```

It tells Render:

- the backend service lives in `backend`
- the service uses Docker
- `/ready` is the health check path
- important environment variables must be configured

## Required backend production variables

```text
ECOMMERCE_APP_ENV=production
ECOMMERCE_APP_NAME=ai-e-commerce-api
ECOMMERCE_HTTP_PORT=8080
ECOMMERCE_DATABASE_URL=<neon-pooled-connection-string>
ECOMMERCE_JWT_SECRET=<long-random-secret>
ECOMMERCE_JWT_ISSUER=ai-e-commerce-api
ECOMMERCE_CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

## Migration strategy

Before using the deployed API, run:

```powershell
go run ./cmd/migrate up
```

against the production database URL.

On a hosted platform, you can run this as:

- a one-off shell command if the platform supports it
- a local command with `ECOMMERCE_DATABASE_URL` set to the hosted Neon URL
- a later CI/CD migration job

For this learning project, running migrations manually is acceptable, but you should document it clearly.

## Backend lesson

Deploying is not just "put code online." You must provide configuration, database connectivity, health checks, secrets, and a repeatable migration process.
