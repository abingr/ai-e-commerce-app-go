# Soft Delete Pattern

Phase 4 uses soft delete for products.

Instead of removing a row from the database, delete does this:

```sql
UPDATE products
SET is_active = false
WHERE id = $1
```

Public product queries only return:

```sql
WHERE is_active = true
```

Why backend teams use soft delete:

- Keeps historical records.
- Avoids breaking references from orders, carts, or analytics.
- Makes accidental deletion easier to recover.

Tradeoffs:

- Queries must remember to filter inactive records.
- Unique constraints may need extra care.
- Old data can grow over time and may need archival.

CV talking point:

> Used soft deletes for product removal to preserve historical product data while hiding inactive items from public catalog APIs.
