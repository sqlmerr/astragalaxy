# 🚀 astragalaxy
![Go 1.23](https://img.shields.io/static/v1?logo=Go&label=&message=Go+1.23&color=FFFFFF)
![Python 3.12](https://img.shields.io/static/v1?logo=Python&label=&message=Python+3.12&color=FFFFFF)

![GitHub commit activity](https://img.shields.io/github/commit-activity/t/sqlmerr/astragalaxy)
![Github Created At](https://img.shields.io/github/created-at/sqlmerr/astragalaxy)
![GitHub last commit](https://img.shields.io/github/last-commit/sqlmerr/astragalaxy)

[![Author](https://img.shields.io/badge/Author-@sqlmerr-purple)](https://sqlmerr.github.io)
[![License](https://img.shields.io/badge/License-MIT-purple)](#license)


# What is this?
**AstraGalaxy** is a space travelling game. Explore different locations. Travel to interesting systems and planets. Gather resources and create your own base (in future updates). You can play it in telegram bot!

# 👨‍💻 Run
## Local
### Specify environment variables
```bash
cp .env.example .env  # or just copy .env and edit it
```

### Install packages
```bash
apt install docker docker-compose -y
```

### And then run in docker.
```bash
docker compose up --build
```

## Docker
```bash
docker run -p "8000:8000" -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_HOST=0.0.0.0 -e POSTGRES_PORT=5432 -e JWT_SECRET=secret -e SECRET_TOKEN=token  ghcr.io/sqlmerr/astragalaxy:latest
```

## Docker Compose
```yaml
services:
  db:
    image: postgres:16.2-alpine
    command: -p 5432
    expose:
        - 5432
    ports:
        - '5432:5432'
    volumes:
        - app-db-data:/var/lib/postgresql/
    env_file:
        - .env
    environment:
        - PGDATA=/var/lib/postgresql/data/pgdata
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
  astragalaxy:
    image: ghcr.io/sqlmerr/astragalaxy:latest
    ports:
      - '8000:8000'
    env_file:
      - .env
    depends_on:
      db:
        condition: service_healthy
volumes:
  app-db-data:
```