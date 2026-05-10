# Phase 1 To-Do-Test

Use this file after Phase 1 changes to verify the foundation is working.

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
- Local database port is `55432`.

Code/config involved:

- `docker-compose.yml` starts the PostgreSQL container.
- `backend/internal/config/config.go:23` uses `127.0.0.1:55432` as the default database URL.

## 2. Run Backend Automated Tests

From the backend folder:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

Expected result:

```text
ok   ai-e-commerce-app-go/backend/internal/handlers
```

What this verifies:

- The health handler returns HTTP `200 OK`.
- The response body is valid JSON.
- The response contains `status: ok`.
- The response contains `service: ai-e-commerce-api`.

Code path involved:

- `backend/internal/handlers/health_test.go:15` runs `TestHealthReturnsOK`.
- `backend/internal/handlers/health_test.go:19` creates `HealthHandler`.
- `backend/internal/handlers/health_test.go:20` registers `GET /health`.
- `backend/internal/handlers/health_test.go:25` sends the fake HTTP request.
- `backend/internal/handlers/health.go:26` runs `HealthHandler.Health`.

## 3. Start The Backend API

From the backend folder:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/api
```

Expected result:

- The command keeps running.
- You should see a JSON log saying the API server is starting on port `8080`.

Code path involved:

- `backend/cmd/api/main.go:16` starts the application.
- `backend/cmd/api/main.go:17` loads configuration with `config.Load`.
- `backend/internal/config/config.go:16` loads `.env` and fallback values.
- `backend/cmd/api/main.go:23` connects to PostgreSQL.
- `backend/internal/database/postgres.go:10` creates the connection pool.
- `backend/internal/database/postgres.go:29` pings PostgreSQL.
- `backend/cmd/api/main.go:30` creates the Gin router.
- `backend/internal/routes/router.go:19` registers middleware and routes.
- `backend/cmd/api/main.go:39` starts listening for HTTP requests.

## 4. Verify Health Endpoint With CLI

Open a second terminal and run:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/health
```

Expected result:

```text
service           status
-------           ------
ai-e-commerce-api ok
```

Code path involved:

- `backend/internal/routes/router.go:31` maps `GET /health`.
- `backend/internal/handlers/health.go:26` runs `HealthHandler.Health`.
- `backend/internal/handlers/health.go:27` returns the JSON response.
- `backend/internal/routes/middleware.go:26` logs the request after the handler runs.

Browser verification:

Open this URL:

```text
http://localhost:8080/health
```

You should see JSON similar to:

```json
{
  "service": "ai-e-commerce-api",
  "status": "ok"
}
```

## 5. Verify Readiness Endpoint With CLI

Run:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/ready
```

Expected result:

```text
database  status
--------  ------
connected ready
```

Code path involved:

- `backend/internal/routes/router.go:32` maps `GET /ready`.
- `backend/internal/handlers/health.go:33` runs `HealthHandler.Ready`.
- `backend/internal/handlers/health.go:37` pings PostgreSQL.
- `backend/internal/handlers/health.go:45` returns the success response.

Browser verification:

Open this URL:

```text
http://localhost:8080/ready
```

You should see JSON similar to:

```json
{
  "database": "connected",
  "status": "ready"
}
```

Failure test:

Stop PostgreSQL:

```powershell
docker compose stop postgres
Invoke-RestMethod -Uri http://localhost:8080/ready
```

Expected result:

- The request should fail with HTTP `503 Service Unavailable`.
- This proves readiness depends on database connectivity.

Start PostgreSQL again afterward:

```powershell
docker compose up -d postgres
```

## 6. Verify OpenAPI Documentation File

Open:

```text
c:\Training\Golang\ai-e-commerce-app-go\backend\docs\api\openapi.yaml
```

Expected result:

- `/health` is documented.
- `/ready` is documented.
- Both endpoints include expected response descriptions.

Code/spec relationship:

- `backend/docs/api/openapi.yaml` documents the routes registered in `backend/internal/routes/router.go:31` and `backend/internal/routes/router.go:32`.

## 7. Start And Verify Frontend

From the frontend folder:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm install
npm run dev -- --port 5173
```

Browser verification:

Open:

```text
http://localhost:5173
```

Expected result:

- The page title says `Backend Learning Storefront`.
- `API status` eventually shows `ok`.

Code path involved:

- `frontend/src/App.jsx:5` renders the React component.
- `frontend/src/App.jsx:8` runs the `useEffect` hook.
- `frontend/src/App.jsx:9` calls `http://localhost:8080/health`.
- `frontend/src/App.jsx:11` stores the returned `status`.
- `frontend/src/App.jsx:26` displays the API status.
- `backend/internal/routes/middleware.go:11` allows requests from `http://localhost:5173`.

## 8. Build And Audit Frontend

From the frontend folder:

```powershell
npm run build
npm audit
```

Expected result:

- `npm run build` completes successfully.
- `npm audit` reports `found 0 vulnerabilities`.

What this verifies:

- The React app compiles.
- The current frontend dependency set has no known npm audit issues.

## Phase 1 Pass Criteria

Phase 1 is considered working when:

- `docker compose ps` shows PostgreSQL running and healthy.
- `go test ./...` passes.
- `go run ./cmd/api` starts the API.
- Browser and CLI checks for `/health` return `status: ok`.
- Browser and CLI checks for `/ready` return `database: connected`.
- React page loads at `http://localhost:5173`.
- React page shows `API status ok`.
- `npm run build` passes.
- `npm audit` reports zero vulnerabilities.
