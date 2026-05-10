# Testing Basics

Backend tests protect behavior while the codebase grows.

Phase 1 includes a handler test for `GET /health` in `internal/handlers/health_test.go`.

The test uses:

- `testing`, Go's built-in test package
- `httptest`, a standard package for testing HTTP handlers
- Gin test mode

The test verifies:

- HTTP status is `200 OK`
- Response is valid JSON
- `status` is `ok`
- `service` is `ai-e-commerce-api`

Testing habit:

Start with small tests around stable behavior. As the project grows, add repository tests, service tests, authentication tests, and integration tests.
