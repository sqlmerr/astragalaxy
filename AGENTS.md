# AGENTS.md

Guidance for AI coding agents and contributors working in this repository.

## Project Overview

AstraGalaxy is an API-based space exploration game inspired by No Man's Sky. A Go backend exposes a
REST API that players interact with via HTTP requests or custom clients. The `client/` directory is
an official (optional) client built with React and TanStack Start, living in this monorepo. Players
manage agents that explore a procedurally generated universe, trade, mine, craft, and build
facilities.

See [`README.md`](README.md) for a full project description and user-facing docs.

## Repository Layout

```
api/openapi.yaml           # OpenAPI 3.0 spec (source of truth; jsons are generated)
cmd/server/main.go         # Go HTTP server entry point
client/                    # Official (optional) client — React + TanStack Start (Bun)
config.yaml                # Game configuration (seed, rules)
data/                      # Static game data (items, resources, recipes, facilities)
internal/                  # Private Go packages
  auth/                    # Authentication (JWT, agent tokens)
  data/                    # Data layer: registry + Postgres repositories/sqlc
  errors/                  # Error codes and types
  game/                    # Core game logic (domain services)
  logger/                  # Structured logging (Zap)
  model/                   # Domain models
  transport/               # HTTP handlers, middleware, server
Justfile                   # Development task runner
sqlc.yaml                  # SQL code generation config
```

## Tech Stack

- **Backend**: Go 1.26, PostgreSQL (pgx/v5) + Redis, JWT auth, OpenAPI 3.0 (Scalar docs), sqlc.
- **Frontend**: React 19 + TypeScript, TanStack Start (SSR), TanStack Router + Query, shadcn/ui +
  Tailwind CSS, PixiJS, Vite, Bun.
- **Tooling**: `just`, Docker Compose, sqlc, redocly.

## Commands

Run via `just` (uses `.env`; `.env` is required — see `.env.example`):

```bash
just env-up            # Start PostgreSQL & Redis (dev)
just env-down          # Stop services
just env-rm            # Stop and remove volumes
just migrate-create NAME
just migrate-up        # Apply migrations
just migrate-down N    # Rollback N migrations
just migrate-force V   # Force migration version
just gen-docs          # Bundle openapi.yaml -> out/openapi.json + regenerate client types
just gen-db            # sqlc generate
just run-server        # go mod tidy && go run cmd/server/main.go
just test              # go test -v ./internal/game/...
just run-client        # client dev server (port 3000)
```

Client checks (run inside `client/`):

```bash
bun test            # vitest run
bun typecheck       # tsc --noEmit
bun lint            # eslint
bun format          # prettier --write
bun check           # prettier --check
```

## Key Workflows

### API Specification

- `api/openapi.yaml` is the single source of truth. Endpoint definitions are split into
  `api/paths/**`. The auth security schemes (`UserAuth`, `AgentAuth`, `AgentId`) live in the
  `components/securitySchemes` block.
- Regenerate the bundled JSON and client types with `just gen-docs` after changing the spec.
  Client types are written to `src/api/schema.d.ts` (used with `openapi-fetch` /
  `openapi-react-query`).

### Domain Logic (Backend)

- `internal/game/**` holds the domain services (agents, ships, inventory, navigation, mining,
  crafting, cooldowns, etc.). Each package may have `*_test.go` files (see `just test`).
- `internal/transport/**` holds HTTP handlers/middleware; `internal/model` holds domain models.
- Error codes/types live in `internal/errors`; use them consistently in handlers.
- `internal/data` is the persistence layer — sqlc-generated Postgres code under
  `internal/data/postgres/database/sqlc`, migrations and queries under
  `internal/data/postgres/database/{migration,query}`.

### Game Data

- Static game data lives in `data/` (e.g. `items.json`, `resources.json`, `recipes.json`,
  `facilities.json`) and is loaded through the registry in `internal/data`.

## Conventions

- Backend is Go. Follow existing package structure and use the established dependencies
  (pgx, koanf for config, envconfig for env vars, zap for logging).
- The client is typed and uses the generated OpenAPI types; keep TypeScript strict.
- Match surrounding code style; do not add code comments unless they add real value.
- Run relevant tests before finishing: `just test` for backend, and
  `bun test` / `bun typecheck` for client changes.
- Write everything in English
