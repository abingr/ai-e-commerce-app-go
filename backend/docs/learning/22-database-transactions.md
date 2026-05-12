# Database Transactions

A database transaction groups multiple SQL statements into one all-or-nothing operation.

Phase 6 needs a transaction because checkout does several related writes:

1. Read the authenticated user's cart items.
2. Create an `orders` row.
3. Create one `order_items` row per cart item.
4. Clear the user's cart.

If step 3 fails after the order is created, we do not want a half-created order. If step 4 fails, we do not want the user to accidentally checkout the same cart again. A transaction protects this workflow.

## Where it is implemented

The transaction lives in `backend/internal/repositories/order_repository.go` inside `CreateFromCart`.

Important calls:

- `r.db.Begin(ctx)` starts the transaction.
- `defer tx.Rollback(ctx)` makes rollback the default if something fails.
- SQL inserts create the order and order items.
- The cart is deleted only after order items are created.
- `tx.Commit(ctx)` makes the changes permanent.

## Why rollback is deferred

The deferred rollback is a defensive pattern. If the function returns early because of an error, rollback runs automatically.

If commit succeeds, the later rollback call has no effect because the transaction is already closed.

## What transaction safety means here

After `POST /api/v1/orders`:

- success means the order exists and the cart is empty
- failure means the checkout changes should not be partially saved

This is one of the most important backend ideas in e-commerce systems.
