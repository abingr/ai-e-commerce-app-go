# Admin Product Management

Phase 4 adds admin-only product write APIs:

```http
POST   /api/v1/admin/products
PUT    /api/v1/admin/products/:id
DELETE /api/v1/admin/products/:id
```

Public users can still browse:

```http
GET /api/v1/products
GET /api/v1/products/:id
```

Why admin endpoints are separate:

- The URL makes the security boundary clear.
- Public and admin workflows can evolve independently.
- Middleware can protect the whole `/admin` route group.

Backend flow:

```text
ProductHandler
-> ProductService
-> ProductRepository
-> PostgreSQL
```

This is the same layered pattern used for reads, now extended to writes.

CV talking point:

> Implemented admin product CRUD operations using a layered handler-service-repository design.
