# harbor-server

[English README](./README.md)

`harbor-server` 是一个面向加密资产交易场景的 Go 后端，提供用户侧 HTTP API、后台管理接口、后台任务调度、WebSocket 推送、上传/CDN 服务，以及少量运维工具。

## 特性

- 所有可执行入口统一收口到 `cmd/*`
- 启动编排统一放在 `internal/bootstrap/*`
- 配置统一以 `.env` 为主
- 基础依赖覆盖 MySQL、MongoDB、Redis、Ethereum、Tron

## 目录结构

```text
.
├── cmd/                    # 服务与工具入口
│   ├── api                 # 用户侧 HTTP API
│   ├── admin               # 后台管理 API
│   ├── task                # 后台任务入口
│   ├── wss                 # WebSocket 服务
│   ├── cdn                 # 上传/静态服务
│   └── tools               # 运维工具入口
├── internal/bootstrap/     # 启动层
├── internal/tools/         # 工具共享实现
├── http/                   # HTTP 运行层与业务模块
├── admin_modules/          # 后台接口层
├── admin_models/           # 后台业务模型层
├── models/                 # 核心业务模型层
├── config/                 # 配置加载与兼容层
├── task/                   # 任务实现
├── lib/                    # 基础设施与链上集成
└── utils/                  # 通用工具
```

## 环境要求

- Go `1.18+`
- MySQL
- MongoDB
- Redis

## 配置方式

项目现在默认使用 `.env` 作为主配置来源。

1. 复制模板：

```bash
cp .env.example .env
```

2. 至少补齐这些配置：

- MySQL：`DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASS`、`DB_NAME`
- MongoDB：`MONGO_URI`、`MONGO_DBNAME`
- Redis：`REDIS_HOST`、`REDIS_PORT`
- 服务端口：`API_PORT`、`ADMIN_PORT`、`WSS_PORT`、`CDN_PORT`
- 可选外部依赖：`ETH_RPC_URL`、`TRON_GRPC_ADDR`、`EXCHANGE_RATE_API_KEY`、`WS_ADMIN_PASS`

3. 兼容说明：

- `config/config.example.json` 仅作为兼容模板保留
- 真实敏感配置应只放在 `.env`，不要提交到仓库

## 服务入口

服务统一通过 `cmd/*` 启动：

```bash
go run ./cmd/api
go run ./cmd/admin
go run ./cmd/task data
go run ./cmd/task approve
go run ./cmd/task task
go run ./cmd/wss
go run ./cmd/cdn
```

### 服务说明

- `cmd/api`：用户侧 API 服务
- `cmd/admin`：后台管理服务
- `cmd/task`：后台任务服务，通过 mode 选择任务类型
- `cmd/wss`：WebSocket 推送服务
- `cmd/cdn`：上传与静态资源服务

## 工具入口

运维工具也统一收口在 `cmd/tools/*`：

```bash
go run ./cmd/tools/mongo-init
go run ./cmd/tools/checkapprove
go run ./cmd/tools/lucky
go run ./cmd/tools/recover-kline
```

### 工具说明

- `mongo-init`：初始化 MongoDB 行情相关集合
- `checkapprove`：授权检查辅助脚本
- `lucky`：轻量定时/演示工具
- `recover-kline`：向 MongoDB 恢复历史 K 线数据

## 开发命令

全量构建：

```bash
go build ./...
```

运行测试：

```bash
go test ./...
```

## 扩展文档

- [部署文档](./docs/deployment.zh-CN.md)
- [API / 服务说明](./docs/services.zh-CN.md)
- [重构路线图](./docs/refactor-roadmap.zh-CN.md)

## 说明

- 当前仓库已经完成一轮入口层与配置层收口
- 历史根目录壳入口与重复工具脚本已被移除
- 后续重构可以继续围绕启动层、配置层、后台拆分和核心模型解耦推进
