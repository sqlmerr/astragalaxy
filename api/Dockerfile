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

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

# Build our application.
RUN CGO_ENABLED=0 GOOS=linux go build -o web .

FROM scratch AS run

WORKDIR /app

COPY --from=build /compiler/.env .
COPY --from=build /compiler/web .

EXPOSE 8000
CMD ["./web"]