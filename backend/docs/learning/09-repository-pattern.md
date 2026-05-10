# Repository Pattern

The repository pattern keeps database code away from HTTP handlers.

In Phase 2, product requests follow this flow:

```text
HTTP request
-> Gin route
-> ProductHandler
-> ProductService
-> ProductRepository
-> PostgreSQL
```

Files involved:

- `internal/handlers/product_handler.go`
- `internal/services/product_service.go`
- `internal/repositories/product_repository.go`
- `internal/models/product.go`

Why this helps:

- Handlers focus on HTTP concerns: query params, route params, status codes, JSON.
- Services hold business workflow. Phase 2 is simple, but future phases will add richer rules.
- Repositories handle SQL queries and database scanning.
- Tests become easier because handlers can use a fake service.

In interviews, avoid saying that every project needs many layers. The key judgment is whether the separation helps the code stay understandable as features grow.

CV talking point:

> Implemented a layered backend structure with handlers, services, repositories, and models to separate HTTP, business, and database responsibilities.
