# User-Scoped Resources

A user-scoped resource belongs to one authenticated user.

The cart is user-scoped. Every cart query includes `user_id`, which comes from the JWT middleware.

This protects users from seeing or changing another user's cart.

Examples:

```sql
WHERE user_id = $1
```

Routes:

```http
GET    /api/v1/cart
POST   /api/v1/cart/items
PATCH  /api/v1/cart/items/:product_id
DELETE /api/v1/cart/items/:product_id
DELETE /api/v1/cart/items
```

Security behavior:

- Missing token returns `401`.
- Invalid token returns `401`.
- Valid token can only access that user's cart.

Backend engineer note:

Never accept `user_id` from the request body for user-owned data when the authenticated identity already tells you who the user is.

CV talking point:

> Implemented user-scoped cart operations using JWT-derived identity instead of trusting user IDs from client input.
