# React API Integration

Phase 8 connects the React frontend to the backend APIs built in Phases 1 to 7.

The main frontend file is:

```text
frontend/src/App.jsx
```

## What the frontend now does

The UI can:

- check API health
- list products
- filter products by category
- search products
- register or login a customer
- store the JWT token in `localStorage`
- send the JWT as a bearer token
- add products to the cart
- update cart quantities
- remove cart items
- checkout and create an order
- list recent orders

## API helper pattern

`apiRequest` is the central helper in `frontend/src/App.jsx`.

It handles:

- building the full API URL
- adding JSON headers
- parsing JSON responses
- converting backend error responses into readable UI messages
- showing validation field errors from Phase 7

This keeps individual UI actions smaller. For example, `addToCart`, `checkout`, and `fetchOrders` can focus on the workflow instead of repeating fetch error handling.

## JWT flow

When login or register succeeds:

1. The backend returns `data.token`.
2. React stores it in `localStorage` as `auth_token`.
3. Protected requests include:

```text
Authorization: Bearer <token>
```

The backend middleware reads this header in:

```text
backend/internal/routes/auth_middleware.go
```

## Backend endpoints used

The frontend calls:

- `GET /health`
- `GET /api/v1/products`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/me`
- `GET /api/v1/cart`
- `POST /api/v1/cart/items`
- `PATCH /api/v1/cart/items/:product_id`
- `DELETE /api/v1/cart/items/:product_id`
- `POST /api/v1/orders`
- `GET /api/v1/orders`

## Backend lesson

Frontend integration is where API design quality becomes obvious. Consistent response shapes, readable errors, clear authentication rules, and predictable endpoint names make the client much easier to build.
