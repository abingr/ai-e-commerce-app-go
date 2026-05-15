# CV Summary

## Project title

AI E-Commerce App Go

## Short CV bullet

Built a production-style mini e-commerce application using Go, Gin, PostgreSQL, JWT authentication, role-based authorization, React, Docker, OpenAPI documentation, Postman collections, and automated tests.

## Expanded CV bullets

- Designed and implemented a layered Go backend with handler, service, repository, and model packages.
- Built REST APIs for product catalog, customer authentication, admin product management, shopping cart, checkout, and order history.
- Implemented JWT authentication, bcrypt password hashing, role-based admin authorization, request ID tracing, structured logging, and consistent API error responses.
- Used PostgreSQL migrations, foreign keys, unique constraints, soft delete, user-scoped resources, transactions, and order item snapshotting.
- Added React frontend integration for product browsing, login/register, cart management, checkout, and order history.
- Created OpenAPI documentation, Postman collection/environment files, phase-by-phase learning notes, and manual test checklists.
- Added automated backend tests, frontend build/audit checks, Dockerfiles, and GitHub Actions CI.

## Technologies

- Go
- Gin
- PostgreSQL
- pgx
- golang-migrate
- JWT
- bcrypt
- React
- Vite
- Docker
- GitHub Actions
- OpenAPI
- Postman

## Interview talking points

- Why the backend is split into handlers, services, repositories, and models.
- Why JWT identity is used instead of accepting `user_id` from request bodies.
- How RBAC protects admin product routes.
- How cart items are converted into orders inside a database transaction.
- Why order items snapshot product name and price.
- How request IDs help connect API responses to backend logs.
- Why fast unit tests run by default while database integration tests are opt-in.
