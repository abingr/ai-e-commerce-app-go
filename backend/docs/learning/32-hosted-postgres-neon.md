# Hosted Postgres with Neon

Phase 11 recommends Neon for the hosted PostgreSQL database.

## Why Neon fits this project

Neon is a managed PostgreSQL provider with a free tier suitable for portfolio projects and learning.

For this app, Neon provides:

- hosted PostgreSQL
- SSL database connections
- connection pooling
- browser dashboard
- no local Docker requirement in production

## Connection string

Neon gives a connection string similar to:

```text
postgresql://user:password@host/dbname?sslmode=require
```

Use the pooled connection string when possible, especially for serverless or small free hosting environments.

Set it as:

```text
ECOMMERCE_DATABASE_URL=<neon-pooled-connection-string>
```

Never commit this value to Git.

## Running migrations

From your local machine:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
$env:ECOMMERCE_DATABASE_URL="<neon-pooled-connection-string>"
go run ./cmd/migrate up
```

Then verify:

```powershell
go run ./cmd/api
```

with the same database URL.

## Local vs hosted database

Local development:

```text
postgres://ecommerce_user:ecommerce_password@127.0.0.1:55432/ecommerce_db?sslmode=disable
```

Hosted Neon:

```text
postgresql://...neon.tech/...?...sslmode=require
```

The code does not need to change. Only the environment variable changes.

## Backend lesson

Using managed databases teaches an important production habit: infrastructure details should live in configuration, not in application code.
