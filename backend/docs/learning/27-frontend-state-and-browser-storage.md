# Frontend State and Browser Storage

Phase 8 uses React state and browser storage to make the frontend feel like a small working app.

## React state

The frontend stores temporary UI data with `useState`.

Important state values in `frontend/src/App.jsx`:

- `health`
- `products`
- `currentUser`
- `cart`
- `orders`
- `actionStatus`
- `authForm`

When a backend request succeeds, the related state is updated. React then re-renders the screen.

## Browser storage

The JWT token is stored in:

```text
localStorage.auth_token
```

This lets the browser remember the login after refresh.

For a learning project this is acceptable, but production systems must think carefully about token storage, XSS risk, refresh tokens, session expiry, and logout behavior.

## Why cart and orders are fetched after login

After login, React calls:

- `fetchCart`
- `fetchOrders`

This mirrors how real apps work. Once the user identity is known, the frontend loads user-scoped resources.

The backend still enforces ownership. The frontend does not send `user_id`; it only sends the JWT.

## Error display

Phase 7 added structured errors. Phase 8 uses those errors in the UI.

If the backend returns validation details like:

```json
{
  "fields": [
    {
      "field": "email",
      "rule": "email"
    }
  ]
}
```

the frontend turns that into a readable status message.

## Backend lesson

The frontend should improve usability, but the backend must remain the source of truth. The backend still validates input, checks JWTs, scopes data by user, and creates orders inside transactions.
