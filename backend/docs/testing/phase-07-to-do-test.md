# Phase 7 To-Do Test: Validation, Error Responses, and Request IDs

This checklist verifies API hardening changes from Phase 7.

## 1. Run automated backend tests

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

Code being tested:

- `backend/internal/handlers/response.go`
- `backend/internal/handlers/auth_handler_test.go`
- `backend/internal/routes/middleware.go`
- `backend/internal/routes/middleware_test.go`

Expected result:

```text
ok
```

for all backend packages.

## 2. Start the backend API

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/api
```

Code being run:

- `backend/cmd/api/main.go`
- `backend/internal/routes/router.go`
- `backend/internal/routes/middleware.go`

## 3. Verify request ID in browser

Open:

```text
http://localhost:8080/health
```

Open browser developer tools, go to the Network tab, click the `/health` request, and check Response Headers.

Verify:

```text
X-Request-ID
```

is present.

Code being called:

- `requestID` in `backend/internal/routes/middleware.go`
- `HealthHandler.Health` in `backend/internal/handlers/health.go`

## 4. Verify custom request ID in Postman

In Postman, run:

```text
System > Health
```

Add a request header:

```text
X-Request-ID: my-postman-request-001
```

Verify in the response headers:

```text
X-Request-ID: my-postman-request-001
```

Code being called:

- `requestID` accepts the client-provided header
- `requestLogger` logs the same request ID after the handler finishes

## 5. Test validation details on registration

In Postman, run:

```text
Auth > Register Customer
```

Use this invalid body:

```json
{
  "name": "",
  "email": "not-an-email",
  "password": "short"
}
```

Expected status:

```text
400 Bad Request
```

Expected response includes:

```json
{
  "code": "VALIDATION_ERROR",
  "error": "invalid registration payload",
  "fields": [
    {
      "field": "name",
      "rule": "required"
    }
  ],
  "request_id": "..."
}
```

The exact order of `fields` can vary.

Code being called:

- `AuthHandler.Register` in `backend/internal/handlers/auth_handler.go`
- `JSONValidationError` in `backend/internal/handlers/response.go`

## 6. Test validation details on cart input

First login as a customer so `customer_token` is populated.

Run:

```text
Cart > Add Item To Cart
```

Use this invalid body:

```json
{
  "product_id": "not-a-uuid",
  "quantity": 0
}
```

Expected status:

```text
400 Bad Request
```

Expected response includes:

```json
{
  "code": "VALIDATION_ERROR",
  "error": "invalid cart item payload",
  "fields": [
    {
      "field": "product_id",
      "rule": "uuid"
    },
    {
      "field": "quantity",
      "rule": "min"
    }
  ]
}
```

Code being called:

- `CartHandler.AddItem` in `backend/internal/handlers/cart_handler.go`
- `JSONValidationError` in `backend/internal/handlers/response.go`

## 7. Test authentication error shape

Run:

```text
Orders > List My Orders
```

Remove the `Authorization` header.

Expected status:

```text
401 Unauthorized
```

Expected response:

```json
{
  "code": "UNAUTHORIZED",
  "error": "authorization header required",
  "request_id": "..."
}
```

Code being called:

- `requireAuth` in `backend/internal/routes/auth_middleware.go`
- `JSONError` in `backend/internal/handlers/response.go`

## 8. Watch structured logs

In the terminal running the backend, inspect the JSON log line after each request.

Verify it includes:

```text
request_id
method
path
status
latency_ms
client_ip
```

For internal server errors, `requestLogger` also records the attached error details from `c.Errors`.

## 9. Run frontend checks

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm install
npm run build
npm audit
```

Expected result:

- build completes successfully
- audit reports no vulnerabilities
