# harbor-server

[中文文档](./README.zh-CN.md)

`harbor-server` is a Go backend for a crypto trading platform. It provides user-facing HTTP APIs, an admin backend, background task runners, WebSocket push services, a CDN/upload service, and a small set of maintenance tools.

## Highlights

- Single executable entrypoint layout under `cmd/*`
- Shared startup orchestration in `internal/bootstrap/*`
- Unified environment-driven configuration with `.env`
- Core integrations with MySQL, MongoDB, Redis, Ethereum, and Tron

## Project Layout

```text
.
├── cmd/                    # service and tool entrypoints
│   ├── api                 # user-facing HTTP API
│   ├── admin               # admin HTTP API
│   ├── task                # background jobs
│   ├── wss                 # websocket server
│   ├── cdn                 # upload/static service
│   └── tools               # maintenance tools
├── internal/bootstrap/     # process startup orchestration
├── internal/tools/         # shared tool implementations
├── http/                   # HTTP runtime and API modules
├── admin_modules/          # admin HTTP handlers
├── admin_models/           # admin-side business logic
├── models/                 # core domain logic and initialization
├── config/                 # configuration loading and compatibility layer
├── task/                   # task implementations
├── lib/                    # infra and chain integration helpers
└── utils/                  # shared helpers
```

## Requirements

- Go `1.18+`
- MySQL
- MongoDB
- Redis

## Configuration

The project now treats `.env` as the primary runtime configuration source.

1. Copy the example file:

```bash
cp .env.example .env
```

2. Fill in the required values:

- MySQL: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
- MongoDB: `MONGO_URI`, `MONGO_DBNAME`
- Redis: `REDIS_HOST`, `REDIS_PORT`
- Service ports: `API_PORT`, `ADMIN_PORT`, `WSS_PORT`, `CDN_PORT`
- Optional integrations: `ETH_RPC_URL`, `TRON_GRPC_ADDR`, `EXCHANGE_RATE_API_KEY`, `WS_ADMIN_PASS`

3. Optional compatibility file:

- `config/config.example.json` is provided only as a compatibility template.
- Real secrets should stay in `.env` and must not be committed.

## Service Entrypoints

Run services only through `cmd/*`:

```bash
go run ./cmd/api
go run ./cmd/admin
go run ./cmd/task data
go run ./cmd/task approve
go run ./cmd/task task
go run ./cmd/wss
go run ./cmd/cdn
```

### Service Roles

- `cmd/api`: user-facing API service
- `cmd/admin`: admin backend service
- `cmd/task`: background workers, selected by task mode
- `cmd/wss`: WebSocket push service
- `cmd/cdn`: upload/static file service

## Tool Entrypoints

Maintenance tools are also normalized under `cmd/tools/*`:

```bash
go run ./cmd/tools/mongo-init
go run ./cmd/tools/checkapprove
go run ./cmd/tools/lucky
go run ./cmd/tools/recover-kline
```

### Tool Roles

- `mongo-init`: initialize MongoDB collections for market data
- `checkapprove`: approval-related helper script
- `lucky`: lightweight timer/demo utility
- `recover-kline`: recover historical kline data into MongoDB

## Development

Build everything:

```bash
go build ./...
```

Run tests:

```bash
go test ./...
```

## Notes

- The repository has been refactored to keep service and tool entrypoints explicit.
- Legacy root-level wrapper mains and duplicated tool scripts have been removed.
- Configuration has been migrated away from committed secrets and toward environment variables.
