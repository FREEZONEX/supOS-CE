# Tier0-Edge 项目深度分析报告

> 分析日期: 2026-03-15  
> 项目: [FREEZONEX/Tier0-Edge](https://github.com/FREEZONEX/Tier0-Edge) — 基于统一命名空间（UNS）的开源工业物联网平台  
> 当前版本: v2.0.0 (Go Refactor)  
> 许可证: Apache 2.0

---

## 目录

1. [项目概况](#1-项目概况)
2. [安全问题（严重）](#2-安全问题严重)
3. [代码质量问题](#3-代码质量问题)
4. [架构与设计问题](#4-架构与设计问题)
5. [部署与运维问题](#5-部署与运维问题)
6. [测试覆盖问题](#6-测试覆盖问题)
7. [文档与项目管理问题](#7-文档与项目管理问题)
8. [建议与改进路线图](#8-建议与改进路线图)

---

## 1. 项目概况

Tier0-Edge 是一个基于 UNS（统一命名空间）方法论的 IIoT 平台，包含:

| 组件 | 技术栈 | 说明 |
|------|--------|------|
| **后端** | Go 1.24 + go-zero | REST API, GORM, MQTT, Keycloak |
| **前端** | React 18 + Vite 7 + pnpm monorepo | Ant Design, CopilotKit AI, MCP |
| **消息中间件** | EMQX 5.8 | MQTT Broker |
| **数据存储** | PostgreSQL + TimescaleDB | 关系型 + 时序数据 |
| **API 网关** | Kong 3.9 | 认证, 路由, 代理 |
| **身份认证** | Keycloak 26 | OAuth2 / OIDC |
| **数据流** | Node-RED 4.0 | 协议转换, 事件流 |
| **部署** | Docker Compose | 全栈容器化部署 |

**关联仓库:**
- [FREEZONEX/supOS-frontend](https://github.com/FREEZONEX/supOS-frontend) (2 stars, 3 forks)
- [FREEZONEX/supOS-backend](https://github.com/FREEZONEX/supOS-backend) (8 stars, 5 forks)

---

## 2. 安全问题（严重）

### 2.1 硬编码凭据 — 严重性: 🔴 Critical

项目中存在大量硬编码的密码、密钥和凭据，分布在多个文件中:

| 文件 | 内容 | 风险等级 |
|------|------|----------|
| `deploy/.env.default:45` | OAuth Client Secret: `VaOS2makbDhJJsLlYPt4Wl87bo9VzXiO` | 🔴 Critical |
| `deploy/.env.default:56-58` | PostgreSQL/TSDB 密码: `postgres` | 🔴 High |
| `deploy/docker-compose.yml:199` | Grafana 管理员密码: `supos` | 🔴 High |
| `deploy/docker-compose.yml:295-296` | Keycloak 管理员: `admin/tier0` | 🔴 High |
| `deploy/docker-compose.yml:420-421` | MinIO: `admin/adminpassword` | 🔴 High |
| `deploy/docker-compose.yml:122` | EMQX 服务密码: `public` | 🔴 High |
| `backend/share/clients/keycloak.go:229` | Keycloak 密码硬编码: `tier0` | 🔴 Critical |
| `frontend/apps/web/src/pages/uns/components/topic-detail/FetchData.tsx:10` | 硬编码 API Key | 🔴 Critical |
| `deploy/mount/emqx/config/default_api_key.conf` | EMQX API Key 明文存储 | 🔴 High |

### 2.2 SQL 注入风险 — 严重性: 🟠 High

| 文件 | 位置 | 问题 |
|------|------|------|
| `backend/internal/repo/relationDB/unsNamespace_query_complex.go:186-195` | 时间戳参数通过 `fmt.Sprintf` 拼接 SQL | 若 `*updateStartTime` 来自用户输入，可被注入 |
| `backend/internal/adapters/postgresql/PgPersistentService.go:216` | 表名和字段名未充分转义 | 标识符注入风险 |
| `backend/internal/adapters/postgresql/table_name_util.go:9-27` | `GetFullTableName()` 未转义双引号 | PostgreSQL 标识符注入 |

### 2.3 Grafana 匿名管理员访问 — 严重性: 🔴 Critical

```yaml
# deploy/docker-compose.yml
GF_AUTH_ANONYMOUS_ENABLED: "true"
GF_AUTH_ANONYMOUS_ORG_ROLE: "Admin"
```

任何人无需认证即可以管理员身份访问 Grafana，可查看/修改所有仪表板和数据源。

### 2.4 EMQX ACL 形同虚设 — 严重性: 🟠 High

```erlang
# deploy/mount/emqx/config/acl.conf
{allow, all}.   # 最后一条规则允许所有客户端访问所有主题
```

ACL 最后一条 `{allow, all}` 使得所有前面的限制都失效，任何 MQTT 客户端都可以发布/订阅任意主题。

### 2.5 容器以 root 运行 — 严重性: 🟡 Medium

以下服务以 `user: root` 运行:
- Node-RED (`docker-compose.yml:135`)
- EventFlow (`docker-compose.yml:159`)
- Grafana (`docker-compose.yml:185`)
- Keycloak (`docker-compose.yml:269`)

Dockerfile 中也未设置非 root 用户（`backend/Dockerfile` 和 `frontend/Dockerfile-CN`）。

### 2.6 Konga 无认证 — 严重性: 🟠 High

```yaml
# deploy/docker-compose.yml
NO_AUTH: "true"   # Konga Kong 管理界面无需登录
```

### 2.7 Docker Socket 挂载 — 严重性: 🟠 High

Portainer 挂载了宿主机的 Docker Socket:
```yaml
- /var/run/docker.sock:/var/run/docker.sock
```
任何通过 Portainer 漏洞的攻击者都可以完全控制宿主机。

---

## 3. 代码质量问题

### 3.1 Express 错误处理中间件签名错误 — 严重性: 🔴 Critical

```typescript
// frontend/apps/services-express/src/middleware/error.ts
export const errorHandler = (req: Request, res: Response) => {  // ❌ 缺少 err 和 next 参数
  console.error('Sup-os-edge-frontend: Server Error:', req);
  res.status(500).json({
    error: config.nodeEnv === 'production' ? 'Internal Server Error' : req,  // ❌ 泄露 req 对象
  });
};
```

**问题:**
1. Express 错误中间件必须有 4 个参数 `(err, req, res, next)`，当前签名导致此中间件**永远不会被调用**
2. 开发环境下将整个 `req` 对象作为响应返回，可能泄露敏感信息（headers, cookies 等）

### 3.2 大量 `console.log` 残留在生产代码中

在 `frontend/apps/services-express/src/utils/mcp-client-manager.ts` 一个文件中就有 **20+ 处** `console.log`，其他文件也有大量调试日志:

| 文件 | 数量 |
|------|------|
| `mcp-client-manager.ts` | 20+ |
| `mcp-client.ts` | 7+ |
| `App.tsx` | 1 |
| `routers/index.tsx` | 2 |
| `utils/request.ts` | 2 |

### 3.3 TypeScript `any` 类型滥用

ESLint 配置中 `@typescript-eslint/no-explicit-any` 被设置为 `off`，导致大量 API 函数使用 `any` 类型:
- `frontend/apps/web/src/apis/inter-api/uns.ts` — 多个函数参数和返回值为 `any`
- `frontend/apps/web/src/App.tsx:96` — `event: any`

### 3.4 React Hooks 依赖缺失

| 文件 | 问题 |
|------|------|
| `frontend/apps/web/src/routers/index.tsx:42-54` | `useEffect` 使用了 `currentUserInfo`、`systemInfo` 等变量，但依赖数组只有 `[params?.isLogin]` |
| `frontend/apps/web/src/App.tsx:48` | `useEffect` 依赖 `systemInfo?.containerMap` 但未加入依赖数组 |

### 3.5 路由配置问题

```typescript
// frontend/apps/web/src/routers/index.tsx
const routes = [
  { path: '/', element: <RootRedirect /> },   // 第一个 '/' 路由
  { path: '/', element: <Layout />, ... },     // 第二个 '/' 路由 — 永远不会匹配
```

- 两个 `path: '/'` 路由，`Layout` 下的子路由可能永远不会被渲染
- 路径命名不一致: `/EventFlow` (PascalCase) vs `/collection-flow` (kebab-case)

### 3.6 Go 后端错误处理不一致

Handler 层存在两种不同的错误处理模式:
1. `result.Http(w, r, resp, err)` — 统一封装
2. `httpx.ErrorCtx` / `httpx.OkJsonCtx` — go-zero 原生

混用导致 API 响应格式不统一，增加前端处理复杂度。

### 3.7 Import Handler 无错误处理

```go
// backend/internal/handler/supos/uns/importHandler.go
importExport.Import()  // 返回值的错误未被捕获或处理
```

### 3.8 Go 后端入口路径依赖

```go
// backend/backend.go:41-43
if info, er := os.Stat("../deploy/"); er == nil && info.IsDir() {
    confFile = "etc/backend-local.yaml"
}
```

依赖相对路径判断运行环境，在不同工作目录下行为不可预测。

---

## 4. 架构与设计问题

### 4.1 子模块配置混乱

- `deploy/.gitmodules` 定义了 `frontend` 和 `backend` 子模块，指向 `supOS-frontend` 和 `supOS-backend` 仓库
- 但这些子模块**从未被初始化**，目录不存在
- 根目录的 `frontend/` 和 `backend/` 是普通目录，不是子模块
- `.gitmodules` 不在仓库根目录，Git 不会处理它
- **结论:** 子模块配置是遗留产物，造成困惑

### 4.2 多仓库与单仓库的混合

项目同时存在:
1. 主仓库 `Tier0-Edge` 包含前后端代码
2. 独立仓库 `supOS-frontend` 和 `supOS-backend`
3. Docker 镜像 `tier0/tier0-frontend:1.0.1-R8` 和 `tier0/tier0-backend:1.0.1-R8`

代码的真实来源（source of truth）不明确，可能导致版本不同步。

### 4.3 数据库迁移无版本控制

```go
// backend/internal/repo/relationDB/modelMigrate.go
fs.WalkDir(sqlFiles, "migrations_sqls", func(path string, d fs.DirEntry, err error) error {
    // SQL 文件按文件系统顺序执行，无版本号
    // 错误只打印警告，不中断启动
    if err != nil {
        log.Println("SQL WARN:", err, sql)  // ⚠️ 失败只是警告
    }
})
```

**问题:**
- 无迁移版本号，无法跟踪已执行的迁移
- SQL 执行失败只打警告，不中断启动
- 无回滚机制
- SQL 按文件名排序执行，无显式依赖管理
- 每次启动都会重新执行所有 SQL

### 4.4 服务间依赖缺乏健康检查

Kong 启动时依赖 PostgreSQL、EMQX、Keycloak，但 `depends_on` 未使用 `condition: service_healthy`:

```yaml
# deploy/docker-compose.yml
kong:
  depends_on:
    - emqx        # 无 healthcheck
    - uns         # 无 healthcheck
    - keycloak    # 有 healthcheck 但未在 depends_on 中使用条件
    - postgresql  # 有 healthcheck 但未在 depends_on 中使用条件
```

Kong 在依赖服务尚未就绪时执行 `kong migrations bootstrap`，会导致启动失败。

### 4.5 中文镜像源硬编码

Docker 构建和前端配置中硬编码了中国镜像源:
- `backend/Dockerfile`: `mirrors.aliyun.com`, `goproxy.cn`
- `frontend/Dockerfile-CN`: `registry.npmmirror.com`

国际用户使用时会遇到访问慢或不可用的问题。

---

## 5. 部署与运维问题

### 5.1 缺少 CI/CD 流水线

项目**完全没有 CI/CD 配置**:
- 无 GitHub Actions
- 无 GitLab CI
- 代码提交后无自动化测试、构建、部署
- 从 PR 记录看，合并前缺乏自动化检查

### 5.2 Docker Compose 缺少资源限制

只有 `marimo` 和 `keycloak` 设置了资源限制，其余 10+ 个服务没有:

| 缺少资源限制的服务 |
|-----|
| frontend, uns (backend), emqx, nodered, eventflow, postgresql, tsdb, kong, portainer, minio |

在生产环境中，单个服务的内存泄漏可能拖垮整台机器。

### 5.3 缺少日志轮转

主 `docker-compose.yml` 没有配置 `logging` driver:
```yaml
# ❌ 主文件无日志配置
# ✅ deploy/docker/run-env/docker-compose.yml 有日志配置（仅用于开发）
```

长期运行后日志文件会无限增长，耗尽磁盘空间。

### 5.4 敏感端口直接暴露

以下端口直接映射到宿主机:

| 端口 | 服务 | 风险 |
|------|------|------|
| 5432 | PostgreSQL | 数据库直接可访问 |
| 2345 | TimescaleDB | 数据库直接可访问 |
| 18083 | EMQX Dashboard | 管理界面暴露 |
| 8001 | Kong Admin API | 管理 API 暴露 |
| 9443 | Portainer | 容器管理暴露 |

### 5.5 缺少健康检查

10+ 个服务缺少 `healthcheck` 配置（见 4.4 节详细列表），无法实现有效的服务编排和故障检测。

### 5.6 安装脚本打印默认凭据

```bash
# deploy/bin/install.sh:120-121
echo "Default username: tier0"
echo "Default password: tier0"
```

---

## 6. 测试覆盖问题

### 6.1 前端零测试

- 未配置任何测试框架（无 Jest、Vitest、Playwright、Cypress）
- **零个测试文件** — 没有 `*.test.ts`、`*.test.tsx`、`*.spec.ts`
- 无测试相关的 npm scripts

### 6.2 后端测试有限

找到 ~20 个 Go 测试文件，但:
- 集中在 `backend/test/`（独立工具测试）和 `backend/internal/adapters/`（适配器层）
- **业务逻辑层（`internal/logic/`）零测试**
- **Handler 层零测试**
- 无集成测试、无 API 端到端测试

### 6.3 无 Lint 检查

- 前端 ESLint 配置存在但关键规则被禁用（`no-explicit-any: off`）
- 后端**无 golangci-lint** 或类似工具配置
- CI 中无自动化 lint 检查

---

## 7. 文档与项目管理问题

### 7.1 文档质量

- `backend/REAME_dev.md` — 文件名拼写错误（应为 `README_dev.md`）
- 中英文混杂，无统一语言标准
- API 文档依赖 Swagger（`/swagger`），但无独立 API 文档
- 部署文档仅覆盖基本 Docker Compose 场景，缺少生产部署指南

### 7.2 Git 工作流

- 从 PR 记录看，开发在 Gitee 进行，再合并到 GitHub
- 缺少 `CONTRIBUTING.md`（后端 README 提到但文件不存在）
- 无 Issue 模板、PR 模板
- 提交信息中英文混杂且不规范

### 7.3 未完成功能（TODO 统计）

项目中有 10+ 处 TODO/FIXME:
- `WebsocketService.go` — 3 处 TODO（文件导入、消息消费、实时推送）
- `expression.go` — 函数提取和布尔结果检测未实现
- `errutil/message.go` — i18n 翻译未集成
- 多个 DTO 验证待实现

---

## 8. 建议与改进路线图

### 🔴 P0 — 立即修复（安全相关）

1. **移除所有硬编码凭据**
   - `.env.default` 和 `.env.example` 只保留占位符
   - Go 和 TypeScript 代码中的硬编码密码/API Key 改为从环境变量读取
   - 轮换已暴露的 OAuth Client Secret

2. **修复 Grafana 匿名管理员访问**
   - 设置 `GF_AUTH_ANONYMOUS_ENABLED: "false"`
   - 或至少将 `GF_AUTH_ANONYMOUS_ORG_ROLE` 改为 `"Viewer"`

3. **修复 EMQX ACL**
   - 移除 `{allow, all}` 规则
   - 实施最小权限原则

4. **修复 SQL 注入风险**
   - `unsNamespace_query_complex.go` 中的时间戳参数使用参数化查询
   - `table_name_util.go` 添加标识符转义

5. **修复 Express 错误中间件**
   - 改为正确的 4 参数签名 `(err, req, res, next)`
   - 移除 `req` 对象在响应中的泄露

### 🟠 P1 — 短期改进（1-2 周）

6. **添加 CI/CD 流水线**
   - GitHub Actions: lint → test → build → push image
   - 至少包含: Go lint (`golangci-lint`), Go test, ESLint, TypeScript 类型检查

7. **完善 Docker 安全配置**
   - 所有容器添加非 root 用户
   - 移除不必要的端口暴露（PostgreSQL、EMQX Dashboard、Kong Admin）
   - 添加资源限制和日志轮转
   - 为所有服务添加 healthcheck

8. **清理代码**
   - 移除所有 `console.log` 调试语句
   - 启用 ESLint `no-console` 规则
   - 配置 `golangci-lint`

9. **修复前端路由**
   - 解决重复 `path: '/'` 问题
   - 统一路径命名风格（全部使用 kebab-case）

### 🟡 P2 — 中期改进（1-2 月）

10. **添加测试**
    - 前端: 配置 Vitest + React Testing Library，核心组件和 API 层达到 50% 覆盖
    - 后端: 业务逻辑层和 Handler 层添加单元测试
    - API 端到端测试

11. **数据库迁移管理**
    - 引入迁移版本控制工具（如 golang-migrate 或 goose）
    - 实现回滚机制
    - 迁移失败应中断启动

12. **清理仓库结构**
    - 删除无用的 `deploy/.gitmodules`
    - 明确 Tier0-Edge vs supOS-frontend/supOS-backend 的关系
    - 统一代码来源

13. **国际化部署支持**
    - Dockerfile 支持通过 build-arg 选择镜像源
    - 提供国际版和中国版两套配置

### 🟢 P3 — 长期改进

14. **可观测性**
    - 集成 Prometheus + Grafana 监控
    - 集成 Loki 或 ELK 日志收集
    - 添加 OpenTelemetry 链路追踪

15. **生产就绪**
    - Kubernetes Helm Chart 部署方案
    - 高可用配置（EMQX 集群、PostgreSQL 主从）
    - 备份和恢复方案
    - 安全审计和渗透测试

16. **开发者体验**
    - 完善 CONTRIBUTING.md
    - 添加 Issue 和 PR 模板
    - 规范提交信息（Conventional Commits）
    - 开发环境一键搭建脚本

---

## 总结

Tier0-Edge 作为一个 IIoT 平台，架构设计合理（UNS 方法论 + EMQX + TimescaleDB），技术栈现代（Go + React + CopilotKit AI）。但项目目前处于快速开发阶段，存在以下核心短板:

| 维度 | 评级 | 说明 |
|------|------|------|
| **安全性** | ⭐ (1/5) | 大量硬编码凭据、SQL 注入风险、宽松的 ACL 和认证配置 |
| **代码质量** | ⭐⭐ (2/5) | 错误处理不完善、类型安全弱、调试代码残留 |
| **测试覆盖** | ⭐ (1/5) | 前端零测试、后端仅有适配器层测试 |
| **部署运维** | ⭐⭐ (2/5) | 无 CI/CD、缺少资源限制和健康检查 |
| **文档** | ⭐⭐ (2/5) | 基础文档存在但不完整，缺少生产部署指南 |
| **架构设计** | ⭐⭐⭐⭐ (4/5) | UNS 架构合理，技术选型现代 |

**最重要的三件事:**
1. 立即修复安全问题（硬编码凭据、SQL 注入、Grafana/EMQX 开放访问）
2. 建立 CI/CD 流水线和基础测试
3. 完善 Docker 部署的安全和运维配置
