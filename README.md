# ai-e-commerce-app-go

A production-style learning project for a mini electronics e-commerce application built with Go, Gin, PostgreSQL, JWT authentication, React, OpenAPI, Postman, Docker, and GitHub Actions.

## Highlights

- Layered backend architecture: handlers, services, repositories, models
- PostgreSQL migrations and seed data
- Product catalog, admin product management, cart, checkout, and order history
- JWT authentication, bcrypt password hashing, and role-based authorization
- User-scoped resources and transaction-backed checkout
- Consistent validation errors, request IDs, and structured logs
- React frontend connected to the backend API
- OpenAPI documentation and Postman collection
- Automated tests, Dockerfiles, and CI workflow

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

Phase 6 adds:

- Checkout without payment
- Order and order item database tables
- Transaction-backed order creation from the cart
- Order history APIs for authenticated users
- Product name and price snapshots on order items

Phase 7 adds:

- Consistent API error responses
- Field-level validation error details
- Request IDs with `X-Request-ID`
- Structured request logging with internal error capture

Phase 8 adds:

- React integration with backend APIs
- Product search and category filtering
- Browser login/register flow with JWT storage
- Cart management from the frontend
- Checkout and order history from the frontend

Phase 9 adds:

- Broader service and middleware tests
- Documented testing strategy
- Opt-in database integration test guidance

Phase 10 adds:

- Backend and frontend Dockerfiles
- GitHub Actions CI
- Final README polish
- CV summary document

Phase 11 adds:

- Deployment-ready backend CORS configuration
- Vite environment variable for the frontend API URL
- Render deployment blueprint
- Vercel frontend configuration
- Neon/Postgres deployment learning docs

## Project layout

```text
backend/       Go API
frontend/      Basic React client
docs/          CV/project summary notes
docker-compose.yml
.github/       GitHub Actions workflow
```

## Run Locally

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

Import Postman files:

```text
postman/ai-e-commerce-app-go.postman_collection.json
postman/ai-e-commerce-local.postman_environment.json
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

Follow the Phase 6 manual verification checklist:

```text
backend/docs/testing/phase-06-to-do-test.md
```

Follow the Phase 7 manual verification checklist:

```text
backend/docs/testing/phase-07-to-do-test.md
```

Follow the Phase 8 manual verification checklist:

```text
backend/docs/testing/phase-08-to-do-test.md
```

Follow the Phase 9 manual verification checklist:

```text
backend/docs/testing/phase-09-to-do-test.md
```

Follow the Phase 10 manual verification checklist:

```text
backend/docs/testing/phase-10-to-do-test.md
```

Follow the Phase 11 manual verification checklist:

```text
backend/docs/testing/phase-11-to-do-test.md
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

Build Docker images:

```powershell
docker build -t ai-ecommerce-backend ./backend
docker build -t ai-ecommerce-frontend ./frontend
```

## CI

GitHub Actions workflow:

```text
.github/workflows/ci.yml
```

It runs backend formatting/tests and frontend install/build/audit checks.

## CV Notes

Use this summary when preparing your resume or interview notes:

```text
docs/cv-summary.md
```

Free-hosting deployment guide:

```text
docs/deployment-free-hosting.md
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
