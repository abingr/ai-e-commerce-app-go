# Authorization And RBAC

Authentication answers:

> Who are you?

Authorization answers:

> Are you allowed to do this?

Phase 4 adds role-based access control, often called RBAC.

Current roles:

- `customer`: can browse products and use normal customer features.
- `admin`: can manage products.

Request flow for admin endpoints:

```text
request
-> requireAuth middleware
-> JWT parsed
-> user role stored in Gin context
-> requireRole("admin") middleware
-> admin handler runs
```

Important status codes:

- `401 Unauthorized`: the request is not authenticated.
- `403 Forbidden`: the request is authenticated, but the user lacks permission.

Code involved:

- `internal/routes/auth_middleware.go`
- `internal/routes/router.go`

CV talking point:

> Added role-based authorization middleware to protect admin-only product management endpoints.
