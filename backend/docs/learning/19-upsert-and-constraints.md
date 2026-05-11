# Upsert And Constraints

Phase 5 uses an upsert when adding an item to the cart.

The cart table has this constraint:

```sql
UNIQUE (user_id, product_id)
```

That prevents duplicate cart rows for the same user and product.

When a user adds a product that is already in the cart, the repository increases the existing quantity instead of creating a duplicate row.

Important SQL idea:

```sql
ON CONFLICT (user_id, product_id)
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
```

Other useful constraints:

- `quantity > 0`
- `user_id` references `users(id)`
- `product_id` references `products(id)`

Backend engineer note:

Database constraints protect data even if application code has a bug.

CV talking point:

> Used PostgreSQL unique constraints and upsert behavior to prevent duplicate cart items while supporting repeated add-to-cart actions.
