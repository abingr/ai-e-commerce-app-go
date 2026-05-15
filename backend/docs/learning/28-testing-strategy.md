# Testing Strategy

Phase 9 broadens the test suite so the project demonstrates more than happy-path coding.

## Test layers in this project

The backend now has:

- handler tests for HTTP status codes and response bodies
- service tests for business workflow delegation
- middleware tests for request IDs, CORS, auth roles, and validation behavior
- repository integration tests that can be enabled when PostgreSQL is running

## Unit tests

Unit tests do not require external services. They use small stub implementations.

Examples:

- `backend/internal/handlers/auth_handler_test.go`
- `backend/internal/services/product_service_test.go`
- `backend/internal/routes/middleware_test.go`

These tests are fast and should run on every commit.

## Integration tests

The product repository integration test is intentionally opt-in:

```powershell
$env:ECOMMERCE_RUN_INTEGRATION_TESTS="true"
go test ./internal/repositories
```

It requires:

- PostgreSQL running
- migrations already applied
- seed data available

Keeping it opt-in prevents normal test runs from failing just because Docker is not open.

## Backend lesson

Senior backend engineers test at multiple levels. Unit tests catch logic mistakes quickly. Integration tests catch database and wiring problems. CI runs the reliable checks automatically so the main branch stays healthy.
