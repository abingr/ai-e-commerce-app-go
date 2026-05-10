# Project Structure

Production Go services usually separate application startup, business code, infrastructure code, and documentation.

In this project:

- `cmd/api` contains the executable entry point.
- `internal/config` loads environment configuration.
- `internal/database` owns PostgreSQL connection setup.
- `internal/handlers` contains HTTP handlers.
- `internal/routes` wires handlers and middleware into Gin routes.
- `docs/learning` contains explanations for each phase.

The `internal` directory is important in Go. Code inside it cannot be imported by other Go modules, which helps keep application internals private.

CV talking point:

> Structured the backend using Go's `internal` package pattern to separate configuration, database access, HTTP handlers, and routing.
