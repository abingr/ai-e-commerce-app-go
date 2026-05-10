# Environment Configuration

Applications should not hard-code environment-specific values such as ports, database URLs, passwords, or JWT secrets.

This project loads configuration in `internal/config/config.go`.

Examples:

- `ECOMMERCE_APP_ENV`
- `ECOMMERCE_APP_NAME`
- `ECOMMERCE_HTTP_PORT`
- `ECOMMERCE_DATABASE_URL`

The `.env.example` file documents which variables are expected. It is safe to commit because it contains local development defaults only.

For local development, copy `.env.example` to `.env`. The backend loads `.env` on startup and falls back to safe development defaults when a variable is missing.

Backend engineer note:

- Development, test, staging, and production should be able to use the same code with different environment variables.
- Secrets should not be committed to Git.
- Configuration should be loaded once at startup and passed into the application.
