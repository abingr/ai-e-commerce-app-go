# Database Schema Design

Phase 2 introduces the `products` table.

Good schema design starts with the business concept. In this project, a product is an electronics item that customers can browse and later add to cart.

Important columns:

- `id`: UUID primary key. UUIDs are safe to expose in URLs because they are hard to guess.
- `name`, `description`, `brand`, `category`: product catalog information.
- `price_cents`: stores money as an integer number of cents.
- `stock_quantity`: current available stock.
- `image_url`: product image shown by the frontend.
- `is_active`: lets the business hide a product without deleting historical data.
- `created_at`, `updated_at`: audit timestamps.

Why `price_cents` instead of `price` as a floating-point number:

Floating-point numbers can introduce rounding issues. Backend systems usually store money as integers in the smallest currency unit, such as cents.

Schema file:

```text
backend/migrations/000001_create_products_table.up.sql
```

CV talking point:

> Designed a normalized product catalog schema using UUID primary keys, integer money storage, stock tracking, active/inactive flags, and indexed query fields.
