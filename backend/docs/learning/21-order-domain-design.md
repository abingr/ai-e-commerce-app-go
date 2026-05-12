# Order Domain Design

Phase 6 adds checkout without payment. In a real e-commerce backend, checkout is the point where temporary cart data becomes a permanent business record.

## What an order represents

An order is a durable record that says:

- which user placed the order
- what products were bought
- what price each product had at checkout time
- what quantity was bought
- what total amount was calculated
- what status the order is currently in

The code represents this in `backend/internal/models/order.go`:

- `Order` is the parent record.
- `OrderItem` is the child record for each product in the order.
- `OrderStatusConfirmed` is the status used in Phase 6.

## Why orders are separate from carts

A cart is editable. A user can add, update, remove, or clear items.

An order is historical. Once created, it should not depend on the current cart or current product price. This is why Phase 6 creates `orders` and `order_items` tables instead of reusing `cart_items`.

## API design

Phase 6 adds these authenticated routes in `backend/internal/routes/router.go`:

- `POST /api/v1/orders`
- `GET /api/v1/orders`
- `GET /api/v1/orders/:id`

The user identity comes from the JWT middleware. The order handler does not accept `user_id` from the request body because clients should not be trusted to choose which user's orders they are reading or creating.

## Senior backend pattern

Notice the same layer pattern from earlier phases:

- Handler: HTTP status codes and request/response shape.
- Service: business use case boundary.
- Repository: database queries and transactions.
- Model: shared domain data structures.

This keeps checkout understandable even as it becomes more complicated later with stock reservation, payment, shipping, cancellation, and refunds.
