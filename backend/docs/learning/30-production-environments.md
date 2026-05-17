# Production Environments

Phase 11 prepares the project for cloud deployment.

Local development and production should not use the same settings. A production app needs different values for:

- database URL
- JWT secret
- allowed frontend origin
- app environment
- public API base URL used by the frontend

## Backend environment variables

The backend reads configuration in:

```text
backend/internal/config/config.go
```

Important variables:

```text
ECOMMERCE_APP_ENV
ECOMMERCE_APP_NAME
ECOMMERCE_HTTP_PORT
ECOMMERCE_DATABASE_URL
ECOMMERCE_JWT_SECRET
ECOMMERCE_JWT_ISSUER
ECOMMERCE_CORS_ALLOWED_ORIGINS
```

## CORS origins

In local development, the React app runs at:

```text
http://localhost:5173
```

In production, the frontend might run at:

```text
https://your-app.vercel.app
```

The backend now accepts a comma-separated list:

```text
ECOMMERCE_CORS_ALLOWED_ORIGINS=http://localhost:5173,https://your-app.vercel.app
```

The CORS middleware is in:

```text
backend/internal/routes/middleware.go
```

It only returns `Access-Control-Allow-Origin` when the request `Origin` header is in the allowed list.

## Frontend environment variables

The frontend reads:

```text
VITE_API_BASE_URL
```

from:

```text
frontend/src/App.jsx
```

Local value:

```text
VITE_API_BASE_URL=http://localhost:8080
```

Production value:

```text
VITE_API_BASE_URL=https://your-backend-url.onrender.com
```

## Backend lesson

Production readiness starts with configuration. Hardcoded local URLs are fine at the beginning, but deployable systems need environment variables so each environment can provide its own safe values.
