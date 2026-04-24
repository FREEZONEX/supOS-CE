---
name: tier0-page-api-integration
description: 在 Tier0 Edge/Enterprise 前端中完成页面与后端 API 的稳定对接。用于将页面动作接入 src/apis/inter-api/、处理请求状态与错误反馈，并完成端到端页面联调。
---

# Tier0 页面 API 对接（Edge/Enterprise）

## 概览

将页面交互可靠接入后端 API 服务层，保证请求行为可预测、可维护。

## 工作流

1. 明确 API 合同与 endpoint（路径、方法、请求/响应字段）。
2. 先确认 `src/apis/inter-api/<module>.ts` 是否已有对应方法；没有则先补。
3. 在页面中接入 service，显式处理：
   - `loading`（加载中状态）
   - `success`（成功反馈）
   - `error`（失败提示）
4. 请求参数映射集中放在提交逻辑附近，不分散在展示组件里。
5. 若页面已使用 Zustand store，沿用该层承载数据获取与变更。
6. 错误提示复用现有模式（`message.error`、`notification`、行内错误等）。
7. 页面对接产生的用户可见文案（成功/失败提示、空态文案、按钮文案）统一走 i18n key，不直接硬编码字符串。

## 对接规则

### API Service 层约束

- **统一使用 `ApiWrapper`** 封装请求，禁止在页面组件内直接使用 `axios` 或 `fetch`。
- Service 函数放在 `src/apis/inter-api/<module>.ts`，所有 API 从 `src/apis/inter-api/index.ts` 统一导出。
- Service 函数保持薄且确定：只做参数整形 + 请求发送 + 响应映射，不含业务逻辑。
- 方法命名与同文件风格一致（如 `getXxxApi`、`createXxxApi`、`updateXxxApi`、`deleteXxxApi`）。

```typescript
// 示例：src/apis/inter-api/feature.ts
import { ApiWrapper } from '@/utils/request';

const api = new ApiWrapper('/inter-api/supos/feature');

export const getFeatureListApi = async (params: GetFeatureListParams) => {
  return api.get('', { params });
};

export const createFeatureApi = async (data: CreateFeatureRequest) => {
  return api.post('', data);
};
```

### 页面层约束

- 展示型深层组件中不直接发请求。
- 禁止在页面/组件内拼接后端 URL，必须通过 `src/apis/inter-api/` 调用。
- 使用 Zustand store 承载列表数据、分页状态、加载状态时，store action 内直接调用 service。
- 使用 `usePagination` hook 统一管理分页逻辑。

### 状态处理规范

```typescript
// 推荐模式：在 store action 中处理请求
const fetchList = async (params) => {
  set({ loading: true });
  try {
    const res = await getFeatureListApi(params);
    set({ list: res.data, total: res.total });
  } catch (e) {
    message.error(formatMessage({ id: 'feature.fetchError' }));
  } finally {
    set({ loading: false });
  }
};
```

## 校验

- 用浏览器 Network 面板确认请求路径、方法与参数是否正确。
- 验证成功与失败分支 UI 是否都有正确反馈。
- 运行 lint 确认无报错：

```bash
pnpm lint
```

- 提交前确认页面新增用户可见文案未出现硬编码（日志/调试输出除外）。
- 提交前确认 `src/locale/index.js` 新增的 key 已执行 `pnpm intl:once` 同步生成。

## i18n 门禁

- 所有新增/修改的用户可见文案（按钮、提示、标签等）必须有 i18n key。
- 对应 key 必须在 `src/locale/index.js` 中定义（源文件）。
- 运行 `pnpm intl:once` 后 `src/locale/zh-CN.json` 与 `src/locale/en-US.json` 必须同步更新。
- 不要只改生成的 JSON 文件而漏掉 `src/locale/index.js` 源文件。
