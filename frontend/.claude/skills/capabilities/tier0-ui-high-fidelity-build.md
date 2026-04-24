---
name: tier0-ui-high-fidelity-build
description: 在 Tier0 Edge/Enterprise 前端中进行页面与组件高还原开发。用于按设计稿/截图/参考页实现高还原 UI、拆分组件、适配响应式与补齐交互状态。
---

# Tier0 UI 高还原开发（Edge/Enterprise）

## 概览

在不破坏现有技术栈与风格的前提下，高还原交付页面与组件。

## 强制规范：组件复用优先（最高优先级）

> 本节规则优先级高于所有其他指导，任何情况下不得绕过。

### 组件选用优先级

1. **项目已有封装组件**（如 `ProTable`、`usePagination`、`AuthWrapper`）— 优先复用，避免重复造轮子。
2. **Ant Design 5.x 组件**— 满足需求的首选 UI 组件库。
3. **Carbon Design System 图标**（`@carbon/icons-react`）— 项目唯一图标库，禁止引入其他图标库。
4. **自定义封装**— 仅在上述均无法满足需求时才新建组件。

### 常用项目组件速查

```
需要表格？         → ProTable（@/components/pro-table）
需要分页？         → usePagination（@/hooks）
需要国际化？       → useTranslate（@/hooks）
需要响应式？       → useMediaSize（@/hooks）
需要 Tab 生命周期？→ useActivate（@/contexts/tabs-lifecycle-context）
需要时间格式化？   → formatTimestamp（@/utils/format）
需要权限控制？     → AuthWrapper（@/components/auth）
需要 API 请求？    → ApiWrapper（@/utils/request）
需要弹窗确认？     → Popconfirm 或 Modal（antd）
```

### 禁止引入新 UI 库

严格禁止在页面和组件开发中引入 Ant Design 以外的 UI 组件库。

## 工作流

1. 先提炼设计约束：间距、字号、颜色、状态、交互、断点。
2. 在同应用中找最近似页面/组件并复用结构。
3. 开发前先核对组件契约：
   - 优先使用项目已有组件（`ProTable`、`AuthWrapper` 等）。
   - 再使用 Ant Design 5.x 原生组件。
   - 实现前先搜索 `src/components/` 确认是否有可复用实现。
4. 分层实现：
   - 页面骨架
   - 功能区块
   - 可复用组件
5. 优先沿用现有栈：
   - Ant Design 5.x 组件
   - SCSS Modules（`.module.scss`）
   - CSS 变量（`var(--supos-*)`）
6. 补齐状态：`loading/empty/error/disabled/success`。
7. 用并排对比方式微调：优先改间距和行高，再改其余细节。
8. 处理可见文案时严格走国际化（强制）：
   - 使用 `useTranslate()` 或 `useTranslate('ComponentName')` 获取格式化函数。
   - 不在 JSX 中直接写用户可见字符串（按钮、标题、说明、空态、报错等）。
   - 不在模块级静态变量里存放"已翻译文案"；静态配置只存 i18n key，渲染时再取翻译结果。
   - 新增 key 先写入 `src/locale/index.js`，再运行 `pnpm intl:once` 同步生成 JSON 文件。

## 样式约束

颜色规则见 CLAUDE.md `### UI 开发` 节。截图/设计稿色值映射使用 `$match-screenshot-colors`（含变量速查表）。

## 字体与图标规范

- **字体**：使用 IBM Plex Sans（全局默认），禁止在组件内手写 `font-family`。
- **图标**：只使用 `@carbon/icons-react`，示例：`import { Search, Close } from '@carbon/icons-react'`；禁止引入其他图标库。

## 组件结构规范

```
pages/
  feature-name/
    index.tsx              # 主页面组件
    index.module.scss      # 页面样式
    components/            # 页面私有组件
      Component.tsx
      Component.module.scss
```

- 页面私有组件放路由下 `components/`。
- 应用级复用组件放 `src/components/`。
- 使用 `@/` 别名引用本应用模块（映射到 `apps/web/src/`）。

## 组件拆分与注释触发条件

- 满足任一条件时，拆分为独立子组件文件：
  - 子组件超过 40 行且具备独立 props 与渲染逻辑。
  - 配置/映射数据超过 5 项（如状态映射、图标映射、字段组）。
  - 工具函数超过 3 个，或单个函数超过 10 行。
  - 同一块 UI 在 2 个及以上页面复用。
- 仅在"非显而易见逻辑"前加简洁注释（复杂条件、业务约束、workaround），避免注释复述代码字面含义。

## Form 与 Modal 规则（强制）

- `initialValues`、字段规则、字段分组、复杂 `label/placeholder` 等配置优先提取为语义化常量，不内联在 `<Form>` JSX。
- 对动态字段显隐、联动校验、提交前数据整形等逻辑，拆成独立函数并命名，避免在 JSX 中堆叠表达式。

## 代码风格约束

- TypeScript 优先，组件用 `*.tsx`。
- SCSS Modules，class 名使用 kebab-case。
- 使用 `classnames` 库处理条件 class。
- 导入分组：React → 第三方 → 本地组件 → 样式。
- `import { type FC }` 方式导入类型。
- 组件使用 `FC` 类型声明。

## 质量标准

- 不引入相邻页面视觉回归。
- 不新增 UI 库。
- 不引入无说明的硬编码色值。
- 不引入新的硬编码业务文案（日志/调试输出除外）。
- 原逻辑未改动时，不删除或改写用户已有注释。
- 提交前清理 `console.log` 调试代码。
