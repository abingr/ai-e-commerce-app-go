# Gin Basics

Gin is a lightweight HTTP web framework for Go. It helps with routing, JSON responses, middleware, path parameters, and request binding.

In Phase 1, Gin is used in `internal/routes/router.go`.

Important ideas:

- `gin.New()` creates a router.
- `router.GET("/health", handler.Health)` registers a GET endpoint.
- `c.JSON(statusCode, payload)` writes a JSON response.
- Middleware runs before or after handlers for cross-cutting concerns like logging and authentication.

Why backend engineers use frameworks:

- Less repetitive HTTP boilerplate
- Clear routing
- Middleware support
- Easier request validation and response formatting

Gin is not the business logic layer. It should mainly translate HTTP requests into application actions and translate results back into HTTP responses.
