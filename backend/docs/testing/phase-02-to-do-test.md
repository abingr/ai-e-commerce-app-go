# Phase 2 To-Do-Test

Use this checklist after Phase 2 changes to verify the product catalog.

## 1. Start PostgreSQL

From the project root:

```powershell
cd c:\Training\Golang\AI_Workspace
docker compose up -d postgres
docker compose ps
```

Expected result:

- `ai_ecommerce_postgres` is `Up`.
- The status should become `healthy`.
- Local database port is `55432`.

Code/config involved:

- `docker-compose.yml` exposes PostgreSQL on `55432`.
- `backend/internal/config/config.go:23` uses that port in the default database URL.

## 2. Run Product Migrations

From the backend folder:

```powershell
cd c:\Training\Golang\AI_Workspace\backend
go run ./cmd/migrate up
```

Expected result:

```text
Migrations completed.
```

or, if you already ran it:

```text
No migration changes to apply.
```

Code path involved:

- `backend/cmd/migrate/main.go:16` starts the migration command.
- `backend/cmd/migrate/main.go:21` loads the database config.
- `backend/cmd/migrate/main.go:23` points the migrator at `backend/migrations`.
- `backend/migrations/000001_create_products_table.up.sql` creates the table.
- `backend/migrations/000002_seed_products.up.sql` inserts electronics products.

## 3. Run Automated Tests

From the backend folder:

```powershell
go test ./...
```

Expected result:

- Handler tests pass.
- Repository integration test is skipped by default unless explicitly enabled.

What this verifies:

- `GET /api/v1/products` can return products.
- `GET /api/v1/products/:id` can return a product.
- Invalid product IDs return `400`.
- Missing products return `404`.
- Handler errors return `500`.

Code path involved:

- `backend/internal/handlers/product_handler_test.go` tests HTTP behavior.
- `backend/internal/handlers/product_handler.go:28` runs product listing.
- `backend/internal/handlers/product_handler.go:47` runs product detail.

Optional database-backed repository test:

```powershell
$env:ECOMMERCE_RUN_INTEGRATION_TESTS="true"
go test ./internal/repositories -v
Remove-Item Env:ECOMMERCE_RUN_INTEGRATION_TESTS
```

Expected result:

- The test connects to PostgreSQL.
- The test confirms seeded products exist.

Code path involved:

- `backend/internal/repositories/product_repository_integration_test.go:15` starts the integration test.
- `backend/internal/repositories/product_repository.go:23` queries the `products` table.

## 4. Start The Backend API

From the backend folder:

```powershell
go run ./cmd/api
```

Expected result:

- The command keeps running.
- API listens on `http://localhost:8080`.

Code path involved:

- `backend/cmd/api/main.go:16` starts the API.
- `backend/cmd/api/main.go:23` connects to PostgreSQL.
- `backend/internal/routes/router.go:32` creates the product repository.
- `backend/internal/routes/router.go:33` creates the product service.
- `backend/internal/routes/router.go:34` creates the product handler.
- `backend/internal/routes/router.go:40` registers `GET /api/v1/products`.
- `backend/internal/routes/router.go:41` registers `GET /api/v1/products/:id`.

## 5. Verify Product List With CLI

Open a second terminal:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/products
```

Expected result:

- Response has a `data` field.
- `data` contains seeded electronics products.

Code path involved:

- `backend/internal/routes/router.go:40` maps the route.
- `backend/internal/handlers/product_handler.go:28` reads query parameters.
- `backend/internal/services/product_service.go:22` calls the repository.
- `backend/internal/repositories/product_repository.go:23` runs the SQL query.

Browser verification:

Open:

```text
http://localhost:8080/api/v1/products
```

You should see JSON with a `data` array.

## 6. Verify Product Filtering

CLI:

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/products?category=Accessories"
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/products?search=phone"
```

Expected result:

- Category filter returns only matching category products.
- Search filter returns products whose name, description, or brand contains the search text.

Code path involved:

- `backend/internal/handlers/product_handler.go:30` reads `category`.
- `backend/internal/handlers/product_handler.go:31` reads `search`.
- `backend/internal/repositories/product_repository.go:32` adds the category SQL condition.
- `backend/internal/repositories/product_repository.go:37` adds the search SQL condition.

Browser verification:

Open:

```text
http://localhost:8080/api/v1/products?category=Accessories
http://localhost:8080/api/v1/products?search=phone
```

## 7. Verify Product Detail

First get a product ID from the list response. Then run:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/products/YOUR_PRODUCT_ID
```

Expected result:

- Response has a `data` object.
- The product ID matches the ID in the URL.

Code path involved:

- `backend/internal/routes/router.go:41` maps the route.
- `backend/internal/handlers/product_handler.go:47` reads the route ID.
- `backend/internal/handlers/product_handler.go:49` validates UUID format.
- `backend/internal/services/product_service.go:26` calls the repository.
- `backend/internal/repositories/product_repository.go:68` queries one product by ID.

Browser verification:

Open:

```text
http://localhost:8080/api/v1/products/YOUR_PRODUCT_ID
```

## 8. Verify Error Cases

Invalid UUID:

```powershell
try { Invoke-RestMethod -Uri http://localhost:8080/api/v1/products/not-a-uuid } catch { $_.Exception.Response.StatusCode.value__ }
```

Expected result:

- HTTP status is `400`.
- Response says `invalid product id`.

Valid UUID but not found:

```powershell
try { Invoke-RestMethod -Uri http://localhost:8080/api/v1/products/00000000-0000-0000-0000-000000000000 } catch { $_.Exception.Response.StatusCode.value__ }
```

Expected result:

- HTTP status is `404`.
- Response says `product not found`.

## 9. Verify Frontend Product Catalog

From the frontend folder:

```powershell
cd c:\Training\Golang\AI_Workspace\frontend
npm install
npm run dev -- --port 5173
```

Browser verification:

Open:

```text
http://localhost:5173
```

Expected result:

- Page shows `Backend Learning Storefront`.
- API status shows `ok`.
- Products status shows `ready`.
- Product cards appear with images, names, descriptions, prices, brands, and stock.

Code path involved:

- `frontend/src/App.jsx:8` stores product loading state.
- `frontend/src/App.jsx:16` calls `GET /api/v1/products`.
- `frontend/src/App.jsx:48` renders product cards.
- `frontend/src/App.jsx:69` formats `price_cents` as dollars.
- `backend/internal/routes/middleware.go:11` allows frontend requests from `localhost:5173`.

## 10. Build And Audit Frontend

```powershell
npm run build
npm audit
```

Expected result:

- Build completes successfully.
- Audit reports `found 0 vulnerabilities`.

## Phase 2 Pass Criteria

Phase 2 is working when:

- PostgreSQL is running.
- `go run ./cmd/migrate up` completes or reports no changes.
- `go test ./...` passes.
- Optional repository integration test passes when enabled.
- `GET /api/v1/products` returns seeded products.
- `category` and `search` filters work.
- `GET /api/v1/products/:id` returns one product.
- Invalid product ID returns `400`.
- Unknown product ID returns `404`.
- React frontend displays product cards.
- `npm run build` passes.
- `npm audit` reports zero vulnerabilities.
