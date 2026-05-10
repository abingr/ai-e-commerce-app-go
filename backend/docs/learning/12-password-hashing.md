# Password Hashing

Passwords must never be stored as plain text.

Phase 3 uses `bcrypt` to hash passwords before saving users.

Important code:

```text
internal/services/auth_service.go
```

Registration uses:

```go
bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
```

Login uses:

```go
bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
```

Why hashing matters:

- If the database leaks, attackers should not immediately see user passwords.
- Bcrypt is intentionally slow, which makes brute-force attacks harder.
- Each bcrypt hash includes a salt, so identical passwords do not produce identical hashes.

Backend engineer note:

- Never log passwords.
- Never return password hashes in JSON.
- Enforce a minimum password length.
- Use trusted libraries instead of writing password hashing code yourself.

CV talking point:

> Secured user credentials using bcrypt password hashing and excluded password hashes from API responses.
