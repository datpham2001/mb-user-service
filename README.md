# mb-user-service

API Gateway service for **Motorbok** — a motorbike hiring platform. This service is the single entry point for all client requests, handling authentication, request routing, rate limiting, and cross-cutting concerns (CORS, tracing, logging) before forwarding to downstream microservices.

## Responsibilities

- **Authentication & Authorization** — JWT-based access/refresh token issuance and validation; Google OAuth2 login
- **Rate limiting** — per-client request throttling via Redis
- **Observability** — structured logging (Logrus), distributed tracing (OpenTelemetry), and request-ID propagation
- **Health checks** — liveness probes for Postgres and Redis dependencies

## Tech Stack

| Concern | Library |
|---|---|
| HTTP framework | Gin |
| ORM | GORM + pgx (Postgres) |
| Cache / Rate limit | go-redis v9 |
| Auth | golang-jwt/jwt v5, Google OAuth2 |
| Config | Viper (YAML) |
| Migrations | Goose |
| Tracing | OpenTelemetry (stdout exporter; swap in prod) |
| Validation | go-playground/validator v10 |

## Architecture

Clean Architecture with four layers:

```
internal/
  domain/           # Entities, repository interfaces, domain errors
  application/      # Use cases (business logic), DTOs
  infrastructure/   # DB, Redis, JWT, OAuth, logging, tracing, config
  presentation/     # HTTP: Gin server, controllers, middlewares
```

Dependency direction: `presentation → application → domain ← infrastructure`

## API Endpoints

| Method | Path | Auth required | Description |
|---|---|---|---|
| POST | `/api/v1/auth/register` | No | Register a new account |
| POST | `/api/v1/auth/login` | No | Login with email/password |
| POST | `/api/v1/auth/token/refresh` | No | Refresh access token |
| POST | `/api/v1/auth/oauth/google/callback` | No | Google OAuth2 login |
| POST | `/api/v1/auth/logout` | Yes | Logout and revoke token |
| GET | `/health` | No | Service health check |

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [`goose`](https://github.com/pressly/goose) (migrations)
- [`golangci-lint`](https://golangci-lint.run) (linting)

### Local Setup

```bash
# 1. Start Postgres (localhost:5433) and Redis (localhost:6380)
make docker-up

# 2. Copy and fill in config
cp configs/env.example.yaml configs/env.local.yaml

# 3. Run migrations
make migrate-up

# 4. Start the server
make run
```

### Configuration

Config is loaded from `configs/env.{APP_ENV}.yaml` (default: `env.local.yaml`). Set `APP_ENV` to switch environments. Any config key can be overridden via environment variables using `_`-separated notation (e.g. `DATABASE_HOST`).

See `configs/env.example.yaml` for all available keys.

## Common Commands

```bash
make run                            # Run locally
make build                          # Compile binary to bin/api
make test                           # Run all tests with race detector
make lint                           # Run golangci-lint
make docker-up / docker-down        # Start / stop Postgres + Redis
make migrate-up / migrate-down      # Apply / rollback DB migrations
make create-migration name=<name>   # Scaffold a new migration file
```

## Adding a New Feature

1. Add entity in `internal/domain/entities/`
2. Add repository method to the interface in `internal/domain/repositories/`
3. Implement it in `internal/infrastructure/persistence/postgresinfra/repositories/`
4. Add use-case logic in `internal/application/usecases/` (DTOs in `internal/application/dto/`)
5. Add controller + routes in `internal/presentation/http/controllers/`
6. Wire it in `cmd/api/main.go` (`initControllers`)
