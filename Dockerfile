# Get Go image from DockerHub.
FROM golang:latest AS build

# Set working directory.
WORKDIR /compiler

# Copy dependency locks so we can cache.
COPY go.mod go.sum ./

# Get all of our dependencies.
RUN go mod download


# Copy all of our remaining application.
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g ./cmd/web/main.go

# Build our application.
RUN CGO_ENABLED=0 GOOS=linux go build -o web ./cmd/web/main.go

FROM scratch AS run

WORKDIR /app

COPY --from=build /compiler/.env .
COPY --from=build /compiler/web .
COPY --from=build /compiler/data/items.json ./data/items.json
COPY --from=build /compiler/data/tags.json ./data/tags.json
COPY --from=build /compiler/data/locations.json ./data/locations.json

EXPOSE 8000
CMD ["./web"]