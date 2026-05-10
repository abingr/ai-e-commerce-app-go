# Phase 4 To-Do-Test

Use this checklist after Phase 4 changes to verify admin product management and role-based authorization.

## 1. Start The Application

From the project root:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go
docker compose up -d postgres
```

From the backend folder:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/migrate up
go run ./cmd/api
```

Expected result:

- PostgreSQL is running.
- Migrations are applied.
- API listens on `http://localhost:8080`.

## 2. Run Automated Tests

From the backend folder:

```powershell
go test ./...
```

Expected result:

- Product handler tests pass.
- RBAC middleware tests pass.
- Existing auth and product tests still pass.

Code path involved:

- `backend/internal/handlers/product_handler_test.go` tests admin product handler behavior.
- `backend/internal/routes/auth_middleware_test.go` tests role authorization.
- `backend/internal/routes/auth_middleware.go:52` checks user role.

## 3. Import Postman Files

Open Postman and import:

```text
postman/ai-e-commerce-app-go.postman_collection.json
postman/ai-e-commerce-local.postman_environment.json
```

Select the environment:

```text
AI E-Commerce Local
```

Environment variables used:

- `base_url`
- `customer_email`
- `admin_email`
- `password`
- `customer_token`
- `admin_token`
- `product_id`

## 4. Register Customer In Postman

Run:

```text
Auth -> Register Customer
```

Expected result:

- Status `201 Created`, or `409 Conflict` if the user already exists.
- If status is `201`, Postman stores `customer_token` automatically.

Then run:

```text
Auth -> Login Customer
```

Expected result:

- Status `200 OK`.
- Postman stores `customer_token`.

## 5. Verify Customer Cannot Create Product

Run:

```text
Admin Products -> Customer Cannot Create Product
```

Expected result:

- Status `403 Forbidden`.
- Response says `insufficient permissions`.

Code path involved:

- `backend/internal/routes/router.go:49` applies auth and admin-role middleware to `/admin`.
- `backend/internal/routes/auth_middleware.go:17` validates JWT.
- `backend/internal/routes/auth_middleware.go:52` checks role.

## 6. Create Or Promote Admin User

Register an admin user as a normal customer first:

In Postman, set `admin_email` to:

```text
admin@example.com
```

Temporarily change the body of:

```text
Auth -> Register Customer
```

to use:

```json
{
  "name": "Admin User",
  "email": "{{admin_email}}",
  "password": "{{password}}"
}
```

Run the request. Then promote that user in PostgreSQL:

```powershell
docker exec ai_ecommerce_postgres psql -U ecommerce_user -d ecommerce_db -c "UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';"
```

Expected result:

```text
UPDATE 1
```

Then run:

```text
Auth -> Login Admin
```

Expected result:

- Status `200 OK`.
- Postman stores `admin_token`.

## 7. Create Product As Admin

Run:

```text
Admin Products -> Create Product
```

Expected result:

- Status `201 Created`.
- Response contains a new product.
- Postman stores `product_id`.

Code path involved:

- `backend/internal/routes/router.go:50` maps admin create product.
- `backend/internal/handlers/product_handler.go:79` handles create.
- `backend/internal/repositories/product_repository.go:86` inserts the product.

## 8. Update Product As Admin

Run:

```text
Admin Products -> Update Product
```

Expected result:

- Status `200 OK`.
- Product name changes to `USB-C Docking Station Pro`.

Code path involved:

- `backend/internal/routes/router.go:51` maps admin update product.
- `backend/internal/handlers/product_handler.go:101` handles update.
- `backend/internal/repositories/product_repository.go:96` updates the product.

## 9. Delete Product As Admin

Run:

```text
Admin Products -> Delete Product
```

Expected result:

- Status `204 No Content`.

Code path involved:

- `backend/internal/routes/router.go:52` maps admin delete product.
- `backend/internal/handlers/product_handler.go:138` handles delete.
- `backend/internal/repositories/product_repository.go:123` soft deletes the product.

## 10. Confirm Soft Delete

Run:

```text
Products -> Get Product Detail
```

Expected result:

- Status `404 Not Found`.

Run:

```text
Products -> List Products
```

Expected result:

- The deleted product is not visible in the public product list.

Why:

- `backend/internal/repositories/product_repository.go:27` lists only active products.
- Delete sets `is_active = false`.

## 11. Optional CLI Smoke Test

Register customer:

```powershell
$customerEmail = "customer+$([DateTimeOffset]::UtcNow.ToUnixTimeSeconds())@example.com"
$body = @{ name = "Customer User"; email = $customerEmail; password = "password123" } | ConvertTo-Json
$customer = Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/register -Method Post -ContentType "application/json" -Body $body
```

Verify customer cannot create product:

```powershell
$productBody = @{
  name = "USB-C Docking Station"
  description = "Multi-port dock for laptops."
  brand = "DockPro"
  category = "Accessories"
  price_cents = 15900
  stock_quantity = 15
  image_url = "https://example.com/dock.jpg"
} | ConvertTo-Json

try {
  Invoke-RestMethod -Uri http://localhost:8080/api/v1/admin/products -Method Post -ContentType "application/json" -Headers @{ Authorization = "Bearer $($customer.data.token)" } -Body $productBody
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

Expected result:

```text
403
```

## 12. Build And Audit Frontend

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

## Phase 4 Pass Criteria

Phase 4 is working when:

- `go test ./...` passes.
- Postman collection imports successfully.
- Customer token can access customer routes.
- Customer token receives `403` on admin product creation.
- Admin token can create a product.
- Admin token can update the product.
- Admin token can delete the product.
- Deleted product returns `404` on public detail endpoint.
- Deleted product is hidden from public list.
- `npm run build` passes.
- `npm audit` reports zero vulnerabilities.
