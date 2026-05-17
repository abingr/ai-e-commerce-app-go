# Free Hosting Deployment Plan

Recommended free-hosting layout:

```text
Vercel
  React/Vite frontend

Render or Koyeb
  Go/Gin backend API

Neon
  Hosted PostgreSQL
```

## Required environment variables

Backend:

```text
ECOMMERCE_APP_ENV=production
ECOMMERCE_APP_NAME=ai-e-commerce-api
ECOMMERCE_HTTP_PORT=8080
ECOMMERCE_DATABASE_URL=<neon-pooled-connection-string>
ECOMMERCE_JWT_SECRET=<long-random-secret>
ECOMMERCE_JWT_ISSUER=ai-e-commerce-api
ECOMMERCE_CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

Frontend:

```text
VITE_API_BASE_URL=https://your-backend-url.onrender.com
```

## Deployment order

1. Create Neon Postgres database.
2. Deploy backend to Render or Koyeb.
3. Run backend migrations against Neon.
4. Deploy frontend to Vercel.
5. Update backend CORS origin with the final Vercel URL.
6. Test health, products, login, cart, checkout, and orders.

## Why not Vercel for the Go backend first?

Vercel is excellent for frontend hosting, but the current backend is a normal long-running Gin HTTP server. Vercel's Go support is function-based. That is useful, but it would require reshaping the backend into serverless handlers.

For this project, deploying the frontend to Vercel and the Go API to Render or Koyeb is a cleaner first production deployment.
