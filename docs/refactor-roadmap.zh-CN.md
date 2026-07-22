# 重构路线图

本文档用于描述 `harbor-server` 当前已经完成的重构、后续建议阶段以及每一阶段的边界与预期收益。

## 1. 当前状态

项目已经从“历史多入口、配置分散、工具脚本散落”的状态，进入到“正在分阶段收口”的过渡期。

目前主路径已经比较明确：

- 可执行入口：`cmd/*`
- 启动编排：`internal/bootstrap/*`
- 用户 API：`http/modules/*`
- 后台接口：`admin_modules/*`
- 核心业务：`models/*`
- 配置加载：`config/*`

## 2. 已完成阶段

### 阶段 1：入口层收口

已完成内容：

- 将服务入口统一收口到 `cmd/api`、`cmd/admin`、`cmd/task`、`cmd/wss`、`cmd/cdn`
- 将工具入口统一收口到 `cmd/tools/*`
- 删除根目录与历史目录中重复/散落的 wrapper main
- 收敛多份重复 Mongo 初始化脚本为一份共享实现

收益：

- 可执行路径清晰
- 代码扫描噪音显著下降
- 后续部署和文档化成本降低

### 阶段 2：配置层收口

已完成内容：

- 以 `.env` 作为主配置来源
- 启动层 `internal/bootstrap/shared` 现在也能读到 `.env`
- 删除仓库内真实 `config/config.json`
- 增加 `config/config.example.json` 作为兼容模板
- 将部分硬编码敏感信息改为环境变量读取

收益：

- 避免敏感信息进入仓库
- 启动层和旧配置层读取逻辑更一致
- 部署方式更标准

## 3. 后续阶段总览

建议按以下顺序继续推进。

### 阶段 3：启动初始化解耦

目标：

- 让 `bootstrap` 真正成为唯一启动编排层
- 把 `httprun.Execute()` 中的初始化职责拆回启动层
- 让运行层只负责 HTTP/WSS 装配，不再负责系统初始化和后台 goroutine 拉起

重点文件：

- `internal/bootstrap/api/runtime.go`
- `internal/bootstrap/admin/runtime.go`
- `internal/bootstrap/wss/runtime.go`
- `http/httprun/httpprocess.go`
- `models/model.go`

建议动作：

- 把 `models.InitData()` 调用位置收口到 `bootstrap`
- 把 `RpcServer`、`CheckApprove`、`GetWalletBalance` 等后台线程装配从 `httprun.Execute()` 中抽离
- 把 API/WSS/CDN/Admin 的启动参数统一成显式 `Options`

收益：

- 启动链路更符合“启动层/路由层/业务层”职责分离
- 单测和局部调试成本降低

### 阶段 4：认证与会话层解耦

目标：

- 把 `http/common` 对 `models.MODEL_USER` 的直接依赖拆开
- 将用户鉴权逻辑抽成独立服务接口

重点文件：

- `http/common/common.go`
- `models/user.go`
- `admin_modules/user.go`

建议动作：

- 提取 `auth/session` 服务
- 用接口替换 `NeedLogin` 中对 `models.MODEL_USER.CheckSessionId` 的直接调用
- 后台 `CheckLogin` 也逐步抽离为复用中间件

收益：

- 路由层不直接耦合核心模型
- 后台和前台的鉴权逻辑可逐步统一

### 阶段 5：后台模块拆分

目标：

- 拆掉 `admin_modules/user.go` 单文件大后台入口
- 按业务主题拆分后台路由与后台模型

建议拆分维度：

- `admin_user`
- `admin_wallet`
- `admin_trade`
- `admin_system`
- `admin_notice`
- `admin_agent`
- `admin_loan`

重点文件：

- `admin_modules/user.go`
- `admin_models/userModel.go`
- `admin_models/systemModel.go`
- `admin_models/tradeModel.go`

收益：

- 后台接口目录更清晰
- 风险面缩小
- 后续权限和菜单体系更易维护

### 阶段 6：核心账本与交易域拆分

目标：

- 逐步拆解 `models` 中的全局单例和交叉依赖
- 优先把账本、交易、资产、充值提现从“大模型层”中拆开

建议优先级：

1. 账本/余额变更
2. 充值提现
3. 资产与闪兑
4. 交易委托与持仓
5. 行情/币种配置

重点文件：

- `models/model.go`
- `models/user.go`
- `models/credit.go`
- `models/assets.go`
- `models/trade.go`
- `models/system.go`

建议动作：

- 以“接口 + service”方式重构，而不是一次性目录搬迁
- 优先提炼余额变更能力，收口到账本服务
- 再逐步让充值、交易、资产调用账本服务而非直接改库

收益：

- 降低核心业务耦合度
- 为单元测试和后续服务化打基础

### 阶段 7：文档与协议补全

目标：

- 让工程重构结果能被前端、测试、运维和后续开发稳定使用

建议补齐：

- 用户 API 字段级接口文档
- 后台 API 字段级接口文档
- WebSocket 协议文档
- 部署与进程管理模板
- 变更日志与迁移说明

## 4. 推荐执行顺序

建议按下面的节奏推进，而不是并行大改：

1. 启动初始化解耦
2. 认证与会话解耦
3. 后台模块拆分
4. 核心账本与交易域拆分
5. 文档和协议补全

原因：

- 前两步是“骨架治理”
- 第三步是“后台治理”
- 第四步才是真正高风险的核心业务拆分
- 第五步用于稳定交付

## 5. 每阶段约束

为了避免重构失控，建议遵守以下约束：

- 每一阶段只解决一个层级的问题
- 不在同一阶段同时改启动层和核心交易逻辑
- 所有新配置统一走 `.env`
- 所有新可执行入口统一放在 `cmd/*`
- 每阶段完成后必须跑 `go build ./...` 和 `go test ./...`
- 能删除的历史壳代码就删除，不做“保留兼容入口”的长期拖延

## 6. 近期建议

如果继续按当前节奏推进，下一步最值得做的是：

### 近期优先级 1

启动初始化解耦。

先把这些逻辑从运行层拔出来：

- RPC 服务启动
- `models.InitData()`
- 授权检测线程
- 钱包余额检测线程
- WSS 数据服务线程

### 近期优先级 2

后台模块拆分设计。

先不直接大改逻辑，只做：

- 路由分组规划
- 文件拆分草案
- `admin_models` 主题边界定义

### 近期优先级 3

账本服务抽象。

围绕余额变更、账变日志、通知副作用，抽出第一个真正的 service 边界。

## 7. 目标形态

重构后的理想形态建议是：

- `cmd/*`：唯一入口
- `internal/bootstrap/*`：唯一启动层
- `internal/app/*` 或等价目录：服务装配层
- `internal/service/*`：业务服务层
- `internal/repository/*`：数据访问层
- `models/*`：逐步退化或被拆分吸收
- `docs/*`：稳定维护部署、协议、路线图和接口文档

## 8. 完成标志

当出现以下结果时，说明本轮重构基本成功：

- 没有根目录历史壳入口
- `.env` 成为唯一主配置入口
- `bootstrap` 独立承担启动职责
- `admin_modules/user.go` 不再是单文件大后台
- `models/model.go` 不再承载过多全局状态与副作用
- API / WSS / Task / Admin 的边界清晰可描述
