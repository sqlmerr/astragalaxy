# AstraGalaxy

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/sqlmerr/astragalaxy)

AstraGalaxy is an API-based space exploration game inspired by No Man's Sky. Players interact with the game through HTTP requests or custom clients, managing agents that explore a procedurally generated universe, trade, mine, craft, and build facilities. This repository (monorepo) contains the game server and an official (optional) client.

## Game Overview

- **Agents**: Each player account can have up to 5 agents. Agents are the key characters that fly between systems, galaxies, and planets, trade with other agents, engage in combat, complete missions, and build bases.
- **Ships**: Each agent commands ships (Trader, Scout, Miner) with different capabilities.
- **Universe**: Procedurally generated galaxy with systems containing planets, waypoints, and asteroid fields.
- **Economy**: Resource mining, trading, crafting, and facility management.

## Technology Stack

**Backend:**
- Go 1.26
- PostgreSQL (pgx) + Redis
- JWT authentication
- OpenAPI 3.0 (Scalar for docs)
- sqlc for type-safe SQL

**Client (official, optional):**
- React 19 + TypeScript
- TanStack Start (SSR)
- TanStack Router + Query
- shadcn/ui + Tailwind CSS
- PixiJS for game rendering
- Vite

## Quick Start

### Prerequisites

- Go 1.26+
- Bun (for the client) — optional
- Docker & Docker Compose
- Just (command runner)

### Setup (Local)

1. **Clone and configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your database/redis credentials and JWT secret
   ```

2. **Start infrastructure:**
   ```bash
   just env-up
   ```

3. **Run migrations:**
   ```bash
   just migrate-up
   # or from the server binary (uses the same .env):
   go run ./cmd/server migrate up
   ```

4. **Generate API docs (optional):**
   ```bash
   just gen-docs
   ```

5. **Run server:**
   ```bash
   just run-server
   # Server starts at http://localhost:8080
   # API docs at http://localhost:8080/docs
   ```

6. **Run client (in another terminal):**
   ```bash
   just run-client
   # Client at http://localhost:3000
   ```

### Docker

The server ships as a container image published to GitHub Container Registry on every release:

```bash
docker build -t astragalaxy .

docker run --rm -p 8080:8080 \
  --env-file .env \
  -e POSTGRES_HOST=host.docker.internal \
  -e REDIS_ADDR=host.docker.internal:6379 \
  astragalaxy
```

Images are tagged with the release version (`v1.2.3`, `1.2`), `latest`, and the commit SHA:

```bash
docker pull ghcr.io/sqlmerr/astragalaxy:latest
```

CI (`.github/workflows/test-and-build.yml`) runs `go vet`, `go test`, and validates the image build on pushes and pull requests.
CD (`.github/workflows/publish.yml`) runs the same checks and pushes the image to GHCR when a release is published.

To use builtin admin cli to manage server instance, do this:
```bash
docker compose exec server admincli 
```

## Configuration

### Environment Variables (`.env`)

| Variable | Description |
|----------|-------------|
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `POSTGRES_DB` | Database name |
| `POSTGRES_HOST` | Database host |
| `POSTGRES_PORT` | Database port (5432) |
| `POSTGRES_TIMEOUT` | Connection timeout |
| `REDIS_USERNAME` | Redis username (optional) |
| `REDIS_PASSWORD` | Redis password |
| `REDIS_ADDR` | Redis address (host:port) |
| `DB` | Redis database number |
| `HTTP_ADDR` | Server listen address (default: 0.0.0.0:8080) |
| `HTTP_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout |
| `LOGGER_LEVEL` | Log level (DEBUG, INFO, WARN, ERROR) |
| `LOGGER_FOLDER` | Log output directory |
| `JWT_SECRET` | Secret for JWT signing |
| `TOKEN_DURATION` | JWT token expiry (e.g., 7d) |

### Game Config (`config.yaml`)

```yaml
game:
  seed: 123456789             # Universe generation seed
  rules:
    disableCooldowns: true    # Disable action cooldowns (dev)
    disableInventoryLimit: false
registry:
  itemsPath: data/items.json
  resourcesPath: data/resources.json
  recipesPath: data/recipes.json
  facilitiesPath: data/facilities.json
```

## API Authentication

Two authentication methods are supported:

### 1. Agent Token (Direct)
Each agent has a unique token prefixed with `ag_agent_` followed by a random hash. Use as Bearer token:
```
Authorization: Bearer ag_agent_<hash>
```
Best for single-agent clients/bots.

### 2. User JWT + Agent ID Header
Use your account JWT with an agent ID header:
```
Authorization: Bearer <user_jwt>
X-Agent-ID: <agent_uuid>
```
Best for multi-agent management interfaces.

### Security Schemes (OpenAPI)
- **UserAuth**: HTTP Bearer (JWT)
- **AgentAuth**: HTTP Bearer (agent token)
- **AgentId**: API Key header `X-Agent-ID` (must be used with UserAuth)

## API

Full API specification is available in [`api/openapi.yaml`](api/openapi.yaml).

## Game Mechanics

### Resources
Resources have tiers (basic, advanced, exotic) and tags (`asteroid_resource`, `planet_resource`). Found on planets based on biome type (habitable, frozen, dead) and in asteroid fields.

### Mining
- **Asteroid mining**: Extract asteroid resources (iron, copper, titanium, helium, uranium, iridium)
- **Planet mining**: Extract planetary resources based on biome deposits

### Crafting
Recipes require specific facilities (smelter, printer) and minimum tier. Example: Steel requires Smelter tier 1, 5 Iron + 5 Carbon → 1 Steel (10s base).

### Facilities
- **Printers**: Craft items (portable, standard, advanced tiers)
- **Smelters**: Process resources (portable, standard tiers)
Higher tiers have better time/cost multipliers.

### Ships
| Type | Role |
|------|------|
| Trader | Cargo capacity, trading |
| Scout | Exploration, radar range |
| Miner | Mining efficiency |

### Cooldowns
Actions have cooldowns (configurable). Warp cooldown depends on distance.

## Development Commands

```bash
just env-up           # Start PostgreSQL & Redis
just env-down         # Stop services
just env-rm           # Stop and remove volumes
just migrate-create NAME  # Create new migration
just migrate-up       # Apply migrations
just migrate-down N   # Rollback N migrations
just migrate-force V  # Force migration version
just gen-docs         # Bundle OpenAPI to JSON
just gen-db           # Generate SQLC code
just run-server       # Run Go server
just test             # Run game tests
just run-client       # Run client dev server
```

### Migrations from the server binary

Migration SQL files are embedded into the server binary, so `migrate` can be run
without Docker:

```bash
go run ./cmd/server migrate up        # Apply all pending migrations (default)
go run ./cmd/server migrate up 2      # Apply only the next 2 migrations
go run ./cmd/server migrate down 1    # Roll back the last migration
go run ./cmd/server migrate version   # Show current migration version
go run ./cmd/server migrate force 1   # Set version manually (dirty state recovery)
```

Configuration comes from `POSTGRES_*` environment variables; a `.env` file in the
current working directory is loaded automatically (existing environment variables
take precedence). In Docker the same command is `server migrate up`.

## Testing

```bash
# Backend tests
just test

# Frontend tests
cd client && bun test

# Frontend typecheck
cd client && bun typecheck
```

## Project Structure Details

### Key Domain Models
- **User**: Account holder (username, password hash)
- **Agent**: Player character (token, inventory, max 5/user)
- **Ship**: Vehicle (type, status, location, modules, inventory)
- **Inventory**: Resource/item storage (capacity limits)
- **Cooldown**: Action rate limiting per agent

### World Generation
Deterministic procedural generation based on seed. Systems at integer coordinates (x, y) contain planets and waypoints with biome-dependent resource deposits.

## Contributing

1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Run `just test` and client checks (if you changed the client)
5. Submit PR

## License

MIT License