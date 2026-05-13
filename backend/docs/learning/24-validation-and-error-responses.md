# Validation and Error Responses

Phase 7 improves how the API tells clients what went wrong.

Before this phase, most errors looked like this:

```json
{
  "error": "invalid registration payload"
}
```

That is readable, but it does not tell the frontend or API tester which field failed.

Phase 7 keeps the simple `error` message and adds machine-friendly fields:

```json
{
  "code": "VALIDATION_ERROR",
  "error": "invalid registration payload",
  "fields": [
    {
      "field": "email",
      "rule": "email"
    }
  ],
  "request_id": "..."
}
```

## Where it is implemented

The shared response helpers live in:

```text
backend/internal/handlers/response.go
```

Important functions:

- `JSONError` writes a consistent non-validation error response.
- `JSONValidationError` writes validation errors with field details.
- `RecordError` attaches internal errors to the Gin context so middleware can log them.

## Why keep `error` as a string

Many APIs use an object like:

```json
{
  "error": {
    "message": "..."
  }
}
```

For this learning project, keeping `error` as a string makes Postman testing easier and avoids breaking the earlier phase examples. The `code`, `fields`, and `request_id` fields add the production-style structure.

## Backend lesson

A good API error should be useful to:

- the user, through a readable message
- the frontend, through a stable error code
- the developer, through field-level validation details
- the backend team, through a request ID that can be found in logs
