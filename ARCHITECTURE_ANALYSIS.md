# Architecture Deep-Dive: AI/MCP Integration, Authentication & System Design

> Generated analysis of the supOS Community Edition platform architecture.

---

## Table of Contents

1. [CopilotKit Integration](#1-copilotkit-integration)
2. [MCP (Model Context Protocol) Architecture](#2-mcp-model-context-protocol-architecture)
3. [Authentication & Authorization Architecture](#3-authentication--authorization-architecture)
4. [Kong as API Gateway](#4-kong-as-api-gateway)
5. [Frontend-Backend Communication](#5-frontend-backend-communication)
6. [Plugin Architecture](#6-plugin-architecture)
7. [End-to-End Auth Flow](#7-end-to-end-auth-flow)
8. [Architectural Decisions Summary](#8-architectural-decisions-summary)

---

## 1. CopilotKit Integration

### 1.1 Frontend Initialization

**File:** `frontend/apps/web/src/main.tsx` (lines 1-29)

CopilotKit is initialized at the **root level** of the React application, wrapping the entire `<App />` component:

```typescript
<CopilotKit
  runtimeUrl="/copilotkit"
  showDevConsole={false}
>
  <App />
</CopilotKit>
```

**Key observations:**
- `runtimeUrl="/copilotkit"` routes all AI requests through a relative URL, which Kong proxies to the Express BFF service on port 4000.
- An alternative direct MCP client agent URL (`/mcpclient/home/api/copilotkit`) is commented out, showing the system evolved from direct-to-agent to proxied architecture.
- `showDevConsole={false}` disables the CopilotKit debug panel in production.
- SSE MCP endpoints and named agents (`sample_agent`) are commented out, indicating they explored both SSE-based MCP and LangGraph agent patterns.

### 1.2 CopilotKit Middleware (Express BFF)

**File:** `frontend/apps/services-express/src/middleware/copilotkit.ts` (lines 1-119)

This is the core AI orchestration layer. It configures **four LLM providers** via LangChain adapters:

| Provider | Model Class | Default Model | Config Key |
|----------|-------------|---------------|------------|
| **Ollama** (local) | `ChatOllama` | `llama3.2` | `config.ollamaModal` |
| **OpenAI** | `ChatOpenAI` | `gpt-4o` | `config.openAiModel` |
| **Anthropic** | `ChatAnthropic` | `claude-3-7-sonnet-thinking` | `config.anthropicAiModel` |
| **Alibaba Tongyi** | `ChatOpenAI` (OpenAI-compatible) | `gpt-4o` | `config.tongyiModal` |

Each provider is wrapped in a `LangChainAdapter` that binds tools and streams:

```typescript
const serviceAdapterByOpenai = new LangChainAdapter({
  chainFn: async ({ messages, tools }) => {
    return openaiModel.bindTools(tools).stream(messages);
  },
});
```

**Provider selection** is runtime-configurable via `config.llmType` (env var `LLM_TYPE`), defaulting to `'openai'`. The selection map at line 61-66:

```typescript
const llmType = {
  ollama: serviceAdapterByllama,
  openai: serviceAdapterByOpenai,
  anthropic: serviceAdapterByAnthropic,
  tongyi: serviceAdapterByTongyi,
};
```

### 1.3 Dynamic CopilotRuntime with MCP Integration

Lines 68-90 implement a **singleton pattern with cache invalidation** for the `CopilotRuntime`:

```typescript
function createOrUpdateCopilotRuntime(): CopilotRuntime {
  const currentMcpServers = mcpManager?.getMCPClientCache()?.map((m) => ({ endpoint: m.endpoint })) || [];
  if (!globalCopilotRuntime || JSON.stringify(currentMcpServers) !== JSON.stringify(oldMcpServers)) {
    oldMcpServers = currentMcpServers;
    globalCopilotRuntime = new CopilotRuntime(
      currentMcpServers?.length > 0
        ? {
            mcpServers: currentMcpServers,
            createMCPClient: async (config) => {
              return await mcpManager.getOrCreateMCPClient(config);
            },
          }
        : {}
    );
  }
  return globalCopilotRuntime;
}
```

The runtime is **lazily recreated** when the MCP server list changes (compared by JSON serialization). This enables dynamic MCP server registration without service restarts.

### 1.4 LLM Configuration

**File:** `frontend/apps/services-express/src/config/index.ts` (lines 1-28)

All LLM settings are loaded from environment variables with fallback defaults:

| Env Variable | Purpose | Default |
|-------------|---------|---------|
| `LLM_TYPE` | Provider selector | `openai` |
| `LLM_MODEL` | Model name | varies per provider |
| `LLM_API_KEY` | API key | placeholder |
| `LLM_BASEURL` | Base URL (Ollama/Tongyi) | provider-specific |
| `AGENT_DEPLOYMENT_URL` | LangGraph deployment | `http://localhost:8123` |
| `LANGSMITH_API_KEY` | LangSmith observability | placeholder |

---

## 2. MCP (Model Context Protocol) Architecture

### 2.1 MCP Client Class

**File:** `frontend/apps/services-express/src/utils/mcp-client.ts` (lines 1-404)

The `MCPClient` class implements the `MCPClientInterface` from `@copilotkit/runtime` and supports **three transport types**:

| Transport | Class | Use Case |
|-----------|-------|----------|
| **SSE** | `SSEClientTransport` | Default; browser-compatible server-sent events |
| **stdio** | `StdioClientTransport` | Local process communication (e.g., npx-based servers) |
| **Streamable HTTP** | `StreamableHTTPClientTransport` | Newer MCP HTTP streaming protocol |

**Key implementation details:**

- **Tool caching** (line 26): Tools are cached after first retrieval (`toolsCache`), with `clearToolsCache()` for forced refresh.
- **Enhanced descriptions** (lines 181-199): Tool descriptions are augmented with required parameters and example inputs to help LLMs format correct calls.
- **Argument normalization** (lines 259-269): Handles common LLM mistakes like double-nested `params` objects: `{ params: { params: { actual_data } } }`.
- **JSON string detection** (lines 302-313): Automatically parses stringified JSON arguments that LLMs sometimes produce.
- **Example input derivation** (lines 319-377): Creates example usage from JSON Schema to guide LLMs.

### 2.2 MCP Client Manager

**File:** `frontend/apps/services-express/src/utils/mcp-client-manager.ts` (lines 1-413)

The `MCPClientManager` is a **singleton** (line 411) that manages the lifecycle of all MCP client connections:

**Cache structure:**
```typescript
interface MCPClientEntry {
  client: MCPClient;
  endpoint: string;
  lastUsed: number;
  isConnected: boolean;
}
```

**Lifecycle operations:**

| Method | Purpose |
|--------|---------|
| `getOrCreateMCPClient(config)` | Cache-first client retrieval with auto-reconnect |
| `removeMCPClient(endpoint)` | Graceful disconnect + cache eviction |
| `refreshMCPClient(endpoint)` | Disconnect/reconnect + tool cache clear |
| `restartMCPClient(endpoint)` | Full destroy + recreate |
| `stopMCPClient(endpoint)` | Disconnect but keep in cache |
| `cleanupExpiredClients()` | Remove clients idle > 30 minutes |
| `getToolsListByEndpoint(endpoint?)` | Get tools for one or all clients |

**Health checking** (lines 68-88): On every cache hit, the manager checks if the client is connected and attempts reconnection if needed, with graceful fallback to eviction if reconnection fails.

**TTL-based cleanup** (line 14): Cache TTL is 30 minutes (`CLIENT_CACHE_TTL = 30 * 60 * 1000`).

### 2.3 Transport URL Protocol

**File:** `frontend/apps/services-express/src/utils/path.ts` (lines 1-167)

MCP endpoints use a **custom URL protocol scheme** for transport type encoding:

| Protocol | Example | Parsed As |
|----------|---------|-----------|
| `sse://` | `sse://http://localhost:3001/sse` | SSE transport to given URL |
| `streamable-http://` | `streamable-http://http://localhost:3001/mcp` | Streamable HTTP transport |
| `stdio://` | `stdio://npx/@modelcontextprotocol/server-weather?env=API_KEY:xxx` | Local subprocess with env vars |

For `stdio://`, the parser handles scoped npm packages (e.g., `@scope/package`) by detecting `@` prefixes and merging path segments (lines 54-65).

### 2.4 MCP REST API (Management Routes)

**File:** `frontend/apps/services-express/src/routes/copilotkit/mcp.ts` (lines 1-257)

A full CRUD REST API for managing MCP connections at runtime:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/list` | GET | List all MCP clients with their tools |
| `/add` | POST | Register one or more MCP servers (batch support) |
| `/delete` | POST | Remove an MCP client by endpoint |
| `/refresh` | POST | Refresh (reconnect + clear tool cache) |
| `/restart` | POST | Full stop + recreate |
| `/stop` | POST | Pause a client without removing |

The `/add` endpoint (lines 37-114) supports both single and batch registration, with per-item error tracking and summary reporting.

### 2.5 Demo MCP Server

**File:** `frontend/mcp/demo-mcp-server/src/index.ts` (lines 1-31)

Entry point that supports all three transport modes via CLI argument:

```bash
node index.js stdio      # Standard I/O
node index.js sse        # Server-Sent Events (default port 3001)
node index.js streamableHttp  # HTTP streaming
```

**File:** `frontend/mcp/demo-mcp-server/src/server/index.ts` (lines 1-66)

The demo server registers a single **weather tool** (`get_weather`):

```typescript
server.registerTool('get_weather', {
  description: '获取指定城市的天气信息',
  inputSchema: {
    city: z.string().describe('城市名称（例如：北京、上海、杭州）'),
  },
}, async ({ city }) => {
  const weatherTool = new WeatherTool();
  const result = await weatherTool.getWeatherByCity(city);
  return { content: [{ type: 'text', text: formatWeatherResult(result) }] };
});
```

Uses `@modelcontextprotocol/sdk` with Zod schema validation.

### 2.6 MCP Type Definitions

**File:** `frontend/apps/services-express/src/types/mcp.ts` (lines 1-42)

Core types:
- `TransportType`: `'sse' | 'stdio' | 'streamable-http'`
- `McpClientOptions`: Server URL, transport type, client name, headers, event callbacks
- `McpServerConfig`: Name + transport type + transport-specific config (`StdioConfig` or `HttpConfig`)

---

## 3. Authentication & Authorization Architecture

### 3.1 Keycloak Client (Go Backend)

**File:** `backend/share/clients/keycloak.go` (lines 1-810)

A comprehensive **Keycloak Admin API client** implementing the full OAuth2/OIDC lifecycle:

**Configuration** (lines 93-113):
```go
type KeycloakConfig struct {
  Realm        string  // default: "supos"
  ClientName   string  // default: "supos"
  ClientID     string  // default: "supos"
  ClientSecret string  // default: "VaOS2makbDhJJsLlYPt4Wl87bo9VzXiO"
  IssuerURI    string  // Keycloak base URL
  RedirectURI  string  // OAuth callback URL
}
```

**Key operations:**

| Method | Lines | Purpose |
|--------|-------|---------|
| `Login(username, password)` | 243-256 | Password grant login |
| `GetKeyCloakTokenByCode(code)` | 258-271 | Authorization code exchange |
| `RefreshToken(refreshToken)` | 281-293 | Token refresh |
| `UserInfo(accessToken)` | 273-279 | Get user profile from token |
| `GetAdminToken()` | 221-241 | Get admin token (cached 10 min) |
| `GetUserExchangeTokenByID(userID)` | 434-462 | Token exchange (impersonation) |
| `CreateUser/DeleteUser/UpdateUser` | 304-366 | User CRUD |
| `CreateRole/DeleteRole/SetUserRoles` | 368-423 | Role management |
| `CreateResource/CreatePolicy/CreatePermission` | 471-638 | Authorization services (UMA) |
| `SetLocale(locale)` | 642-689 | Realm i18n configuration |

**Singleton pattern** (lines 147-150): Uses `sync.Once` for thread-safe initialization.

**Admin token caching** (lines 221-241): Admin tokens are cached with `go-cache` (10-minute TTL, 15-minute cleanup) to avoid repeated auth requests.

**Note:** The client UUID resolution (lines 164-174) is currently hardcoded to `Tier0ClientID = "a7b53e5e-3567-470a-9da1-94cc0c7f18e6"` instead of being dynamically resolved.

### 3.2 Backend Auth Middleware (Go)

**File:** `backend/internal/middleware/checktokenwareMiddleware.go` (lines 1-98)

The `CheckTokenWareMiddleware` validates every request:

**Flow:**
1. Check `SYS_OS_AUTH_ENABLE` env var -- if `false`/empty, auth is disabled (guest access allowed).
2. Extract `supos_community_token` from cookies (line 45).
3. Look up token in `cache.TokenCache` -- a server-side token cache.
4. If cache hit, refresh the token's TTL (line 72).
5. Call `authsvc.FetchUserInfo()` with the cached access token to get the full user profile.
6. Inject the user into the request context via `apiutil.SetUserInContext()`.

**Graceful degradation:** When auth is disabled (`SYS_OS_AUTH_ENABLE=false`), all requests proceed as a "Guest" user (`vo.Guest()`). This pattern repeats at every validation step (lines 47-93).

**File:** `backend/internal/middleware/initctxswareMiddleware.go` (lines 1-19)

A stub middleware (placeholder for future context initialization).

### 3.3 Frontend Auth State Management

**File:** `frontend/apps/web/src/utils/auth.ts` (lines 1-61)

Cookie-based token management using `js-cookie`:

- **Token key:** `supos_community_token` (defined in `constans.ts` line 14)
- `getToken()`, `setToken()`, `removeToken()` -- standard cookie CRUD
- `hasPermission(auth)` -- checks button-level permissions against `useBaseStore.buttonList`
  - Dev mode (`import.meta.env.DEV`) bypasses all permission checks (line 31)
  - Permissions use format `button:<permission_name>`
- `filterPermissionToList()` -- filters UI element arrays based on user permissions

**File:** `frontend/apps/web/src/stores/base/index.ts` (lines 1-357)

The `useBaseStore` (Zustand) is the central auth/permissions store:

**Initialization flow** (`updateBaseStore`, lines 147-288):
1. `Promise.allSettled` fetches routes/resources, user info, and system config in parallel.
2. Parses user's `resourceList` and `denyResourceList` into button and page permissions.
3. Deny permissions take priority (`filterObjectArrays(denyOthers, others)`).
4. Builds menu trees, home trees, and button permission lists.
5. If `authEnable === false` or user is `superAdmin`, grants all permissions (`button:*`).
6. Stores everything in Zustand for reactive UI updates.

**File:** `frontend/apps/web/src/apis/keycloak/auth.ts` (lines 1-16)

Direct Keycloak token endpoint access for frontend token exchange:

```typescript
export const getKeycloakToken = async (data) =>
  api.post(`/keycloak/home/auth/realms/${realm}/protocol/openid-connect/token`, data, {
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  });
```

---

## 4. Kong as API Gateway

### 4.1 Kong Configuration

**File:** `deploy/mount/kong/kong_config.yml.tpl` (3797 lines)

This is a **declarative Kong configuration** (format version 3.0) that defines the entire API gateway topology.

**Services (26 total):**

| Service | Host | Port | Protocol | Purpose |
|---------|------|------|----------|---------|
| `frontend` | `frontend` | 3000 | HTTP | Main React SPA |
| `backend` | `uns` | 8080 | HTTP | Go backend (inter-api), 5min timeouts |
| `backend-service-api` | `uns` | 8080 | HTTP | Service API layer |
| `backend-open-api` | `uns` | 8080 | HTTP | Public API layer |
| `GenerativeUI` | `copilotkit` | 4000 | HTTP | CopilotKit/Express BFF, **10min timeouts** |
| `mcpclient` | `mcpclient` | 3000 | HTTP | Dedicated MCP client service |
| `keycloak` | `keycloak` | 8080 | HTTP | Keycloak auth server |
| `nodered` | `nodered` | 1880 | HTTP | Node-RED flow engine |
| `EventFlow` | `eventflow` | 1889 | HTTP | Event flow processing |
| `grafana` | `grafana` | 3000 | HTTP | Dashboards |
| `gitea` | `gitea` | 3000 | HTTP | Git repository (CI/CD) |
| `emqx` | `emqx` | 18083 | HTTP | MQTT broker management |
| `portainer` | `portainer` | 9443 | HTTPS | Container management |
| `minio` | `minio` | 9001 | HTTP | Object storage |
| `gateway` | `gateway` | 8070 | HTTP | Internal gateway service |
| `plugin-frontend` | `plugin` | 3001 | HTTP | Plugin frontend server |
| `konga` | `konga` | 1337 | HTTP | Kong admin UI |
| `apm` | `apm` | 8080 | HTTP | Application performance monitoring |
| `swagger` | `uns` | 8080 | HTTP | API documentation |
| `marimo` | `marimo` | 8080 | HTTP | Notebook service |
| `login` | `keycloak` | 8080 | HTTP | Login redirect service |

**Routes (~50):** Map URL paths to services. Key routes:

| Route | Path | Service | Notes |
|-------|------|---------|-------|
| `frontend` | `/` | frontend | Catch-all for SPA |
| `backend` | `/inter-api/supos/` | backend | Go backend API |
| `copilotkit` | `/copilotkit` | GenerativeUI | AI endpoint, `strip_path: false` |
| `McpClient` | `/mcpclient/home` | mcpclient | MCP management UI |
| `login` | `/tier0-login` | login (keycloak) | OAuth login redirect |
| `open-backend-api` | `/open-api/` | backend-open-api | Public API with key-auth |

### 4.2 Kong Plugins

**Active plugins (7):**

1. **`supos-auth-checker`** (custom, lines 1932-1970): Global auth enforcement.
   - Enabled via `${KONG_AUTH_ENABLED}` (templated).
   - Whitelist regex patterns for public paths (auth endpoints, assets, locale, Keycloak, etc.).
   - Enables both `enable_deny_check` and `enable_resource_check`.
   - Login redirect URL includes Keycloak OAuth flow with `client_id=tier0`.

2. **`key-auth`** (lines 1973-1997): API key authentication on the `open-backend-api` route only.
   - Keys accepted via header, query, or body (`apikey` parameter).
   - Associated with the `open-api` consumer.

3. **`request-transformer`** (lines 1999-2039): On the login service.
   - Appends OAuth query parameters (`client_id`, `response_type`, `scope`, `redirect_uri`).

4. **`response-transformer`** (lines 2041-2077): On the login service.
   - Sets security headers: `X-Frame-Options: DENY`, `Content-Security-Policy: frame-ancestors 'none'`.

5. **`cors`** (lines 2079-2108): On the open-api route.
   - Allows all origins (`*`), all standard methods.

6. **`response-transformer`** (2nd instance, lines 2110-2144): On MinIO/object storage.
   - Sets `X-Frame-Options: SAMEORIGIN` (allows embedding within same origin).

7. **`supos-url-transformer`** (custom, lines 2146-2161): On the login route.
   - Redirects to `/?isLogin=true` after successful auth.

### 4.3 Custom Kong Plugins (Lua)

**Auth Checker Plugin:**

**File:** `backend/internal/adapters/kong/kong-plugins/kong-plugin-auth-checker/handler.lua` (lines 1-179)

This is the **core authentication gateway plugin** (priority 1000):

1. **Whitelist check:** Matches request path against regex whitelist; whitelisted paths bypass auth.
2. **Cookie extraction:** Reads `supos_community_token` from cookies.
3. **Backend validation:** Makes HTTP GET to `http://uns:8080/inter-api/supos/auth/userinfo` with the token cookie.
4. **Deny policy check** (if `enable_deny_check`): Compares request path against user's `denyResourceList` URIs.
5. **Resource permission check** (if `enable_resource_check`): Verifies request path is in user's `resourceList`.
6. On failure: Redirects to `login_url` (401) or `forbidden_url` (403).

**URL Transformer Plugin:**

**File:** `backend/internal/adapters/kong/kong-plugins/kong-plugin-url-transformer/handler.lua` (lines 1-64)

Post-authentication redirect: If user has a valid `supos_community_token` and the backend validates it, redirects to `conf.home_url` (`/?isLogin=true`).

### 4.4 Kong Adapter (Go Backend)

**File:** `backend/internal/adapters/kong/logic/kongLogic.go` (lines 1-489)

A Go client for the **Kong Admin API** using `resty`:

- `QueryRoutes()`: Fetches routes, parses tags for menu metadata (parent names, descriptions, sort order).
- `CreateService()` / `CreateRoute()`: Programmatic service/route creation.
- `MarkMenu()`: Persists menu selections to a local properties file.
- `AddAPIKey()`: Attaches `key-auth` plugin to a route.
- Menu items are identified by the `"menu"` tag in Kong routes.

**File:** `backend/internal/adapters/kong/route/route.go` (lines 1-52)

Registers Go backend HTTP handlers:
- `GET /inter-api/supos/kong/routeList` -- Kong route listing
- `POST /open-api/menu` -- Menu configuration
- `GET /test/nodered` -- NodeRed proxy

---

## 5. Frontend-Backend Communication

### 5.1 Proxy Configuration

**File:** `frontend/apps/web/config/supos.dev.ts` (lines 1-122)

Development proxy setup using Vite's built-in proxy:

**Proxied paths** (line 27-34):
```typescript
export const proxyList = [
  'inter-api',
  'gateway',
  'copilotkit',
  'chat2db/api',
  'minio/inter/supos',
  'files/system/resource',
];
```

All paths are proxied to a configurable base URL (default: `http://office.unibutton.com:11488`). The proxy supports:
- WebSocket upgrades (`ws: true`)
- CORS header rewriting (`changeOrigin: true`)
- Optional VPN/HTTPS proxy agent (commented out)
- Separate `singleList` for paths with a different proxy target

### 5.2 Communication Architecture

The system has a **three-tier** frontend communication model:

```
Browser  ─── Kong (port 8000/8443) ───┬─── frontend (port 3000) ──── React SPA
                                       ├─── uns (port 8080) ───────── Go Backend
                                       ├─── copilotkit (port 4000) ── Express BFF (AI/MCP)
                                       ├─── keycloak (port 8080) ──── Auth Server
                                       ├─── nodered (port 1880) ───── Flow Engine
                                       ├─── eventflow (port 1889) ─── Event Processing
                                       └─── plugin (port 3001) ────── Plugin Frontend
```

**Express BFF role** (`services-express` on port 4000):
- Hosts the CopilotKit runtime endpoint (`/copilotkit`)
- Manages MCP client connections (`/api/copilotkit/mcp/*`)
- Acts as an AI orchestration layer between the frontend and LLM providers
- Has **10-minute timeouts** in Kong (vs 1 minute for most other services) due to LLM response times

**Go Backend role** (`uns` on port 8080):
- Core business logic (UNS - Unified Namespace)
- Auth endpoints (`/inter-api/supos/auth/*`)
- System configuration, user management
- Kong admin operations
- File management
- Open API (key-auth protected)

**Direct service access** (via Kong routing):
- Keycloak: `/keycloak/home/` -- Auth admin UI
- Grafana: `/grafana/home/` -- Dashboards
- Node-RED: `/nodered/home/` -- Flow editor
- EventFlow: `/eventflow/home/` -- Event processing
- MinIO: `/minio/home/` -- Object storage
- Portainer: `/portainer/home/` -- Container management
- Gitea: `/gitea/home/` -- Git/CI/CD
- EMQX: `/emqx/home/` -- MQTT broker

### 5.3 Upstream Configuration

Kong uses **named upstreams** with round-robin load balancing. Key target mappings:

| Upstream | Target |
|----------|--------|
| `frontend` | `frontend:3000` |
| `backend` | `uns:8080` |
| `copilotkit` | `frontend:4000` |
| `keycloak` | `keycloak:8080` |
| `nodered` | `nodered:1880` |
| `plugin` | `frontend:3002` |

Note: The `copilotkit` upstream targets `frontend:4000`, meaning the Express BFF service runs as a sidecar on the same container/host as the frontend, but on port 4000.

---

## 6. Plugin Architecture

### 6.1 Module Federation Configuration

**File:** `frontend/plugins/pluginConfig.ts` (lines 1-38)

Shared dependencies for Module Federation:

```typescript
export const shared = {
  react: { singleton: true, requiredVersion: '18.3.1' },
  'react-dom': { singleton: true, requiredVersion: '18.3.1' },
  'react-router': { singleton: true, requiredVersion: '7.9.4' },
  antd: { singleton: true, requiredVersion: '5.27.4' },
  '@ant-design/icons': { singleton: true, requiredVersion: '6.1.0' },
  ahooks: { singleton: true, requiredVersion: '3.9.5' },
  '@carbon/icons-react': { singleton: true, requiredVersion: '11.69.0' },
  'lodash-es': { singleton: true, requiredVersion: '4.17.21' },
  sass: { singleton: true, requiredVersion: '1.93.2' },
};
```

All major dependencies are **singletons** to prevent duplicate React instances and ensure consistent UI behavior across host and remote modules.

### 6.2 Plugin Vite Base Configuration

**File:** `frontend/plugins/vite.base.ts` (lines 1-60)

Every plugin uses `@module-federation/vite` for **Vite-based Module Federation**:

```typescript
federation({
  name: `supos-ce/${name}`,
  manifest: true,
  exposes: {
    './index': './src/App.tsx',      // Main component
    './enUS': './src/locale/en-US.json',  // English locale
    './zhCN': './src/locale/zh-CN.json',  // Chinese locale
    ...exposes,
  },
  remotes: {
    '@supos_host': `supos-ce/host@${hostOrigin}/mf-manifest.json`,
  },
  shared: shared,
})
```

**Key design decisions:**
- **Manifest-based loading** (`manifest: true`): Plugins publish a `mf-manifest.json` for dynamic discovery.
- **Standard exports**: Every plugin exposes `./index` (App component) and locale files.
- **Host dependency**: Plugins consume `@supos_host` for shared hooks (`useTranslate`, `usePagination`), components (`ComLayout`, `ProTable`, `AuthButton`), and utilities.
- **Base path**: Plugins are served from `/plugin/${name}` in production.
- **Build output**: Each plugin builds to `dist/${name}/`.

### 6.3 Alert Plugin Example

**File:** `frontend/plugins/alert/` directory (27 files)

The Alert plugin demonstrates the full plugin pattern:

**Package:** `@supos-ce/plugin-alert` (line 3 of `package.json`)

**REMOTE_NAME:** `'Alert'` (from `variables.ts`)

**Environment config** (`env.ts`, lines 1-36): Dynamic environment detection for dev vs production, with auto-generated config files.

**App Component** (`src/App.tsx`, lines 1-459):

- Uses host-provided hooks: `usePagination`, `useTranslate`, `useActivate` (tab lifecycle)
- Uses host-provided components: `ComLayout`, `ComContent`, `ComDrawer`, `ProTable`, `ProCard`, `AuthButton`, `ComSearch`, `ComSegmented`
- Uses host-provided utilities: `validInputPattern`
- Implements `ButtonPermission` for auth-gated actions (edit, show, delete, add)
- Full CRUD with drawer-based forms, card/table view toggle, pagination
- i18n via `useTranslate(REMOTE_NAME)` for plugin-scoped translations

**Standalone entry** (`src/main.tsx`): Plugins can run independently for development:
```typescript
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <App title="" />
    </BrowserRouter>
  </StrictMode>
);
```

### 6.4 Plugin Integration in Kong

Plugins are served via the `plugin-frontend` service (host: `plugin`, port: 3001) with route `/plugin/`. The plugin frontend server runs on port 3001 on the `plugin` upstream (target: `frontend:3002`).

Additionally, plugin routes like `Alert` (`/Alert`) are mapped to the main frontend service, allowing the host app to render plugin content within its shell.

---

## 7. End-to-End Auth Flow

### 7.1 Initial Login Flow

```
1. User visits /  (any protected path)
         │
2. Kong auth-checker plugin intercepts
         │
3. No supos_community_token cookie?
    ├─ YES → Redirect to Keycloak login:
    │    /keycloak/home/auth/realms/tier0/protocol/openid-connect/auth
    │    ?client_id=tier0&redirect_uri=${BASE_URL}/inter-api/supos/auth/token
    │    &response_type=code&scope=openid
    │         │
    │  4. User authenticates with Keycloak
    │         │
    │  5. Keycloak redirects to /inter-api/supos/auth/token?code=XXX
    │         │
    │  6. Go backend exchanges code for tokens via KeycloakClient.GetKeyCloakTokenByCode()
    │         │
    │  7. Backend sets supos_community_token cookie
    │         │
    │  8. URL Transformer plugin redirects to /?isLogin=true
    │
    └─ NO (cookie exists) →
         │
    9. Kong auth-checker calls http://uns:8080/inter-api/supos/auth/userinfo
       with the cookie
         │
   10. Backend CheckTokenWareMiddleware:
       a. Reads supos_community_token from cookie
       b. Looks up in TokenCache (server-side)
       c. Calls authsvc.FetchUserInfo() with stored access_token
       d. Returns user info + resourceList + denyResourceList
         │
   11. Kong auth-checker checks:
       a. Deny policy: Is path in denyResourceList? → 403
       b. Resource check: Is path in resourceList? → Allow / 403
         │
   12. Request proceeds to upstream service
```

### 7.2 Frontend Auth State Load

```
1. App.tsx useEffect → fetchBaseStore(true)
         │
2. Promise.allSettled([
     getRoutesResourceApi(),   // GET /inter-api/supos/resource/routes
     getUserInfo(),            // GET /inter-api/supos/auth/user
     getSystemConfig(),        // GET /inter-api/supos/systemConfig
   ])
         │
3. Parse user resources:
   - resourceList → pages user CAN access
   - denyResourceList → pages user CANNOT access (deny overrides allow)
   - buttonGroup → button-level permissions (format: "button:Module.action")
         │
4. Build navigation:
   - menuTree (sidebar navigation)
   - homeTree (home page cards)
   - homeTabGroup (home page tabs)
         │
5. Store in Zustand (useBaseStore):
   - currentUserInfo, buttonList, menuGroup, systemInfo
         │
6. UI renders based on permissions
```

---

## 8. Architectural Decisions Summary

### AI/MCP Layer
- **Decision:** Express BFF as AI orchestration layer, not the Go backend.
  - **Rationale:** Node.js ecosystem has better LangChain/CopilotKit support; streaming-first design; 10-minute timeouts configured in Kong for LLM responses.
- **Decision:** Dynamic MCP server registration via REST API.
  - **Rationale:** Users can add/remove MCP servers at runtime without service restarts; supports all three MCP transport types.
- **Decision:** Custom URL protocol for MCP endpoints (`sse://`, `stdio://`, `streamable-http://`).
  - **Rationale:** Encodes transport type + connection params in a single string; enables batch management.

### Authentication
- **Decision:** Kong custom Lua plugin for auth enforcement (not JWT validation).
  - **Rationale:** Centralizes auth at the gateway; delegates token validation to the Go backend; supports both cookie-based sessions and API key auth.
- **Decision:** Cookie-based sessions (`supos_community_token`) rather than JWT in Authorization header.
  - **Rationale:** Simpler for iframe-embedded services (Grafana, Node-RED, etc.); single cookie works across all proxied services.
- **Decision:** Server-side token cache with Keycloak token exchange.
  - **Rationale:** Avoids exposing Keycloak tokens to the browser; backend manages token lifecycle.

### Plugin System
- **Decision:** Vite Module Federation for plugin architecture.
  - **Rationale:** True runtime module loading; shared dependencies prevent bloat; plugins can develop independently with their own dev servers.
- **Decision:** Host provides shared component library and hooks.
  - **Rationale:** Consistent UI across plugins; reduces plugin bundle sizes; centralized auth and i18n.

### API Gateway
- **Decision:** Kong as the single entry point with declarative configuration.
  - **Rationale:** Consolidates 25+ microservices behind one gateway; tag-based menu system; programmatic route management via Go adapter.
- **Decision:** Template variables (`${BASE_URL}`, `${KONG_AUTH_ENABLED}`, `${ENABLE_MCP}`) in Kong config.
  - **Rationale:** Environment-specific deployment configuration; feature flags for optional services.

### Communication Patterns
- **Decision:** Three API path prefixes: `/inter-api/supos/` (internal), `/service-api/supos/` (service), `/open-api/` (public).
  - **Rationale:** Clear separation of internal/service/public APIs; different auth policies per prefix.
- **Decision:** Frontend dev proxy mirrors production Kong routing.
  - **Rationale:** Development environment closely matches production; same URL patterns work in both environments.
