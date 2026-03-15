# Tier0-Edge 重写方案

> 作者: Architecture Analysis  
> 日期: 2026-03-15  
> 状态: 提案

---

## 目录

1. [重写原则：保留什么、砍掉什么、重做什么](#1-重写原则)
2. [整体架构](#2-整体架构)
3. [技术选型](#3-技术选型)
4. [后端设计](#4-后端设计)
5. [数据管道与实时引擎](#5-数据管道与实时引擎)
6. [前端设计](#6-前端设计)
7. [部署架构](#7-部署架构)
8. [开发规范与工程化](#8-开发规范与工程化)
9. [分阶段实施路线](#9-分阶段实施路线)

---

## 1. 重写原则

### 保留（核心价值）

| 保留项 | 原因 |
|--------|------|
| UNS 方法论 | 统一命名空间是正确的 IIoT 架构范式，是项目的核心差异化价值 |
| Go 后端 | 适合高并发消息处理、低内存占用、单二进制部署，是 edge 场景的正确选择 |
| EMQX | 工业级 MQTT broker，Erlang 高可用，已有成熟的 IIoT 生态 |
| TimescaleDB | 基于 PostgreSQL 的时序扩展，同时满足关系和时序查询需求 |
| React 前端 | 生态成熟，组件库丰富，团队已有经验 |
| 接口抽象 | `IPersistentService`、`IDataSinkService` 等接口设计是好的 |
| 事件驱动思想 | 适配器间通过事件通信的松耦合设计正确 |

### 砍掉（过度复杂或错误选择）

| 砍掉项 | 原因 | 替代 |
|--------|------|------|
| Kong API 网关 | 对 edge 场景过重，配置复杂（26 服务、50 路由、7 插件），引入额外 PostgreSQL 依赖 | Go 内置反向代理 |
| Keycloak | 资源消耗大（需要 512MB+ 内存），配置复杂，对边缘部署不友好 | 内置轻量 OAuth2 + 可选外部 OIDC |
| Express BFF | 仅服务 AI/MCP，多了一个服务节点和技术栈 | 合并到 Go 后端，AI 部分用子进程 |
| Konga | Kong 的管理 UI，Kong 砍掉后不再需要 | - |
| Portainer | Docker 管理工具，不是平台核心 | 可选 profile |
| Marimo | Python notebook，不是核心功能 | 可选 profile |
| 自研 Spring 容器 | 双 DI 系统的根源，init() 隐式注册不可预测 | Go 标准 DI（uber/fx） |
| 5 个缓存库 | 碎片化，维护负担 | 统一用一个（otter） |

### 重做（方向对但实现有问题）

| 重做项 | 当前问题 | 新方案 |
|--------|----------|--------|
| 事件总线 | 同步阻塞、无错误隔离 | 基于 channel 的异步事件总线 |
| 数据管道 | WS 阻塞消息消费、磁盘队列无回压 | Go channel pipeline + 背压机制 |
| 数据库迁移 | 无版本化、失败不中断 | golang-migrate |
| 告警系统 | 完全未实现 | 从第一天就作为核心模块 |
| 实时计算 | 完全未实现 | 内置表达式引擎 |
| 前端状态管理 | God Store、无类型安全 | 分域 Store + TanStack Query |
| 前端路由 | 零代码分割 | 全部 lazy loading |
| API 类型安全 | 全部 any | OpenAPI 生成 + 端到端类型 |
| Node-RED 集成 | 深度耦合、单点故障 | 松耦合 + 可选依赖 |

---

## 2. 整体架构

### 2.1 核心设计理念

**"少即是多"** — 当前版本 14 个 Docker 容器，重写目标是核心功能只需 **3 个容器**。

```
最小部署:
┌──────────────────────────────────────────────────────┐
│  Container 1: tier0                                   │
│  (Go binary: API + WebSocket + 反向代理 + 告警引擎)     │
├──────────────────────────────────────────────────────┤
│  Container 2: emqx                                    │
│  (MQTT Broker)                                        │
├──────────────────────────────────────────────────────┤
│  Container 3: timescaledb                             │
│  (PostgreSQL + TimescaleDB: 关系 + 时序 + 认证)         │
└──────────────────────────────────────────────────────┘

可选扩展:
  + node-red       (SourceFlow/EventFlow, 按需启用)
  + grafana        (可视化, 按需启用)
  + minio          (对象存储, 按需启用)
```

### 2.2 架构拓扑

```
                    ┌─────────────────────────────────────────────┐
                    │              tier0 (单 Go 二进制)              │
                    │                                               │
  HTTP/WS          │  ┌──────────┐  ┌──────────┐  ┌─────────────┐ │
  ────────────────►│  │ HTTP      │  │ WebSocket│  │ Reverse     │ │
  Browser/API      │  │ Router    │  │ Hub      │  │ Proxy       │ │
                    │  └────┬─────┘  └────┬─────┘  └──────┬──────┘ │
                    │       │             │               │         │
                    │  ┌────▼─────────────▼───────────────▼──────┐ │
                    │  │              Core Engine                  │ │
                    │  │  ┌──────┐ ┌────────┐ ┌───────┐ ┌──────┐│ │
                    │  │  │ UNS  │ │ Alarm  │ │ Calc  │ │ Auth ││ │
                    │  │  │Mgr   │ │Engine  │ │Engine │ │      ││ │
                    │  │  └──┬───┘ └───┬────┘ └───┬───┘ └──────┘│ │
                    │  │     │         │          │              │ │
                    │  │  ┌──▼─────────▼──────────▼────────────┐│ │
                    │  │  │         Event Bus (async)           ││ │
                    │  │  └──┬───────────────┬─────────────┬───┘│ │
                    │  └─────┼───────────────┼─────────────┼────┘ │
                    │        │               │             │       │
                    │  ┌─────▼─────┐  ┌──────▼──────┐ ┌───▼────┐ │
                    │  │ MQTT      │  │ TimescaleDB │ │ WS     │ │
                    │  │ Consumer  │  │ + PostgreSQL│ │ Push   │ │
                    │  └─────┬─────┘  └──────┬──────┘ └───┬────┘ │
                    └────────┼───────────────┼────────────┼───────┘
                             │               │            │
                     ┌───────▼───────┐ ┌─────▼─────┐     │
                     │    EMQX       │ │TimescaleDB │     ▼
                     │  MQTT Broker  │ │  Database  │  Browsers
                     └───────────────┘ └───────────┘
```

### 2.3 设计原则

1. **单二进制**: 一个 Go binary 包含 API Server、WebSocket Hub、反向代理、告警引擎、静态文件服务。前端构建产物嵌入 Go 二进制中，真正的单文件部署。

2. **数据库合一**: TimescaleDB 本身就是 PostgreSQL 扩展，不需要两个独立的数据库实例。一个 TimescaleDB 同时存时序数据（hypertable）和关系数据（普通表），以及用户认证数据。

3. **最少外部依赖**: 核心运行只需 EMQX 和 TimescaleDB。Node-RED、Grafana 等按需启用，缺少它们系统照常工作。

4. **类型安全端到端**: Go struct → OpenAPI spec → TypeScript types，一次定义，处处使用。

---

## 3. 技术选型

### 3.1 后端

| 领域 | 选择 | 理由 |
|------|------|------|
| 语言 | Go 1.24+ | 保持不变，高并发、低内存、单二进制 |
| HTTP 框架 | `net/http` + `chi` router | 不用 go-zero。chi 是标准 `net/http` 兼容的轻量路由器，无供应商锁定，中间件生态丰富 |
| DI | `uber-go/fx` | 成熟的 DI 框架，显式模块声明，可测试，启动时验证完整依赖图 |
| ORM | `sqlc` + 手写 SQL | 替代 GORM。sqlc 从 SQL 生成类型安全的 Go 代码，性能好，可控性强。复杂查询用手写 SQL。对于时序数据和复杂聚合查询，ORM 是错误的抽象层 |
| 迁移 | `golang-migrate` | 版本化迁移，支持 up/down，CLI + 库两种用法 |
| 验证 | `go-playground/validator` | 结构体标签验证，覆盖所有 DTO 输入 |
| 日志 | `slog` (Go 标准库) | Go 1.21+ 内置结构化日志，不需要第三方 |
| 缓存 | `maypok86/otter` | 单一缓存库，高性能并发安全，支持 TTL |
| MQTT | `eclipse/paho.golang` v2 | Paho v2 支持 MQTT 5.0，连接池，自动重连 |
| 告警引擎 | `expr-lang/expr` | Go 表达式引擎，编译后执行快，支持自定义函数，适合告警规则和实时计算 |
| WebSocket | `coder/websocket` | 基于 `net/http` 的高性能 WebSocket 库 |
| API 文档 | `swaggo/swag` → OpenAPI 3.0 | 从注释生成 OpenAPI spec，同时生成前端类型 |
| 嵌入前端 | `embed.FS` | Go 1.16+ 原生支持，将前端构建产物嵌入二进制 |

**为什么不用 go-zero?**

go-zero 是优秀的微服务框架，但对 Tier0-Edge 的 edge 部署场景来说:
- 它引入了 etcd、Redis 等微服务基础设施依赖
- 代码生成（goctl）产生大量样板代码（当前 198 个 handler 文件、269 个 logic 文件）
- `ServiceContext` 的依赖注入方式与 Spring 容器冲突
- 对于单体+边缘部署，`chi` + `fx` 更轻量灵活

### 3.2 前端

| 领域 | 选择 | 理由 |
|------|------|------|
| 框架 | React 19 + Vite 7 | 保持不变，但升级 React 19（useOptimistic、Actions） |
| 路由 | TanStack Router | 文件系统路由、类型安全参数、内置代码分割、比 react-router 更好的类型 |
| 数据获取 | TanStack Query v5 | 替代手动 axios + useEffect。缓存、去重、后台刷新、乐观更新 |
| 状态管理 | Zustand v5 (分域) | 保持 Zustand，但拆为 5+ 个独立 store |
| UI | Ant Design 5 | 保持不变，组件库成熟 |
| 表单 | React Hook Form + Zod | 替代 antd Form 的重型方案，性能更好，校验逻辑可前后端复用 |
| 类型 | OpenAPI → `openapi-typescript` | 从后端 OpenAPI spec 自动生成 TypeScript 类型 |
| 图表 | AntV G2/X6 | 保持不变 |
| i18n | `react-intl` | 保持不变 |
| 构建 | pnpm + Turborepo | 保持不变 |
| 实时 | 原生 WebSocket + TanStack Query 整合 | WebSocket 消息更新 Query cache |

### 3.3 基础设施

| 领域 | 选择 | 理由 |
|------|------|------|
| MQTT Broker | EMQX 5.8 | 保持不变 |
| 数据库 | TimescaleDB (= PostgreSQL + 时序扩展) | 合并为一个实例 |
| 容器编排 | Docker Compose | 保持不变，edge 场景不需要 K8s |
| API 网关 | 去掉，Go 内置 | Go binary 自带反向代理功能，转发 Node-RED、Grafana 等 |
| 认证 | Go 内置 + OIDC 协议 | 内置本地认证，可选对接外部 OIDC Provider（Keycloak、Authentik 等） |
| CI/CD | GitHub Actions | 从第一天就有 |

---

## 4. 后端设计

### 4.1 项目结构

```
tier0/
├── cmd/
│   └── tier0/
│       └── main.go              # 入口: fx.New() 启动
├── internal/
│   ├── domain/                   # 领域模型（纯 Go struct，无外部依赖）
│   │   ├── namespace.go          # UNS 命名空间实体
│   │   ├── alarm.go              # 告警规则 & 告警事件
│   │   ├── user.go               # 用户、角色
│   │   ├── flow.go               # Source/Event Flow
│   │   └── dashboard.go          # 仪表板
│   │
│   ├── port/                     # 端口（接口定义）
│   │   ├── inbound/              # 入站端口
│   │   │   ├── namespace_service.go    # UNS 业务接口
│   │   │   ├── alarm_service.go        # 告警业务接口
│   │   │   ├── calc_service.go         # 计算引擎接口
│   │   │   └── auth_service.go         # 认证接口
│   │   └── outbound/             # 出站端口
│   │       ├── namespace_repo.go       # UNS 数据访问
│   │       ├── timeseries_repo.go      # 时序数据访问
│   │       ├── mqtt_publisher.go       # MQTT 发布
│   │       ├── websocket_hub.go        # WS 推送
│   │       └── nodered_client.go       # Node-RED（可选）
│   │
│   ├── app/                      # 应用服务（编排层，实现 inbound port）
│   │   ├── namespace_app.go
│   │   ├── alarm_app.go
│   │   ├── calc_app.go
│   │   ├── pipeline_app.go       # MQTT → 处理 → 持久化 管道编排
│   │   └── auth_app.go
│   │
│   ├── adapter/                  # 适配器（实现 outbound port）
│   │   ├── postgres/             # PostgreSQL/TimescaleDB
│   │   │   ├── namespace_repo.go
│   │   │   ├── timeseries_repo.go
│   │   │   ├── alarm_repo.go
│   │   │   ├── user_repo.go
│   │   │   ├── migrations/       # SQL 迁移文件 (001_init.up.sql ...)
│   │   │   └── queries/          # sqlc 查询定义
│   │   │       ├── namespace.sql
│   │   │       ├── alarm.sql
│   │   │       └── sqlc.yaml
│   │   ├── mqtt/
│   │   │   ├── consumer.go       # MQTT 共享订阅消费
│   │   │   └── publisher.go      # MQTT 发布
│   │   ├── ws/
│   │   │   └── hub.go            # WebSocket 连接管理 & 推送
│   │   ├── nodered/
│   │   │   └── client.go         # Node-RED HTTP API 客户端
│   │   └── proxy/
│   │       └── reverse_proxy.go  # 反向代理（Grafana, Node-RED, MinIO）
│   │
│   ├── handler/                  # HTTP Handler（入口层）
│   │   ├── namespace_handler.go
│   │   ├── alarm_handler.go
│   │   ├── flow_handler.go
│   │   ├── auth_handler.go
│   │   ├── ws_handler.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── recovery.go
│   │   └── router.go             # 路由注册
│   │
│   ├── engine/                   # 核心引擎
│   │   ├── alarm/
│   │   │   ├── engine.go         # 告警规则引擎
│   │   │   ├── rule.go           # 规则编译与评估
│   │   │   └── notifier.go       # 告警通知（WS + 持久化 + webhook）
│   │   ├── calc/
│   │   │   ├── engine.go         # 实时计算引擎
│   │   │   └── functions.go      # 内置函数（avg, max, min, 滑动窗口）
│   │   └── pipeline/
│   │       ├── pipeline.go       # 消息处理管道
│   │       ├── stage_parse.go    # Stage 1: JSON 解析
│   │       ├── stage_validate.go # Stage 2: 字段校验
│   │       ├── stage_calc.go     # Stage 3: 实时计算
│   │       ├── stage_alarm.go    # Stage 4: 告警检测
│   │       ├── stage_ws.go       # Stage 5: WS 推送（异步）
│   │       └── stage_sink.go     # Stage 6: 持久化（批量）
│   │
│   └── config/
│       └── config.go             # 配置结构（从 env / yaml 加载）
│
├── web/                          # 前端构建产物（嵌入到 Go 二进制）
│   └── dist/                     # vite build output
│
├── migrations/                   # 数据库迁移（golang-migrate 格式）
│   ├── 001_init_schema.up.sql
│   ├── 001_init_schema.down.sql
│   ├── 002_uns_namespace.up.sql
│   ├── 003_timeseries.up.sql
│   ├── 004_alarm.up.sql
│   └── 005_auth.up.sql
│
├── api/
│   └── openapi.yaml              # OpenAPI 3.0 spec（生成来源）
│
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

### 4.2 依赖注入 (uber/fx)

```go
// cmd/tier0/main.go
func main() {
    fx.New(
        // 配置
        config.Module,

        // 适配器层
        postgres.Module,     // PostgreSQL + TimescaleDB repos
        mqtt.Module,         // MQTT consumer + publisher
        ws.Module,           // WebSocket hub
        nodered.Module,      // Node-RED client (可选)
        proxy.Module,        // 反向代理

        // 引擎层
        pipeline.Module,     // 消息处理管道
        alarm.Module,        // 告警引擎
        calc.Module,         // 计算引擎

        // 应用层
        app.Module,          // 业务服务

        // 入口层
        handler.Module,      // HTTP handlers + router

        // 启动
        fx.Invoke(startServer),
    ).Run()
}
```

**与当前方案的对比:**

| 维度 | 当前 (go-zero + Spring) | 重写 (chi + fx) |
|------|------------------------|-----------------|
| 依赖声明 | `init()` 隐式注册，运行时才知道缺失 | 显式模块声明，启动时验证完整依赖图 |
| 可测试性 | 需要 mock 全局 Spring 容器 | `fx.Replace` / `fx.Decorate` 替换任意依赖 |
| 生命周期 | `RefreshBeanContext()` 手动调用 | fx 自动管理 `OnStart` / `OnStop` |
| 错误处理 | `GetBean` panic | fx 启动失败时返回错误，不 panic |
| 复杂度 | 双 DI 系统交织 | 单一、一致的 DI 机制 |

### 4.3 六边形架构 (Hexagonal Architecture)

```
                              ┌─────────────────┐
                              │    Handler       │
                              │  (HTTP/WS 入口)   │
                              └───────┬─────────┘
                                      │
                              ┌───────▼─────────┐
                   ┌──────────│   Inbound Port   │──────────┐
                   │          │   (Service 接口)   │          │
                   │          └───────┬─────────┘          │
                   │                  │                     │
          ┌────────▼──────┐  ┌───────▼─────────┐  ┌───────▼───────┐
          │  App Service  │  │  App Service    │  │  App Service  │
          │ (Namespace)   │  │  (Alarm)        │  │  (Pipeline)   │
          └────────┬──────┘  └───────┬─────────┘  └───────┬───────┘
                   │                 │                     │
          ┌────────▼─────────────────▼─────────────────────▼──────┐
          │                    Domain Model                        │
          │            (纯 Go struct，无外部依赖)                    │
          └────────┬─────────────────┬─────────────────────┬──────┘
                   │                 │                     │
          ┌────────▼──────┐  ┌──────▼──────┐  ┌───────────▼─────┐
          │ Outbound Port │  │Outbound Port│  │  Outbound Port  │
          │  (Repo 接口)   │  │ (MQTT 接口)  │  │  (WS 接口)      │
          └────────┬──────┘  └──────┬──────┘  └───────────┬─────┘
                   │                │                     │
          ┌────────▼──────┐  ┌──────▼──────┐  ┌───────────▼─────┐
          │   Adapter     │  │   Adapter   │  │    Adapter      │
          │  (PostgreSQL) │  │   (EMQX)    │  │   (WebSocket)   │
          └───────────────┘  └─────────────┘  └─────────────────┘
```

**关键约束:**
- `domain/` 不导入任何外部包，是纯粹的业务逻辑
- `port/` 只定义接口，不包含实现
- `app/` 通过接口编排业务，不直接依赖具体适配器
- `adapter/` 实现 `port/outbound/` 中的接口
- `handler/` 调用 `port/inbound/` 中的接口

这确保了 **任何适配器都可以被替换**。例如:
- 测试时用 SQLite 替代 TimescaleDB
- 未来可以换掉 EMQX（换成 NanoMQ 或 HiveMQ）
- 可以用 Kafka 替代 MQTT 作为消息传输

### 4.4 认证设计: 从重到轻

**当前**: Kong (Lua 插件做认证) → Go Backend → Keycloak → PostgreSQL
**总计 3 个组件、2 次网络跳转**完成一次认证。

**重写后**:

```go
// 模式 1: 本地认证（默认，零外部依赖）
// 用户数据存在 TimescaleDB 的普通表中
// JWT 由 Go 签发和验证，密钥存数据库

type LocalAuthProvider struct {
    userRepo   port.UserRepository
    jwtSigner  *jwt.Signer
}

func (p *LocalAuthProvider) Login(email, password string) (*Token, error) { ... }
func (p *LocalAuthProvider) Verify(token string) (*Claims, error) { ... }

// 模式 2: 外部 OIDC（可选，对接 Keycloak / Authentik / Azure AD）
// 配置 OIDC discovery URL 即可

type OIDCAuthProvider struct {
    discoveryURL string
    clientID     string
    verifier     *oidc.IDTokenVerifier
}

func (p *OIDCAuthProvider) Verify(token string) (*Claims, error) { ... }

// 在配置中选择:
// auth:
//   provider: local            # 或 oidc
//   oidc:
//     discovery_url: https://keycloak.example.com/realms/tier0
//     client_id: tier0
```

**好处:**
- 最小部署不需要 Keycloak（-1 容器，-512MB 内存）
- 本地认证开箱即用，0 配置
- 企业用户可对接已有的身份系统

### 4.5 反向代理: 替代 Kong

```go
// internal/adapter/proxy/reverse_proxy.go

type ProxyConfig struct {
    Routes []ProxyRoute `yaml:"routes"`
}

type ProxyRoute struct {
    PathPrefix string `yaml:"path_prefix"`  // e.g. "/nodered"
    Target     string `yaml:"target"`       // e.g. "http://nodered:1880"
    StripPrefix bool  `yaml:"strip_prefix"` // 是否移除前缀
    Auth        bool  `yaml:"auth"`         // 是否需要认证
}

// proxy:
//   routes:
//     - path_prefix: /nodered
//       target: http://nodered:1880
//       strip_prefix: true
//       auth: true
//     - path_prefix: /grafana
//       target: http://grafana:3000
//       strip_prefix: true
//       auth: false
```

Go 标准库的 `httputil.ReverseProxy` 完全满足需求，不需要一个独立的 Kong 实例。

---

## 5. 数据管道与实时引擎

### 5.1 管道架构

这是重写中**最核心的改动**。当前是一条同步线性管道；重写为异步多阶段管道:

```
MQTT Message
     │
     ▼
┌────────────────┐
│  Parse Stage   │  JSON → map[string]string
│  (goroutine)   │  topic → UnsDefinition (LRU cached)
└───────┬────────┘
        │ channel (buffered, 背压)
        ▼
┌────────────────┐
│ Validate Stage │  类型检查、范围检查、质量码
│  (goroutine)   │
└───────┬────────┘
        │ fan-out (3 个 channel)
        ├─────────────────────────────────────────┐
        │                                          │
        ▼                                          ▼
┌────────────────┐  ┌────────────────┐  ┌────────────────┐
│  Calc Stage    │  │  Alarm Stage   │  │  WS Push Stage │
│  (goroutine)   │  │  (goroutine)   │  │  (goroutine)   │
│  实时计算       │  │  规则评估       │  │  异步推送       │
└───────┬────────┘  └───────┬────────┘  └────────────────┘
        │                   │
        ▼                   ▼
┌────────────────┐  ┌────────────────┐
│  Sink Stage    │  │ Alarm Notify   │
│  (goroutine)   │  │  (goroutine)   │
│  批量写入 TSDB  │  │  WS + Webhook  │
└────────────────┘  └────────────────┘
```

**关键设计:**

```go
// internal/engine/pipeline/pipeline.go

type Pipeline struct {
    parseIn    chan MQTTMessage       // MQTT → Parse
    validateIn chan ParsedMessage     // Parse → Validate
    calcIn     chan ValidatedMessage  // Validate → Calc
    alarmIn    chan ValidatedMessage  // Validate → Alarm (fan-out)
    wsIn       chan ValidatedMessage  // Validate → WS Push (fan-out)
    sinkIn     chan SinkMessage       // Calc → Sink (批量)
}

type PipelineConfig struct {
    ParseWorkers    int  // 解析 worker 数
    ValidateWorkers int  // 校验 worker 数
    CalcWorkers     int  // 计算 worker 数
    AlarmWorkers    int  // 告警 worker 数
    WSWorkers       int  // WS 推送 worker 数
    SinkWorkers     int  // 持久化 worker 数
    SinkBatchSize   int  // 批量写入大小
    SinkFlushInterval time.Duration  // 最大等待时间
    ChannelBuffer   int  // channel 缓冲区大小
}
```

**与当前方案的对比:**

| 维度 | 当前 | 重写 |
|------|------|------|
| WS 推送 | 同步，阻塞消息处理 | 异步，独立 goroutine |
| 告警检测 | 不存在 | 管道内置阶段 |
| 实时计算 | 空实现 | 管道内置阶段 |
| 持久化 | 磁盘队列 → DB | Channel → 批量 DB（省去磁盘 IO） |
| 背压 | 无，磁盘满则丢数据 | Channel 满则阻塞上游，MQTT QoS 1 保证不丢 |
| 并发 | 单 goroutine 串行 | 每阶段可配置 worker 数 |

### 5.2 告警引擎

```go
// internal/engine/alarm/engine.go

type AlarmEngine struct {
    rules    map[int64]*CompiledRule  // unsId → 编译后的告警规则
    repo     port.AlarmRepository
    notifier AlarmNotifier
}

type AlarmRule struct {
    ID          int64
    NamespaceID int64
    Name        string
    Expression  string     // "temperature > 80 && humidity < 20"
    Severity    string     // critical, warning, info
    Duration    time.Duration  // 持续多久才触发（防抖）
    Cooldown    time.Duration  // 两次告警之间的最小间隔
    Actions     []Action       // webhook, email, ws_notify
}

// 使用 expr-lang/expr 编译和执行规则
type CompiledRule struct {
    Rule     AlarmRule
    Program  *vm.Program        // 编译后的表达式
    State    *AlarmState        // 当前状态（触发中/已恢复/静默中）
}

// 告警评估 — 在管道的 Alarm Stage 中调用
func (e *AlarmEngine) Evaluate(nsID int64, data map[string]any) []AlarmEvent {
    rules := e.getRules(nsID)
    var events []AlarmEvent
    for _, rule := range rules {
        result, _ := expr.Run(rule.Program, data)
        if triggered, ok := result.(bool); ok && triggered {
            if event := rule.State.Trigger(data); event != nil {
                events = append(events, *event)
            }
        } else {
            if event := rule.State.Recover(); event != nil {
                events = append(events, *event)
            }
        }
    }
    return events
}
```

### 5.3 实时计算引擎

```go
// internal/engine/calc/engine.go

type CalcEngine struct {
    definitions map[int64]*CalcDefinition
    windows     map[int64]*SlidingWindow  // 滑动窗口数据
}

type CalcDefinition struct {
    SourceID   int64     // 源数据的命名空间 ID
    TargetID   int64     // 计算结果写入的命名空间 ID
    Expression string    // "avg(source.temperature, '5m')"
    Interval   time.Duration
}

// 内置函数
// avg(field, window)  — 滑动窗口平均
// max(field, window)  — 滑动窗口最大
// min(field, window)  — 滑动窗口最小
// rate(field, window) — 变化率
// delta(field)        — 与上一个值的差
// abs(value)          — 绝对值
// clamp(value, min, max) — 限幅
```

---

## 6. 前端设计

### 6.1 项目结构

```
frontend/
├── apps/
│   └── web/
│       ├── src/
│       │   ├── routes/              # TanStack Router 文件系统路由
│       │   │   ├── __root.tsx       # 根布局
│       │   │   ├── index.tsx        # 首页 (/)
│       │   │   ├── uns/
│       │   │   │   ├── index.tsx    # UNS 页面
│       │   │   │   └── $nsId.tsx    # 命名空间详情
│       │   │   ├── alarms/
│       │   │   │   ├── index.tsx
│       │   │   │   └── rules.tsx
│       │   │   ├── flows/
│       │   │   │   ├── source.tsx
│       │   │   │   └── event.tsx
│       │   │   ├── dashboards/
│       │   │   │   └── index.tsx
│       │   │   └── settings/
│       │   │       ├── users.tsx
│       │   │       └── system.tsx
│       │   │
│       │   ├── features/            # 按功能域组织的模块
│       │   │   ├── uns/
│       │   │   │   ├── components/  # UNS 特有组件
│       │   │   │   │   ├── UnsTree.tsx          # <300 行
│       │   │   │   │   ├── UnsTreeNode.tsx
│       │   │   │   │   ├── UnsTreeToolbar.tsx
│       │   │   │   │   ├── TopicDetail.tsx
│       │   │   │   │   ├── NamespaceForm.tsx
│       │   │   │   │   └── FieldEditor.tsx
│       │   │   │   ├── hooks/       # UNS 业务 hooks
│       │   │   │   │   ├── useUnsTree.ts
│       │   │   │   │   ├── useNamespace.ts
│       │   │   │   │   └── useRealtimeData.ts
│       │   │   │   ├── api.ts       # UNS API (类型安全)
│       │   │   │   └── types.ts     # UNS 类型 (从 OpenAPI 生成)
│       │   │   │
│       │   │   ├── alarm/
│       │   │   │   ├── components/
│       │   │   │   ├── hooks/
│       │   │   │   ├── api.ts
│       │   │   │   └── types.ts
│       │   │   │
│       │   │   ├── flow/
│       │   │   └── auth/
│       │   │
│       │   ├── shared/              # 共享组件和工具
│       │   │   ├── components/      # 通用 UI 组件
│       │   │   ├── hooks/           # 通用 hooks
│       │   │   ├── stores/          # 全局 store（拆分后）
│       │   │   │   ├── auth.ts      # ~50 行: 用户信息、登录状态
│       │   │   │   ├── system.ts    # ~50 行: 系统配置、版本
│       │   │   │   ├── menu.ts      # ~60 行: 菜单树
│       │   │   │   ├── permission.ts# ~40 行: 权限映射
│       │   │   │   └── theme.ts     # ~40 行: 主题
│       │   │   └── utils/
│       │   │
│       │   ├── generated/           # 自动生成（不手动编辑）
│       │   │   └── api-types.ts     # 从 OpenAPI spec 生成
│       │   │
│       │   └── main.tsx
│       │
│       ├── vite.config.ts
│       └── package.json
│
├── packages/
│   ├── ui/                          # 共享 UI 组件库
│   └── tsconfig/
│
├── pnpm-workspace.yaml
└── turbo.json
```

### 6.2 状态管理策略

**三层数据管理:**

```
┌─────────────────────────────────────────────────────┐
│  Layer 1: TanStack Query (服务端状态)                  │
│  — API 数据的获取、缓存、刷新、乐观更新                  │
│  — 替代所有 useEffect + axios 模式                     │
│  — 自动处理 loading/error 状态                         │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│  Layer 2: Zustand (客户端全局状态)                      │
│  — 只存纯客户端状态：主题、语言、侧边栏展开等             │
│  — 5 个小 store，每个 < 60 行                          │
│  — 不存任何服务端数据（那是 TanStack Query 的事）        │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│  Layer 3: WebSocket + Query Invalidation (实时数据)    │
│  — WebSocket 收到消息后 invalidate 相关 Query          │
│  — 或直接通过 setQueryData 更新缓存                     │
│  — 无需额外 store                                     │
└─────────────────────────────────────────────────────┘
```

**具体实现:**

```typescript
// features/uns/hooks/useRealtimeData.ts

function useRealtimeData(namespaceId: number) {
  const queryClient = useQueryClient();

  useEffect(() => {
    const ws = new WebSocket(`/api/ws?id=${namespaceId}`);
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // 直接更新 Query 缓存，触发组件重渲染
      queryClient.setQueryData(
        ['namespace', namespaceId, 'realtime'],
        data
      );
    };
    return () => ws.close();
  }, [namespaceId]);

  return useQuery({
    queryKey: ['namespace', namespaceId, 'realtime'],
    queryFn: () => api.getLatestData(namespaceId),
    staleTime: Infinity,  // WebSocket 负责更新，不需要轮询
  });
}
```

```typescript
// features/uns/api.ts — 类型安全，从 OpenAPI 生成

import type { paths } from '@/generated/api-types';
import { queryOptions } from '@tanstack/react-query';

type NamespaceTree = paths['/api/uns/tree']['get']['responses']['200']['content']['application/json'];
type CreateNamespace = paths['/api/uns/namespace']['post']['requestBody']['content']['application/json'];

export const unsQueries = {
  tree: () => queryOptions({
    queryKey: ['uns', 'tree'],
    queryFn: () => fetch('/api/uns/tree').then(r => r.json()) as Promise<NamespaceTree>,
  }),

  detail: (id: number) => queryOptions({
    queryKey: ['uns', 'detail', id],
    queryFn: () => fetch(`/api/uns/namespace/${id}`).then(r => r.json()),
  }),
};

// 在组件中使用
function UnsTree() {
  const { data: tree, isLoading } = useQuery(unsQueries.tree());
  // tree 类型自动推导，IDE 有完整补全
}
```

### 6.3 组件设计规范

**规则: 单个组件文件不超过 300 行。**

以 UNS 树组件为例，当前 1082 行的 `uns-tree/index.tsx` 拆分为:

```
当前: uns-tree/index.tsx (1082 行)
  ↓ 拆分为:
UnsTree.tsx (120 行)
├── 组件骨架、布局、Context Provider
│
├── UnsTreeToolbar.tsx (80 行)
│   └── 搜索框、过滤按钮、新建按钮
│
├── UnsTreeList.tsx (150 行)
│   └── 虚拟化列表渲染（大量节点时性能关键）
│
├── UnsTreeNode.tsx (100 行)
│   └── 单个节点的渲染（图标、名称、类型标签）
│
├── UnsTreeContextMenu.tsx (80 行)
│   └── 右键菜单（创建子节点、重命名、删除）
│
├── useUnsTree.ts (100 行)
│   └── 树的展开/折叠/选中/搜索/拖拽逻辑
│
└── useUnsTreeDnd.ts (80 行)
    └── 拖拽排序逻辑
```

### 6.4 代码分割

TanStack Router 内置代码分割:

```typescript
// routes/uns/index.tsx
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/uns/')({
  // 自动代码分割：此路由的组件和 loader 只在访问 /uns 时加载
  component: () => import('@/features/uns/UnsPage'),
  loader: ({ context }) => context.queryClient.ensureQueryData(unsQueries.tree()),
});
```

---

## 7. 部署架构

### 7.1 最小部署 (3 容器)

```yaml
# docker-compose.yml
services:
  tier0:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://tier0:${DB_PASSWORD}@db:5432/tier0
      - MQTT_URL=tcp://emqx:1883
      - AUTH_PROVIDER=local                    # 或 oidc
      - AUTH_JWT_SECRET=${JWT_SECRET}
    depends_on:
      db:
        condition: service_healthy
      emqx:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/health"]
      interval: 10s
      timeout: 5s
      retries: 3
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '2.0'
    logging:
      driver: json-file
      options:
        max-size: "50m"
        max-file: "3"

  emqx:
    image: emqx/emqx:5.8
    environment:
      - EMQX_LISTENERS__TCP__DEFAULT__BIND=0.0.0.0:1883
    healthcheck:
      test: ["CMD", "emqx", "ctl", "status"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          memory: 512M

  db:
    image: timescale/timescaledb:2.20.0-pg17
    environment:
      - POSTGRES_DB=tier0
      - POSTGRES_USER=tier0
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U tier0"]
      interval: 5s
      timeout: 3s
      retries: 5
    deploy:
      resources:
        limits:
          memory: 1G

volumes:
  db_data:
```

**对比当前部署:**

| 维度 | 当前 (14 容器) | 重写 (3 容器) |
|------|---------------|--------------|
| 最小内存 | ~4GB+ | ~2GB |
| 容器数 | 14 | 3 |
| 端口暴露 | 15+ | 1 (只暴露 tier0 的 3000) |
| 配置文件 | .env.default (60+ 变量) | .env (5 变量) |
| 首次启动时间 | 3-5 分钟 | 30 秒 |

### 7.2 扩展部署

```yaml
# docker-compose.override.yml — 按需添加
services:
  nodered:
    image: nodered/node-red:4.0
    profiles: ["flows"]

  grafana:
    image: grafana/grafana:11.5
    profiles: ["monitoring"]

  minio:
    image: minio/minio
    profiles: ["storage"]
```

```bash
# 最小部署
docker compose up -d

# 带 Node-RED
docker compose --profile flows up -d

# 全功能
docker compose --profile flows --profile monitoring --profile storage up -d
```

### 7.3 Dockerfile

```dockerfile
# === 前端构建 ===
FROM node:22-alpine AS frontend
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml frontend/package.json ./
RUN pnpm fetch
COPY frontend/ .
RUN pnpm install --offline && pnpm build

# === 后端构建 ===
FROM golang:1.24-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/apps/web/dist ./web/dist
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /tier0 ./cmd/tier0

# === 运行 ===
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata wget
RUN adduser -D -u 1000 tier0
USER tier0
COPY --from=backend /tier0 /usr/local/bin/tier0
COPY migrations/ /migrations/
EXPOSE 3000
ENTRYPOINT ["tier0"]
```

**特点:**
- 前端嵌入 Go 二进制，单文件部署
- 非 root 用户运行
- 最终镜像 ~30MB（Alpine + 静态 Go binary）
- 无中国镜像源硬编码（CI 中通过 build-arg 按需设置）

---

## 8. 开发规范与工程化

### 8.1 从第一天就有的 CI/CD

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: timescale/timescaledb:2.20.0-pg17
        env:
          POSTGRES_PASSWORD: test
        options: --health-cmd pg_isready
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go vet ./...
      - run: golangci-lint run
      - run: go test -race -coverprofile=coverage.out ./...
      - run: go build ./cmd/tier0

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: pnpm
          cache-dependency-path: frontend/pnpm-lock.yaml
      - run: pnpm install --frozen-lockfile
        working-directory: frontend
      - run: pnpm lint
        working-directory: frontend
      - run: pnpm typecheck
        working-directory: frontend
      - run: pnpm test
        working-directory: frontend
      - run: pnpm build
        working-directory: frontend

  docker:
    needs: [backend, frontend]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: docker/build-push-action@v5
        with:
          push: true
          tags: tier0/tier0:${{ github.sha }}
```

### 8.2 API 类型同步流水线

```makefile
# Makefile

# 1. 从 Go 注释生成 OpenAPI spec
api-spec:
	swag init -g cmd/tier0/main.go -o api/

# 2. 从 OpenAPI spec 生成 TypeScript 类型
api-types: api-spec
	npx openapi-typescript api/swagger.json -o frontend/apps/web/src/generated/api-types.ts

# 3. 从 SQL 生成 Go 查询代码
db-codegen:
	cd internal/adapter/postgres && sqlc generate

# 4. 一键全部生成
generate: db-codegen api-spec api-types
```

**端到端类型安全链路:**
```
SQL (migrations/*.sql)
     ↓ sqlc generate
Go structs + queries (adapter/postgres/queries/*.go)
     ↓ swag init
OpenAPI 3.0 spec (api/swagger.json)
     ↓ openapi-typescript
TypeScript types (generated/api-types.ts)
     ↓ TanStack Query
React 组件 (完整类型推导)
```

**修改一个字段，从数据库到前端的类型全部自动更新。**

### 8.3 测试策略

```
后端:
├── 单元测试 (go test)
│   ├── domain/ — 纯业务逻辑，无需 mock
│   ├── app/ — mock outbound port 接口
│   └── engine/ — 告警规则评估、计算引擎
├── 集成测试 (testcontainers)
│   ├── adapter/postgres/ — 真实 TimescaleDB
│   └── adapter/mqtt/ — 真实 EMQX
└── API 测试 (httptest)
    └── handler/ — 完整 HTTP 请求/响应

前端:
├── 单元测试 (vitest)
│   ├── features/*/hooks/ — 业务 hooks
│   └── shared/utils/ — 工具函数
├── 组件测试 (vitest + testing-library)
│   └── features/*/components/ — 核心组件
└── E2E 测试 (playwright)
    └── 关键用户流程（登录→创建命名空间→查看数据）
```

### 8.4 代码规范

```
后端:
- golangci-lint (exhaustive, govet, errcheck, staticcheck)
- 每个 PR 必须通过 lint + test
- 组件不超过 300 行

前端:
- ESLint: no-console (error), no-explicit-any (error), react-hooks/exhaustive-deps (error)
- 组件不超过 300 行
- 必须使用 TanStack Query 获取数据（禁止裸 useEffect + fetch）
- 每个 feature 必须有 hooks/、components/、api.ts、types.ts
```

---

## 9. 分阶段实施路线

### Phase 1: 地基 (4 周)

```
Week 1-2: 后端骨架
├── 项目结构搭建（六边形架构）
├── fx DI 配置
├── chi 路由 + 中间件（auth, logging, recovery）
├── golang-migrate 数据库迁移
├── sqlc 代码生成
├── 本地认证（JWT 签发/验证 + 用户 CRUD）
└── CI/CD（GitHub Actions: lint + test + build）

Week 3-4: 前端骨架
├── TanStack Router 文件系统路由
├── TanStack Query 配置
├── OpenAPI → TypeScript 类型生成
├── Zustand 分域 store
├── Ant Design 主题 + 布局
└── 认证流程（登录/登出/JWT 刷新）
```

**交付物:** 能登录、看到空 dashboard 的可运行系统。

### Phase 2: UNS 核心 (4 周)

```
Week 5-6: UNS 数据模型 + CRUD
├── 命名空间树（创建/查询/更新/删除/移动）
├── TimescaleDB 时序表管理（hypertable, 视图）
├── 前端 UnsTree 组件 + TopicDetail 组件
└── UNS 导入/导出

Week 7-8: 数据管道
├── MQTT 消费者（EMQX 共享订阅）
├── 6 阶段异步管道（parse → validate → calc → alarm → ws → sink）
├── WebSocket 实时数据推送
├── 前端 useRealtimeData hook
└── 批量持久化（TimescaleDB）
```

**交付物:** 能创建 UNS 命名空间、接收 MQTT 数据、实时查看、持久化存储。

### Phase 3: 告警与计算 (3 周)

```
Week 9-10: 告警引擎
├── 告警规则 CRUD（后端 + 前端）
├── expr-lang 表达式编译和评估
├── 告警触发、恢复、静默
├── 告警通知（WebSocket + Webhook）
└── 告警历史查询

Week 11: 实时计算
├── 计算字段定义
├── 滑动窗口聚合（avg, max, min）
├── 引用字段解析
└── 计算结果回写管道
```

**交付物:** 完整的告警和实时计算功能。

### Phase 4: 集成与扩展 (3 周)

```
Week 12: Node-RED 集成
├── SourceFlow（设备接入流）
├── EventFlow（事件处理流）
├── 流的 CRUD + 自动创建
└── 反向代理（/nodered → Node-RED）

Week 13: 仪表板与可视化
├── Grafana 集成（反向代理 + 数据源配置）
├── 内置简易仪表板（基于 AntV G2）
└── 拓扑图

Week 14: 收尾
├── OIDC 外部认证对接
├── 数据备份/恢复
├── 性能测试 + 调优
├── 文档完善
└── Release v1.0
```

**交付物:** 功能完整的 Tier0-Edge v2 重写版。

---

## 附: 关键决策记录 (ADR)

### ADR-001: 为什么砍掉 Kong

**背景:** 当前使用 Kong 作为 API 网关（26 服务、50 路由、7 插件、1 个自定义 Lua 插件），需要独立的 PostgreSQL 数据库。

**决策:** 用 Go `httputil.ReverseProxy` 替代。

**理由:**
1. Tier0 是 edge 部署，不是云端微服务，不需要独立 API 网关
2. Kong 的认证插件需要网络调用 Go 后端的 `/auth/userinfo`，不如直接在 Go 中做
3. Go 内置 reverse proxy 可以完成路由转发
4. 减少 2 个容器（Kong + Konga）和 1 个数据库依赖

**风险:** 失去 Kong 的限流、熔断等高级功能。缓解: 后续可通过 Go 中间件实现。

### ADR-002: 为什么默认不用 Keycloak

**背景:** 当前 Keycloak 是必装组件（512MB+ 内存），但大多数 edge 部署只需要基本的用户名/密码认证。

**决策:** 内置本地认证，Keycloak 作为可选 OIDC Provider。

**理由:**
1. Edge 场景通常 2-5 个用户，不需要企业级身份管理
2. Keycloak 的内存占用对 8GB 的 edge 设备是显著负担
3. 保留 OIDC 协议支持，企业用户仍可对接
4. 减少 1 个容器和大量配置

### ADR-003: 为什么用 sqlc 替代 GORM

**背景:** 当前使用 GORM，存在 N+1 查询不明显、复杂 SQL 用字符串拼接、性能不可预测等问题。

**决策:** 核心查询用 sqlc，简单 CRUD 可选 GORM。

**理由:**
1. IIoT 的查询模式以时序聚合为主，需要精确控制 SQL
2. sqlc 从 SQL 生成类型安全的 Go 代码，编译期发现问题
3. TimescaleDB 的高级功能（hypertable, continuous aggregates）不能通过 ORM 表达
4. SQL 可以被 DBA 审查和优化

### ADR-004: 为什么用异步管道替代同步处理

**背景:** 当前 MQTT 消息处理是同步线性的：解析 → WS 推送（阻塞）→ 持久化。

**决策:** 基于 Go channel 的多阶段异步管道。

**理由:**
1. IIoT 数据管道天然适合流式处理
2. WS 推送不应阻塞持久化（慢客户端不应影响数据完整性）
3. 告警检测和实时计算需要独立的处理阶段
4. Go channel 提供天然的背压和并发控制
5. 每个阶段可独立伸缩 worker 数

---

## 总结: 重写的核心价值

| 维度 | 当前版本 | 重写版本 |
|------|----------|----------|
| 最小部署 | 14 容器，4GB+ 内存 | 3 容器，2GB 内存 |
| 启动时间 | 3-5 分钟 | 30 秒 |
| 配置复杂度 | 60+ 环境变量 | 5 个环境变量 |
| 功能完整度 | ~50 个空壳 Logic | 告警 + 计算 从第一天就有 |
| 类型安全 | Go + `any` TypeScript | SQL → Go → OpenAPI → TypeScript 端到端 |
| 测试 | 前端零测试 | CI 强制 lint + test |
| 消息吞吐 | 同步阻塞，WS 拖慢管道 | 异步管道，每阶段独立伸缩 |
| 新人上手 | 理解双 DI + 50 个空壳 + 14 容器 | 六边形架构 + 统一 DI + 3 容器 |
