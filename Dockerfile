FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/admincli ./cmd/admincli

FROM node:22-alpine AS docs

RUN npm i -g @redocly/cli@2

WORKDIR /spec

COPY api ./api

RUN redocly bundle api/openapi.yaml -o /out/openapi.json


FROM alpine:3.22

RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app

WORKDIR /app

COPY --from=build /out/server /usr/local/bin/server
COPY --from=build /out/admincli /usr/local/bin/admincli
COPY --from=docs /out/openapi.json ./out/openapi.json
COPY data ./data/
COPY config.example.yaml ./config.yaml

RUN mkdir -p /app/out/logs && chown -R app:app /app

ENV HTTP_ADDR=0.0.0.0:8080 \
    HTTP_SHUTDOWN_TIMEOUT=10s \
    LOGGER_LEVEL=INFO \
    LOGGER_FOLDER=/app/out/logs

EXPOSE 8080
USER app

ENTRYPOINT ["server"]
