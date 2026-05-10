# JWT And Protected Routes

JWT means JSON Web Token.

In Phase 3, login and registration return a signed JWT. The frontend stores the token and sends it on protected requests:

```http
Authorization: Bearer <token>
```

The token contains claims:

- `user_id`
- `email`
- `role`
- issuer
- issued time
- expiry time

Protected route flow:

```text
GET /api/v1/me
-> requireAuth middleware
-> parse and validate JWT
-> set user identity in Gin context
-> AuthHandler.Me
-> AuthService.GetUser
-> UserRepository.FindByID
```

Why JWT is useful:

- The API can authenticate requests without storing session data in memory.
- It works well for stateless REST APIs.
- It is easy for frontend clients to send in an `Authorization` header.

Tradeoffs:

- A stolen token can be used until it expires.
- Logout is client-side unless you add a token denylist or server-side session store.
- The signing secret must be protected.

Backend engineer note:

- Use HTTPS in production.
- Keep JWT expiry reasonably short.
- Store JWT secrets in environment variables, not source code.
- Validate the signing method and issuer.

CV talking point:

> Added JWT middleware to protect authenticated endpoints and propagate user identity through request context.
