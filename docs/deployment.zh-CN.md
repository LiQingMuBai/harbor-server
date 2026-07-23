# 部署文档

本文档描述 `harbor-server` 的部署依赖、配置准备、服务启动顺序以及常见发布方式。

## 1. 部署概览

`harbor-server` 当前包含 5 个常驻服务和 4 个工具入口：

- `cmd/api`：用户侧 HTTP API，同时启动内部 RPC 服务
- `cmd/admin`：后台管理服务，会消费 Redis 队列并调用 API RPC
- `cmd/task`：后台任务服务，按模式执行
- `cmd/wss`：WebSocket 推送服务
- `cmd/cdn`：上传与静态资源服务

建议把它部署为多进程服务，而不是把所有角色塞进同一个进程。

## 2. 环境要求

- Go `1.18+`
- MySQL
- MongoDB
- Redis
- Linux 或 macOS

## 3. 配置准备

项目当前以 `.env` 作为主配置来源。

1. 复制模板：

```bash
cp .env.example .env
```

2. 至少补齐以下配置：

- MySQL：`DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASS`、`DB_NAME`
- MongoDB：`MONGO_URI`、`MONGO_DBNAME`
- Redis：`REDIS_HOST`、`REDIS_PORT`
- API：`API_PORT`、`API_LOCAL_IP`、`API_RPC_PORT`
- Admin：`ADMIN_PORT`、`RPC_CLIENTS`
- WSS：`WSS_PORT`
- CDN：`CDN_PORT`、`CDN_DOMAIN`

3. 推荐补齐以下外部依赖配置：

- `ETH_RPC_URL`
- `TRON_GRPC_ADDR`
- `EXCHANGE_RATE_API_KEY`
- `WS_ADMIN_PASS`
- `SMS_API_URL`
- `SMS_API_ID`
- `SMS_API_KEY`

## 4. 关键配置项说明

### 基础依赖

- `DB_*`：MySQL 连接配置
- `MONGO_*`：MongoDB 配置，优先使用 `MONGO_URI`
- `REDIS_*`：Redis 配置

### 服务配置

- `API_PORT`：用户 API 监听端口，默认 `9001`
- `API_LOCAL_IP`：API 内部 RPC 绑定地址，默认 `127.0.0.1`
- `API_RPC_PORT`：API 内部 RPC 端口，默认 `9010`
- `ADMIN_PORT`：后台服务端口，默认 `8080`
- `RPC_CLIENTS`：Admin 连接的 RPC 客户端列表，例如 `9010=127.0.0.1,9020=127.0.0.1`
- `WSS_PORT`：WebSocket 服务端口，默认 `9088`
- `CDN_PORT`：上传服务端口，默认 `9999`
- `CDN_DOMAIN`：上传后返回的静态资源前缀域名
- `TASK_MODE`：任务模式，可选 `data`、`approve`、`task`

### 可选集成

- `ETH_RPC_URL`：Ethereum RPC 节点地址
- `TRON_GRPC_ADDR`：Tron gRPC 地址，默认 `grpc.trongrid.io:50052`
- `EXCHANGE_RATE_API_KEY`：汇率接口 Key
- `WS_ADMIN_PASS`：WebSocket 管理员登录口令
- `RECOVER_MONGO_URI`：`recover-kline` 工具专用 Mongo URI

## 5. 构建

全量构建：

```bash
go build ./...
```

按服务编译：

```bash
go build -o bin/api ./cmd/api
go build -o bin/admin ./cmd/admin
go build -o bin/task ./cmd/task
go build -o bin/wss ./cmd/wss
go build -o bin/cdn ./cmd/cdn
```

## 6. 启动顺序

推荐启动顺序如下：

1. MySQL / MongoDB / Redis
2. `api`
3. `admin`
4. `task`
5. `wss`
6. `cdn`

原因：

- `api` 会启动内部 RPC 服务，`admin` 依赖它处理后台广播任务
- `task`、`wss`、`admin` 都依赖数据库和缓存初始化
- `cdn` 相对独立，可以最后启动

## 7. 启动命令

### API

```bash
go run ./cmd/api
```

### Admin

```bash
go run ./cmd/admin
```

### Task

三种模式分别运行：

```bash
go run ./cmd/task data
go run ./cmd/task approve
go run ./cmd/task task
```

也可以通过环境变量控制：

```bash
TASK_MODE=task go run ./cmd/task
```

### WSS

```bash
go run ./cmd/wss
```

### CDN

```bash
go run ./cmd/cdn
```

## 8. 工具命令

### 初始化 Mongo 行情集合

```bash
go run ./cmd/tools/mongo-init
```

### 授权检查辅助脚本

```bash
go run ./cmd/tools/checkapprove
```

### 定时演示工具

```bash
go run ./cmd/tools/lucky
```

### 恢复 K 线历史

```bash
go run ./cmd/tools/recover-kline
```

## 9. Docker 部署

当前仓库已经内置以下 Docker 部署文件：

- `Dockerfile`
- `docker-compose.yml`
- `.env.docker.example`
- `.dockerignore`

### 9.1 适用场景

当前 Docker 方案默认适合以下部署结构：

- `admin`、`http`、`task` 运行在 Docker 容器中
- `Redis` 使用 compose 内置容器
- `MySQL` 使用外部实例，例如阿里云 RDS
- `MongoDB` 使用外部实例，例如阿里云 MongoDB

说明：

- `http` 容器实际运行 `cmd/api`
- `admin` 容器实际运行 `cmd/admin`
- `task` 容器实际运行 `cmd/task`

### 9.2 准备容器配置

复制模板：

```bash
cp .env.docker.example .env.docker
```

至少补齐以下配置：

- `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASS`、`DB_NAME`
- `MONGO_URI`、`MONGO_DBNAME`
- `API_PORT`、`API_RPC_PORT`
- `ADMIN_PORT`
- `TASK_MODE`

如果使用 compose 内置 Redis，推荐保持：

- `REDIS_HOST=redis`
- `REDIS_PORT=6379`
- `REDIS_PASSWORD=` 可空，或按需设置

Docker 网络下有两个关键配置必须注意：

- `API_LOCAL_IP=0.0.0.0`
- `RPC_CLIENTS=9010=http`

原因：

- `cmd/api` 在容器内需要对其他容器暴露 RPC，所以不能继续绑定 `127.0.0.1`
- `cmd/admin` 不能再通过本机回环地址访问 API RPC，必须改成 compose 服务名 `http`

### 9.3 启动命令

构建并启动：

```bash
docker compose --env-file .env.docker up -d --build redis http admin task
```

查看状态：

```bash
docker compose --env-file .env.docker ps
```

查看日志：

```bash
docker compose --env-file .env.docker logs -f redis
docker compose --env-file .env.docker logs -f http
docker compose --env-file .env.docker logs -f admin
docker compose --env-file .env.docker logs -f task
```

停止服务：

```bash
docker compose --env-file .env.docker down
```

### 9.4 服务与端口

- `redis`：仅在 compose 网络中使用，数据持久化到 `redis_data`
- `http`：对外暴露 `${API_PORT}`，容器内同时暴露 `${API_RPC_PORT}` 给 `admin`
- `admin`：对外暴露 `${ADMIN_PORT}`
- `task`：不对外暴露端口

### 9.5 Task 模式

默认：

```env
TASK_MODE=task
```

如果你要跑其他模式，修改 `.env.docker` 后重建 `task` 即可：

```env
TASK_MODE=data
```

或：

```env
TASK_MODE=approve
```

然后执行：

```bash
docker compose --env-file .env.docker up -d --build task
```

### 9.6 Docker 部署注意事项

- `admin` 依赖 `http` 的内部 RPC，所以 compose 内保留了启动依赖关系
- `http`、`admin`、`task` 都会等待 `redis` 健康后再启动
- 日志目录挂载到宿主机 `./logs`
- 当前 compose 不负责初始化 MySQL、MongoDB
- 如果你的 MySQL 和 MongoDB 已经放在阿里云，只需要把 `.env.docker` 填成真实地址即可
- 宿主机需要提前安装 Docker 和 Docker Compose

## 10. 生产部署建议

### 进程管理

建议使用以下方式之一托管进程：

- `systemd`
- `supervisor`
- `pm2` 仅作为简单托管层
- 容器化部署

### 日志

- 标准输出日志建议交给进程管理器采集
- `http/common` 中仍保留本地文件日志逻辑，生产环境建议统一重定向或收口

### 静态目录

`cdn` 服务默认依赖以下目录：

- `./static`
- `./pdf`
- `./whitepaper`

发布时要确保这些目录存在并具备写权限。

### 安全

- 不要提交 `.env`
- 不要提交 `config/config.json`
- `WS_ADMIN_PASS` 必须单独配置
- 生产环境必须替换所有默认地址与示例值

## 11. 发布检查清单

- `.env` 已填写且未提交
- Docker 部署时 `.env.docker` 已填写且未提交
- MySQL / MongoDB / Redis 可连通
- `API_LOCAL_IP` 与 `RPC_CLIENTS` 配置一致
- `CDN_DOMAIN` 指向真实外部访问域名
- `ETH_RPC_URL`、`TRON_GRPC_ADDR`、短信与汇率相关配置已按业务需要补齐
- `go build ./...` 和 `go test ./...` 已执行
