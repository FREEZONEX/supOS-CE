---
name: tier0-delivery-orchestrator
description: Tier0 Edge/Enterprise 前端需求总控编排技能。用于端到端交付场景：从后端接口确认、service 层封装、页面 API 对接、UI 高还原开发到提交门禁检查。用户提到"总控 skill""端到端交付""全流程开发"时使用。
---

# Tier0 交付总控（Edge/Enterprise）

## 概览

将一个需求按固定阶段串起来执行，默认按以下子技能顺序调度：

1. `$tier0-writing-plans-lite`（复杂需求先产出轻量计划）
2. 接口层确认（后端 go-zero REST API 定义确认）
3. `$tier0-page-api-integration`（service 封装 + 页面对接）
4. `$tier0-ui-high-fidelity-build`（UI 高还原 + i18n）
5. `$tier0-systematic-debugging`（仅在异常或回归失败时触发）

## 输入契约

开始执行前先确认最小输入：

- 目标应用（通常是 `apps/web`）
- 目标页面或路由路径
- 后端 API 接口路径与方法（若接口是新增，先确认 go-zero `.api` 定义）
- UI 参考来源（设计稿、截图或现有页面）
- 验收标准（业务 + 视觉）

若信息不足，只追问缺失的最小字段。

## 通用编辑约束

- 在"代码逻辑未改变"的情况下，不删除、不改写用户已有注释。
- 若因逻辑调整必须改注释，只做最小必要改动，不顺手清洗或重写其他注释。
- 新增注释与文件写回统一保持 UTF-8 编码，避免中文注释乱码。

## 主流程

### 阶段 0：基线同步（强制）

- 拉取最新代码后先确认依赖安装：

```bash
cd frontend && pnpm install
```

- 确认开发环境可正常启动：

```bash
pnpm dev:web
```

### 阶段 1：轻量计划（可选但推荐）

- 需求跨多层（后端接口 + service + 页面 + UI）时，先调用 `$tier0-writing-plans-lite`。
- 计划经确认后再进入开发阶段。

### 阶段 2：接口层确认

- 确认后端已有对应 REST API（路径、方法、请求/响应结构）。
- 若接口是新增，先在 `backend/http/supos/<module>.api` 中确认定义，并由后端同学执行 goctl 生成 handler。
- 确认接口文档或 Swagger 可访问（`http://localhost:8080/swagger/`）。
- 将接口约定落到 `src/apis/inter-api/<module>.ts` 中的函数签名。

### 阶段 3：页面 API 对接

- 调用 `$tier0-page-api-integration`。
- 在 `src/apis/inter-api/<module>.ts` 中封装 API 函数。
- 在页面/store 中接入 service，补齐 `loading/empty/error/success` 状态。
- 对接产生的用户可见文案（成功提示、失败提示、空态等）统一走国际化 key。

### 阶段 4：UI 实现

- 调用 `$tier0-ui-high-fidelity-build`。
- 根据设计源实现或修正页面与组件。
- 保持与同模块页面一致的组件与样式风格。
- 检查新增用户可见文案是否全部走国际化，并同步补齐 `src/locale/index.js`（源文件）。
- 若存在用户可见硬编码文案（日志/调试输出除外），不得进入下一阶段。

### 阶段 5：异常处理（按需）

- 任何阶段出现失败或回归异常，调用 `$tier0-systematic-debugging`。
- 必须先给出根因证据链，再进入修复。

## 输出要求

输出结果必须包含：

- 改动文件列表
- 执行命令列表
- 校验结果
- 未决假设或阻塞项

## 提交前门禁（强制）

- 提交前必须通过 lint 检查：

```bash
pnpm lint
```

- 提交前确认构建不报错（快速验证）：

```bash
pnpm build:web
```

- i18n 门禁：
  - 新增/修改的用户可见文案必须在 `src/locale/index.js` 中有对应 key。
  - 运行 `pnpm intl:once` 确认 `zh-CN.json` 与 `en-US.json` 已同步生成。
- 提交前清理 `console.log` 调试代码。

## 停机条件

出现以下情况先停下并询问用户：

- API 合同不明确，继续会造成破坏性改动。
- 设计信息不足，无法承诺高还原。
- 用户要求与当前 Monorepo 边界冲突。
- 排障连续 3 次假设失败，需先复盘数据流设计。
