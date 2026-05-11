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

Before continuing, confirm these environment values:

```text
base_url       = http://localhost:8080
customer_email = customer01@gmail.com
admin_email    = admin01@gmail.com
password       = 12345678
```

Important Postman variable note:

- Use `{{customer_email}}` only when `customer_email` exists as a Postman environment variable.
- Do not write `{{customer01@gmail.com}}`.
- Do not write `{{12345678}}`.

Correct variable-based JSON:

```json
{
  "name": "customer01",
  "email": "{{customer_email}}",
  "password": "{{password}}"
}
```

Correct literal JSON:

```json
{
  "name": "customer01",
  "email": "customer01@gmail.com",
  "password": "12345678"
}
```

Before testing auth, run:

```text
System -> Health
```

Request details:

```text
Method: GET
URL:    {{base_url}}/health
```

Expected result:

- Status `200 OK`.
- Response contains `status: ok`.

If this returns `404`, check that the URL is using backend port `8080`, not frontend port `5173`.

## 4. Register Customer In Postman

Run this request:

```text
Auth -> Register Customer
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/auth/register
```

Headers:

```text
Content-Type: application/json
```

Body tab:

```text
raw -> JSON
```

Body:

```json
{
  "name": "Customer User",
  "email": "{{customer_email}}",
  "password": "{{password}}"
}
```

Expected result:

- Status `201 Created`, or `409 Conflict` if the user already exists.
- If status is `201`, Postman stores `customer_token` automatically.
- Response includes `data.user.role = customer`.
- Response includes `data.token`.
- Response must not include `password_hash`.

If you get `404 page not found`:

- Confirm method is `POST`.
- Confirm URL is `{{base_url}}/api/v1/auth/register`.
- Confirm `base_url` is `http://localhost:8080`.
- Confirm the backend API is running.

Then run:

```text
Auth -> Login Customer
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/auth/login
```

Headers:

```text
Content-Type: application/json
```

Body:

```json
{
  "email": "{{customer_email}}",
  "password": "{{password}}"
}
```

Expected result:

- Status `200 OK`.
- Postman stores `customer_token`.
- Response includes `data.user.email`.
- Response includes `data.token`.

Then run:

```text
Auth -> Get Me
```

Request details:

```text
Method: GET
URL:    {{base_url}}/api/v1/me
```

Headers:

```text
Authorization: Bearer {{customer_token}}
```

Expected result:

- Status `200 OK`.
- Response user email matches `customer_email`.

## 5. Verify Customer Cannot Create Product

Run:

```text
Admin Products -> Customer Cannot Create Product
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/admin/products
```

Headers:

```text
Content-Type: application/json
Authorization: Bearer {{customer_token}}
```

Body:

```json
{
  "name": "USB-C Docking Station",
  "description": "Multi-port dock for laptops.",
  "brand": "DockPro",
  "category": "Accessories",
  "price_cents": 15900,
  "stock_quantity": 15,
  "image_url": "https://example.com/dock.jpg"
}
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

In the selected Postman environment, set `admin_email` to:

```text
admin01@gmail.com
```

Option A: temporarily change the body of:

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

Run the request.

Option B: duplicate `Auth -> Register Customer`, rename it `Register Admin`, and use the admin body above.

Then promote that user in PostgreSQL:

```powershell
docker exec ai_ecommerce_postgres psql -U ecommerce_user -d ecommerce_db -c "UPDATE users SET role = 'admin' WHERE email = 'admin01@gmail.com';"
```

Expected result:

```text
UPDATE 1
```

Then run:

```text
Auth -> Login Admin
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/auth/login
```

Headers:

```text
Content-Type: application/json
```

Body:

```json
{
  "email": "{{admin_email}}",
  "password": "{{password}}"
}
```

Expected result:

- Status `200 OK`.
- Postman stores `admin_token`.
- Response includes `data.user.role = admin`.

If response still says `customer`, run the SQL promotion command again and then login again.

## 7. Create Product As Admin

Run:

```text
Admin Products -> Create Product
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/admin/products
```

Headers:

```text
Content-Type: application/json
Authorization: Bearer {{admin_token}}
```

Body:

```json
{
  "name": "USB-C Docking Station",
  "description": "Multi-port dock for laptops.",
  "brand": "DockPro",
  "category": "Accessories",
  "price_cents": 15900,
  "stock_quantity": 15,
  "image_url": "https://example.com/dock.jpg"
}
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

Request details:

```text
Method: PUT
URL:    {{base_url}}/api/v1/admin/products/{{product_id}}
```

Headers:

```text
Content-Type: application/json
Authorization: Bearer {{admin_token}}
```

Body:

```json
{
  "name": "USB-C Docking Station Pro",
  "description": "Updated multi-port dock for laptops.",
  "brand": "DockPro",
  "category": "Accessories",
  "price_cents": 17900,
  "stock_quantity": 20,
  "image_url": "https://example.com/dock-pro.jpg"
}
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

Request details:

```text
Method: DELETE
URL:    {{base_url}}/api/v1/admin/products/{{product_id}}
```

Headers:

```text
Authorization: Bearer {{admin_token}}
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

Request details:

```text
Method: GET
URL:    {{base_url}}/api/v1/products/{{product_id}}
```

Expected result:

- Status `404 Not Found`.

Run:

```text
Products -> List Products
```

Request details:

```text
Method: GET
URL:    {{base_url}}/api/v1/products
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
