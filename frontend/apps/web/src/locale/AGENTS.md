# 前端 locale fallback 规则

本目录的 `zh-CN.json` 和 `en-US.json` 只是应用启动和后端语言包暂不可用时的有限 fallback，不是业务语言包主源。

## 硬规则

- 两份 JSON 由仓库根 `scripts/generate-i18n-fallback.mjs` 生成，禁止手工新增、修改或删除词条。
- 禁止将 `backend/etc/i18n/*.json` 全量复制到本目录。
- 普通业务词条只修改 `backend/etc/i18n/zh_CN.json` 和 `backend/etc/i18n/en_US.json`。
- fallback key 只能来自 `frontend/apps/web/scripts/i18n-fallback-keys.mjs` 的明确白名单，并受 160 个 key 的硬上限约束。
- 不要为了让检查通过而把普通业务 key 加入白名单。

## 正确修改流程

普通业务词条：

1. 更新后端中英文语言包。
2. 从仓库根执行 `node scripts/check-i18n.mjs`。
3. 不修改本目录文件。

确实需要公共启动或通用 UI fallback：

1. 更新后端中英文语言包。
2. 更新 `frontend/apps/web/scripts/i18n-fallback-keys.mjs` 白名单。
3. 从仓库根执行 `node scripts/generate-i18n-fallback.mjs`。
4. 从仓库根执行 `node scripts/check-i18n.mjs`。

若检查报告 `unexpected keys`，应删除越界词条或重新运行生成器；不得恢复旧的全量同步方式。
