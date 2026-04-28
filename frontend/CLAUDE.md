# CLAUDE.md

> **同步维护要求**
>
> - 本目录中的 `CLAUDE.md` 与 `AGENTS.md` 是配对文档
> - AI 修改其中任意一份时，必须同步检查并更新另一份的对应说明
> - `.claude/skills/` 是本仓库 skills 的唯一事实来源，Codex 使用 `.agents/skills/`
> - 若映射缺失，请在 `Tier0` 根仓库执行 `bash shell/map-agent-skills.sh`
>
> 本节为当前有效约束；若与后文历史内容冲突，以本节为准。所有文字类描述一律使用中文。

## AI 开发规范（Canonical）

### 语言

- 所有输出（注释、文档、commit 信息、对话）必须使用中文。

### Skills

- Skills 单一来源：`.claude/skills/`；Codex 映射目录：`.agents/skills/`。
- 始终通过 `Skill` 工具调用，不要手动读取 skill 文件。
- 若 Codex 未识别到 skills，在 `Tier0` 根仓库执行 `bash shell/map-agent-skills.sh`，然后重开会话。
- Skills 导航详见 [AGENTS.md](./AGENTS.md)。

### 复用优先

- 实现前先搜索 `src/components/`、`src/hooks/`、`src/utils/` 是否已有可用实现。
- 高频易被忽略的工具：`ProTable`（`@/components/pro-table`）、`usePagination`、`useTranslate`、`ApiWrapper`（`@/utils/request`）。

### UI 开发

- 还原截图/设计稿时，**必须先调用 `match-screenshot-colors` skill** 映射颜色到 CSS 变量。
- 颜色必须使用 `var(--supos-*)` CSS 变量，禁止硬编码色值（`#161616`、`#f4f4f4` 等）。
- 优先使用应用级语义变量（`--supos-*`，定义于 `apps/web/src/index.scss :root`），禁止直接使用主题级变量（`--supos-t-*`）。
- 图标只使用 `@carbon/icons-react`，禁止引入其他图标库。
- 路由 import 使用 `react-router`，禁止 `react-router-dom`。
- 样式使用 SCSS Modules，class 名 kebab-case。

### i18n

- 所有用户可见文案必须国际化，禁止硬编码中英文字符串。
- React 组件内使用 `useTranslate()` 或 `useTranslate('ComponentName')`；非 React 场景使用 `getIntl`。
- 词条源文件是 `src/locale/index.js`（中文源），**禁止只改生成的 JSON 文件**。
- 新增词条后执行 `pnpm intl:once` 或 `pnpm intl:watch` 生成 `zh-CN.json` / `en-US.json`。
- i18n key 末级名控制在 25 字符以内，优先短语义命名，禁止整句式 key。

### API 层

- 禁止在页面组件内直接使用 `axios` 或 `fetch`，统一通过 `src/apis/inter-api/<module>.ts` 封装。
- 使用 `ApiWrapper`（`@/utils/request`）封装请求。
- 所有 API 从 `src/apis/inter-api/index.ts` 统一导出。

### 禁止项

- 禁止在生产代码中遗留 `console.log`（提交前必须清理）。
- 禁止使用 `@ts-ignore` / `@ts-expect-error` 作为永久方案。
- 禁止硬编码颜色值、业务常量（角色名、状态码等）。
- 禁止引入 `@carbon/icons-react` 以外的图标库。
- 禁止引入 Ant Design 以外的 UI 组件库。
- 排障禁止"先试试这个改动"式猜修，必须先定位根因。

### 提交前检查

- `pnpm lint` 无报错。
- 若改动涉及 i18n：`src/locale/index.js` 已更新 → `pnpm intl:once` → 确认 JSON 同步。
- 清理 `console.log` 和临时 `@ts-ignore`。
- 机密配置放 `.env`，禁止提交。

### 停机条件

出现以下情况必须暂停并询问用户：

- API 合同不明确，继续会造成破坏性改动。
- 设计信息不足，无法确认高还原交付。
- 排障连续 3 次假设失败，需先复盘数据流设计。
