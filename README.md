# Task Tracker Server (Go)

## Quick start

1. Start Postgres (adjust DSN in config if needed).

2. Run server:
```powershell
go run ./cmd/server
```

## Structure
```
/cmd/server
/internal/api
/internal/auth
/internal/repo
/internal/service
/internal/domain
/internal/storage
/internal/notify
/migrations
/openapi
```

## OpenAPI
File: `openapi/openapi.yaml`

## Migrations
Use your preferred tool (e.g., golang-migrate). The migration files live under `migrations/`.
