# Docker 部署

本文档提供 `admin`、`http`、`task` 三个服务的 Docker 部署方式。

## 1. 说明

- `http` 容器实际运行 `cmd/api`
- `admin` 容器运行 `cmd/admin`
- `task` 容器运行 `cmd/task`
- `redis` 由 compose 内置提供
- MySQL 和 MongoDB 使用外部实例，例如阿里云 RDS / MongoDB
- 你只需要把外部 MySQL、MongoDB 地址写入 `.env.docker`

## 2. 准备配置

复制容器环境模板：

```bash
cp .env.docker.example .env.docker
```

至少补齐以下配置：

- `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASS`、`DB_NAME`
- `MONGO_URI`、`MONGO_DBNAME`
- `API_PORT`、`API_RPC_PORT`
- `ADMIN_PORT`
- `TASK_MODE`

如果你使用这套 compose 默认 Redis：

- `REDIS_HOST=redis`
- `REDIS_PORT=6379`
- `REDIS_PASSWORD` 可留空，也可以自行设置

容器部署时有两个关键差异：

- `API_LOCAL_IP` 必须是 `0.0.0.0`
- `RPC_CLIENTS` 不能写 `127.0.0.1`，要写成 `9010=http`

## 3. 启动

构建并启动四个容器：

```bash
docker compose --env-file .env.docker up -d --build redis http admin task
```

查看状态：

```bash
docker compose --env-file .env.docker ps
```

查看日志：

```bash
docker compose --env-file .env.docker logs -f http
docker compose --env-file .env.docker logs -f admin
docker compose --env-file .env.docker logs -f task
docker compose --env-file .env.docker logs -f redis
```

停止服务：

```bash
docker compose --env-file .env.docker down
```

## 4. 端口

- `http`：`${API_PORT}`
- `admin`：`${ADMIN_PORT}`
- `task`：无对外端口
- `redis`：默认仅在 compose 网络内使用
- `api` 内部 RPC：`${API_RPC_PORT}`，仅在容器网络中暴露给 `admin`

## 5. 文件说明

- `Dockerfile`：统一的多阶段构建文件
- `docker-compose.yml`：`redis + http + admin + task` 编排
- `.env.docker.example`：容器部署配置模板
- `.dockerignore`：减少构建上下文

## 6. Task 模式

当前 compose 默认：

```env
TASK_MODE=task
```

如果你要跑其他模式，可以改 `.env.docker` 后重启：

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

## 7. 注意事项

- `admin` 依赖 `http` 的内部 RPC，所以 compose 里保留了 `depends_on`
- `http`、`admin`、`task` 都会等待 `redis` 健康后再启动
- 三个容器都会把日志写到挂载目录 `./logs`
- `redis` 数据会持久化到 compose volume `redis_data`
- 这套 compose 不负责初始化 MySQL、MongoDB
- 如果你的 MySQL 和 MongoDB 已经放在阿里云，只需要把 `.env.docker` 填成真实地址即可
- 如果宿主机没有 Docker 或 Docker Compose，本方案无法直接执行
