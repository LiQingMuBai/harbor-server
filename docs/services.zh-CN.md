# API / 服务说明

本文档描述 `harbor-server` 当前的服务职责、接口风格、主要模块以及后台接口分组。

## 1. 服务划分

### `api`

用户侧 API 服务。

职责：

- 提供用户登录、注册、资料、认证等接口
- 提供交易、资产、充值提现、消息、矿机、系统配置等业务接口
- 启动内部 RPC 服务供后台广播调用

入口：

- `cmd/api`
- `internal/bootstrap/api`

### `admin`

后台管理服务。

职责：

- 后台登录、权限、角色、菜单
- 用户、钱包、提现、充值审核
- 币种、汇率、公告、规则、矿机、控制器配置
- 通过 Redis 队列消费系统指令并调用 API RPC 客户端

入口：

- `cmd/admin`
- `internal/bootstrap/admin`

### `task`

后台任务服务。

职责：

- 行情/K 线数据处理
- 授权充值任务
- 资产更新、清理委托、矿机收益、贷款统计、汇率更新

入口：

- `cmd/task`
- `internal/bootstrap/task`

### `wss`

WebSocket 推送服务。

职责：

- 用户侧 WebSocket 登录
- 行情、消息、服务消息推送
- 可选管理员 WebSocket 登录

入口：

- `cmd/wss`
- `internal/bootstrap/wss`

### `cdn`

上传与静态资源服务。

职责：

- 接收文件上传
- 暴露 `/static`、`/pdf`、`/whitepaper` 目录

入口：

- `cmd/cdn`
- `internal/bootstrap/cdn`

## 2. HTTP 协议约定

### 返回格式

统一返回结构：

```json
{
  "code": 0,
  "data": {}
}
```

常见状态码：

- `0`：成功
- `10001`：未登录
- `10002`：参数错误

### 参数格式

当前接口不是标准 REST JSON body 风格，而是通过参数 `p` 传入一个 JSON 字符串。

典型结构：

```json
{
  "sid": "session-id",
  "uid": 10001,
  "data": {
    "key": "value"
  }
}
```

说明：

- `sid`：会话标识
- `uid`：用户 ID
- `data`：业务参数对象

### 鉴权方式

用户侧接口中，带 `NeedLogin` 的路由会校验：

- `sid` 对应的会话
- `uid` 与会话用户是否一致

后台接口中，大多数管理路由会经过 `CheckLogin`。

## 3. 用户侧 API 模块

用户 API 主要注册在 `http/modules/*`。

### 用户模块

前缀：`/user/*`

主要能力：

- 登录注册：`/user/login`、`/user/register`、`/user/register/code`
- 用户资料：`/user/userinfo`、`/user/profile`
- 安全能力：`/user/changepass`、`/user/set_cash_password`、`/user/update_cash_password`
- 认证能力：`/user/auth1`、`/user/auth2`
- Google 验证：`/user/google/secret`、`/user/google/bind`、`/user/google/auth_login`
- 手机邮箱绑定：`/user/phone/send`、`/user/phone/bind`、`/user/email/send`、`/user/email/bind`
- 钱包登录与授权：`/user/login/wallet`、`/user/approve`
- 公告：`/user/notice/list`、`/user/notice/detail`、`/user/notice/read`

### 交易模块

前缀：`/trade/*`

主要能力：

- 委托下单：`/trade/delegate`
- 委托列表：`/trade/delegate/list`
- 持仓列表：`/trade/opend/list`
- 平仓列表：`/trade/close/list`

### 资产模块

前缀：`/assets/*`

主要能力：

- 资产兑换：`/assets/exchange`
- 资产列表：`/assets/list`
- 闪兑：`/assets/quickexchange`

### 充值提现模块

前缀：`/credit/*`

主要能力：

- 充值：`/credit/recharge`
- 提现：`/credit/withdraw`
- 充值提现记录：`/credit/recharge/list`、`/credit/withdraw/list`
- 账变日志：`/credit/logs`
- 钱包管理：`/credit/wallet/add`、`/credit/wallet/list`、`/credit/wallet/del`
- 划转：`/credit/wallet/transfer`、`/credit/wallet/transfer/logs`
- 账户转换：`/credit/wallet/exchange`、`/credit/wallet/exchange2`

### 消息模块

前缀：`/message/*`

主要能力：

- 消息列表：`/message/list`
- 未读数量：`/message/unread`
- 已读：`/message/read`
- 客服消息：`/message/service/list`

### 系统模块

前缀：`/system/*`、`/loan/*`、`/coin/*`

主要能力：

- 站点配置：`/system/config`、`/system/config2`
- 公告：`/system/notice`、`/system/notice/detail`
- 币种/汇率：`/system/coinlist`、`/system/currency`
- 规则文案：`/system/rule`、`/system/rule/detail`
- 借贷：`/loan/product`、`/loan/order`、`/loan/order/list`
- 新币/申购：`/coin/buy/list`、`/coin/buy/order`、`/coin/new/list`

### 矿机模块

前缀：`/mining/*`

主要能力：

- 产品列表：`/mining/list`
- 购买：`/mining/buy`
- 解锁：`/mining/unlock`
- 订单列表：`/mining/order/list`
- 明细与统计：`/mining/detail`、`/mining/count`、`/mining/accepts`

## 4. 后台 API 模块

当前后台接口集中在 `admin_modules/user.go` 一个大模块中，路径前缀为 `/admin/*`。

### 认证与权限

- `/admin/login`
- `/admin/logout`
- `/admin/token_info`
- `/admin/admin_list`
- `/admin/add_manage`
- `/admin/role_list`
- `/admin/handler_role`
- `/admin/mean_list`
- `/admin/mean_router`
- `/admin/auth_router`

### 用户与资产

- `/admin/userlist`
- `/admin/userinfo`
- `/admin/op_user`
- `/admin/user_asset`
- `/admin/usercoinlog`
- `/admin/opuserAssetWallet`
- `/admin/save_parentmemo`

### 充值、提现与钱包

- `/admin/recharge_list`
- `/admin/recharge_op`
- `/admin/withdraw_list`
- `/admin/withdraw_op`
- `/admin/save_withdraw`
- `/admin/addr_list`
- `/admin/addr_op`
- `/admin/addr_del`
- `/admin/user_wallet`
- `/admin/del_userwallet`
- `/admin/collect_wallet`
- `/admin/collect_list`

### 认证审核

- `/admin/uauth_list`
- `/admin/uauth_op`
- `/admin/uauth_del`
- `/admin/uauth2_list`
- `/admin/uauth2_op`
- `/admin/uauth2_del`
- `/admin/save_uauth`

### 交易运营

- `/admin/close_trade_list`
- `/admin/hold_trade_list`
- `/admin/delegate_now`
- `/admin/delegate_history`
- `/admin/delegate_del`
- `/admin/manual_delegate_trade`
- `/admin/spot_delegate`
- `/admin/opspot`

### 币种、汇率与系统配置

- `/admin/coin_list`
- `/admin/save_coin`
- `/admin/del_coin`
- `/admin/coindesc_list`
- `/admin/op_coindesc`
- `/admin/currency_list`
- `/admin/op_currency`
- `/admin/del_currency`
- `/admin/setting`
- `/admin/sitecount`
- `/admin/kline_config`

### 公告、消息与通知

- `/admin/notice_list`
- `/admin/notice_op`
- `/admin/notice_del`
- `/admin/msg_list`
- `/admin/send_msg`
- `/admin/send_user_notice`
- `/admin/del_usernotice`
- `/admin/notify_list`
- `/admin/clearNotify`
- `/admin/chatlist`

### 风控、规则与控制器

- `/admin/rulelist`
- `/admin/rulehandler`
- `/admin/del_rule`
- `/admin/controller_list`
- `/admin/del_controller`
- `/admin/kline_controller`
- `/admin/explode_controller`
- `/admin/explode_list`
- `/admin/op_explode`
- `/admin/del_explode`

### 矿机、贷款、代理

- `/admin/minner_open`
- `/admin/minner_list`
- `/admin/minner_op`
- `/admin/minner_del`
- `/admin/loan`
- `/admin/op_loan`
- `/admin/loanorder_list`
- `/admin/audit_loan`
- `/admin/agent_count`
- `/admin/agent_list`
- `/admin/employer_list`

## 5. WebSocket 说明

入口：

- `/wss`

用途：

- 用户登录后接收行情、消息和服务消息
- 管理员可通过 `admin_pass` 进行特殊登录，口令由 `WS_ADMIN_PASS` 控制

说明：

- WSS 不是独立业务协议文档，此处只描述服务角色
- 若要前后端联调，建议后续补一版单独的 WebSocket 消息协议说明

## 6. CDN 说明

上传入口：

- `POST /`

行为：

- 通过 query 参数校验 `sid`
- 表单字段使用 `file`
- 支持类型：`jpg`、`jpeg`、`png`、`gif`、`pdf`
- 返回值中包含最终资源地址

静态路径：

- `/static`
- `/pdf`
- `/whitepaper`

## 7. 当前接口文档边界

本文档目前是“服务说明 + 模块级接口目录”，不是逐字段 API Reference。

后续如果需要继续细化，建议按以下顺序补充：

1. 登录注册流程文档
2. 充值提现流程文档
3. 交易下单与持仓流程文档
4. WebSocket 消息协议文档
5. Admin 权限与菜单文档
