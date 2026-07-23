# Docker Deployment

This document covers Docker deployment for the `admin`, `http`, and `task` services.

## Overview

- `http` runs `cmd/api`
- `admin` runs `cmd/admin`
- `task` runs `cmd/task`
- `redis` is included in the compose stack
- MySQL and MongoDB are expected to be external services, such as Alibaba Cloud instances
- You only need to point `.env.docker` to the external MySQL and MongoDB endpoints

## Prepare Config

Copy the Docker environment template:

```bash
cp .env.docker.example .env.docker
```

Fill in at least:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
- `MONGO_URI`, `MONGO_DBNAME`
- `API_PORT`, `API_RPC_PORT`
- `ADMIN_PORT`
- `TASK_MODE`

When using the bundled Redis service:

- `REDIS_HOST=redis`
- `REDIS_PORT=6379`
- `REDIS_PASSWORD` is optional

Container deployment requires two important differences from local startup:

- `API_LOCAL_IP=0.0.0.0`
- `RPC_CLIENTS=9010=http`

## Start Services

```bash
docker compose --env-file .env.docker up -d --build redis http admin task
```

Check service status:

```bash
docker compose --env-file .env.docker ps
```

Tail logs:

```bash
docker compose --env-file .env.docker logs -f http
docker compose --env-file .env.docker logs -f admin
docker compose --env-file .env.docker logs -f task
docker compose --env-file .env.docker logs -f redis
```

Stop services:

```bash
docker compose --env-file .env.docker down
```

## Ports

- `http`: `${API_PORT}`
- `admin`: `${ADMIN_PORT}`
- `task`: no public port
- `redis`: internal to the compose network by default
- API internal RPC: `${API_RPC_PORT}`, exposed only inside the compose network

## Files

- `Dockerfile`: shared multi-stage build
- `docker-compose.yml`: orchestration for `redis + http + admin + task`
- `.env.docker.example`: Docker environment template
- `.dockerignore`: reduced build context

## Task Mode

The compose file defaults to:

```env
TASK_MODE=task
```

To switch modes, update `.env.docker` and rebuild the task service:

```bash
docker compose --env-file .env.docker up -d --build task
```

## Notes

- `admin` depends on the internal RPC exposed by `http`
- `http`, `admin`, and `task` wait for Redis to become healthy
- Application logs are written into the mounted `./logs` directory
- Redis data is persisted in the `redis_data` volume
- The compose file does not initialize MySQL or MongoDB for you
