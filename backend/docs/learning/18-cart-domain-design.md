# Cart Domain Design

Phase 5 adds a shopping cart.

A cart is user-scoped. That means each authenticated user sees and changes only their own cart.

This project uses one table:

```text
cart_items
```

Important columns:

- `user_id`: owner of the cart item.
- `product_id`: product in the cart.
- `quantity`: how many units the user wants.
- `created_at`, `updated_at`: audit timestamps.

There is no separate `carts` table yet. For this project, a user's cart is simply all rows in `cart_items` for that user.

Backend flow:

```text
JWT
-> user_id from auth middleware
-> CartHandler
-> CartService
-> CartRepository
-> PostgreSQL
```

CV talking point:

> Built authenticated shopping cart APIs scoped by user identity, backed by PostgreSQL foreign keys and quantity constraints.
