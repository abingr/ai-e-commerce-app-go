# Phase 6 To-Do Test: Orders and Checkout

This checklist verifies checkout without payment. Use Postman for HTTP requests and the browser for quick API visibility.

## 1. Start PostgreSQL

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go
docker compose up -d postgres
```

Verify:

```powershell
docker ps
```

You should see the `ai_ecommerce_postgres` container running.

## 2. Run migrations

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/migrate up
```

Code being run:

- `backend/cmd/migrate/main.go`
- `backend/migrations/000005_create_orders_tables.up.sql`

Verify that the command prints successful migration output and no error.

## 3. Run automated tests

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

Code being tested:

- `backend/internal/handlers/order_handler.go`
- `backend/internal/handlers/order_handler_test.go`
- `backend/internal/services/order_service.go`
- `backend/internal/services/order_service_test.go`

Expected result:

```text
ok
```

for all backend packages.

## 4. Start the backend API

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/api
```

Keep this terminal open.

Code being run:

- `backend/cmd/api/main.go`
- `backend/internal/routes/router.go`
- `backend/internal/handlers/order_handler.go`
- `backend/internal/repositories/order_repository.go`

## 5. Import Postman files

Import:

```text
postman/ai-e-commerce-app-go.postman_collection.json
postman/ai-e-commerce-local.postman_environment.json
```

Select the `AI E-Commerce Local` environment.

## 6. Login or register a customer in Postman

Run:

```text
Auth > Register Customer
```

If the email already exists, run:

```text
Auth > Login Customer
```

Verify:

- status is `201 Created` for register, or `200 OK` for login
- environment variable `customer_token` is populated

Code being called:

- `backend/internal/handlers/auth_handler.go`
- `backend/internal/services/auth_service.go`
- `backend/internal/routes/auth_middleware.go` for later protected requests

## 7. Select a product for the cart

Run:

```text
Products > List Products
```

Verify:

- status is `200 OK`
- environment variable `cart_product_id` is populated by the Postman test script

Code being called:

- `backend/internal/handlers/product_handler.go`
- `backend/internal/repositories/product_repository.go`

## 8. Add an item to the cart

Run:

```text
Cart > Add Item To Cart
```

Body:

```json
{
  "product_id": "{{cart_product_id}}",
  "quantity": 2
}
```

Verify:

- status is `200 OK`
- response `data.items` contains one product
- response `data.total_cents` is greater than `0`

Code being called:

- `CartHandler.AddItem` in `backend/internal/handlers/cart_handler.go`
- `CartService.AddItem` in `backend/internal/services/cart_service.go`
- `CartRepository.AddItem` in `backend/internal/repositories/cart_repository.go`

## 9. Create an order from the cart

Run:

```text
Orders > Create Order From Cart
```

Verify:

- status is `201 Created`
- response `data.status` is `confirmed`
- response `data.items[0].product_name` is present
- response `data.items[0].unit_price_cents` is present
- environment variable `order_id` is populated

Code being called:

- `OrderHandler.CreateFromCart` in `backend/internal/handlers/order_handler.go`
- `OrderService.CreateFromCart` in `backend/internal/services/order_service.go`
- `OrderRepository.CreateFromCart` in `backend/internal/repositories/order_repository.go`

Important database behavior:

- inserts into `orders`
- inserts into `order_items`
- deletes checked-out rows from `cart_items`
- commits all changes together

## 10. Verify the cart was cleared

Run:

```text
Cart > Get Cart
```

Verify:

- status is `200 OK`
- response `data.items` is an empty array
- response `data.total_cents` is `0`

This confirms checkout did not leave old cart items behind.

## 11. List your orders

Run:

```text
Orders > List My Orders
```

Verify:

- status is `200 OK`
- response `data` contains at least one order
- the latest order has the same ID as `{{order_id}}`

Code being called:

- `OrderHandler.List`
- `OrderService.ListByUser`
- `OrderRepository.ListByUser`

## 12. Get one order detail

Run:

```text
Orders > Get Order Detail
```

Verify:

- status is `200 OK`
- response `data.id` matches `{{order_id}}`
- response `data.items` contains snapshot fields

Code being called:

- `OrderHandler.GetByID`
- `OrderService.FindByIDForUser`
- `OrderRepository.FindByIDForUser`

## 13. Verify empty cart checkout fails

Run again:

```text
Orders > Create Order From Cart
```

Expected response:

```json
{
  "error": "cart is empty"
}
```

Expected status:

```text
400 Bad Request
```

Code being called:

- `OrderRepository.CreateFromCart` returns `repositories.ErrEmptyCart`
- `OrderHandler.CreateFromCart` converts it to HTTP `400`

## 14. Verify in browser

Open:

```text
http://localhost:8080/health
```

You should see:

```json
{
  "service": "ai-e-commerce-api",
  "status": "ok"
}
```

Browser note:

Protected order endpoints require a JWT bearer token, so use Postman for `/api/v1/orders`.

## 15. Run frontend checks

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm install
npm run build
npm audit
```

Expected result:

- build completes successfully
- audit reports no critical issue that blocks this learning phase
