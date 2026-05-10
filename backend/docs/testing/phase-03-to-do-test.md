# Phase 3 To-Do-Test

Use this checklist after Phase 3 changes to verify registration, login, JWT authentication, and the protected `/me` endpoint.

## 1. Start PostgreSQL

From the project root:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go
docker compose up -d postgres
docker compose ps
```

Expected result:

- `ai_ecommerce_postgres` is `Up`.
- The status should become `healthy`.

## 2. Run User Migration

From the backend folder:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/migrate up
```

Expected result:

```text
Migrations completed.
```

or:

```text
No migration changes to apply.
```

Code path involved:

- `backend/cmd/migrate/main.go:16` starts the migration command.
- `backend/migrations/000003_create_users_table.up.sql` creates the `users` table.
- `backend/internal/config/config.go:26` loads the JWT secret.
- `backend/internal/config/config.go:27` loads the JWT issuer.

## 3. Run Automated Tests

From the backend folder:

```powershell
go test ./...
```

Expected result:

- Auth handler tests pass.
- Auth service tests pass.
- Existing health, product, and repository tests still pass.

What this verifies:

- Registration returns `201 Created`.
- Duplicate email maps to `409 Conflict`.
- Login with invalid credentials returns `401 Unauthorized`.
- Passwords are hashed with bcrypt.
- JWTs are generated and parsed.
- `/me` returns the current user when authentication context exists.

Code path involved:

- `backend/internal/handlers/auth_handler_test.go` tests HTTP behavior.
- `backend/internal/services/auth_service_test.go` tests hashing and JWT behavior.
- `backend/internal/services/auth_service.go:45` registers a user.
- `backend/internal/services/auth_service.go:74` logs a user in.
- `backend/internal/services/auth_service.go:121` parses a JWT.

## 4. Start The Backend API

From the backend folder:

```powershell
go run ./cmd/api
```

Expected result:

- The API listens on `http://localhost:8080`.

Code path involved:

- `backend/internal/routes/router.go:35` creates the user repository.
- `backend/internal/routes/router.go:36` creates the auth service.
- `backend/internal/routes/router.go:37` creates the auth handler.
- `backend/internal/routes/router.go:43` registers `POST /api/v1/auth/register`.
- `backend/internal/routes/router.go:44` registers `POST /api/v1/auth/login`.
- `backend/internal/routes/router.go:47` registers protected `GET /api/v1/me`.

## 5. Register A User With CLI

Use a unique email each time, or reuse the one below and expect `409` after the first success.

```powershell
$body = @{
  name = "Backend Learner"
  email = "learner@example.com"
  password = "password123"
} | ConvertTo-Json

$register = Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/register -Method Post -ContentType "application/json" -Body $body
$register.data.user
$register.data.token
```

Expected result:

- Response contains `data.user`.
- User role is `customer`.
- Response contains `data.token`.
- Response does not contain `password_hash`.

Code path involved:

- `backend/internal/handlers/auth_handler.go:29` handles registration.
- `backend/internal/services/auth_service.go:46` hashes the password.
- `backend/internal/repositories/user_repository.go:23` inserts the user.
- `backend/internal/services/auth_service.go:103` generates the JWT.

Browser verification:

Open:

```text
http://localhost:5173
```

Use the Account panel to register with:

```text
Backend Learner
learner@example.com
password123
```

Expected result:

- Auth status changes to `customer`.
- The panel shows the authenticated email and role.

Frontend code involved:

- `frontend/src/App.jsx:45` handles the form submit.
- `frontend/src/App.jsx:48` chooses register or login endpoint.
- `frontend/src/App.jsx:71` saves the JWT in `localStorage`.
- `frontend/src/App.jsx:124` displays auth status.

## 6. Login With CLI

```powershell
$body = @{
  email = "learner@example.com"
  password = "password123"
} | ConvertTo-Json

$login = Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/login -Method Post -ContentType "application/json" -Body $body
$token = $login.data.token
$login.data.user
```

Expected result:

- Response contains the user.
- Response contains a JWT token.

Code path involved:

- `backend/internal/handlers/auth_handler.go:58` handles login.
- `backend/internal/repositories/user_repository.go:42` finds the user by email.
- `backend/internal/services/auth_service.go:84` compares the bcrypt hash.
- `backend/internal/services/auth_service.go:103` generates a JWT.

## 7. Verify Protected `/me`

Use the token from login:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/me -Headers @{ Authorization = "Bearer $token" }
```

Expected result:

- Response contains the current user.
- The email matches the logged-in user.

Code path involved:

- `backend/internal/routes/auth_middleware.go:17` starts auth middleware.
- `backend/internal/routes/auth_middleware.go:27` reads the bearer token.
- `backend/internal/routes/auth_middleware.go:35` parses the JWT.
- `backend/internal/routes/auth_middleware.go:43` stores `user_id` in Gin context.
- `backend/internal/handlers/auth_handler.go:87` handles `/me`.
- `backend/internal/repositories/user_repository.go:61` loads the user by ID.

Browser verification:

After registering or logging in at:

```text
http://localhost:5173
```

Refresh the page.

Expected result:

- The frontend keeps the session using the stored token.
- Account panel still shows the signed-in user.

Frontend code involved:

- `frontend/src/App.jsx:32` reads the stored token.
- `frontend/src/App.jsx:78` calls `/api/v1/me`.
- `frontend/src/App.jsx:98` signs out by removing the token.

## 8. Verify Auth Error Cases

Missing token:

```powershell
try { Invoke-RestMethod -Uri http://localhost:8080/api/v1/me } catch { $_.Exception.Response.StatusCode.value__ }
```

Expected result:

```text
401
```

Invalid token:

```powershell
try { Invoke-RestMethod -Uri http://localhost:8080/api/v1/me -Headers @{ Authorization = "Bearer invalid-token" } } catch { $_.Exception.Response.StatusCode.value__ }
```

Expected result:

```text
401
```

Wrong password:

```powershell
$body = @{
  email = "learner@example.com"
  password = "wrong-password"
} | ConvertTo-Json

try { Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/login -Method Post -ContentType "application/json" -Body $body } catch { $_.Exception.Response.StatusCode.value__ }
```

Expected result:

```text
401
```

## 9. Build And Audit Frontend

From the frontend folder:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm install
npm run build
npm audit
```

Expected result:

- Build completes successfully.
- Audit reports `found 0 vulnerabilities`.

## Phase 3 Pass Criteria

Phase 3 is working when:

- `go run ./cmd/migrate up` creates or confirms the `users` table.
- `go test ./...` passes.
- Register creates a user and returns a token.
- Duplicate registration returns `409`.
- Login returns a token.
- Wrong password returns `401`.
- `/api/v1/me` rejects missing or invalid tokens.
- `/api/v1/me` returns the current user with a valid token.
- Frontend can register/login and persist session after refresh.
- `npm run build` passes.
- `npm audit` reports zero vulnerabilities.
