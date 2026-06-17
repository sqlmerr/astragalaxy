export PROJECT_ROOT := `pwd`

set dotenv-load := true
set dotenv-path := ".env"
set dotenv-required := true

env-up:
    docker compose -f docker-compose.dev.yaml up postgres -d

env-down:
    docker compose -f docker-compose.dev.yaml down postgres

env-rm:
    docker compose -f docker-compose.dev.yaml down postgres -v && echo "Done"

migrate-create NAME:
    docker compose -f docker-compose.dev.yaml run --rm \
        -u $(id -u):$(id -g) \
        postgres-migrate \
        create \
        -dir /migrations \
        -ext sql \
        -seq "{{ NAME }}"

migrate-up:
    docker compose -f docker-compose.dev.yaml run --rm \
        -u $(id -u):$(id -g) \
        postgres-migrate \
        -path /migrations \
        -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable" \
        up

migrate-down:
    docker compose -f docker-compose.dev.yaml run --rm \
        -u $(id -u):$(id -g) \
        postgres-migrate \
        -path /migrations \
        -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable" \
        down

run-backend:
    @go mod tidy && go run cmd/server/main.go