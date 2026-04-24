---
name: tier0-systematic-debugging
description: 在 Tier0 Edge/Enterprise 项目中进行系统化排障。用于构建报错、接口异常、页面行为异常、联调问题等场景，要求先定位根因再修复，禁止直接猜修。
---

# Tier0 系统化排障（Edge/Enterprise）

## 概览

目标：先找根因，再修复。  
原则：没有证据链，不提交修复方案。

## 触发时机

- lint 或 build 报错
- 页面与后端接口联调异常
- 接口返回不符合预期
- 运行时报错或页面白屏
- "修了一个又冒出另一个"的反复问题

## 四阶段流程

### 阶段 1：根因定位

1. 固定复现路径（输入、环境、步骤），并先做基线检查：

```bash
pnpm install
pnpm lint
```

2. 逐行读错误信息和调用栈。
3. 对比近期变更（git diff、依赖、配置、环境变量）。
4. 先按错误类别分流，避免混修：
   - 依赖缺失（`Cannot find module xxx`）
   - TypeScript 类型错误（编译阶段）
   - 接口路径/参数错误（Network 面板确认）
   - 状态管理数据流异常（Zustand store 检查）
   - i18n key 缺失或格式错误
   - 真实业务逻辑缺陷
5. 在关键边界加日志，形成证据链：
   - 页面 → store/hook
   - store/hook → API service
   - API service → 后端接口
   - 接口响应 → 数据映射
6. 明确"坏数据从哪一层开始出现"。

### 阶段 2：模式对比

1. 在同仓库找一个"已工作的同类实现"。
2. 对比 API service 封装、请求参数、响应映射、错误处理。
3. 列出差异并标记最可疑项。

### 阶段 3：假设验证

1. 一次只验证一个假设。
2. 只做最小改动验证假设。
3. 记录验证结果：
   - 成立 → 进入修复
   - 不成立 → 回到阶段 1

### 阶段 4：修复与回归

1. 修根因，不修症状。
2. 运行最小回归集：
   - `pnpm lint`
   - 必要时 `pnpm build:web` 确认无构建错误
   - 浏览器手动验证核心路径
3. 说明修复点与证据链闭环。

## 常见问题排查剧本

以下剧本只用于快速建立证据，不替代"四阶段流程"。

### 1. 依赖/安装问题（pnpm）

- 先确认在 `frontend/` 根目录执行命令：

```bash
pnpm install
```

- 若锁文件/缓存可疑，再执行：

```bash
pnpm clean
pnpm install
```

### 2. 构建或 lint 报错

```bash
pnpm lint
pnpm build:web
```

- 若错误集中在某个文件/模块，先在该范围收敛修复，再整体回归。
- 禁止用 `@ts-ignore/@ts-expect-error` 作为临时止血提交。

### 3. 接口联调失败

- 用 Network 面板确认请求路径、方法、参数是否与后端约定一致。
- 确认 `.env` 或 `config/supos.dev.ts` 中的代理配置是否正确。
- 确认后端 Swagger（`http://localhost:8080/swagger/`）与实际请求是否匹配。
- 确认 `src/apis/inter-api/<module>.ts` 中的接口路径与后端 `.api` 定义一致。

### 4. 开发服务器无法启动/无响应

- 端口冲突：检查 `3000` 等端口占用。
- 环境变量：核对 `.env` 文件的必填项（`API_PROXY_URL` 等）。
- 代理配置：检查 `config/supos.dev.ts` 中的 proxy 设置。

### 5. 页面显示异常/样式错乱

- 检查 CSS 变量是否正确引用（`var(--supos-*)`，不要硬编码色值）。
- 检查 SCSS Module 是否正确导入（`import styles from './index.module.scss'`）。
- 检查 Ant Design 组件版本兼容性。

### 6. i18n key 缺失或文案不正确

- 确认 `src/locale/index.js` 中已定义对应 key。
- 运行 `pnpm intl:once` 重新生成 `zh-CN.json` 与 `en-US.json`。
- 确认 `useTranslate()` 的 prefix 与 key 路径匹配。
- 检查是否混入了硬编码文案（非 i18n key 的中英文字符串）。

## 红线

- 禁止"先试试这个改动"式猜修。
- 禁止一次提交混入多个不相关修复。
- 连续 3 次假设失败，先暂停并复盘数据流或组件设计。
