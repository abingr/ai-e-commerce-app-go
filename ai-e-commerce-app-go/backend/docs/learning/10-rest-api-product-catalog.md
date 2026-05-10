# REST API Product Catalog

Phase 2 adds product browsing endpoints.

Endpoints:

```http
GET /api/v1/products
GET /api/v1/products/:id
```

`GET /api/v1/products` supports optional query parameters:

```text
category=Accessories
search=keyboard
```

Examples:

```text
http://localhost:8080/api/v1/products
http://localhost:8080/api/v1/products?category=Accessories
http://localhost:8080/api/v1/products?search=phone
```

Response shape:

```json
{
  "data": [
    {
      "id": "uuid",
      "name": "KeyForge Mechanical Keyboard",
      "brand": "KeyForge",
      "category": "Accessories",
      "price_cents": 12900,
      "stock_quantity": 40
    }
  ]
}
```

Why the response uses a `data` wrapper:

- It gives us room to add pagination metadata later.
- Example future response: `{ "data": [...], "meta": { "page": 1 } }`.

Status codes:

- `200 OK`: product list or product detail found.
- `400 Bad Request`: invalid product UUID.
- `404 Not Found`: valid UUID format, but no active product exists.
- `500 Internal Server Error`: unexpected backend/database failure.

CV talking point:

> Built RESTful product catalog APIs with filtering, UUID route validation, consistent JSON responses, and OpenAPI documentation.
