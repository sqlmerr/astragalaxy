FROM lukemathwalker/cargo-chef:latest-rust-1.79.0 as chef
WORKDIR /app

FROM chef AS planner

COPY . .

RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder

COPY --from=planner /app/recipe.json recipe.json

# Установите необходимые зависимости
RUN apt update && apt install -y \
    build-essential \
    cmake \
    libsodium-dev \
    libsecp256k1-dev \
    lz4 \
    liblz4-dev \
    zlib1g-dev \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*

RUN cargo chef cook --release --recipe-path recipe.json

COPY . .

RUN cargo build --release

FROM debian:bookworm-slim AS runtime

RUN apt update && apt install -y \
    libssl-dev \
    libsodium23 \
    libsecp256k1-dev \
    zlib1g \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*


COPY --from=builder /app/.env .env
COPY --from=builder /app/target/release/api api

EXPOSE 8000
CMD ["./api"]
