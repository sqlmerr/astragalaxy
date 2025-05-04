FROM --platform=$BUILDPLATFORM golang:latest AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /compiler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o web ./cmd/web/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS run

LABEL authors="sqlmerr"
LABEL org.opencontainers.image.source="https://github.com/sqlmerr/astragalaxy"

WORKDIR /app

COPY --from=build /compiler/web .
COPY --from=build /compiler/data/items.json ./data/items.json
COPY --from=build /compiler/data/tags.json ./data/tags.json
COPY --from=build /compiler/data/locations.json ./data/locations.json

USER nonroot

EXPOSE 8000
ENTRYPOINT ["./web"]