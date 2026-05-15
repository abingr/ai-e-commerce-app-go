# Phase 10 To-Do Test: Docker, CI, and Final Project Checks

This checklist verifies the final project polish.

## 1. Run backend tests

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

Expected:

```text
ok
```

## 2. Run frontend checks

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm ci
npm run build
npm audit --audit-level=high
```

Expected:

- dependencies install cleanly
- build passes
- audit does not report high or critical vulnerabilities

## 3. Build backend Docker image

Start Docker Desktop first.

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go
docker build -t ai-ecommerce-backend ./backend
```

Expected:

```text
Successfully built
```

or Docker's newer build success output.

## 4. Build frontend Docker image

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go
docker build -t ai-ecommerce-frontend ./frontend
```

Expected:

```text
Successfully built
```

or Docker's newer build success output.

## 5. Run local full-stack manual test

Start PostgreSQL:

```powershell
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

- products load
- login/register works
- cart works
- checkout creates an order
- order history updates

## 6. Verify CI file

Open:

```text
.github/workflows/ci.yml
```

Verify it includes:

- backend `go test ./...`
- backend formatting check
- frontend `npm ci`
- frontend `npm run build`
- frontend `npm audit --audit-level=high`

## 7. Final GitHub check

After pushing, open your GitHub repository and verify:

- all files are at the repository root
- GitHub Actions tab shows the CI workflow
- latest commit is visible on `main`
