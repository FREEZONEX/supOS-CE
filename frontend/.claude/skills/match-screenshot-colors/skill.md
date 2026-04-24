---
name: match-screenshot-colors
description: 分析截图或设计稿中的颜色，将其映射到项目的两级 CSS 变量体系（--supos-t-* 主题级 / --supos-* 应用级）。用户提供截图还原需求、UI 色值匹配、颜色一致性修正时使用。
---

# 截图颜色映射到 CSS 变量

## 使用时机

- 用户提供截图要求还原页面样式
- 需要将设计稿色值对应到项目 CSS 变量
- 页面存在硬编码色值，需替换为 CSS 变量
- 检查颜色与现有主题是否一致

## 两级变量体系

项目使用两级 CSS 变量：

1. **主题级变量**（`--supos-t-*`）— 定义于 `apps/web/src/theme/theme.scss`，基础色盘定义
2. **应用级变量**（`--supos-*`）— 定义于 `apps/web/src/index.scss :root` / `.dark`，语义化映射

**组件样式只允许使用应用级变量（`--supos-*`），禁止直接使用主题级变量（`--supos-t-*`）。**

## 工作流程

### 第 1 步：分析截图颜色

使用 Read 工具读取截图，识别各区域的色值：

- 背景色
- 文字色
- 边框色
- 卡片/容器背景
- 交互元素色（按钮、链接、悬停状态）
- 状态/强调色

### 第 2 步：在主题级找到对应色

读取 `apps/web/src/theme/theme.scss`，找到色盘中与截图色值最接近的变量：

```
截图色 #f4f4f4 → --supos-t-gray-color-10
截图色 #e0e0e0 → --supos-t-gray-color-20
截图色 #161616 → --supos-t-gray-color-100
```

### 第 3 步：找到对应的应用级变量（关键步骤）

**此步骤决定最终使用哪个变量，必须严格执行。**

1. 读取 `apps/web/src/index.scss`
2. **只在 `:root` 段**查找匹配的应用级变量
3. **只在 `.dark` 段**确认暗色主题映射
4. **禁止直接使用** `--supos-t-*` 变量
5. 若无合适的应用级变量，考虑是否需要在 `index.scss` 中新增

建立映射表，只使用应用级变量：

```
背景 #f4f4f4 → var(--supos-craft-bg-color)
边框 #e0e0e0 → var(--supos-border-color)
文字 #161616 → var(--supos-text-color)
卡片背景 #fff → var(--supos-card-color)
```

### 第 4 步：定位页面文件

确认对应的页面文件和 SCSS Module，读取后识别需要更新颜色的元素。

### 第 5 步：应用 CSS 变量

更新 SCSS Module，将硬编码色值替换为 CSS 变量：

```scss
// ✅ 正确
.myComponent {
  color: var(--supos-text-color);
  background: var(--supos-bg-color);
  border: 1px solid var(--supos-border-color);
}

// ❌ 错误
.myComponent {
  color: #161616;
  background: #ffffff;
  border: 1px solid #e0e0e0;
}
```

## 颜色匹配优先级

替换色值时，按以下优先级选择变量：

1. **精确色值匹配** — 找到与原色值完全相同的应用级变量
2. **近似色值匹配** — 找到色值相近的应用级变量（同色阶范围内）
3. **语义匹配** — 仅在前两步均无合适变量时才考虑语义含义

**反例（错误做法）：**

- 原色：`--supos-t-gray-color-50`（#8d8d8d，中灰）
- 错误选择：`--supos-description-card-color`（映射到 gray-70 = #525252，明显偏深）
- 原因：名称语义相似，但实际色值差异过大

**正例（正确做法）：**

- 原色：`--supos-t-gray-color-50`（#8d8d8d，中灰）
- 正确选择：`--supos-icon-disabled-color`（映射到 gray-40 = #a8a8a8，色值相近）

## 处理已有主题级变量

如果发现组件内已直接使用了 `--supos-t-*` 变量：

1. 确认该颜色的用途（背景、文字、边框等）
2. 在 `:root` 中找到语义和色值都匹配的应用级变量
3. 替换为应用级变量

常见替换：

- `var(--supos-t-gray-color-10)` → `var(--supos-craft-bg-color)` 或 `var(--supos-card-bg)`
- `var(--supos-t-gray-color-100)` → `var(--supos-text-color)`
- `var(--supos-t-white-color)` → `var(--supos-card-color)` 或 `var(--supos-bg-color)`

## 常用应用级变量速查

**以下变量均定义于 `apps/web/src/index.scss :root` 段，组件样式应优先使用这些变量。**

### 背景

- `--supos-bg-color` — 页面主背景（亮色主题为白色）
- `--supos-card-color` — 卡片/面板背景（白色）
- `--supos-craft-bg-color` — 工作区背景（#f4f4f4）
- `--supos-header-bg-color` — 头部背景
- `--supos-card-bg` — 备用卡片背景

### 文字

- `--supos-text-color` — 主文字色（亮色主题为 #161616）
- `--supos-anti-text-color` — 反色文字（深色背景上的白字）

### 边框

- `--supos-border-color` — 标准边框
- `--supos-home-border-color` — 首页边框

### 交互状态

- `--supos-menu-hover-color` — 菜单悬停
- `--supos-card-hover-color` — 卡片悬停
- `--supos-menu-active-color` — 菜单选中
- `--supos-theme-color` — 主题色（蓝或绿）

### 其他

- `--supos-nocheck-color` — 未选中/非激活状态
- `--supos-boxshadow-color` — 阴影色

## 输出格式

执行完成后输出：

1. **颜色分析** — 截图中识别到的色值列表
2. **主题变量映射** — 色盘中匹配到的变量
3. **应用变量映射** — 最终使用的语义变量
4. **实施结果** — 更新后的 SCSS，已全部替换为 CSS 变量
5. **验证** — 确认无硬编码色值残留
