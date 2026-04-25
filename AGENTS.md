# Go Fiber CRUD - Agent Instructions

## Project Overview

A high-performance Todo application built with Go Fiber v3 and PostgreSQL. Features user authentication, session management, and CSRF protection.

## Tech Stack

- **Language**: Go 1.26.0
- **Framework**: Fiber v3
- **Database**: PostgreSQL with pgx/v5
- **Templates**: html/v2
- **Auth**: bcrypt + session-based

## Architecture

```
├── cmd/web/main.go        # Entry point, route setup, middleware config
├── internal/
│   ├── config/            # Environment config loading
│   ├── database/          # DB connection (pgx)
│   ├── handlers/          # HTTP handlers (auth.go, todo.go)
│   ├── models/            # Data models (User, Todo)
│   └── repository/        # Data access layer
└── views/                 # HTML templates
```

## Key Patterns

- **Auth**: Session stored user_id; auth middleware redirects to /login if missing
- **CSRF**: All POST forms include `_csrf` token; extracted from form field
- **User Model Note**: User model only has `ID` and `Password` fields - no `Username` field in model (username is used directly in queries)

## Common Commands

```bash
# Run the application
go run cmd/web/main.go

# Install dependencies
go mod tidy

# Build
go build -o app ./cmd/web
```

## Important Notes

- Database requires `users` and `todos` tables (see README.md for schema)
- Session secret configured via SESSION_SECRET env var
- CSRF token passed to all templates via `{{.CSRFToken}}`

## Conventions

- Handlers receive `*session.Store` for session operations
- Repository methods handle their own SQL errors
- Models are simple structs (no business logic)
- Passwords hashed with bcrypt before storage

## Testing

```bash
# Run unit tests
go test ./...

# Run tests with verbose output
go test -v ./...
```

## Linting & Type Checking

```bash
# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```