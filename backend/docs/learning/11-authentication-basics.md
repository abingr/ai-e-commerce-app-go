# Authentication Basics

Authentication answers the question:

> Who is making this request?

Phase 3 adds three endpoints:

```http
POST /api/v1/auth/register
POST /api/v1/auth/login
GET /api/v1/me
```

Registration creates a user account. Login verifies credentials. The `/me` endpoint proves the API can identify the current user from a JWT.

Request flow:

```text
HTTP request
-> AuthHandler
-> AuthService
-> UserRepository
-> PostgreSQL
```

Files involved:

- `internal/handlers/auth_handler.go`
- `internal/services/auth_service.go`
- `internal/repositories/user_repository.go`
- `internal/routes/auth_middleware.go`
- `migrations/000003_create_users_table.up.sql`

Backend engineer note:

- Authentication should never return the password hash.
- Login failure should not reveal whether the email or password was wrong.
- Protected routes should reject missing, malformed, invalid, or expired tokens.

CV talking point:

> Implemented JWT-based authentication with registration, login, protected routes, and role-aware user records.
