---
name: tier0-writing-plans-lite
description: 为 Tier0 Edge/Enterprise 多步骤需求编写轻量可执行计划。用于需求较复杂、跨后端接口/service/页面/组件改动、或需要先对齐实施路径再编码的场景。
---

# Tier0 轻量计划编写（Edge/Enterprise）

## 概览

输出一个短而可执行的实施计划，重点是"下一步怎么做"，不是写长文档。  
适用：跨层改动、风险较高、多人协作、需求边界不清时。

## 输出格式

计划建议保存到：`docs/plans/YYYY-MM-DD-<feature>.md`

文档头固定包含：

```markdown
# <功能名> 实施计划

**目标：** <1-2 句>
**范围：** <改哪些，不改哪些>
**风险：** <最多 3 条>
```

正文保持 5-9 步，每步包含：

1. 改动目标
2. 具体文件路径
3. 完成标准
4. 校验命令

## 计划粒度

- 每步 5-20 分钟可完成。
- 每步只做一件事，避免混合任务。
- 路径必须写到文件级（例如 `apps/web/src/apis/inter-api/feature.ts` 或 `apps/web/src/pages/feature/index.tsx`）。

## Tier0 Edge/Enterprise 推荐阶段

1. 接口确认（后端 API 路径/方法/字段确认）
2. Service 层封装（`src/apis/inter-api/`）
3. 页面与状态流对接（store/hook + service）
4. UI 还原与交互细节（组件、样式、空态）
5. i18n 补齐与回归检查

## 约束

- 不强制 TDD、worktree、子代理流程。
- 不写超长"百科式"计划，优先可执行性。
- 计划批准后再进入编码。
