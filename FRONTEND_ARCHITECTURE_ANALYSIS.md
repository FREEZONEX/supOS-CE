# Frontend Architecture Deep-Dive Analysis

> **Project**: supOS Edge Frontend  
> **Analysis Date**: 2026-03-15  
> **Root**: `/workspace/frontend/`

---

## Table of Contents

1. [Monorepo Structure](#1-monorepo-structure)
2. [State Management](#2-state-management)
3. [Component Architecture](#3-component-architecture)
4. [API Layer Design](#4-api-layer-design)
5. [Routing Architecture](#5-routing-architecture)
6. [Plugin / Micro-Frontend Architecture](#6-plugin--micro-frontend-architecture)
7. [Services Layer (services-express BFF)](#7-services-layer-services-express-bff)
8. [i18n Implementation](#8-i18n-implementation)
9. [Build & Development](#9-build--development)
10. [Summary of Architectural Issues & Recommendations](#10-summary-of-architectural-issues--recommendations)

---

## 1. Monorepo Structure

### 1.1 Workspace Layout

**File**: `frontend/pnpm-workspace.yaml` (lines 1-6)

```
packages:
  - apps/*
  - packages/*
  - plugins/*
  - mcp/*
  - '!apps/services-hono'    # explicitly excluded
```

| Workspace Area | Contents | Purpose |
|---|---|---|
| `apps/web` | Main SPA (React + Vite) | Primary user-facing frontend |
| `apps/services-express` | Express v5 Node server | BFF for AI/CopilotKit, Docker health, MCP management |
| `apps/services-hono` | (excluded from workspace) | Likely deprecated or experimental |
| `packages/scripts` | `@supos-os-edge/scripts` CLI | i18n extraction, JSON-to-properties conversion |
| `packages/typescript-config` | `@supos-os-edge/typescript-config` | Shared tsconfig presets (base, react, node, vite) |
| `plugins/alert` | `@supos-ce/plugin-alert` | Module Federation remote plugin |
| `mcp/demo-mcp-server` | MCP demo server | AI/MCP integration demo |

### 1.2 Turborepo Configuration

**File**: `frontend/turbo.json` (lines 1-62)

- Build tasks (`build`, `build:web`, `build:server`, `build:servicesExpress`) all use `dependsOn: ["^build"]` for proper dependency ordering.
- Dev tasks disable caching (`cache: false`) and are marked `persistent: true`.
- The `intl:once`, `intl:watch`, and `properties:convert` tasks support i18n workflow.

### 1.3 Dependency Management via Catalogs

**File**: `frontend/pnpm-workspace.yaml` (lines 7-123)

pnpm catalogs organize dependencies into semantic groups:
- **framework**: React 18.3.1, react-router 7.9.4
- **build**: Vite 7.2.0, turbo, tsdown, Module Federation
- **util**: zustand 5.0.8, axios, ahooks, lodash-es, dayjs, immer, etc.
- **ui**: antd 5.27.4, @antv/x6, Carbon icons
- **lint**: ESLint 9, Prettier, Stylelint, commitlint, husky
- **ai**: CopilotKit 1.10.6, OpenAI, LangChain (Anthropic/Ollama/OpenAI), MCP SDK
- **server**: Express 5, dockerode

**Assessment**: The catalog system is a **strong pattern** that centralizes version management across the monorepo. This prevents version drift between packages.

### 1.4 Dependency Graph & Circular Dependencies

```
apps/web ─────────────> packages/scripts (devDep, for intl CLI)
apps/web ─────────────> packages/typescript-config (devDep)
apps/services-express ─> (no internal deps, standalone)
plugins/alert ────────> packages/scripts (devDep)
plugins/alert ────────> packages/typescript-config (devDep)
packages/scripts ─────> packages/typescript-config (devDep)
```

**No circular dependencies detected.** The dependency graph is a clean DAG. `services-express` is fully isolated, which is good for independent deployment.

### 1.5 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **No shared UI package** | Medium | `apps/web` exports components via Module Federation (`mfConfig.ts` exposes `./components`, `./utils`, `./hooks`, `./apis`). This tightly couples plugins to the host app's internal structure. A dedicated `packages/ui` would be cleaner. |
| **`services-hono` excluded but present** | Low | `apps/services-hono` exists on disk but is excluded from the workspace. Should be removed or documented. |
| **Single plugin only** | Info | Only `plugins/alert` exists. The plugin architecture is set up for expansion but currently underutilized. |
| **`mcp/` workspace area** | Info | Contains only a demo server. This workspace area may be premature. |

---

## 2. State Management

### 2.1 Overview

The app uses **Zustand** exclusively for global state. All stores use `createWithEqualityFn` from `zustand/traditional` with `shallow` equality for render optimization.

### 2.2 Global Stores

| Store | File | Lines | Purpose |
|---|---|---|---|
| `useBaseStore` | `src/stores/base/index.ts` | 357 | User info, system config, menu trees, permissions, theme plugin config |
| `useI18nStore` | `src/stores/i18n-store.ts` | 146 | Language, messages, antd locale, intl instance |
| `useThemeStore` | `src/stores/theme-store.ts` | 162 | Theme (light/dark/system), primary color, menu layout |
| `useAiStore` | `src/stores/ai-store.ts` | 29 | AI operation results (CopilotKit integration) |
| `useErrorStore` | `src/stores/error-store.ts` | 18 | Global error info |

### 2.3 Page-Level Stores

| Store | File | Lines | Purpose |
|---|---|---|---|
| `treeStore` | `src/pages/uns/store/treeStore.tsx` | 893 | UNS tree state: tree data, pagination, expand/collapse, search, WebSocket-driven lazy loading |

### 2.4 Store Architecture Patterns

**Positive patterns:**
- External-to-React state updates via exported functions (e.g., `setAiResult()`, `setTheme()`, `initI18n()`), keeping stores accessible outside React components.
- `subscribeWithSelector` middleware on `useI18nStore` for fine-grained subscriptions.
- `TreeStoreProvider` pattern with `createStore` + React Context for scoped page-level state — avoids polluting global state.
- `immer` middleware on treeStore for ergonomic nested state updates.
- `persist` middleware on treeStore for `lazyTree` setting survivng page reloads.

**Anti-patterns & Issues:**

| Issue | Severity | Location | Detail |
|---|---|---|---|
| **God store** | High | `stores/base/index.ts` (357 lines) | `useBaseStore` holds user info, system config, menu trees, button permissions, container lists, database types, MQTT broker type, dashboard types, and more. This should be split into `useAuthStore`, `useMenuStore`, `useSystemConfigStore`. |
| **Mega store function** | High | `stores/base/index.ts:147-288` | `updateBaseStore()` is a 140-line function that does 3 parallel API calls, complex permission processing, tree building, theme loading, guide config, UNS tree setup, and i18n initialization — all in one function. |
| **Excessive `any` types** | High | Multiple stores | `useAiStore` uses `[key: string]: any`, `useErrorStore` uses `errorInfo?: any`, base store has `pluginList: any[]`, `buttonGroup?: any[]`. |
| **Overly large page store** | Medium | `pages/uns/store/treeStore.tsx` (893 lines) | The tree store contains ~40 methods including pagination, abort controllers, WebSocket-related logic, and recursive data loading. Should be split into smaller composable stores or use a more modular approach. |
| **Direct localStorage coupling** | Medium | `stores/theme-store.ts`, `stores/i18n-store.ts` | Stores directly call `storageOpt.setOrigin()` for multiple third-party tools (chat2db, node-red, EMQ). This creates hidden coupling to external services. |
| **Module-level subscription** | Low | `stores/i18n-store.ts:104-106` | `useI18nStore.subscribe()` runs at module load time, creating a global subscription that lives forever (cleanup function exported but likely never called). |

---

## 3. Component Architecture

### 3.1 Page Structure (UNS as Case Study)

**Entry**: `src/pages/uns/index.tsx` (179 lines)

The UNS page is composed of:
```
WrapperModule
  └─ TreeStoreProvider (scoped zustand store)
     └─ UnsContextProvider (React context for WS data)
        └─ Module
           ├─ LeftDom (tree panel, 300 lines)
           ├─ ComContent
           │   ├─ TopDom (breadcrumb, import/export, 161 lines)
           │   └─ DetailDom (detail panel)
           └─ ModalContext (modal orchestration, 97 lines)
```

**Assessment**: The composition pattern is reasonable — layout regions separated into `LeftDom`, `TopDom`, `DetailDom` with modals isolated in `ModalContext`.

### 3.2 Component Size Analysis

| Component | File | Lines | Assessment |
|---|---|---|---|
| `treeStore.tsx` | `pages/uns/store/treeStore.tsx` | 893 | **Too large** — should be decomposed |
| `EditButton.tsx` | `pages/uns/components/EditButton.tsx` | 554 | **Too large** — mixes form logic, validation, conditional rendering for multiple data types, save workflow |
| `LeftDom.tsx` | `pages/uns/LeftDom.tsx` | 300 | Acceptable but has unused/commented code |
| `TopDom.tsx` | `pages/uns/TopDom.tsx` | 161 | Good size |
| `index.tsx` (UNS) | `pages/uns/index.tsx` | 179 | Good size |
| `ModalContext.tsx` | `pages/uns/ModalContext.tsx` | 97 | Good size |

### 3.3 Shared Components Library

**File**: `src/components/index.ts` (106 lines of exports)

The shared components library contains **~60 components** organized into directories with a `com-` prefix naming convention. Notable components:

- **Layout**: `com-layout`, `ComContent`, `ComLeft`, `ComRight`
- **Pro components**: `pro-table`, `pro-modal`, `pro-search`, `pro-tree`, `pro-tree-select`, `pro-card`, `pro-codemirror`
- **Primitives**: `com-button`, `com-input`, `com-select`, `com-checkbox`, `com-radio`, `com-text`
- **Domain**: `copilotkit/*`, `dynamic-mf-component`, `error-boundary`, `shepherd`
- **Auth**: `withAuth` HOC, `AuthButton`, `AuthWrapper`

The barrel export file (`index.ts`) uses **both** named re-exports and default re-exports, which doubles the surface area. Example:
```typescript
export * from './com-button';           // named exports
export { default as ComButton } from './com-button';  // default
```

### 3.4 Auth Component Pattern

**File**: `src/components/auth/index.tsx` (43 lines)

Uses a `withAuth` HOC that checks `hasPermission(auth)` and conditionally renders. `AuthButton` is `withAuth(Button)`. `AuthWrapper` wraps children and adds `data-button-auth` attributes.

**Assessment**: Clean pattern. However, permission checks happen at render time only — there's no route-level auth guard.

### 3.5 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **Giant components** | High | `EditButton.tsx` (554 lines) combines form definition, validation rules, API calls, conditional rendering for 3+ data types, and modal management. Should be split into `EditForm`, `EditModal`, and per-data-type renderers. |
| **Excessive commented-out code** | Medium | `LeftDom.tsx`, `EditButton.tsx`, `index.tsx` contain large blocks of commented JSX (Splitter panels, field selectors, etc.). |
| **Inconsistent naming** | Medium | Mix of `PascalCase` (`TopDom.tsx`, `LeftDom.tsx`, `DetailDom.tsx`) and `kebab-case` directories (`com-button`, `pro-table`). Page-level components use `Dom` suffix which is unconventional. |
| **Prop drilling** | Medium | `LeftDom` receives 7 props, many of which are passed through. `currentUnusedTopicNode`, `setCurrentUnusedTopicNode`, `unusedTopicBreadcrumbList` etc. could be absorbed into the tree store or context. |
| **Dual barrel exports** | Low | `components/index.ts` exports both `*` and `default` for each component, creating a large bundle entry point. |

---

## 4. API Layer Design

### 4.1 HTTP Client

**File**: `src/utils/request.ts` (284 lines)

Built on **axios** with a custom `ApiWrapper` class:

- **Interceptors**: Request interceptor adds `X-Sa-Token` header; response interceptor handles business codes (0, 200 = success), error notifications, 401/403/404 handling.
- **Business code pattern**: Responses use `{code, msg, data}` envelope. Success extracts `data`; errors show `msg` via `antd.message`.
- **Custom config flags**: `_businessResponse` (return full envelope), `_noMessage` (suppress error toast), `_noCode` (bypass code checking).
- **Timeout**: 600 seconds (10 minutes) — very generous.
- **ApiWrapper**: Provides `get`, `post`, `put`, `delete`, `upload`, `uploads` methods with URL base prefixing.

### 4.2 API Module Organization

**Directory**: `src/apis/inter-api/` (29 files)

| Module | File | Lines | Endpoints |
|---|---|---|---|
| `uns.ts` | UNS namespace | 244 | Tree, models, templates, labels, subscriptions, import/export, schemas, paste |
| `flow.ts` | Collection/Event flows | 66 | CRUD flows, deploy, save, bind |
| `i18n.ts` | Internationalization | ~50 | Language list, messages |
| `dashboard.ts` | Dashboards | ~50 | CRUD dashboards |
| `user-manage.ts` | User management | ~50 | CRUD users |
| `auth.ts` | Authentication | 10 | getUserInfo, logout |
| `alarm.ts` | Alarms | ~20 | Alarm operations |
| `apps.ts` | App management | ~30 | App operations |
| `global.ts` | Global operations | 26 | Import/export/download |
| (others) | Various | Varies | resource, role, plugin, kong, external, etc. |

**Re-export barrel**: `src/apis/inter-api/index.ts` re-exports 14 modules.

### 4.3 Type Safety Assessment

| Aspect | Rating | Detail |
|---|---|---|
| Request params | Poor | Most functions use `Record<string, unknown>` or `any` for params/data |
| Response types | Poor | Nearly all API functions return untyped `Promise<any>` (via `ApiWrapper`) |
| `ApiWrapper` generics | Unused | `post<T>` and `put<T>` accept generic `T` for request body but return `Promise<any>` |
| Business logic in API layer | Anti-pattern | `getAlertList` transforms response (adding `pageNo`, `pageSize`, `total`). `getAlertForSelect` maps to `{label, value}`. This transformation belongs in the calling layer. |

### 4.4 Caching

**No request/response caching is implemented.** There is no SWR, React Query, or custom cache layer. Every component re-fetches data independently.

### 4.5 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **No response type safety** | High | All API functions return `Promise<any>`. Defining response interfaces and typing `ApiWrapper` return values would catch bugs at compile time. |
| **No caching layer** | High | No SWR/React Query/tanstack-query. Repeated navigation triggers fresh API calls. Consider adding a caching layer for frequently-accessed data. |
| **Business logic in API layer** | Medium | `uns.ts` has functions like `getAlertList` that reshape response data. API modules should be pure transport; transformations belong in hooks/services. |
| **`uns.ts` is too large** | Medium | 244 lines with 40+ exports covering UNS tree, models, templates, labels, subscriptions, import/export, schemas. Should be split by domain. |
| **URL string building** | Low | Several functions manually build query strings (e.g., `updateModelSubscribe` constructs `?id=${id}&enable=${enable}`). Should use axios params. |
| **10-minute timeout** | Low | 600,000ms timeout is extremely long. Most API calls should use a shorter default with opt-in longer timeouts for upload/export operations. |

---

## 5. Routing Architecture

### 5.1 Route Configuration

**File**: `src/routers/index.tsx` (403 lines)

Routes are defined as a static array (`childrenRoutes`) with 20+ routes:

```
/ (root) ─── RootRedirect (login/redirect logic)
/ (layout) ─── Layout wrapper
   ├─ /home
   ├─ /uns
   ├─ /todo
   ├─ /grafana-design
   ├─ /collection-flow
   ├─ /collection-flow/flow-editor (child)
   ├─ /EventFlow
   ├─ /EventFlow/Editor (child)
   ├─ /dashboards
   ├─ /dashboards/preview (child)
   ├─ /account-management
   ├─ /aboutus
   ├─ /Localization
   ├─ /MenuConfiguration
   ├─ /advanced-use
   ├─ /dev-page
   ├─ /plugin-management
   ├─ /OpenData
   ├─ /AppManagement
   ├─ /403 (no permission)
   ├─ /404 (not found)
   └─ [dynamic plugin/iframe routes]
/freeLogin ─── FreeLoginLoader
/share ─── Share
* ─── NotPage
```

### 5.2 Auth/Permissions in Routing

**File**: `src/routers/index.tsx:280-390` (`getRoutesDom`)

The `getRoutesDom()` function dynamically constructs the route tree at runtime based on:
1. **`menuGroup`** (from backend resource API) — determines which routes are visible
2. **`systemInfo.authEnable`** — toggles auth enforcement
3. **`currentUserInfo.pageList`** — user's authorized page resources

Routes without matching backend resources are filtered out (except in dev mode where all routes are accessible). Plugin/iframe routes are dynamically appended from backend configuration.

**No route-level guards/middleware.** Auth is checked during route construction, not at navigation time. Once routes are built, there are no `beforeEnter` guards.

### 5.3 Code Splitting

**No lazy loading.** All page components are statically imported at the top of `routers/index.tsx`:
```typescript
import Uns from '@/pages/uns';
import Todo from '@/pages/todo';
import GrafanaDesign from '@/pages/grafana-design';
// ... 15+ more static imports
```

This means the **entire application is in a single bundle** (aside from Vite's manual chunking for vendor libs). For an app with 20+ pages, this is a significant missed optimization.

### 5.4 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **No code splitting / lazy loading** | High | All 20+ page components are eagerly imported. Use `React.lazy()` or react-router's lazy route loading for significant bundle size reduction. |
| **No route guards** | Medium | Auth is baked into route construction. If a user's permissions change mid-session, stale routes remain accessible until a full page reload. |
| **Inconsistent path casing** | Medium | Mix of `/EventFlow`, `/OpenData`, `/Localization`, `/MenuConfiguration` (PascalCase) with `/collection-flow`, `/account-management`, `/advanced-use` (kebab-case). |
| **Complex route building logic** | Medium | `getRoutesDom()` is a 110-line function with nested maps, filters, and conditional logic. Hard to maintain. |
| **Commented-out code** | Low | `RootRedirect` and `FreeLoginLoader` contain commented-out data router patterns. |

---

## 6. Plugin / Micro-Frontend Architecture

### 6.1 Technology

The plugin system uses **Module Federation** (via `@module-federation/vite` and `@module-federation/enhanced`).

### 6.2 Host Configuration

**File**: `apps/web/config/mfConfig.ts` (57 lines)

The host (`supos-ce/host`) exposes:
```typescript
exposes: {
  './components': './src/components/index.ts',
  './utils': './src/utils/index.ts',
  './hooks': './src/hooks/index.ts',
  './apis': './src/apis/inter-api/index.ts',
  './button-permission': './src/common-types/button-permission.ts',
  './constans': './src/common-types/constans.ts',
  './i18nStore': './src/stores/i18n-store.ts',
  './baseStore': './src/stores/base/index.ts',
  './tabs-lifecycle-context': './src/contexts/tabs-lifecycle-context.ts',
  './useTabsContext': './src/contexts/tabs-context.ts',
}
```

**Note**: Module Federation is currently **commented out** in `vite.config.ts` (line 31: `// federation(mfConfig)`), suggesting it's not active in the main build.

### 6.3 Plugin (Remote) Configuration

**File**: `plugins/vite.base.ts` (60 lines)

Each plugin remote:
- Exposes `./index` (App.tsx), `./enUS` and `./zhCN` (locale files)
- Shares React, react-dom, react-router, antd, @ant-design/icons, ahooks, carbon-icons, lodash-es, sass as singletons
- Consumes the host via `@supos_host` remote pointing to the host's `mf-manifest.json`

**File**: `plugins/pluginConfig.ts` (38 lines) — shared singleton config.

### 6.4 Dynamic Loading

**File**: `src/components/dynamic-mf-component/index.tsx` (51 lines)

`DynamicMFComponent` uses the `useRemote` hook to dynamically load plugin modules:

**File**: `src/hooks/useRemote.ts` (96 lines)

The loading flow:
1. `registerPlugins()` with lifecycle hooks
2. `registerRemotes()` with force: true
3. `loadRemote()` for the target module
4. Load plugin i18n resources and merge with host messages

Wrapped in `ErrorBoundary` with retry capability.

### 6.5 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **MF disabled in production** | High | Module Federation is commented out in `vite.config.ts`. Plugins only work when MF is enabled. |
| **Tight coupling** | High | Plugins import host internals (`components`, `utils`, `hooks`, `apis`, `stores`). Any refactoring of the host breaks all plugins. Should define a stable plugin API contract. |
| **Single plugin** | Info | Only `plugins/alert` exists. The architecture supports multiple plugins but is underutilized. |
| **No version negotiation** | Medium | Shared dependencies use exact `requiredVersion` but no fallback strategy if versions diverge. |

---

## 7. Services Layer (services-express BFF)

### 7.1 Purpose

The `services-express` app is a **Backend-for-Frontend (BFF)** that runs alongside the web app, providing:

1. **AI/CopilotKit runtime** — LLM proxy supporting OpenAI, Ollama, Anthropic, and Alibaba Tongyi
2. **MCP (Model Context Protocol) management** — client lifecycle management for AI tool integration
3. **Docker health monitoring** — container status reporting via Docker API

### 7.2 Architecture

**File**: `apps/services-express/src/index.ts` (32 lines)

```
Express App (port 4000)
├─ POST /copilotkit/* ──── CopilotKit runtime handler
├─ GET /copilotkit/mcp/* ── MCP client management CRUD
├─ GET /open-api/health ─── Docker container health
└─ Error handling middleware
```

### 7.3 CopilotKit Integration

**File**: `src/middleware/copilotkit.ts` (119 lines)

Supports 4 LLM providers via LangChain adapters:
- OpenAI (`ChatOpenAI`)
- Ollama (`ChatOllama`)
- Anthropic (`ChatAnthropic`)
- Tongyi/DashScope (`ChatOpenAI` with custom baseURL)

The `llmType` is selected at startup from `config.llmType` env var. The CopilotKit runtime is lazily recreated when MCP server configurations change.

### 7.4 MCP Client Manager

**File**: `src/utils/mcp-client-manager.ts` (413 lines)

Full lifecycle management for MCP clients:
- **Cache**: In-memory `Map<string, MCPClientEntry>` with 30-minute TTL
- **Operations**: Create, connect, disconnect, refresh, restart, stop, cleanup
- **Health checks**: Auto-reconnect on access if disconnected
- **Tool discovery**: List tools from connected MCP servers

**File**: `src/routes/copilotkit/mcp.ts` (256 lines) — REST API for MCP CRUD operations.

### 7.5 Why BFF is Needed

The Go backend handles core business logic (UNS, auth, flows, etc.). The BFF exists because:
1. **AI/LLM integration** requires Node.js libraries (CopilotKit, LangChain)
2. **MCP protocol** has best support in the JS/TS ecosystem
3. **Docker API access** for container health monitoring
4. Separation of AI concerns from business logic backend

### 7.6 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **Error handler signature** | High | `errorHandler` in `src/middleware/error.ts:7` has `(req, res)` signature (2 args) instead of Express error handler's required `(err, req, res, next)` (4 args). It will **never** be called as an error handler. |
| **Health endpoint reads certs** | High | `src/routes/open-api/health.ts:26-28` reads `/certs/*.pem` on every request. Should be read once at startup. |
| **Chinese log messages** | Low | Console messages are in Chinese (e.g., `获取MCP列表时发生错误`). Should use English for consistency. |
| **No authentication** | Medium | BFF endpoints have no auth middleware. Anyone with network access can manage MCP clients. |
| **Singleton MCP manager** | Low | `mcpManager` is a module-level singleton, which is fine for a single-process server but doesn't scale horizontally. |

---

## 8. i18n Implementation

### 8.1 Architecture

The i18n system uses **react-intl** with a multi-source message loading strategy:

```
Local JSON files (zh-CN.json, en-US.json)
    └─ merged with ─┐
Backend API (/inter-api/supos/uns/i18n/messages)
    └─ merged with ─┐
Plugin i18n (loaded via Module Federation)
    └─ Final merged messages stored in useI18nStore
```

### 8.2 Store

**File**: `src/stores/i18n-store.ts` (146 lines)

- **State**: `lang`, `langMessages`, `antMessages`, `intl` (IntlShape instance), `langList`
- **`initI18n()`**: Loads dayjs locale, antd locale, and merged messages. Stores everything in both Zustand and localStorage.
- **`getIntl()`**: Non-React utility function for formatting messages outside components (falls back to store's intl instance).
- **`connectI18nMessage()`**: Merges additional messages (from plugins) into existing store.

### 8.3 Message Loading

**File**: `src/utils/i18ns.ts` (122 lines)

- **Local**: Static JSON imports for zh-CN and en-US
- **Backend**: Dynamic fetch from `/inter-api/supos/uns/i18n/messages?lang=xx`
- **Antd**: Dynamic imports of all antd locale files (70+ locales)
- **Merge strategy**: `{ ...localMessages, ...backEndMessages }` — backend overrides local

### 8.4 Multi-System Language Sync

The i18n store syncs language preferences across multiple embedded services:
- `editor-language` for Node-RED
- `language` for EMQ/EMQX
- `lang` for Chat2DB
- `SUPOS_LANG` for supOS itself

### 8.5 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **Multiple localStorage writes per language change** | Medium | `initI18n()` performs 5+ `storageOpt.setOrigin()` calls for different services. Should be batched or abstracted. |
| **70+ antd locale loaders defined** | Low | All antd locales are defined in `antSources` map even if the app only supports a few. Not a runtime issue (dynamic import) but adds code bulk. |
| **Module-level subscription leak** | Low | `useI18nStore.subscribe()` at module level (line 104) creates a subscription that's never cleaned up in practice. |
| **Inconsistent translation access** | Medium | Components use `useTranslate()` hook, non-React code uses `getIntl()`, and some places access `useI18nStore` directly. Should standardize. |

---

## 9. Build & Development

### 9.1 Vite Configuration

**File**: `apps/web/vite.config.ts` (78 lines)

Key features:
- **Legacy browser support**: `@vitejs/plugin-legacy` targeting Chrome 89+, Safari 15+, Firefox 89+, Edge 89+
- **Build optimizations**: `drop: ['debugger']`, `pure: ['console.log']` in esbuild
- **Manual chunks**: vendor (react/react-dom/react-router), antd, charts (@antv/x6), utils (ahooks/lodash/dayjs)
- **Module Federation**: Commented out (`// federation(mfConfig)`)
- **Environment variables**: `process.env` injection, `VITE_APP_VERSION`, `VITE_APP_BUILD_TIMESTAMP`
- **Env prefix**: Supports `REACT_APP_`, `VITE_`, `OPENAI_` prefixed variables

### 9.2 Dev Proxy Configuration

**File**: `apps/web/config/supos.dev.ts` (122 lines)

Proxied paths:
- `inter-api`, `gateway`, `copilotkit`, `chat2db/api`, `minio/inter/supos`, `files/system/resource` → API_PROXY_URL
- `/copilotkit`, `/open-api` → localhost:4000 (services-express)
- `/plugin/` → API_PROXY_URL (remote plugins)
- `/iframe` → API_PROXY_URL (with path rewriting)

Configuration via `.env` and `.env.local` files using dotenv.

### 9.3 services-express Build

Uses **tsdown** (fast TypeScript bundler):
- Dev mode: `tsdown --watch`
- Build: `tsdown` → outputs to `dist/`
- Runtime: `node dist/index.js`

### 9.4 Issues & Recommendations

| Issue | Severity | Detail |
|---|---|---|
| **No code splitting** | High | All pages are eagerly imported. Vite supports dynamic imports natively. Adding `React.lazy()` wrappers would enable automatic code splitting. |
| **Large manual chunks may be counterproductive** | Medium | The `vendor` chunk bundles react+react-dom+react-router together, and `antd` is a separate chunk. Since antd is 1MB+, consider tree-shaking or splitting further. |
| **Legacy build doubles output** | Medium | `@vitejs/plugin-legacy` generates both modern and legacy bundles. Verify that Chrome 89+ target actually requires legacy polyfills. |
| **`process.env` injection** | Low | `define: { 'process.env': { ...devInfo } }` replaces all `process.env` references, which can cause issues with libraries that check `process.env.NODE_ENV`. |

---

## 10. Summary of Architectural Issues & Recommendations

### Critical Issues (High Severity)

1. **No code splitting / lazy loading** — All 20+ pages are eagerly imported, creating a monolithic bundle. Use `React.lazy()` + `Suspense`.

2. **God store pattern** — `useBaseStore` (357 lines) combines user, system, menu, permission, and config state. Split into focused stores.

3. **No API response typing** — All API functions return `Promise<any>`. Define response interfaces for compile-time safety.

4. **No data caching layer** — No SWR/React Query. Every navigation triggers fresh API calls. Add a caching strategy for frequently accessed data.

5. **Express error handler broken** — `errorHandler` middleware has wrong signature (2 params instead of 4). Will never catch errors.

6. **Module Federation disabled** — MF is commented out in vite.config.ts. Plugin architecture is non-functional in production.

### Important Issues (Medium Severity)

7. **Giant components** — `EditButton.tsx` (554 lines), `treeStore.tsx` (893 lines) need decomposition.

8. **Excessive `any` usage** — Store types, API params/responses, and component props heavily use `any`.

9. **Prop drilling** — UNS page passes 7+ props through multiple component layers. Use context or store.

10. **Inconsistent naming conventions** — Route paths, component files, and directory names use mixed casing strategies.

11. **No route-level auth guards** — Permissions only checked at route construction time, not at navigation time.

12. **No BFF authentication** — services-express endpoints are unprotected.

13. **Direct localStorage coupling** — Stores directly manage localStorage for multiple external services.

### Recommended Improvements

| Priority | Recommendation | Impact |
|---|---|---|
| P0 | Add `React.lazy()` to all page route imports | 50%+ reduction in initial bundle |
| P0 | Fix Express error handler signature to `(err, req, res, next)` | Error handling actually works |
| P1 | Split `useBaseStore` into `useAuthStore` + `useMenuStore` + `useSystemStore` | Maintainability, render performance |
| P1 | Add React Query / TanStack Query for API data management | Caching, deduplication, refetch control |
| P1 | Type all API response interfaces | Compile-time safety, IDE support |
| P2 | Enable Module Federation or remove dead code | Plugin system works or codebase is cleaner |
| P2 | Create `packages/ui` for shared components (decouple from host internals) | Plugin stability, independent versioning |
| P2 | Add route-level auth middleware | Security, real-time permission enforcement |
| P3 | Standardize naming conventions (kebab-case for paths, PascalCase for components) | Consistency |
| P3 | Remove all commented-out code blocks | Code hygiene |
| P3 | Add auth middleware to services-express | Security |
