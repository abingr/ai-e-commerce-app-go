# Phase 9 To-Do Test: Broader Automated Testing

This checklist verifies the broader test coverage added in Phase 9.

## 1. Run all default backend tests

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

Expected:

```text
ok
```

for all tested packages.

Code being tested:

- `backend/internal/handlers/*_test.go`
- `backend/internal/services/*_test.go`
- `backend/internal/routes/*_test.go`

## 2. Review service-level tests

Open:

```text
backend/internal/services/product_service_test.go
```

Verify these cases exist:

- product filters are passed to the repository
- create product delegates input correctly
- delete product delegates the ID correctly

## 3. Review middleware tests

Open:

```text
backend/internal/routes/middleware_test.go
```

Verify these cases exist:

- request ID is generated
- incoming `X-Request-ID` is reused
- CORS preflight returns `204 No Content`

## 4. Run repository integration test manually

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

Enable integration tests:

```powershell
$env:ECOMMERCE_RUN_INTEGRATION_TESTS="true"
go test ./internal/repositories
```

Expected:

```text
ok
```

If Docker Desktop is not running, this integration test cannot run yet. That is expected.

## 5. Run frontend checks

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm run build
npm audit
```

Expected:

- build passes
- audit reports no vulnerabilities
