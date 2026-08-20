# Anchor 子应用

Anchor（模型库 + 场景）子应用，自 tier0-frontend（云端 SaaS）迁移，由 backend 静态承载在同源 `/anchor/` 路径、主应用菜单以 iframe 嵌入。

## 技术栈说明（与 apps/web 的差异是有意的）

- 本子应用保留源端的 3D 技术栈：**React 19 + three.js**（React 19 是 R3F 9 生态的硬依赖）。
- 不与 apps/web 共享 React 运行时；通过同源 iframe 集成（菜单 `urlType=2` + `openType=0`）。
- `@types/react`/`@types/react-dom` 在本包内固定为 19.x，**不要改回根 catalog（那是主站的 React 18 类型）**。

## UI 规范（必须与 Enterprise 主应用统一）

- **组件库**：Ant Design（`catalog:ui` 同版本）+ `@ant-design/v5-patch-for-react-19`；图标用 `@carbon/icons-react`。
- **主题变量**：`src/theme/ui-vars.css` 与 `src/theme/theme-token.ts` 从 apps/web 同步（文件头有源头标注，改 web 侧后需同步）；业务样式一律使用 `var(--ui-*)`，禁止硬编码色值。
- **明暗主题/语言**：`ThemeBridge` 读取主应用写入的 `localStorage`（`APP_THEME_V2`/`APP_PRIMARY_COLOR`/`APP_LANG`），套用与 web `setThemeRoot` 一致的 html class，并监听 `storage` 事件跟随主应用切换。
- **子页面标题栏**：使用 `.page-title-bar`（对齐 web Flow 详情页：左侧 outlined `Back` + 标题，右侧操作/主按钮）。

## 承载方式

- 生产：`deploy/bin/prepare-web-artifacts.sh` 在 backend 镜像构建前统一构建主前端和 Anchor，并将产物写入 `backend/.build/web`；Dockerfile 继续使用 `backend/` context 将该目录复制进镜像，由 backend 网关静态承载（深链回退到子应用自己的 `index.html`）。
- Vite `base` 固定为 `/anchor/`，react-router `basename` 为 `/anchor`，两者必须与承载路径保持一致。
- `/viewer` 路由是扫码分享页：免登录、可无宿主直开，配套的免鉴权接口见 backend `anchor.api` 的 qr-config/model-file。

## 开发

```bash
pnpm dev:anchor     # 仓库 frontend/ 下执行；http://localhost:5174/anchor/，API 代理到 127.0.0.1:8088
pnpm build:anchor
```
