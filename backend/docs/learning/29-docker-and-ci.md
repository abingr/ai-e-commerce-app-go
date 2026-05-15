# Docker and CI

Phase 10 adds Dockerfiles and a GitHub Actions CI workflow.

## Dockerfiles

The backend Dockerfile is:

```text
backend/Dockerfile
```

It uses a multi-stage build:

1. `golang:1.24-alpine` downloads modules and builds the API and migrate binaries.
2. `alpine:3.21` runs only the compiled binaries and migrations folder.

The frontend Dockerfile is:

```text
frontend/Dockerfile
```

It also uses a multi-stage build:

1. `node:20-alpine` installs dependencies and builds the Vite app.
2. `nginx:1.27-alpine` serves the static files.

## Docker Compose

`docker-compose.yml` still provides PostgreSQL for local development.

The database runs on local port:

```text
55432
```

The backend default database URL points to that port in development.

## CI workflow

The workflow is:

```text
.github/workflows/ci.yml
```

It runs on pushes to `main` and on pull requests.

Backend checks:

- Go setup
- `gofmt` check
- `go test ./...`

Frontend checks:

- Node setup
- `npm ci`
- `npm run build`
- `npm audit --audit-level=high`

## Backend lesson

CI is a safety net. It does not replace local testing, but it prevents broken code from quietly landing on the main branch.
