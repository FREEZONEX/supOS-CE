# Go Backend Architecture Deep-Dive Analysis

**Analyzed:** `/workspace/backend`  
**Date:** 2026-03-15  
**Go version:** 1.24.2  
**Framework:** go-zero (v1.9.2) + custom Spring-like DI container

---

## Table of Contents

1. [Layered Architecture](#1-layered-architecture)
2. [Dependency Injection & Service Context](#2-dependency-injection--service-context)
3. [API Design](#3-api-design)
4. [Data Access Layer](#4-data-access-layer)
5. [Adapter Pattern](#5-adapter-pattern)
6. [Business Logic](#6-business-logic)
7. [Configuration & Initialization](#7-configuration--initialization)
8. [Error Handling Patterns](#8-error-handling-patterns)
9. [go.mod Dependency Analysis](#9-gomod-dependency-analysis)
10. [Summary: Anti-Patterns & Recommendations](#10-summary-anti-patterns--recommendations)

---

## 1. Layered Architecture

### 1.1 Directory Structure

```
backend/
├── backend.go                          # Main entry point
├── etc/                                # Configuration files (YAML)
├── internal/
│   ├── config/config.go                # Config struct definition
│   ├── middleware/                      # HTTP middleware (auth, context init)
│   ├── handler/supos/                  # HTTP handlers (controller layer)
│   │   ├── auth/
│   │   ├── app/
│   │   ├── uns/uns/
│   │   ├── uns/dashboard/
│   │   ├── uns/label/
│   │   └── ... (20+ sub-packages)
│   ├── logic/supos/                    # Business logic layer
│   │   ├── auth/
│   │   ├── uns/uns/                    # UNS logic
│   │   ├── uns/uns/service/            # UNS service classes
│   │   ├── uns/dashboard/service/
│   │   ├── uns/label/service/
│   │   └── ... (20+ sub-packages)
│   ├── repo/                           # Data access layer
│   │   ├── relationDB/                 # GORM-based PostgreSQL DAOs
│   │   ├── keycloak/                   # Keycloak DB access
│   │   └── event/subDev/              # MQTT event repo
│   ├── adapters/                       # External service adapters
│   │   ├── grafana/
│   │   ├── postgresql/
│   │   ├── timescaledb/
│   │   ├── kong/
│   │   └── msg_consumer/              # MQTT message consumer
│   ├── svc/serviceContext.go           # Manual service context wiring
│   ├── common/                         # Shared utilities
│   │   ├── adapter/                    # Adapter interfaces
│   │   ├── serviceApi/                 # Service interfaces (IPersistentService, IWebsocketSender, etc.)
│   │   ├── errors/                     # Custom error types
│   │   ├── event/                      # Event bus event types
│   │   ├── cache/
│   │   ├── config/
│   │   ├── constants/
│   │   ├── dto/
│   │   ├── vo/
│   │   ├── enums/
│   │   └── utils/
│   └── types/                          # API request/response types (goctl-generated)
├── share/                              # Shared libraries
│   ├── spring/                         # Custom Spring-like DI container + Event bus
│   ├── base/                           # Base utilities (StringBuilder, etc.)
│   ├── clients/                        # External clients (Keycloak, NodeRed, Kong)
│   ├── result/
│   └── ctxs/
├── http/                               # HTTP client DTOs
└── test/                               # Test files
```

### 1.2 Layer Communication Flow

```
HTTP Request
  → Middleware (CheckTokenWare → InitCtxsWare)
    → Handler (parse request, validate, call logic)
      → Logic (thin orchestration layer)
        → Service (business logic, registered as Spring beans)
          → Repo (GORM queries)
          → Adapters (external services via pgxpool, resty, MQTT)
          → Event Bus (publish async events)
```

### 1.3 Separation of Concerns Assessment

**Strengths:**
- Clear separation between handler/logic/repo/adapter layers
- Interfaces defined for cross-cutting services (`serviceApi/` package: `IPersistentService`, `IWebsocketSender`, `IDataSinkService`, `IUnsDefinitionService`)
- Adapter pattern properly abstracts external systems (PostgreSQL, TimescaleDB, Grafana, Kong)

**Weaknesses:**
- **Dual DI systems**: The codebase uses BOTH `svc.ServiceContext` (go-zero standard) AND a custom `spring` DI container simultaneously, creating confusion about which system owns what
- **Logic layer is often a thin passthrough**: Many logic files (e.g., `searchTreeLogic.go` at line 32-34) just delegate directly to a Spring bean service with zero logic, making the layer redundant
- **Handler → Logic → Service triple-hop**: The go-zero `goctl` scaffolding creates a Handler → Logic → Service chain where Logic adds no value in many cases

---

## 2. Dependency Injection & Service Context

### 2.1 Two Competing DI Systems

#### System 1: go-zero `ServiceContext` (`internal/svc/serviceContext.go`, lines 28-78)

Manual wiring of infrastructure dependencies:

```go
type ServiceContext struct {
    Config         config.Config
    InitCtxsWare   rest.Middleware
    CheckTokenWare rest.Middleware
    SnowFlake      *utils.SnowFlake
    OssClient      *oss.Client
    Keycloak       *clients.KeycloakClient
    SourceNodeRed  *noderedclient.Client
    EventNodeRed   *noderedclient.Client
    I18n           map[language.Tag]*i18n.MessageFile
}
```

This is passed through the handler layer via closure (`Handler(svcCtx)`) and into logic constructors.

#### System 2: Custom Spring-like IoC Container (`share/spring/beanFactory.go`)

A full-featured DI container built on top of `gitee.com/supos-community-edition/di/v2`:

- **`spring.RegisterLazy[T]`** (line 15-31): Lazy singleton registration via `init()` functions
- **`spring.RegisterBean[T]`** (line 50-57): Eager singleton registration
- **`spring.GetBean[T]`** (line 69-75): Service resolution (panics on not found)
- **`spring.GetBeansOfType[T]`** (line 76-85): Multi-bean resolution by interface type
- **`spring.RefreshBeanContext()`** (line 114-123): Triggers all lazy bean instantiation

Most business services use this system via `init()` functions:

```go
// From UnsAddService (logic/supos/uns/uns/service/UnsAddService.go:44-52)
func init() {
    spring.RegisterLazy[*UnsAddService](func() *UnsAddService {
        return &UnsAddService{
            log:             logx.WithContext(context.Background()),
            sysConfig:       spring.GetBean[*sysconfig.SystemConfig](),
            unsLabelService: spring.GetBean[*service.UnsLabelService](),
            removeService:   spring.GetBean[*UnsRemoveService](),
        }
    })
}
```

#### System 2b: Event Bus (`share/spring/eventBus.go`)

Reflection-based event bus that auto-registers methods prefixed with `OnEvent`:

- Methods like `OnEventBatchCreateTable300` are auto-detected via reflection (line 40-65)
- Trailing numbers determine execution priority (lower = higher priority)
- Events are published via `spring.PublishEvent(eventObj)` (line 86-115)

### 2.2 Initialization Sequence (`backend.go`, lines 36-88)

```
1. Load config
2. Create go-zero REST server
3. svc.NewServiceContext(c)        → DB init, cache init, Keycloak init, OSS init
4. handler.RegisterHandlers()      → Register HTTP routes
5. spring.RegisterBean(ctx)        → Register ServiceContext in Spring container
6. spring.RefreshBeanContext()     → Instantiate ALL lazy beans
7. PublishEvent(ContextRefreshedEvent) → Bootstrap adapters (MQTT, Grafana datasources, etc.)
8. server.Start()
```

### 2.3 Architectural Issues

| Issue | Location | Severity |
|-------|----------|----------|
| **Dual DI systems** create confusion about ownership | `svc/serviceContext.go` vs `spring/beanFactory.go` | HIGH |
| **`init()` ordering** is fragile and implicit | All adapter/service files | HIGH |
| **Blank imports** required to trigger `init()` registration | `backend.go` lines 12-21 | MEDIUM |
| **`spring.GetBean` panics** on resolution failure | `spring/beanFactory.go:71-74` | HIGH |
| **No interface abstraction** for most services; concrete types registered | Throughout | MEDIUM |

---

## 3. API Design

### 3.1 URL Structure (`internal/handler/routes.go`)

The routing file is **auto-generated by goctl** (1392 lines). URL prefixes:

| Prefix | Purpose | Auth |
|--------|---------|------|
| `/inter-api/supos/app` | Application management | CheckToken + InitCtxs |
| `/inter-api/supos/auth` | Authentication (token, logout) | None |
| `/inter-api/supos/uns/*` | UNS namespace operations | CheckToken + InitCtxs |
| `/inter-api/supos/uns/dashboard/*` | Dashboard CRUD | CheckToken + InitCtxs |
| `/inter-api/supos/event/*` | Event flow management | CheckToken + InitCtxs |
| `/inter-api/supos/group` | Group management | CheckToken + InitCtxs |
| `/inter-api/supos/userManage/*` | User management | CheckToken + InitCtxs |
| `/inter-api/supos/kong/*` | Kong route management | None |
| `/open-api/uns/*` | Public UNS access | None |
| `/service-api/supos/*` | Service-to-service APIs | None |

**RESTfulness Assessment:** Mixed. Some routes are RESTful (`GET /installed/:name`, `DELETE /uninstall/:name`), but many use RPC-style paths (`POST /uns/json2fs`, `POST /file/batchQuery`, `POST /alarm/pageList`). The use of POST for read operations is an anti-pattern.

### 3.2 Handler Pattern (goctl-generated)

Every handler follows the same pattern (`handler/supos/uns/uns/createModelInstanceHandler.go`):

```go
func CreateModelInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.CreateTopicDto
        if err := httpx.Parse(r, &req); err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
            return
        }
        l := uns.NewCreateModelInstanceLogic(r.Context(), svcCtx)
        resp, err := l.CreateModelInstance(&req)
        if err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
        } else {
            httpx.OkJsonCtx(r.Context(), w, resp)
        }
    }
}
```

### 3.3 Validation

- **go-zero `httpx.Parse`** handles basic binding from JSON/form/path params
- **Struct tags** with `validate:"required,max=63"` exist on request types (e.g., `CreateTopicDto.Name`)
- **No centralized validation middleware** - validation is scattered
- **Some handlers have manual validation** (`tokenHandler.go` line 18: `strings.TrimSpace(req.Code) == ""`)

### 3.4 Request/Response Types

All types are in `internal/types/types.go` (auto-generated by goctl, 1938 lines). Issues:
- **Massive `CreateTopicDto` struct** (70+ fields, lines 248-319) - a God Object
- **Mixed concerns**: `CreateTopicDto` contains DB fields (`TableName`, `DataSrcID`), UI fields (`AddFlow`, `AddDashBoard`), internal state (`FieldsChanged`, `CompileExpression`), and transport fields - violating single responsibility
- **Inconsistent ID handling**: Some fields use `json:"id,string"`, others use plain `json:"id"` - risks JSON parse errors with large int64 values

---

## 4. Data Access Layer

### 4.1 Repository Pattern (`internal/repo/relationDB/`)

Each entity has a "Dao" file (repository) and a "Model" file:

| Entity | Repo File | Model File |
|--------|-----------|------------|
| UnsNamespace | `unsNamespace.go` + `unsNamespace_query.go` + `unsNamespace_query_tree.go` + `unsNamespace_update.go` + `unsNamespace_export.go` | `modelUNS.go` |
| Dashboard | `dashboardDao.go` | `dashboardModel.go` |
| UnsLabel | `unsLabel.go` | (in `modelUNS.go`) |
| AppKey | `appKeyDao.go` | `appKeyModel.go` |
| Group | `groupDao.go` | `groupModel.go` |

### 4.2 GORM Usage

The repo uses **GORM v1.31.0** with a custom DB connection retrieval pattern:

```go
// unsNamespace.go:29-44
func GetDb(ctx context.Context) *gorm.DB {
    if connObj := ctx.Value("db"); connObj != nil {
        if db, is := connObj.(*gorm.DB); is {
            return db
        }
    }
    db := stores.GetCommonConn(ctx)
    // ... debug/silent logger selection
    return db
}
```

**Context-based transaction propagation** (`unsNamespace.go:45-47`):
```go
func SetDb(ctx context.Context, db *gorm.DB) context.Context {
    return context.WithValue(ctx, "db", db)
}
```

### 4.3 UnsNamespaceRepo Pattern

```go
type UnsNamespaceRepo struct{}  // Stateless struct - methods operate on injected *gorm.DB

func (p UnsNamespaceRepo) Insert(db *gorm.DB, data *UnsNamespace) error
func (p UnsNamespaceRepo) MultiInsert(db *gorm.DB, data []*UnsNamespace) error
func (p UnsNamespaceRepo) Update(db *gorm.DB, data *UnsNamespace) error
func (p UnsNamespaceRepo) Delete(db *gorm.DB, id int64) error
func (p UnsNamespaceRepo) SelectById(db *gorm.DB, id int64) (*UnsNamespace, error)
func (p UnsNamespaceRepo) FindOneByAlias(db *gorm.DB, alias string) (*UnsNamespace, error)
// ... 40+ query methods
```

### 4.4 Data Access Issues

| Issue | Location | Severity |
|-------|----------|----------|
| **String-interpolated SQL** (SQL injection risk) | `unsNamespace_query.go:198` (`WHERE lay_rec like '"+layRec+"/%'`) | CRITICAL |
| **More string-interpolated SQL** | `unsNamespace_query.go:369` (`lay_rec like '"+escapeSQL(layRec)+"%'`) | HIGH |
| **Raw SQL with manual string building** | `unsNamespace_query.go:206-226` (CountByParentAliasAndNames) | HIGH |
| **`context.WithValue` with string key "db"** for transactions | `unsNamespace.go:30,46` | MEDIUM |
| **Module-level `sync.Once` for column caching** | `unsNamespace.go:74-94` (MultiUpdate) | LOW |
| **No repo interface** - all repos are concrete structs | All repo files | MEDIUM |
| **Soft-delete filter `WHERE status=1`** repeated in EVERY query** | Throughout `unsNamespace_query.go` | MEDIUM |
| **Magic numbers** (`id>1000`, `id>10`, `data_type<>5`) scattered in queries | `unsNamespace_query.go:192,239,285` | HIGH |

### 4.5 Model Design (`modelUNS.go`)

The `UnsNamespace` model (lines 31-71) is a **massive struct with 30+ fields** including:
- JSON columns (`Fields`, `RefUns`, `Refers`, `Extend`, `LabelIds`) with custom `Scan`/`Value` implementations
- Transient fields marked with `gorm:"-"` (`ModelAlias`, `PathName`, `OldPath`, `CountExistsSiblings`)
- Custom JSON deserialization for `Fields` that handles legacy data format quirks (lines 241-292)

---

## 5. Adapter Pattern

### 5.1 Overview

Adapters live in `internal/adapters/` and abstract external service interactions:

| Adapter | Purpose | Registration |
|---------|---------|-------------|
| `postgresql/PgPersistentService` | Write UNS data to PostgreSQL | `spring.RegisterLazy` in `init()` |
| `timescaledb/TsdbPersistentService` | Write UNS time-series data to TimescaleDB | `spring.RegisterLazy` in `init()` |
| `grafana/GrafanaEventHandler` | Create/delete Grafana dashboards | `spring.RegisterLazy` in `init()` |
| `kong/logic/KongLogic` | Kong API gateway management | Singleton via `sync.Once` |
| `msg_consumer/UnsMessageConsumer` | MQTT message processing pipeline | `spring.RegisterBean` in `init()` |

### 5.2 Interface Abstraction

Service interfaces are defined in `internal/common/serviceApi/`:

```go
// IPersistentService.go
type IPersistentService interface {
    Persistent(unsData []UnsData)
    GetDataSrcId() types.SrcJdbcType
    GetDataSourceProperties() DataSourceProperties
    FillLastRecord(uns *types.UnsDefinition)
    Save(creates []types.UnsInfo) error
    Remove(topics []types.UnsInfo) error
}

// IWebsocketSender.go
type IWebsocketSender interface {
    SendMessage(msg WebsocketMessage)
    HasTopologies() bool
}

// IDataSinkService.go
type IDataSinkService interface {
    Sink(ctx context.Context, unsData []TopicMessage)
}
```

Both `PgPersistentService` and `TsdbPersistentService` implement `IPersistentService`. They are resolved at runtime via `spring.GetBeansOfType[serviceApi.IPersistentService]()`.

### 5.3 PostgreSQL Adapter (`adapters/postgresql/PgPersistentService.go`)

- Uses **pgxpool** (not GORM) for direct connection pooling (lines 21-28)
- Handles table creation, SQL batch execution, and data persistence
- Event-driven: responds to `RemoveTopicsEvent` via naming convention `OnEventRemoveTopicsEvent7`

### 5.4 TimescaleDB Adapter (`adapters/timescaledb/TsdbPersistentService.go`)

- Also uses **pgxpool** directly
- Manages TimescaleDB-specific views and hypertables
- Complex view migration logic with error recovery (lines 157-274)
- Has detailed retry/rename logic for conflicting table/view names

### 5.5 Grafana Adapter (`adapters/grafana/GrafanaEventHandler.go`)

- Event-driven: `OnEventBatchCreateTable300`, `OnEventRemoveTopicsEvent300`, `OnEventContextRefreshedEvent300`
- Creates datasources on startup with exponential retry (lines 113-142)
- Delegates to grafana utility functions for actual API calls

### 5.6 Kong Adapter (`adapters/kong/logic/kongLogic.go`)

- Uses **go-resty** HTTP client with retry configuration (line 88-96)
- Full Kong Admin API wrapper (services, routes, plugins)
- **Singleton via `sync.Once`** (not using Spring container) - inconsistent with other adapters

### 5.7 MQTT Message Consumer (`adapters/msg_consumer/UnsMessageConsumer.go`)

The central data pipeline:

```
MQTT Message → OnMsg() → parseJsonList() → procDataAndSendWs() → sendData()
                              ↓                      ↓                ↓
                         JSON parsing          WebSocket push     DB persistence
                                                                 (via IDataSinkService)
```

**Notable:** Uses `unsafe.Pointer` for zero-copy string↔bytes conversion (lines 89-100) - an optimization that is fragile and could break with GC changes.

### 5.8 Adapter Issues

| Issue | Location | Severity |
|-------|----------|----------|
| **Kong adapter uses different DI pattern** (`sync.Once` singleton) vs Spring container | `kong/logic/kongLogic.go:79-96` | MEDIUM |
| **`unsafe.Pointer` for string conversion** | `msg_consumer/UnsMessageConsumer.go:89-100` | MEDIUM |
| **Deprecated `reflect.SliceHeader`/`reflect.StringHeader`** | `msg_consumer/UnsMessageConsumer.go:94-98` | MEDIUM |
| **Goroutine leaks possible** in event handlers with `go func()` | `grafana/GrafanaEventHandler.go:64,103` | MEDIUM |
| **No circuit breaker** for external service calls | All adapters | LOW |
| **Hardcoded retry logic** with exponential backoff but no max timeout | `grafana/GrafanaEventHandler.go:118-134` | LOW |

---

## 6. Business Logic

### 6.1 Logic Layer Structure

The logic layer follows the goctl pattern with per-endpoint logic structs:

```go
// searchTreeLogic.go
type SearchTreeLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (l *SearchTreeLogic) SearchTree(req *types.SearchTreeReq) (*types.SearchTreeResp, err) {
    resp, err = spring.GetBean[*service.UnsQueryService]().SearchTree(l.ctx, req)
    return
}
```

Most logic files are **thin passthroughs** to Spring-registered services. The real business logic lives in `logic/supos/*/service/` packages.

### 6.2 UnsAddService (`logic/supos/uns/uns/service/UnsAddService.go`)

The most complex business service (590 lines). Key operations:

- `CreateModelInstance` (lines 454-524): Single UNS item creation
- `CreateModelAndInstancesInner` (lines 53-260): Batch creation with dependency resolution
- `saveBatchAndSendEvent` (lines 344-450): Transactional save with event publication

**Complexity concerns:**
- `CreateModelAndInstancesInner` is 200+ lines with 8+ parameters and multiple nested data structures
- Transaction management is manual (`tx.Begin()` / `tx.Commit()` / `tx.Rollback()`)
- Error collection via `map[string]string` is ad-hoc

### 6.3 WebsocketService (`logic/supos/uns/uns/service/WebsocketService.go`)

Manages real-time WebSocket subscriptions (615 lines):

- Uses multiple `sync.Map` instances for concurrent session tracking
- Supports subscriptions by UNS ID, topic path, and alias
- Event-driven topology updates via `OnEventUnsTopologyChangeEvent`

**Architecture:**
```
WebsocketService
├── sessions           (sessionId → WsSubscription)
├── idToSessionsMap    (unsId → set of sessionIds)
├── topicToSessionsMap (topic → set of sessionIds)
├── aliasToSessionsMap (alias → set of sessionIds)
└── topologySessions   (set of sessionIds subscribed to topology)
```

### 6.4 Anti-Patterns in Business Logic

| Anti-Pattern | Location | Description |
|-------------|----------|-------------|
| **God Object** | `types.CreateTopicDto` (70+ fields) | Used for create, update, batch, import, and internal state |
| **Logic layer as passthrough** | Most `*Logic.go` files | Adds indirection without value |
| **`spring.GetBean` called in hot paths** | `websocketLogic.go:54`, `searchTreeLogic.go:33` | Should be injected at construction |
| **Chinese comments/error messages in code** | Throughout | Hampers internationalization and readability for non-Chinese speakers |
| **Manual transaction management** | `UnsAddService.go:355-448` | Error-prone; should use a `WithTransaction` pattern |
| **Panic recovery in business logic** | `UnsAddService.go:358-361` | Masks bugs instead of failing fast |
| **`TODO` comments left in production code** | `WebsocketService.go:128,256`, `modelUNS.go:114` | Incomplete implementations |
| **Mixed return patterns** | Some return `(result, nil)` with error in result.Code, others return `(nil, err)` | Inconsistent error signaling |

---

## 7. Configuration & Initialization

### 7.1 Config Structure (`internal/config/config.go`)

```go
type Config struct {
    rest.RestConf                              // go-zero base config (host, port, log, etc.)
    Database       conf.Database               // Primary PostgreSQL DSN
    OssConf        conf.OssConf                // Object storage (local/minio)
    Export         ExportConfig                 // Import/export settings
    GrafanaUrl     string                      // Grafana base URL
    PostgresqlUrl  string                      // Sink PostgreSQL URL (env: SINK_PG_URL)
    TimescaledbUrl string                      // Sink TimescaleDB URL (env: SINK_TSDB_URL)
    DevLink        conf.EventConf              // MQTT broker config
    CacheRedis     cache.ClusterConf           // Redis config
    KeycloakDSN    string                      // Keycloak DB DSN
    OAuthKeyCloak  clients.KeycloakConfig      // Keycloak OAuth config
    NodeRed        nodered.NodeRedConfig       // Node-RED instances
    Kong           clients.KongConfig          // Kong API gateway
}
```

### 7.2 Config Files

| File | Purpose |
|------|---------|
| `etc/backend.yaml` | Production config (Docker) |
| `etc/backend-local.yaml` | Local development |
| `etc/backend-dev.yaml` | Development environment |
| `etc/backend-api.yaml` | API-specific config |

### 7.3 Config Issues

| Issue | Location | Severity |
|-------|----------|----------|
| **Secrets in config files** (Keycloak ClientSecret, EMQX keys) | `etc/backend.yaml:48-50,63` | CRITICAL |
| **Hardcoded max request body to 2GB** overrides config | `backend.go:47` | MEDIUM |
| **`ElasticsearchConfig` defined but unused** | `config/config.go:37-42` | LOW |
| **Config auto-detection via filesystem** (`../deploy/` check) | `backend.go:42-44` | LOW |
| **Environment variable fallbacks** only for some fields | `config.go:18-19` | LOW |

---

## 8. Error Handling Patterns

### 8.1 Custom Error Types (`internal/common/errors/errors.go`)

```go
type BuzError struct {
    Code   int
    Msg    string
    Params []any
}

type AppError struct { *BuzError }
type NodeRedError struct { *BuzError }
```

Helper constructors: `BadRequest()`, `NotFound()`, `InternalError()`, etc.

### 8.2 Response Helpers (`internal/common/errors/response.go`)

```go
type ResultVO struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, msg string, data interface{})
func Fail(w http.ResponseWriter, httpStatusCode, code int, msg string)
```

### 8.3 Error Propagation

**Repo → Logic → Handler flow:**

```
Repo:    return nil, stores.ErrFmt(err)     // Wraps GORM errors
Logic:   resp.Code = 500; resp.Msg = err    // Business errors in response body
Handler: httpx.ErrorCtx(ctx, w, err)        // go-zero error handler
```

### 8.4 Error Handling Issues

| Issue | Location | Description |
|-------|----------|-------------|
| **GlobalErrorHandler is completely commented out** | `errors/handler.go` (entire file) | No panic recovery middleware |
| **Mixed error signaling**: some use Go errors, some embed in response structs | `createModelInstanceLogic.go:39-43` vs `searchTreeLogic.go:33` | Inconsistent |
| **Error codes embedded in response body** with HTTP 200 | Multiple logic files (e.g., `SearchPaged` returns `resp.Code = 500` with HTTP 200) | Anti-pattern |
| **`stores.ErrFmt`** wraps all DB errors uniformly, losing context | All repo files | Makes debugging harder |
| **Some errors silently swallowed** | `svc/serviceContext.go:53-57` (Keycloak init logs but continues) | Could mask critical failures |

---

## 9. go.mod Dependency Analysis

### 9.1 Key Dependencies

| Dependency | Version | Purpose | Notes |
|-----------|---------|---------|-------|
| `zeromicro/go-zero` | v1.9.2 | HTTP framework | Core framework |
| `gorm.io/gorm` | v1.31.0 | ORM (primary DB) | Latest |
| `jackc/pgx/v5` | v5.7.4 | PostgreSQL driver (adapters) | Used alongside GORM |
| `gorilla/websocket` | v1.5.3 | WebSocket support | Stable |
| `eclipse/paho.mqtt.golang` | v1.5.0 | MQTT client | |
| `kong/go-kong` | v0.72.1 | Kong Admin API SDK | Not actually used (resty used instead) |
| `go-resty/resty/v2` | v2.7.0 | HTTP client | For Kong, external APIs |
| `expr-lang/expr` | v1.17.6 | Expression evaluation | For calculated fields |
| `nicksnyder/go-i18n/v2` | v2.5.1 | Internationalization | |
| `docker/docker` | v28.5.2 | Docker API | For app management |
| `dgraph-io/ristretto/v2` | v2.3.0 | Cache | |
| `maypok86/otter/v2` | v2.2.1 | Cache | Duplicate cache library? |
| `goburrow/cache` | v0.1.4 | Cache | Third cache library? |
| `karlseguin/ccache/v2` | v2.0.8 | Cache | Fourth cache library! |
| `patrickmn/go-cache` | v2.1.0 | Cache | Fifth cache library! |

### 9.2 Dependency Concerns

| Concern | Details | Severity |
|---------|---------|----------|
| **5 different caching libraries** | ristretto, otter, goburrow, ccache, go-cache | HIGH |
| **`kong/go-kong` imported but resty used instead** | Kong adapter uses resty directly | LOW |
| **Custom DI library from gitee** | `gitee.com/supos-community-edition/di/v2` - non-standard | MEDIUM |
| **Deprecated packages** | `golang-jwt/jwt/v4` alongside `jwt/v5`, deprecated `reflect.SliceHeader` | LOW |
| **Two JSON libraries** | `json-iterator/go` + standard `encoding/json` | LOW |
| **`parnurzeal/gorequest`** is archived/unmaintained | Used alongside resty | LOW |

---

## 10. Summary: Anti-Patterns & Recommendations

### Critical Issues

1. **SQL Injection via String Interpolation**
   - `unsNamespace_query.go:198`: `WHERE lay_rec like '"+layRec+"/%'`
   - `unsNamespace_query.go:369`: `WHERE lay_rec like '"+escapeSQL(layRec)+"%'`
   - `unsNamespace_query.go:394`: Same pattern
   - `unsNamespace_query.go:413`: `lay_rec like '"+layRec+"%'` with NO escaping
   - **Recommendation:** Use parameterized queries exclusively

2. **Secrets in Version Control**
   - `etc/backend.yaml` contains Keycloak secrets, EMQX API keys
   - **Recommendation:** Use environment variables or secret management

3. **`spring.GetBean` Panics on Missing Beans**
   - `share/spring/beanFactory.go:71-74`
   - **Recommendation:** Use `GetBeanOrErr` everywhere, or add startup validation

### High-Priority Issues

4. **Dual DI Systems** (ServiceContext + Spring Container)
   - **Recommendation:** Migrate fully to one system; preferably the Spring container since it's already dominant

5. **God Object `CreateTopicDto`** (70+ fields, used everywhere)
   - **Recommendation:** Split into domain-specific DTOs (CreateFolderReq, CreateFileReq, BatchImportReq, InternalUnsState)

6. **Magic Numbers in Queries** (`id>1000`, `data_type<>5`, etc.)
   - **Recommendation:** Use named constants from `constants` package consistently

7. **No Global Error Handler** (commented out in `errors/handler.go`)
   - **Recommendation:** Implement go-zero's `httpx.SetErrorHandlerCtx` for uniform error responses

### Medium-Priority Issues

8. **Thin Logic Layer** adds indirection without value
9. **Manual Transaction Management** instead of declarative patterns
10. **Chinese-only Error Messages** in some code paths
11. **5 Caching Libraries** should be consolidated
12. **`unsafe.Pointer` usage** in hot path (deprecated API)
13. **No repository interfaces** - makes testing harder
14. **`context.WithValue` with string keys** instead of typed keys

### Architecture Strengths

- **Clean adapter abstraction** with `IPersistentService` interface
- **Event-driven architecture** enables loose coupling between components
- **WebSocket subscription management** is well-structured with proper concurrency handling
- **Code generation** via goctl ensures consistency in handler/type patterns
- **Proper middleware chain** for auth and context initialization
- **SQL escaping utilities** exist (though not always used)
- **Comprehensive GORM model** with custom JSON serialization for backward compatibility
