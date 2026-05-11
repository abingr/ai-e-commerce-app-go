# Phase 5 To-Do-Test

Use this checklist after Phase 5 changes to verify the authenticated shopping cart.

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
- Migration `000004_create_cart_items_table.up.sql` is applied.
- API listens on `http://localhost:8080`.

## 2. Run Automated Tests

From the backend folder:

```powershell
go test ./...
```

Expected result:

- Cart handler tests pass.
- Cart service tests pass.
- Existing auth, product, and RBAC tests still pass.

Code path involved:

- `backend/internal/handlers/cart_handler_test.go` tests HTTP cart behavior.
- `backend/internal/services/cart_service_test.go` tests service behavior.
- `backend/internal/routes/router.go:52` protects the cart route group.

## 3. Import Or Refresh Postman Files

Open Postman and import the latest files:

```text
postman/ai-e-commerce-app-go.postman_collection.json
postman/ai-e-commerce-local.postman_environment.json
```

Select the environment:

```text
AI E-Commerce Local
```

Confirm these environment values:

```text
base_url       = http://localhost:8080
customer_email = customer01@gmail.com
password       = 12345678
```

New Phase 5 variable:

```text
cart_product_id
```

`Products -> List Products` saves the first public product ID into `cart_product_id`.

## 4. Login Customer In Postman

Run:

```text
Auth -> Login Customer
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/auth/login
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

If customer does not exist yet, run:

```text
Auth -> Register Customer
```

then run login again.

## 5. Choose A Product For Cart Testing

Run:

```text
Products -> List Products
```

Expected result:

- Status `200 OK`.
- Response contains products.
- Postman stores the first product ID in `cart_product_id`.

You may also manually copy any active product ID into:

```text
cart_product_id
```

## 6. Get Empty Or Existing Cart

Run:

```text
Cart -> Get Cart
```

Request details:

```text
Method: GET
URL:    {{base_url}}/api/v1/cart
```

Headers:

```text
Authorization: Bearer {{customer_token}}
```

Expected result:

- Status `200 OK`.
- Response contains `data.items`.
- Response contains `data.total_cents`.

Code path involved:

- `backend/internal/routes/router.go:53` maps `GET /cart`.
- `backend/internal/handlers/cart_handler.go:31` handles the request.
- `backend/internal/repositories/cart_repository.go:19` joins cart items with products.

## 7. Add Item To Cart

Run:

```text
Cart -> Add Item To Cart
```

Request details:

```text
Method: POST
URL:    {{base_url}}/api/v1/cart/items
```

Headers:

```text
Content-Type: application/json
Authorization: Bearer {{customer_token}}
```

Body:

```json
{
  "product_id": "{{cart_product_id}}",
  "quantity": 2
}
```

Expected result:

- Status `200 OK`.
- Cart contains the product.
- Item quantity is `2`.
- `line_total_cents` equals `unit_price_cents * quantity`.

Code path involved:

- `backend/internal/routes/router.go:54` maps add item.
- `backend/internal/handlers/cart_handler.go:50` handles add item.
- `backend/internal/services/cart_service.go:29` adds item and returns fresh cart.
- `backend/internal/repositories/cart_repository.go:70` performs the upsert.

Repeat the same request.

Expected result:

- The same product row remains.
- Quantity increases by `2`.
- This proves the upsert works.

## 8. Update Cart Item Quantity

Run:

```text
Cart -> Update Cart Item Quantity
```

Request details:

```text
Method: PATCH
URL:    {{base_url}}/api/v1/cart/items/{{cart_product_id}}
```

Headers:

```text
Content-Type: application/json
Authorization: Bearer {{customer_token}}
```

Body:

```json
{
  "quantity": 3
}
```

Expected result:

- Status `200 OK`.
- Item quantity is exactly `3`.

Code path involved:

- `backend/internal/routes/router.go:55` maps update item.
- `backend/internal/handlers/cart_handler.go:84` handles update item.
- `backend/internal/services/cart_service.go:37` passes exact quantity.
- `backend/internal/repositories/cart_repository.go:92` updates the cart row.

## 9. Remove Cart Item

Run:

```text
Cart -> Remove Cart Item
```

Request details:

```text
Method: DELETE
URL:    {{base_url}}/api/v1/cart/items/{{cart_product_id}}
```

Headers:

```text
Authorization: Bearer {{customer_token}}
```

Expected result:

- Status `200 OK`.
- Product is removed from `data.items`.

Code path involved:

- `backend/internal/routes/router.go:56` maps remove item.
- `backend/internal/handlers/cart_handler.go:126` handles remove item.
- `backend/internal/repositories/cart_repository.go:114` deletes the row.

## 10. Clear Cart

First add one or more items again.

Then run:

```text
Cart -> Clear Cart
```

Request details:

```text
Method: DELETE
URL:    {{base_url}}/api/v1/cart/items
```

Headers:

```text
Authorization: Bearer {{customer_token}}
```

Expected result:

- Status `204 No Content`.

Then run:

```text
Cart -> Get Cart
```

Expected result:

- Status `200 OK`.
- `data.items` is empty.
- `data.total_cents` is `0`.

Code path involved:

- `backend/internal/routes/router.go:57` maps clear cart.
- `backend/internal/handlers/cart_handler.go:160` handles clear cart.
- `backend/internal/repositories/cart_repository.go:130` deletes all rows for the user.

## 11. Verify Error Cases

Missing token:

```text
GET {{base_url}}/api/v1/cart
```

without the `Authorization` header.

Expected result:

- Status `401 Unauthorized`.

Invalid quantity:

```json
{
  "product_id": "{{cart_product_id}}",
  "quantity": 0
}
```

Expected result:

- Status `400 Bad Request`.

Missing product:

```json
{
  "product_id": "00000000-0000-0000-0000-000000000000",
  "quantity": 1
}
```

Expected result:

- Status `404 Not Found`.

## 12. Optional CLI Smoke Test

Login customer:

```powershell
$body = @{ email = "customer01@gmail.com"; password = "12345678" } | ConvertTo-Json
$login = Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/login -Method Post -ContentType "application/json" -Body $body
$token = $login.data.token
```

Get product ID:

```powershell
$products = Invoke-RestMethod -Uri http://localhost:8080/api/v1/products
$productID = $products.data[0].id
```

Add item:

```powershell
$body = @{ product_id = $productID; quantity = 2 } | ConvertTo-Json
Invoke-RestMethod -Uri http://localhost:8080/api/v1/cart/items -Method Post -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" } -Body $body
```

Get cart:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/cart -Headers @{ Authorization = "Bearer $token" }
```

## 13. Build And Audit Frontend

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

## Phase 5 Pass Criteria

Phase 5 is working when:

- `go run ./cmd/migrate up` creates or confirms `cart_items`.
- `go test ./...` passes.
- Postman collection imports successfully.
- Customer can get cart.
- Customer can add a product to cart.
- Repeated add increases quantity.
- Customer can update item quantity.
- Customer can remove one item.
- Customer can clear the cart.
- Missing token returns `401`.
- Invalid quantity returns `400`.
- Missing product returns `404`.
- `npm run build` passes.
- `npm audit` reports zero vulnerabilities.
