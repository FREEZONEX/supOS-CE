# Tier0-Edge 功能实现与架构深度分析报告

> 分析日期: 2026-03-15  
> 项目: [FREEZONEX/Tier0-Edge](https://github.com/FREEZONEX/Tier0-Edge)  
> 版本: v2.0.0 (Go Refactor)

---

## 目录

1. [整体架构评价](#1-整体架构评价)
2. [后端架构分析](#2-后端架构分析)
3. [前端架构分析](#3-前端架构分析)
4. [UNS 核心数据流分析](#4-uns-核心数据流分析)
5. [功能完成度分析](#5-功能完成度分析)
6. [跨模块架构问题](#6-跨模块架构问题)
7. [改进建议](#7-改进建议)

---

## 1. 整体架构评价

### 1.1 架构亮点

Tier0-Edge 的整体架构设计有几个值得肯定的点:

**UNS 方法论落地**  
项目真正围绕 Unified Namespace（统一命名空间）方法论构建，EMQX 作为语义消息中间件、TimescaleDB 作为时序存储、PostgreSQL 作为关系存储，形成了清晰的 Source → UNS → Sink 数据管道。这在开源 IIoT 项目中是比较少见的完整实现。

**接口抽象设计**  
后端定义了一组干净的服务接口（`IPersistentService`、`IDataSinkService`、`IWebsocketSender`、`IUnsDefinitionService`），使得 PostgreSQL 和 TimescaleDB 可以作为可互换的持久化后端，这是优秀的策略模式应用。

**事件驱动的松耦合**  
自研的 Spring-like 事件总线（`spring/eventBus.go`）通过反射自动发现 `OnEvent*` 方法，支持优先级排序，实现了适配器层的松耦合。27 种事件类型覆盖了命名空间的全生命周期。

**AI/MCP 集成前瞻性**  
CopilotKit + MCP（Model Context Protocol）的集成方案是超前的，支持 4 种 LLM 提供商（Ollama、OpenAI、Anthropic、阿里通义）、3 种 MCP 传输方式（SSE、stdio、Streamable HTTP），以及运行时动态注册 MCP 服务器。

---

### 1.2 核心架构问题总览

| 问题 | 严重程度 | 影响范围 |
|------|----------|----------|
| 双重 DI 系统并存 | 🔴 架构级 | 后端全局 |
| ~50 个空壳 Logic 文件（功能未实现） | 🔴 功能级 | 后端业务逻辑 |
| 前端零代码分割、God Store | 🟠 性能级 | 前端全局 |
| 实时计算/告警/API 管理完全未实现 | 🔴 功能级 | 核心功能 |
| WebSocket → MQTT 反向通道断开 | 🟠 功能级 | 数据双向流 |
| 事件总线同步执行、无错误恢复 | 🟠 架构级 | 后端稳定性 |
| 数据库迁移无版本化 | 🟠 运维级 | 数据完整性 |

---

## 2. 后端架构分析

### 2.1 分层架构

后端采用 go-zero 框架，分层如下:

```
Handler (HTTP 入口)
   ↓
Logic (业务编排)
   ↓
Service (业务逻辑，spring 容器管理)
   ↓
Repo/Adapter (数据访问 / 外部服务适配)
```

**问题一: 双重 DI 系统并存**

这是后端最大的架构问题。项目同时使用两套依赖注入机制:

1. **go-zero 标准的 `ServiceContext`** — 通过闭包传递给 Handler 和 Logic
2. **自研 Spring 容器** (`share/spring/beanFactory.go`) — 通过 `init()` 函数全局注册 Bean

```go
// 系统 1: go-zero ServiceContext — 手动构建依赖
// backend/internal/svc/serviceContext.go
type ServiceContext struct {
    Config         config.Config
    Keycloak       *clients.KeycloakClient
    SourceNodeRed  *noderedclient.Client
    ...
}

// 系统 2: Spring 容器 — init() 自动注册
// backend/internal/logic/supos/uns/uns/service/WebsocketService.go
func init() {
    spring.RegisterLazy[*WebsocketService](func() *WebsocketService {
        return &WebsocketService{
            unsQueryService: spring.GetBean[*UnsQueryService](),  // 从 Spring 容器获取
            topologyService: spring.GetBean[*service.UnsTopologyService](),
        }
    })
}
```

这导致:
- Handler 层从 `svcCtx` 获取 Keycloak、NodeRed 等基础设施依赖
- Service 层从 `spring.GetBean` 获取业务服务依赖  
- **两个系统互不可见**: ServiceContext 不知道 Spring 容器中注册了什么，反之亦然
- **依赖关系不可预测**: `init()` 执行顺序依赖 Go 包导入顺序
- **测试困难**: 无法轻松替换 Spring 容器中的 Bean 进行单元测试

**问题二: `GetBean` 在找不到 Bean 时 panic**

```go
// backend/share/spring/beanFactory.go:69-74
func GetBean[Component any]() Component {
    rs, er := GetBeanOrErr[Component]()
    if er != nil {
        panic(er)  // 🔴 任何可选依赖缺失都会导致进程崩溃
    }
    return rs
}
```

在生产环境中，一个可选服务的缺失不应该导致整个应用崩溃。且该 panic 没有被任何 recover 中间件捕获（全局 error handler 被注释掉了）。

**问题三: 事件总线是同步阻塞的**

```go
// backend/share/spring/eventBus.go:86-114
func PublishEvent(eventObj any) error {
    // ...
    for _, handler := range listeners {
        err := handler.listener(eventObj)  // 同步调用，阻塞当前 goroutine
        if err != nil {
            return err  // 第一个错误就中断，后续 listener 不执行
        }
    }
}
```

- 所有 event listener 串行同步执行
- 第一个 listener 返回错误就中断后续所有 listener
- 在高吞吐 MQTT 消息处理路径中，这会成为严重瓶颈
- 没有超时、重试、错误隔离机制

### 2.2 数据访问层

**优点:**
- GORM 使用规范，模型定义清晰
- 物化路径（`lay_rec`）实现树形结构查询效率高
- `IPersistentService` 接口让 PostgreSQL 和 TimescaleDB 可互换

**问题:**
- 查询构建器大量使用字符串拼接（`unsNamespace_query_complex.go` 超过 200 行条件构建）
- 缺少 Repository 模式的完整抽象，部分逻辑层直接操作 GORM `*gorm.DB`
- 时间戳过滤直接拼接字符串到 SQL，存在注入风险

### 2.3 God Object: `CreateTopicDto`

`types/types.go` 中的 `CreateTopicDto` 承担了过多职责:

```
CreateTopicDto (70+ 字段)
├── 创建请求
├── 更新请求
├── 批量导入数据
├── 内部状态追踪 (ParentID, PathId 等)
├── 事件发布载荷
└── Node-RED 流配置
```

应当按职责拆分为 `CreateRequest`、`UpdateRequest`、`ImportRecord`、`InternalState` 等独立结构体。

### 2.4 缓存系统碎片化

项目同时引入了 **5 个缓存库**:
- `dgraph-io/ristretto`
- `maypok86/otter`
- `goburrow/cache`
- `karlseguin/ccache`
- `patrickmn/go-cache`

这表明缓存策略在不同开发阶段由不同人引入，没有统一收敛。应当选择一个作为标准（推荐 `otter`，它是最现代的 Go 缓存库）。

---

## 3. 前端架构分析

### 3.1 Monorepo 结构

```
frontend/
├── apps/
│   ├── web/           # 主 React SPA (20000+ 行页面代码)
│   ├── services-express/  # BFF for AI/CopilotKit/MCP
│   └── services-hono/     # 备选 BFF (excluded from workspace)
├── packages/
│   ├── scripts/       # 构建脚本
│   └── typescript-config/  # 共享 TS 配置
├── plugins/
│   └── alert/         # 告警插件 (Module Federation)
└── mcp/
    └── demo-mcp-server/   # MCP 演示服务
```

**优点:**
- pnpm workspace + Turborepo 的选型合理
- pnpm catalog 统一版本管理是好实践
- 依赖图无循环依赖

**问题: BFF 职责不清晰**

Express BFF (`services-express`) 的存在理由不够充分:
- 主要功能: CopilotKit 代理 + MCP 管理 + Docker 健康检查
- 这些功能完全可以集成到 Go 后端，减少一个服务节点
- 当前架构下，一个请求可能经过: Browser → Kong → Express BFF → Go Backend → DB，链路过长

### 3.2 状态管理

**God Store 问题**

`useBaseStore`（357 行）混合了完全不相关的状态:

```typescript
// frontend/apps/web/src/stores/base/index.ts
useBaseStore = create<BaseState>((set, get) => ({
    // 用户信息
    currentUserInfo, setCurrentUserInfo,
    // 系统配置
    systemInfo, setSystemInfo,
    // 菜单树
    menuTreeNode, setMenuTreeNode,
    // 权限
    pagePermissionMap, buttonPermissionMap,
    // 容器配置
    containerMap,
    // 路由
    routes,
    // 初始化
    fetchBaseInfo,
    // ... 还有更多
}))
```

这导致:
- 任何一个状态变化都可能触发大量无关组件重新渲染
- 虽然使用了 `shallow` equality，但状态粒度太粗
- 难以测试和维护

**页面级 Store 过大**

UNS 页面的 `treeStore`（893 行）包含 ~40 个方法和状态字段，应当拆分为:
- `unsTreeNavStore` — 树的展开/选中/搜索
- `unsDataStore` — 命名空间数据的 CRUD
- `unsUIStore` — 对话框/面板显示状态

### 3.3 组件设计

**UNS 模块代码量分析:**

| 文件 | 行数 | 问题 |
|------|------|------|
| `uns-tree/index.tsx` | 1082 | 树组件过大，混合渲染/状态/交互逻辑 |
| `use-create-modal/form-content/index.tsx` | 896 | 表单组件过大 |
| `store/treeStore.tsx` | 893 | Store 过大 |
| `reverse-modal/source-form/json/index.tsx` | 631 | JSON 编辑器组件 |
| `EditButton.tsx` | 554 | 编辑按钮组件内嵌完整表单和 API 调用 |
| `FieldsFormList.tsx` | 522 | 字段列表表单 |
| `topic-detail/index.tsx` | 517 | 主题详情页 |

UNS 模块总计 **20,140 行**，占前端页面代码的 **60%**，说明核心功能复杂度集中但缺乏合理拆分。

**问题模式:**
- 组件超过 500 行，违反单一职责
- `EditButton.tsx` 一个"按钮"组件包含完整的 Form、校验、API 调用和渲染逻辑
- 缺少 custom hooks 来提取业务逻辑（如 `useUnsCreate`、`useUnsTree`）

### 3.4 零代码分割

```typescript
// frontend/apps/web/src/routers/index.tsx
import UNS from '@/pages/uns';               // 同步导入
import EventFlowPage from '@/pages/event-flow';  // 同步导入
import Home from '@/pages/home';              // 同步导入
// ... 全部 20+ 页面都是同步导入
```

**所有页面在首次加载时全部打包进初始 bundle**，没有使用 `React.lazy()` 或任何动态导入。对于一个包含 20+ 页面、33,000+ 行页面代码的应用，这意味着:
- 首屏加载包含大量用户不需要的代码
- UNS 模块（20,000 行）即使用户只访问首页也会被加载
- 预计使用 `React.lazy()` 可减少 50%+ 的初始 bundle 大小

### 3.5 API 层缺少类型安全和缓存

```typescript
// frontend/apps/web/src/apis/inter-api/uns.ts
export const addModel = (data: any) => request.post<any>('/inter-api/supos/uns/...', data);
export const editModel = (data: any) => request.put<any>('/inter-api/supos/uns/...', data);
```

- 几乎所有 API 函数参数和返回类型都是 `any`
- 没有 API 响应类型定义
- 没有使用 React Query / SWR 等数据获取库
- 没有请求去重和缓存，相同数据可能被重复请求
- 无乐观更新（Optimistic Update）支持

### 3.6 Module Federation 插件未完全启用

```typescript
// frontend/apps/web/vite.config.ts
// Module Federation 配置存在但被注释掉
```

插件系统（Module Federation）的基础设施已搭建，但在主应用中被注释掉。Alert 插件直接导入宿主的内部组件和 hooks，耦合度过高。

---

## 4. UNS 核心数据流分析

### 4.1 完整数据管道

```
设备/协议 → Node-RED (SourceFlow)
                ↓ MQTT Publish
          EMQX Broker ($share/uns/#)
                ↓ QoS 1 共享订阅
     UnsMessageConsumer.OnMsg()
                ↓
    ┌───────────────────────┐
    │ 1. 主题解析 (缓存 LRU) │  ← UnsDefinitionService (110万条缓存)
    │ 2. JSON 解析           │
    │ 3. 字段校验(类型/范围)  │
    │ 4. 质量码设置           │
    └───────────────────────┘
                ↓
    ┌───────────────────────┐
    │ WebSocket 推送         │  → 前端实时显示
    │ (同步，阻塞消息处理)    │
    └───────────────────────┘
                ↓
    ┌───────────────────────┐
    │ 磁盘队列 (二进制编码)   │  ← 64MB 文件，顺序写入
    │ QueueDataSinkService  │
    └───────────────────────┘
                ↓ 批量消费
    ┌───────────────────────┐
    │ IPersistentService    │
    │ ├── PgPersistent      │  → PostgreSQL (关系型)
    │ └── TsdbPersistent    │  → TimescaleDB (时序)
    └───────────────────────┘
```

### 4.2 数据流问题

**问题一: WebSocket 推送阻塞 MQTT 消息处理**

```go
// backend/internal/adapters/msg_consumer/UnsMessageConsumer.go:81
msgList := u.procDataAndSendWs(ctx, def, data, strPayload, nil)  // 包含 WS 推送
t1 = time.Now()
u.sendData(ctx, msgList)  // 数据持久化
```

`procDataAndSendWs` 内部同步调用 WebSocket 推送:

```go
// 同步推送到所有订阅该主题的 WebSocket 连接
sessions.Range(func(key, value any) bool {
    subscription.WriteLock.Lock()
    subscription.conn.Write(msg)  // 如果某个连接慢，阻塞整个消息处理
    subscription.WriteLock.Unlock()
    return true
})
```

在高频数据场景下（工业传感器可以每秒产生数百条消息），WebSocket 推送的延迟会直接影响消息消费速率。

**问题二: WebSocket → MQTT 反向通道未实现**

```go
// backend/internal/logic/supos/uns/uns/service/WebsocketService.go:217
// TODO: Call topicMessageConsumer.onMessageByAlias(alias, body)
```

前端通过 WebSocket 发送的 `/send?t=alias&body=payload` 命令无法真正写入 MQTT，也就是说 **前端到设备的双向通信是断开的**。

**问题三: 磁盘队列满时静默丢数据**

磁盘队列在写满时没有告警机制，也没有回压（backpressure）策略。在 MQTT 消息洪峰时可能导致数据静默丢失。

**问题四: 实时计算完全未实现**

```go
// backend/internal/adapters/msg_consumer/UnsRealtimeCalcService.go
func (c UnsRealtimeCalcService) TryCalculate(...) (...) {
    if def == nil || len(def.RefUns) == 0 {
        return
    }
    //TODO  实时计算
    return
}
```

虽然消息处理管道中有实时计算的钩子，但实现体完全是空的。这意味着:
- 计算型字段（如移动平均、阈值比较）不工作
- 引用类型字段不工作
- 聚合类型字段不工作

### 4.3 TimescaleDB 集成设计

**优点:**
- 单一超级表 `uns_timeserial` + 动态列 (`double_1`, `str_1`...) + SQL 视图的方案在 IIoT 场景下是合理的
- 2 小时分块、50 分区、1 小时压缩、2 年保留策略配置合理
- 冲突处理使用临时表合并策略

**问题:**
- 视图与已有表名冲突时的自动迁移逻辑复杂，缺少单元测试
- 动态列的数量有上限，当字段数超出预定义列数时的处理不明确
- 缺少数据降采样（downsampling）策略

---

## 5. 功能完成度分析

### 5.1 空壳 Logic 文件统计

通过搜索后端 `internal/logic/` 目录中包含 `// todo: add your logic here` 的文件，发现 **约 50 个 Logic 文件完全未实现**:

| 模块 | 未实现文件数 | 影响 |
|------|-------------|------|
| **告警管理** (`alarm/`) | 4/4 | 🔴 告警创建、查询、更新、确认全部不可用 |
| **Kong API 管理** (`uns/kong/`) | 4/4 | 🔴 API 路由的 CRUD 全部不可用 |
| **仪表板** (`dashboard/`) | 3/4 | 🟠 仪表板详情和查询大部分不可用 |
| **文件管理** (`file/`) | 3/3 | 🔴 文件批量查询、更新、Blob 全部不可用 |
| **全局导入导出** (`global/`) | 5/5 | 🔴 系统级数据导入导出全部不可用 |
| **应用管理** (`app/`) | 4/4 | 🔴 应用搜索、更新、卸载全部不可用 |
| **示例** (`example/`) | 7/7 | 🟡 演示功能全部不可用 |
| **UNS 扩展** (external, batch, parser) | 7/7 | 🟠 外部树、批量查询、解析器不可用 |
| **挂载/开发工具** (`mount/`, `devtools/`) | 3/3 | 🟡 辅助功能不可用 |
| **附件/事件流代理** | 2/2 | 🟠 附件删除、事件流代理不可用 |

**已实现且功能完整的核心模块:**
- ✅ UNS 命名空间 CRUD（创建/查询/更新/删除/树管理）
- ✅ MQTT 消息消费 → 数据持久化管道
- ✅ WebSocket 实时数据推送（ID/Topic/Alias 三种订阅模式）
- ✅ UNS 数据导入导出（Excel/JSON）
- ✅ SourceFlow / EventFlow 的 Node-RED 集成
- ✅ Keycloak 身份认证和授权
- ✅ Kong API 网关配置（声明式）
- ✅ 拓扑图实时更新

### 5.2 关键未实现功能的影响

**告警管理 — 🔴 关键缺失**

对于 IIoT 平台来说，告警是最核心的功能之一。当前:
- 前端有告警插件 UI（`plugins/alert/`）
- 后端 4 个告警 Logic 文件全部为空壳
- 后端 `UnsQueryService` 中有 `//TODO Alarm` 注释
- 没有任何告警规则引擎、触发器或通知机制

这意味着 **系统无法在传感器数据异常时发出任何告警**，这对工业场景是致命的。

**实时计算 — 🔴 关键缺失**

`UnsRealtimeCalcService.TryCalculate()` 完全为空。这意味着:
- 计算字段（如公式 `temperature_celsius = (raw - 32) * 5/9`）不工作
- 引用字段不工作
- 聚合字段不工作
- UNS 数据模型中定义的这些高级类型形同虚设

**WebSocket 写回 — 🟠 重要缺失**

`HandleCmdMsg` 中的 `/send` 命令处理有 TODO，意味着前端无法通过 WebSocket 向设备写入数据。

**引用完整性检查 — 🟠 重要缺失**

```go
// backend/internal/logic/supos/uns/uns/service/uns_remove_helper.go:179
//TODO 引用检查
```

删除命名空间节点时不检查是否被其他节点引用，可能导致悬挂引用和数据不一致。

---

## 6. 跨模块架构问题

### 6.1 前后端数据契约无类型共享

后端使用 `goctl` 生成 `types/types.go`，前端 API 层返回 `Promise<any>`。两端之间没有:
- 共享的 API Schema（如 OpenAPI/Swagger 生成的类型）
- 自动化的类型同步机制
- API 版本管理

虽然后端有 Swagger UI（`/swagger`），但前端没有使用生成的类型。

### 6.2 三层服务调用链

```
Browser → Kong (auth + proxy) → Go Backend (business logic) → EMQX/PostgreSQL/TimescaleDB
                               → Express BFF (AI/MCP only)
```

**问题:**
- Kong → Go Backend 之间的超时是 300 秒，过长
- Express BFF 只服务 AI/MCP 功能，但部署为独立服务，增加运维负担
- Kong 同时转发给 Go Backend 和 Express BFF，路由规则复杂（~50 条规则）

### 6.3 Node-RED 深度耦合

SourceFlow 和 EventFlow 的实现完全依赖 Node-RED:
- 流的创建/删除通过 HTTP API 操作 Node-RED
- 流的定义存储在 Node-RED 的文件系统中
- 导入导出需要直接读写 Node-RED 的流配置

这意味着:
- Node-RED 是单点故障
- 无法水平扩展数据采集能力
- 流的版本控制依赖 Node-RED，而不是 Git 等标准工具

### 6.4 `unsafe` 包使用

```go
// backend/internal/adapters/msg_consumer/UnsMessageConsumer.go:89-100
func b2s(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
func s2b(s string) (b []byte) {
    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
    // ...
}
```

在消息消费的热路径上使用 `unsafe` 进行零拷贝字符串转换。虽然是性能优化，但:
- `reflect.SliceHeader` 和 `reflect.StringHeader` 在 Go 1.21+ 中已弃用
- 如果 `[]byte` 后续被修改，会导致 string 不可变性被破坏
- 应使用 Go 1.22+ 的 `unsafe.String()` 和 `unsafe.Slice()` 替代

---

## 7. 改进建议

### 7.1 架构层面

**A. 统一 DI 系统**

当前双重 DI 系统（go-zero ServiceContext + Spring 容器）是最需要解决的架构债务。建议:

```
方案 1: 全部收敛到 Spring 容器
- 将 ServiceContext 的成员全部注册到 Spring 容器
- Handler 层通过 spring.GetBean 获取依赖
- 优势: 一致性好，支持懒加载和事件
- 风险: 偏离 go-zero 标准用法

方案 2: 全部收敛到 ServiceContext
- Service 层的依赖在 ServiceContext 中构建，通过参数传递
- 移除 Spring 容器的 init() 注册模式
- 优势: 符合 go-zero 惯例，显式依赖
- 风险: ServiceContext 会变得很大

推荐: 方案 1，因为 Spring 容器的事件总线和懒加载是核心架构
```

**B. 事件总线改为异步**

```go
// 当前: 同步阻塞
func PublishEvent(eventObj any) error {
    for _, handler := range listeners {
        handler.listener(eventObj)  // 阻塞
    }
}

// 建议: 关键路径异步化
func PublishEventAsync(eventObj any) {
    go func() {
        for _, handler := range listeners {
            func() {
                defer recover()  // 错误隔离
                handler.listener(eventObj)
            }()
        }
    }()
}
```

至少在 MQTT 消息处理路径中，WebSocket 推送和其他非关键事件应当异步化。

**C. 解耦 Express BFF**

将 CopilotKit/MCP 功能迁移到 Go 后端（作为独立 HTTP handler），消除 Express BFF 节点:
- 减少一个服务的部署和维护
- 减少请求链路长度
- 减少技术栈复杂度

或者如果确实需要 Node.js 运行时（因为 CopilotKit SDK 是 JS），则考虑将其作为 sidecar 而非独立服务。

### 7.2 功能层面

**优先级 P0: 实现告警系统**

```
1. 定义告警规则模型（阈值、表达式、持续时间）
2. 在 MQTT 消息处理管道中集成规则引擎
3. 实现告警触发 → 通知（WebSocket + 持久化）
4. 对接已有的前端告警插件 UI
```

**优先级 P0: 实现实时计算**

```
1. 完善 UnsRealtimeCalcService.TryCalculate()
2. 支持基础计算（四则运算、比较、范围）
3. 支持引用类型字段
4. 支持滑动窗口聚合
```

**优先级 P1: WebSocket 双向通道**

```
1. 实现 HandleCmdMsg 中的 /send 命令
2. 调用 UnsMessageConsumer.OnMessageByAlias
3. 通过 MQTT 发布到设备
```

**优先级 P1: 引用完整性检查**

在删除命名空间节点前，检查被引用关系，拒绝或级联处理。

### 7.3 前端层面

**A. 代码分割 — 最高优先级、最低成本**

```typescript
// 当前
import UNS from '@/pages/uns';

// 改为
const UNS = React.lazy(() => import('@/pages/uns'));
```

仅此一项改动即可将初始 bundle 减小 50%+。

**B. 拆分 God Store**

```typescript
// 当前: 1 个 357 行的 baseStore
// 建议: 拆为 4 个独立 store

const useAuthStore = create<AuthState>(...);      // 用户信息、登录状态
const useSystemStore = create<SystemState>(...);   // 系统配置、容器信息
const useMenuStore = create<MenuState>(...);       // 菜单树、路由
const usePermissionStore = create<PermState>(...); // 权限映射
```

**C. 引入 TanStack Query**

替代手动 axios + useEffect 模式:
- 自动请求去重和缓存
- 后台刷新（stale-while-revalidate）
- 乐观更新
- 类型安全的 API hooks

**D. 大组件拆分**

以 `uns-tree/index.tsx`（1082 行）为例:
```
拆分为:
├── UnsTree.tsx (主容器，~100 行)
├── useUnsTreeData.ts (数据获取 hook)
├── useUnsTreeActions.ts (操作 hook: 展开/折叠/选中/拖拽)
├── UnsTreeNode.tsx (节点渲染)
├── UnsTreeToolbar.tsx (搜索/过滤工具栏)
└── UnsTreeContextMenu.tsx (右键菜单)
```

### 7.4 数据流层面

**A. MQTT 消息消费改进**

```
当前:  MQTT msg → 解析 → WS推送(同步) → 磁盘队列 → 持久化
建议:  MQTT msg → 解析 → Channel → [WS推送(异步), 磁盘队列] → 持久化
```

使用 Go channel 将消息处理拆分为并行管道，WebSocket 推送不再阻塞持久化。

**B. 磁盘队列增加回压和告警**

- 队列使用量超过 80% 时触发告警
- 队列满时拒绝新消息而非静默丢弃
- 提供队列深度的 metrics 接口

**C. 数据降采样**

TimescaleDB 支持 Continuous Aggregates，建议为高频数据配置自动降采样:
```sql
CREATE MATERIALIZED VIEW uns_hourly
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', time), tag, avg(double_1), max(double_1), min(double_1)
FROM uns_timeserial
GROUP BY 1, 2;
```

---

## 总结

Tier0-Edge 在架构层面有清晰的愿景（UNS 方法论 + 事件驱动 + 接口抽象），技术选型现代且合理。但作为 v2.0.0 版本:

**核心矛盾: 架构设计超前于功能实现**

项目定义了完整的接口、事件、数据模型，但约 50 个 Logic 文件仍然是空壳。告警、实时计算、API 管理等 IIoT 平台的关键功能完全缺失。

| 维度 | 评级 | 总结 |
|------|------|------|
| **架构设计** | ⭐⭐⭐⭐ (4/5) | UNS 落地好，接口抽象干净，事件驱动合理 |
| **核心功能 (UNS CRUD + 数据流)** | ⭐⭐⭐⭐ (4/5) | 消息消费→持久化管道完整，WebSocket 实时推送可用 |
| **功能完成度** | ⭐⭐ (2/5) | ~50 个空壳文件，告警/实时计算/API 管理缺失 |
| **前端工程质量** | ⭐⭐ (2/5) | 零代码分割、God Store、无类型安全、无测试 |
| **后端工程质量** | ⭐⭐⭐ (3/5) | DI 混乱和缓存碎片化是主要债务 |
| **数据流完整性** | ⭐⭐⭐ (3/5) | 单向管道完整，反向通道断开，缺少回压 |

**建议的投入方向:**
1. 🔴 先补齐关键功能（告警、实时计算），而非继续添加新模块
2. 🟠 统一 DI 系统、异步化事件总线，稳固后端架构
3. 🟠 前端代码分割 + Store 拆分 + TanStack Query，立竿见影提升体验
4. 🟡 清理 50 个空壳文件 — 要么实现，要么移除，避免给人"功能已完成"的错觉
