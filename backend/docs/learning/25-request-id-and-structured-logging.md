# Request ID and Structured Logging

Phase 7 adds request IDs and better structured logging.

## What is a request ID?

A request ID is a unique value attached to one HTTP request. It helps connect:

- the HTTP response returned to the client
- the backend log line written by the server
- the specific error that happened during that request

If a client sends an `X-Request-ID` header, the API uses it. If not, the API creates one.

## Where it is implemented

Request ID middleware lives in:

```text
backend/internal/routes/middleware.go
```

Important functions:

- `requestID` creates or accepts an `X-Request-ID`
- `requestLogger` logs method, path, status, latency, client IP, request ID, and internal error details

The middleware is registered in:

```text
backend/internal/routes/router.go
```

## Why structured logs matter

Structured logs use named fields instead of plain text. This makes logs easier to search later.

Example fields:

- `request_id`
- `method`
- `path`
- `status`
- `latency_ms`
- `client_ip`
- `error`

In production, these fields are usually sent to tools like CloudWatch, Datadog, Grafana Loki, or Elasticsearch.

## Backend lesson

Do not expose raw internal errors to users. Instead:

- return a safe message to the client
- record the real error on the Gin context
- let middleware include that internal error in backend logs
