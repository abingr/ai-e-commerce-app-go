# Phase 8 To-Do Test: React Frontend Integration

This checklist verifies that the React frontend can use the backend features built so far.

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
- all SQL files in `backend/migrations`

## 3. Start the backend API

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go run ./cmd/api
```

Keep this terminal open.

Verify in browser:

```text
http://localhost:8080/health
```

Expected:

```json
{
  "service": "ai-e-commerce-api",
  "status": "ok"
}
```

## 4. Start the frontend

Open a second terminal:

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm install
npm run dev
```

Open:

```text
http://localhost:5173
```

Code being run:

- `frontend/src/App.jsx`
- `frontend/src/styles.css`

## 5. Verify API health in the UI

In the page header, verify:

```text
API ok
```

Code being called:

- `checkHealth` in `frontend/src/App.jsx`
- `HealthHandler.Health` in `backend/internal/handlers/health.go`

## 6. Verify product catalog

In the browser, verify:

- product cards are visible
- product image, name, brand, category, price, and stock are shown
- search input filters products
- category dropdown filters products

Code being called:

- `fetchProducts` in `frontend/src/App.jsx`
- `ProductHandler.List` in `backend/internal/handlers/product_handler.go`
- `ProductRepository.List` in `backend/internal/repositories/product_repository.go`

## 7. Register or login

Use the account panel.

For a new user, switch to:

```text
Register
```

Example:

```text
Name: Phase Eight Customer
Email: phase8@example.com
Password: password123
```

If the email already exists, switch to:

```text
Login
```

Verify:

- Auth status changes to signed in
- role shows `customer`
- browser DevTools Application tab shows `localStorage.auth_token`

Code being called:

- `submitAuth` in `frontend/src/App.jsx`
- `AuthHandler.Register` or `AuthHandler.Login`
- `AuthService` creates or validates the JWT

## 8. Add product to cart

Click:

```text
Add
```

on any product card.

Verify:

- cart total increases
- cart line item appears
- action status says product was added

Code being called:

- `addToCart` in `frontend/src/App.jsx`
- `CartHandler.AddItem` in `backend/internal/handlers/cart_handler.go`
- `CartRepository.AddItem` in `backend/internal/repositories/cart_repository.go`

## 9. Update cart quantity

Use:

```text
+
-
```

buttons in the cart.

Verify:

- quantity changes
- line total changes
- cart total changes

Code being called:

- `updateCartItem` in `frontend/src/App.jsx`
- `CartHandler.UpdateItem`
- `CartRepository.UpdateItem`

## 10. Remove cart item

Click:

```text
Remove
```

Verify:

- item disappears
- cart total updates

Code being called:

- `removeCartItem` in `frontend/src/App.jsx`
- `CartHandler.RemoveItem`
- `CartRepository.RemoveItem`

## 11. Checkout

Add at least one product again, then click:

```text
Checkout
```

Verify:

- cart becomes empty
- orders count increases
- order appears in the Orders section

Code being called:

- `checkout` in `frontend/src/App.jsx`
- `OrderHandler.CreateFromCart`
- `OrderRepository.CreateFromCart`

Important backend behavior:

- order is created
- order items snapshot product name and price
- cart is cleared
- database transaction commits all changes together

## 12. Refresh browser

Refresh:

```text
http://localhost:5173
```

Verify:

- user remains signed in
- cart and orders reload from backend

Code being called:

- `fetchCurrentUser`
- `fetchCart`
- `fetchOrders`

## 13. Sign out

Click:

```text
Sign out
```

Verify:

- user is cleared
- cart is cleared from the UI
- orders are cleared from the UI
- `localStorage.auth_token` is removed

## 14. Run automated checks

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\backend
go test ./...
```

```powershell
cd c:\Training\Golang\ai-e-commerce-app-go\frontend
npm run build
npm audit
```

Expected result:

- backend tests pass
- frontend build passes
- audit reports no vulnerabilities
