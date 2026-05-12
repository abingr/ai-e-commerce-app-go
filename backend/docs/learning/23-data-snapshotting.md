# Data Snapshotting

Snapshotting means copying important data at the time an event happens.

In Phase 6, order items store:

- `product_id`
- `product_name`
- `unit_price_cents`
- `quantity`
- `line_total_cents`

Even though the product already exists in the `products` table, the order copies the product name and price into `order_items`.

## Why this matters

Imagine a customer buys a keyboard for 12900 cents. Tomorrow, an admin changes that product price to 14900 cents.

The old order should still show the price paid at checkout time. It should not change just because the product catalog changed.

That is why `backend/internal/repositories/order_repository.go` inserts `product_name` and `unit_price_cents` into `order_items`.

## Product ID is still stored

The order keeps `product_id` so the system can still trace the order item back to the product catalog.

But the order display should use the snapshot fields for historical accuracy:

- use `product_name` from `order_items`
- use `unit_price_cents` from `order_items`

## Backend lesson

Do not always normalize everything away. In business systems, some duplication is intentional because history matters.
