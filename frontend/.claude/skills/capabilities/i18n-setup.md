---
name: i18n-setup
description: 设置并验证 Tier0 Edge/Enterprise 前端的 i18n 国际化能力，包含新页面接入、词条维护，以及 react-intl + useTranslate 的正确使用方式。
---

# i18n-setup（Edge/Enterprise）

用于在本仓库中新增、改造、校验国际化能力，确保文案可维护、可协作、可回归。

## 适用场景

- 新页面接入 i18n
- 现有页面补齐 i18n 文案
- 插件暴露语言包
- 确认语言切换行为正确

## 技术栈说明

本仓库使用 **react-intl + Zustand + Ant Design ConfigProvider** 组合：

- 核心状态：`apps/web/src/stores/i18n-store.ts`
- Provider：`apps/web/src/LanguageProvider.tsx`
- 初始化函数：`initI18n`（先加载本地 JSON，再合并后端词条）
- 词条源文件：`apps/web/src/locale/index.js`（中文源）
- 词条生成文件：`apps/web/src/locale/zh-CN.json`、`apps/web/src/locale/en-US.json`

## 必守约束

- React 组件内使用 `useTranslate()` 或 `useTranslate('ComponentName')` 获取格式化函数。
- 非 React 场景使用 `getIntl`，避免在组件里绕过响应式更新。
- 用户可见文案必须国际化，禁止在组件内硬编码中英文字符串。
- `src/locale/index.js` 是源文件，必须先在这里维护词条，再生成 JSON。
- 不要只改 `zh-CN.json` 或 `en-US.json` 而漏掉 `index.js`。
- AI 提取或新增 i18n key 时，末级 key 名尽量控制在 25 个字符以内，优先短语义命名，避免整句式 key。

## 本仓库 i18n 调用约定

### React 组件内

```tsx
import { useTranslate } from '@/hooks';

// 无前缀：直接使用完整 key
const formatMessage = useTranslate();
<h1>{formatMessage({ id: 'feature.pageTitle' })}</h1>;

// 有前缀：ComponentName 作为前缀
const formatMessage = useTranslate('MyComponent');
<button>{formatMessage({ id: 'save' })}</button>;
// 实际查找 key：MyComponent.save
```

### 非 React 场景

```ts
import { getIntl } from '@/stores/i18n-store';

const intl = getIntl();
const text = intl.formatMessage({ id: 'feature.pageTitle' });
```

## 词条维护流程

### 1) 在源文件中新增词条

编辑 `apps/web/src/locale/index.js`，添加中文源词条：

```js
// src/locale/index.js
export default {
  // ... 现有词条
  'feature.pageTitle': '功能页面',
  'feature.save': '保存',
  'feature.cancel': '取消',
  'feature.deleteConfirm': '确定要删除吗？',
};
```

### 2) 生成 JSON 文件

```bash
pnpm intl:once
# 或监听模式（开发期间）
pnpm intl:watch
```

生成结果：

- `apps/web/src/locale/zh-CN.json`（中文）
- `apps/web/src/locale/en-US.json`（英文，需要翻译的留空或填英文）

### 3) 补充英文词条

在 `en-US.json` 中补充英文翻译：

```json
{
  "feature.pageTitle": "Feature Page",
  "feature.save": "Save",
  "feature.cancel": "Cancel",
  "feature.deleteConfirm": "Are you sure you want to delete?"
}
```

### 4) 验证

- 切换语言后确认 UI 文案正确
- 确认无 missing key 警告（浏览器 Console）
- 运行 lint：`pnpm lint`

## 新页面接入 i18n 完整示例

```tsx
import { type FC } from 'react';
import { Button } from 'antd';
import { useTranslate } from '@/hooks';

const FeaturePage: FC = () => {
  const formatMessage = useTranslate('FeaturePage');

  return (
    <div>
      <h1>{formatMessage({ id: 'title' })}</h1>
      <Button onClick={handleSave}>{formatMessage({ id: 'save' })}</Button>
    </div>
  );
};
```

对应 `src/locale/index.js` 词条：

```js
'FeaturePage.title': '功能页面',
'FeaturePage.save': '保存',
```

## 插件语言包

模块联邦插件必须暴露 `./zhCN` 与 `./enUS` 语言包，参考 `plugins/vite.base.ts`。插件内部调用 `useTranslate(REMOTE_NAME)` 确保 key 自动带插件前缀，避免与主应用 key 冲突。

## 注意事项

- `dayjs` 语言在 `initI18n` 中同步切换；新增日期/时间格式化能力时确认是否跟随当前语言。
- Ant Design 组件语言由 `loadAntdLocale` 动态加载，分页、弹窗、空态等默认文案复用已有 key。
- 若扩展新语言（除 `zh-CN`/`en-US` 外），需同时检查：前端本地语言文件、后端语言列表接口、Ant Design locale 映射、插件语言包。

## 验收清单

- `src/locale/index.js` 已定义所有新增词条
- 运行 `pnpm intl:once` 后 `zh-CN.json` 与 `en-US.json` 已同步
- `en-US.json` 已补充英文翻译
- 切换语言可见文案完整正确
- 无 missing key 报错
- 无临时调试代码（如 `console.log`）遗留
