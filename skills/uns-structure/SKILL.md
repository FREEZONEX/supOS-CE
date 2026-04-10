---
name: UNS Structure Design
description: Guidelines for designing and validating Unified Namespace (UNS) namespace tree structures — node types, type folders, topic fields, and data type inference.
---

# UNS Structure Design Skill

This skill defines the rules for designing UNS (Unified Namespace) namespace tree structures in the Tier0 platform. The focus is on the **namespace hierarchy**: how PATH and TOPIC nodes are organized, how `topicType` drives data type inference, and how fields are defined.

---

## Node Types

Every node in the namespace tree has a `type` field:

| `type` Value | Role | Description |
|--------------|------|-------------|
| `PATH` | **Folder** | Container node; organizes hierarchy via `children` |
| `TOPIC` | **Leaf** | Terminal data node; carries `fields` for data schema |

---

## Namespace Hierarchy

The namespace is a recursive tree of PATH and TOPIC nodes. The typical pattern:

```
{Version} / {Site} / {Area} / {TypeFolder} / {Topic}
```

Example tree:

```
v1
└── Plant_Name
    └── SMT-Area-1
        ├── State        (PATH, topicType: STATE)
        │   ├── current_job   (TOPIC)
        │   └── productA      (TOPIC)
        ├── Action       (PATH, topicType: ACTION)
        │   ├── start_job     (TOPIC)
        │   └── stop_job      (TOPIC)
        └── Metric       (PATH, topicType: METRIC)
            ├── board_cycle_time  (TOPIC)
            └── boards_metrics    (TOPIC)
```

---

## Type Folders (Core Concept)

A **Type Folder** is a PATH node with `topicType` set. It determines how all descendant TOPICs store data:

| `topicType` | Purpose | Inferred Data Type | Storage |
|-------------|---------|-------------------|---------|
| `STATE` | State/status data with complex payloads | `JSONB_TYPE` | JSONB |
| `ACTION` | Command/action triggers | `JSONB_TYPE` | JSONB |
| `METRIC` | Time-series numeric data | `TIME_SEQUENCE_TYPE` | Time-series DB |

```json
{
  "name": "State",
  "type": "PATH",
  "topicType": "STATE",
  "children": [...]
}
```

### Data Type Inference Rules

TOPICs **do not** need an explicit `dataType`. The backend (`node2vo`) infers it automatically:

1. **Parent `topicType` is `STATE` or `ACTION`** → topic gets `JSONB_TYPE`; declared `fields` become JSON fields internally.
2. **Parent `topicType` is `METRIC`** → topic gets `TIME_SEQUENCE_TYPE`; fields are stored as time-series columns.
3. **`topicType` propagates downward**: nested PATH children inherit the parent's data type.

---

## PATH Node

### Regular PATH (pure folder)

```json
{
  "name": "SMT-Area-1",
  "type": "PATH",
  "children": [...]
}
```

### Type Folder PATH

```json
{
  "name": "Metric",
  "type": "PATH",
  "topicType": "METRIC",
  "children": [...]
}
```

### All PATH Fields

| JSON Key | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | `string` | **Yes** | Folder name (max 63 chars) |
| `type` | `string` | **Yes** | Must be `"PATH"` |
| `topicType` | `string` | No | `STATE` / `ACTION` / `METRIC` — makes this a type folder |
| `children` | `array` | No | Nested PATH or TOPIC nodes |
| `labels` | `string` | No | Comma-separated label names (e.g. `"alpha,beta"`) |
| `alias` | `string` | No | Auto-generated from path if omitted |
| `displayName` | `string` | No | Human-readable display name |

### Constraints

- Names `label` and `template` are **reserved** and cannot be used as folder names.
- Folder name max length: **63 characters**.

---

## TOPIC Node

A TOPIC is a leaf node that defines a data point with typed fields.

### Basic TOPIC (with fields)

```json
{
  "name": "current_job",
  "type": "TOPIC",
  "fields": [
    { "name": "job_id", "type": "INTEGER" },
    { "name": "status", "type": "FLOAT" },
    { "name": "created_at", "type": "DATETIME" }
  ]
}
```

### TOPIC with optional properties

```json
{
  "name": "productA",
  "type": "TOPIC",
  "description": "desc",
  "enableHistory": "TRUE",
  "extendProperties": { "ext1": "value1" },
  "fields": [
    { "name": "product_id", "type": "INTEGER" },
    { "name": "status", "type": "STRING" }
  ]
}
```

### TOPIC referencing a template

```json
{
  "name": "board_cycle_time",
  "type": "TOPIC",
  "templateAlias": "basic_order",
  "enableHistory": "TRUE",
  "mockData": "FALSE"
}
```

### All TOPIC Fields

| JSON Key | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | `string` | **Yes** | Topic name |
| `type` | `string` | **Yes** | Must be `"TOPIC"` |
| `fields` | `array` | Conditional | Field definitions; required if no `templateAlias` |
| `templateAlias` | `string` | No | Reference to a template's `alias`; provides fields |
| `description` | `string` | No | Human-readable description |
| `enableHistory` | `string` | No | `"TRUE"` / `"FALSE"` — persist historical values |
| `mockData` | `string` | No | `"TRUE"` / `"FALSE"` — simulated data flow (METRIC only; forced FALSE for JSONB) |
| `writeData` | `string` | No | `"TRUE"` / `"FALSE"` — allow external write (maps to `accessLevel`) |
| `extendProperties` | `object` | No | Custom key-value extension properties |
| `labels` | `string` | No | Comma-separated label names |
| `alias` | `string` | No | Auto-generated from path if omitted |
| `displayName` | `string` | No | Human-readable display name |

> Boolean-like fields (`enableHistory`, `mockData`, `writeData`) are **strings**, not JSON booleans.

---

## Field Definitions

Each field has at minimum `name` and `type`:

```json
{ "name": "temperature", "type": "DOUBLE" }
```

### Supported Field Types

| Type | Description |
|------|-------------|
| `INTEGER` | 32-bit integer |
| `LONG` | 64-bit integer |
| `FLOAT` | 32-bit floating point |
| `DOUBLE` | 64-bit floating point |
| `BOOLEAN` | true / false |
| `DATETIME` | Timestamp |
| `STRING` | Text (supports optional `maxLen`) |

> Field type lookup is **case-insensitive**. `"int"` is accepted as alias for `INTEGER`.

### Optional Field Properties

| Property | Type | Description |
|----------|------|-------------|
| `maxLen` | `int` | Max length, STRING fields only (ignored on other types) |
| `unique` | `bool` | Unique constraint |
| `index` | `string` | Index definition |
| `displayName` | `string` | Human-readable field name |
| `remark` | `string` | Field documentation |
| `unit` | `string` | Measurement unit |
| `upperLimit` | `float` | Upper bound |
| `lowerLimit` | `float` | Lower bound |
| `decimal` | `int` | Decimal precision |

---

## Alias Generation

If `alias` is omitted (recommended), the backend auto-generates it from the full path: `{parent_path}/{name}`, processed by `PathUtil.GenerateFileAlias()`. This ensures global uniqueness.

---

## Optional: Templates & Labels

Templates and labels are **optional** top-level sections. They are not required for namespace structure definition.

### Templates

Reusable field schemas that TOPICs can reference via `templateAlias`:

```json
{
  "templates": [
    {
      "alias": "basic_order",
      "name": "BasicTemplate",
      "fields": [
        { "name": "order_id", "type": "STRING" }
      ]
    }
  ]
}
```

### Labels

Simple tags for categorizing nodes:

```json
{
  "labels": [
    { "name": "alpha" },
    { "name": "beta" }
  ]
}
```

---

## Import JSON Format

The full import file wraps everything in a single JSON object. Only `namespace` is required:

```json
{
  "namespace": [
    {
      "name": "v1",
      "type": "PATH",
      "children": [
        {
          "name": "Plant_Name",
          "type": "PATH",
          "children": [
            {
              "name": "SMT-Area-1",
              "type": "PATH",
              "children": [
                {
                  "name": "State",
                  "type": "PATH",
                  "topicType": "STATE",
                  "children": [
                    {
                      "name": "current_job",
                      "type": "TOPIC",
                      "extendProperties": { "ext1": "value1" },
                      "fields": [
                        { "name": "job_id", "type": "INTEGER" },
                        { "name": "status", "type": "FLOAT" },
                        { "name": "created_at", "type": "DATETIME" }
                      ]
                    },
                    {
                      "name": "productA",
                      "type": "TOPIC",
                      "description": "desc",
                      "enableHistory": "TRUE",
                      "fields": [
                        { "name": "product_id", "type": "INTEGER" },
                        { "name": "status", "type": "STRING" }
                      ]
                    }
                  ]
                },
                {
                  "name": "Action",
                  "type": "PATH",
                  "topicType": "ACTION",
                  "children": [
                    {
                      "name": "start_job",
                      "type": "TOPIC",
                      "fields": [
                        { "name": "light_turn_on", "type": "BOOLEAN" },
                        { "name": "height_adjust", "type": "FLOAT" }
                      ]
                    },
                    {
                      "name": "stop_job",
                      "type": "TOPIC",
                      "fields": [
                        { "name": "water_pump_off", "type": "BOOLEAN" },
                        { "name": "window_close", "type": "BOOLEAN" }
                      ]
                    }
                  ]
                },
                {
                  "name": "Metric",
                  "type": "PATH",
                  "topicType": "METRIC",
                  "children": [
                    {
                      "name": "board_cycle_time",
                      "type": "TOPIC",
                      "enableHistory": "TRUE",
                      "mockData": "FALSE",
                      "templateAlias": "basic_order"
                    },
                    {
                      "name": "boards_metrics",
                      "type": "TOPIC",
                      "enableHistory": "TRUE",
                      "mockData": "FALSE",
                      "fields": [
                        { "name": "water", "type": "LONG" },
                        { "name": "temp", "type": "LONG" }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

---

## Validation Checklist

When reviewing or creating UNS structures, verify:

- [ ] Every node has `type` = `PATH` or `TOPIC` and a non-empty `name`
- [ ] Type folder PATH nodes set `topicType` to `STATE`, `ACTION`, or `METRIC`
- [ ] TOPIC nodes have `fields` array or `templateAlias` (one is required)
- [ ] METRIC topics use numeric/datetime field types only
- [ ] STATE/ACTION topics can use any field types (stored as JSONB internally)
- [ ] `enableHistory`, `mockData`, `writeData` are strings `"TRUE"` / `"FALSE"`, not booleans
- [ ] `extendProperties` is a flat key-value map
- [ ] Folder names avoid reserved words (`label`, `template`) and are <= 63 chars
- [ ] If `templateAlias` is used, corresponding template must be defined

---

## Backend Source Reference

| File | Role |
|------|------|
| `backend/.../importExport/service/FileData.go` | `FileData` DTO, `node2vo` (import mapping), `uns2DataVo` (export mapping) |
| `backend/.../importExport/service/UnsImport.go` | Import orchestration and error handling |
| `backend/.../importExport/service/jsonstream/StreamedJsonDecoder.go` | Streaming tree→flat JSON parser |
| `backend/.../importExport/service/templates/all-namespace.json` | Canonical import template |
| `backend/internal/types/field_type.go` | `FieldType` enum (INTEGER, STRING, ...) |
| `backend/internal/common/enums/FolderDataType.go` | `FolderDataType` enum (STATE, ACTION, METRIC) |
