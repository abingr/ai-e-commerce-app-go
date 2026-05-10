# Health And Readiness Checks

Health and readiness endpoints are common in production systems.

## Health

`GET /health` answers the question:

> Is the HTTP server alive?

It returns:

```json
{
  "status": "ok",
  "service": "ai-e-commerce-api"
}
```

## Readiness

`GET /ready` answers the question:

> Can the application serve real traffic?

In Phase 1, readiness checks whether PostgreSQL is reachable.

This distinction matters because an API process can be running while its database connection is broken.

CV talking point:

> Added health and readiness checks to support production monitoring and container orchestration patterns.
