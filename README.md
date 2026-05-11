# ai-e-commerce-app-go

A production-style learning project for a mini electronics e-commerce application.

Phase 1 creates the foundation:

- Go + Gin backend API
- PostgreSQL local database through Docker Compose
- Health and readiness endpoints
- Environment-based configuration
- First HTTP test
- Swagger/OpenAPI learning notes
- Basic React frontend placeholder

Phase 2 adds:

- Product catalog database schema
- Versioned PostgreSQL migrations with seed data
- Product repository, service, and handler layers
- Product list and detail REST APIs
- Basic React product catalog display

Phase 3 adds:

- User registration and login
- Bcrypt password hashing
- JWT token generation and validation
- Protected `GET /api/v1/me` endpoint
- Basic React authentication panel

Phase 4 adds:

- Role-based admin authorization
- Admin product create/update/delete endpoints
- Soft delete for product removal
- Postman collection and environment for API testing

Phase 5 adds:

- Authenticated shopping cart APIs
- Cart item table with foreign keys and unique constraints
- Add/update/remove/clear cart operations
- User-scoped cart access using JWT identity

## Project layout

```text
backend/       Go API
frontend/      Basic React client
docker-compose.yml
```

## Run Phase 1

Start PostgreSQL:

```powershell
docker compose up -d postgres
```

The database is exposed on local port `55432` to avoid conflicts with an existing PostgreSQL installation.

Run the backend:

```powershell
cd backend
copy .env.example .env
go run ./cmd/api
```

Run database migrations:

```powershell
cd backend
go run ./cmd/migrate up
```

Check the API:

```powershell
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

Read the Phase 1 OpenAPI document:

```text
backend/docs/api/openapi.yaml
```

Follow the Phase 1 manual verification checklist:

```text
backend/docs/testing/phase-01-to-do-test.md
```

Follow the Phase 2 manual verification checklist:

```text
backend/docs/testing/phase-02-to-do-test.md
```

Follow the Phase 3 manual verification checklist:

```text
backend/docs/testing/phase-03-to-do-test.md
```

Follow the Phase 4 manual verification checklist:

```text
backend/docs/testing/phase-04-to-do-test.md
```

Follow the Phase 5 manual verification checklist:

```text
backend/docs/testing/phase-05-to-do-test.md
```

Run backend tests:

```powershell
cd backend
go test ./...
```

Run the frontend:

```powershell
cd frontend
npm install
npm run dev
```

## Phase roadmap

1. Foundation: Gin, PostgreSQL, health checks, docs, first test
2. Product catalog
3. User registration and JWT login
4. Admin role and product management
5. Cart
6. Orders and checkout without payment
7. Validation, error handling, and logging improvements
8. React integration
9. Broader tests
10. Dockerfile, CI, README polish, CV summary
