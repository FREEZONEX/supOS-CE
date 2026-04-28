# AGENTS.md

> **同步维护要求**
>
> - 本目录中的 `CLAUDE.md` 与 `AGENTS.md` 是配对文档
> - AI 修改其中任意一份时，必须同步检查并更新另一份的对应说明
> - 前端规范的 Canonical 细则仍以 `CLAUDE.md` 为主，`AGENTS.md` 负责 skills 导航与执行入口
> - `.claude/skills/` 是本仓库 skills 的唯一事实来源，Codex 使用 `.agents/skills/`
> - 若映射缺失，请在 `Tier0` 根仓库执行 `bash shell/map-agent-skills.sh`
> - 根仓库 `update-*` 脚本会在同步代码后自动补齐 skills 映射

## 1) Canonical 规范来源

- 执行任何开发任务前，先遵循 [CLAUDE.md](./CLAUDE.md) 的 Canonical 约束。
- `AGENTS.md` 仅维护：skills 导航、调用顺序、少量补充约定。
- 后续新增/修改规则时，先更新 `CLAUDE.md`，并同步检查 `AGENTS.md` 的导航与入口说明是否需要更新。

## 2) Skills 导航

以下是项目内可用 skills（版本化到仓库）。

### 场景层（scenarios）

- `tier0-delivery-orchestrator`：端到端交付编排（接口确认 → service 封装 → 页面对接 → UI 高还原）  
  file: `.claude/skills/scenarios/tier0-delivery-orchestrator.md`
- `tier0-systematic-debugging`：系统化排障（先证据后修复，禁止猜修）  
  file: `.claude/skills/scenarios/tier0-systematic-debugging.md`

### 能力层（capabilities）

- `tier0-page-api-integration`：页面与 `src/apis/inter-api` 的标准对接  
  file: `.claude/skills/capabilities/tier0-page-api-integration.md`
- `tier0-ui-high-fidelity-build`：页面与组件高还原开发（含严格 i18n 约束）  
  file: `.claude/skills/capabilities/tier0-ui-high-fidelity-build.md`
- `tier0-writing-plans-lite`：轻量可执行计划编写  
  file: `.claude/skills/capabilities/tier0-writing-plans-lite.md`
- `i18n-setup`：i18n 配置与校验流程  
  file: `.claude/skills/capabilities/i18n-setup.md`
- `match-screenshot-colors`：截图颜色映射到项目 CSS 变量体系  
  file: `.claude/skills/match-screenshot-colors/skill.md`

## 3) Skill 调用规则

- 用户显式提到 skill 名称（如 `$tier0-ui-high-fidelity-build`）时，优先使用该 skill。
- 多阶段需求优先从 `tier0-delivery-orchestrator` 开始。
- 优先场景层编排，再调用能力层。
- 优先通过 `Skill` 工具调用 skills；若 skill 缺失/不可读，需明确说明并采用最近 fallback。
- 本仓库使用小写 `skill.md` 作为部分 skill 入口，映射脚本会自动转换为 Codex 可识别的 `SKILL.md` 包装目录。
- 执行映射脚本后，需要重开 Codex 会话刷新 skills 列表。

### 用户意图 → Skill 映射

| 你想做什么                             | 用哪个 Skill                              |
| -------------------------------------- | ----------------------------------------- |
| 新建一个完整页面（接口 + 样式 + 状态） | `tier0-delivery-orchestrator`             |
| 照着设计稿还原一个页面                 | `tier0-ui-high-fidelity-build`            |
| 后端给了新接口，要对接到页面           | `tier0-page-api-integration`              |
| 加国际化 / 多语言支持                  | `i18n-setup`                              |
| 截图颜色映射到 CSS 变量                | `match-screenshot-colors`                 |
| 遇到报错 / 功能行为不对                | `tier0-systematic-debugging`              |
| 写一个执行计划                         | `tier0-writing-plans-lite`                |
| **加个新功能 / 不确定用哪个**          | **`tier0-delivery-orchestrator`（兜底）** |

---

## 项目快速参考

### 目录结构

- `apps/web`：主 React 18 + Vite 前端
- `apps/services-express`：Express v5 API（CopilotKit + 健康检查）
- `packages/`：共享脚本与 TS 配置
- `plugins/`：模块联邦插件

### 常用命令

```bash
pnpm install
pnpm dev:web                  # 启动前端
pnpm dev:servicesExpress      # 启动 AI 服务（按需）
pnpm build:web
pnpm lint
pnpm intl:once                # 生成 i18n JSON
```

### 业务模块速查

- **核心**：UNS、采集流/事件流、仪表盘、插件管理、开放数据、账号/角色
- **主要路由**：`/home`、`/uns`、`/collection-flow`、`/EventFlow`、`/dashboards`、`/plugin-management`
- **权限**：`pageList/buttonList` 控制菜单与按钮；`AuthWrapper` 包裹权限组件
- **API 前缀**：`/inter-api/supos/...`；CopilotKit：`/copilotkit`

### 环境变量（常用）

- `API_PROXY_URL`、`VITE_ASSET_PREFIX`、`VITE_REMOTE_PREFIX`、`PORT`、`LLM_*`
- 登录：`loginPath` 或 `/tier0-login`；免登录 `/freeLogin?token=...&redirectUri=...`

### 常见排查

- 登录循环/403：检查 `loginPath`、`authEnable`、资源启用与用户 `pageList`
- 接口报错：检查 `.env` 代理配置
- 插件加载失败：确认插件安装与联邦配置
