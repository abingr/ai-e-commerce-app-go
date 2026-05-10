# Migrations And Seed Data

Database migrations are versioned changes to the database schema.

Instead of manually creating tables, backend teams keep migration files in source control so every developer and environment can apply the same changes.

This project uses `golang-migrate`.

Phase 2 migrations:

- `000001_create_products_table.up.sql` creates the `products` table.
- `000001_create_products_table.down.sql` drops the `products` table.
- `000002_seed_products.up.sql` inserts sample electronics products.
- `000002_seed_products.down.sql` removes the seeded products.

Run migrations from the backend folder:

```powershell
go run ./cmd/migrate up
```

Rollback one migration:

```powershell
go run ./cmd/migrate down
```

Backend engineer note:

- `up` files move the database forward.
- `down` files undo a migration when safe.
- Seed data makes local development and manual testing easier.
- Production seed data should be handled carefully because it changes real data.

CV talking point:

> Added versioned PostgreSQL migrations and seed data using golang-migrate to make database setup repeatable.
