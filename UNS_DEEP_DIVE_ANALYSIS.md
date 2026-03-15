# UNS (Unified Namespace) Deep-Dive Analysis

## Table of Contents

1. [UNS Data Model](#1-uns-data-model)
2. [MQTT Integration](#2-mqtt-integration)
3. [TimescaleDB Integration](#3-timescaledb-integration)
4. [UNS Business Logic](#4-uns-business-logic)
5. [WebSocket / Real-time Data](#5-websocket--real-time-data)
6. [EventFlow & SourceFlow](#6-eventflow--sourceflow)
7. [Kong Integration](#7-kong-integration)
8. [Skills / UNS Validation](#8-skills--uns-validation)
9. [Data Flow Summary](#9-complete-data-flow)
10. [Missing Features & TODOs](#10-missing-features--todos)
11. [Data Integrity Issues](#11-data-integrity-issues)
12. [Error Handling Analysis](#12-error-handling-analysis)

---

## 1. UNS Data Model

### 1.1 Core Table: `uns_namespace`

**File:** `backend/internal/repo/relationDB/modelUNS.go` (lines 31-71)

The central data structure is `UnsNamespace`, mapped to the `uns_namespace` PostgreSQL table:

| Column | Type | Description |
|--------|------|-------------|
| `id` | `bigint` (PK) | Snowflake-generated ID |
| `lay_rec` | `text NOT NULL` | Layer record path (e.g., `0/1234/5678`) for tree traversal |
| `alias` | `varchar(128) UNIQUE` | Globally unique alias identifier |
| `parent_alias` | `varchar(128)` | Parent node alias |
| `name` | `varchar(512)` | Display name |
| `path` | `text NOT NULL` | Full topic path (e.g., `AMS/ERP/State/material`) |
| `path_type` | `int2 NOT NULL` | 0=Folder, 1=Template, 2=File |
| `data_type` | `int2` | 1=TimeSeries, 2=Relational, 3=Calculation, 5=Alarm, 6=Aggregate, 7=Reference |
| `fields` | `json` | Array of `FieldDefine` objects (schema definition) |
| `status` | `smallint DEFAULT 1` | 1=Active, used for soft-delete |
| `protocol` | `varchar(2000)` | Protocol configuration JSON |
| `data_src_id` | `int2` | Data source type ID (TimescaleDB, PostgreSQL, etc.) |
| `ref_uns` | `jsonb DEFAULT '{}'` | Reference to other UNS entries (map of id→count) |
| `refers` | `json` | Instance field references for calculations |
| `expression` | `varchar(255)` | Calculation expression |
| `table_name` | `varchar(190)` | Physical table name in TimescaleDB |
| `parent_id` | `bigint` | Parent node ID |
| `model_id` | `bigint` | Template/model ID |
| `with_flags` | `integer DEFAULT 0` | Bitfield: dashboard, flow, save2db, access level, subscribe |
| `extend` | `jsonb DEFAULT '{}'` | Extensible metadata |
| `label_ids` | `jsonb` | Label associations (map of labelId→labelName) |
| `mount_type` | `int2 DEFAULT 0` | Mount source type (e.g., collector) |
| `mount_source` | `varchar(256)` | External mount source identifier |
| `pathash` | `INTEGER` | Hash of path for fast lookups |
| `fields_text` | `TEXT` | Flattened fields text with GIN trigram index |

**Migration SQL:** `backend/internal/repo/relationDB/migrations_sqls/uns.sql` (lines 1-340)

### 1.2 Tree Structure: `lay_rec` Pattern

The tree hierarchy is encoded via `lay_rec` — a materialized path using ID segments separated by `/`:

```
Root:    lay_rec = "0"
Child:   lay_rec = "0/12345"
Grand:   lay_rec = "0/12345/67890"
```

This enables efficient subtree queries using `LIKE 'prefix/%'` with PostgreSQL indexes. The custom SQL function `nextIdLong(prefix, lay_rec)` extracts the next-level child ID from the `lay_rec` path for paginated tree queries.

**File:** `backend/internal/repo/relationDB/unsNamespace_query_tree.go` (lines 83-136)

### 1.3 Node Types

| `path_type` | Constant | Description |
|-------------|----------|-------------|
| 0 | `PathTypeDir` | Folder/directory node |
| 1 | `PathTypeTemplate` | Model/template definition |
| 2 | `PathTypeFile` | Leaf data node (topic) |

### 1.4 Data Types

| `data_type` | Constant | Storage |
|-------------|----------|---------|
| 1 | `TIME_SEQUENCE_TYPE` | TimescaleDB hypertable with numeric fields |
| 2 | `RelationType` | PostgreSQL relational table |
| 3 | `CalculationHistType` | Calculated from other topics (historical) |
| 5 | `AlarmRuleType` | Alarm rule definition |
| 6 | `AggregateType` | Aggregated data |
| 7 | `JsonbType` | JSONB storage (STATE/ACTION payloads) |

### 1.5 Supporting Tables

| Table | File Location | Purpose |
|-------|--------------|---------|
| `uns_label` | `modelUNS.go:406-419` | Tag/label definitions |
| `uns_label_ref` | `modelUNS.go:468-478` | Many-to-many: labels ↔ UNS nodes |
| `uns_attachment` | `modelUNS.go:483-496` | File attachments per UNS alias |
| `uns_alarms_data` | `modelUNS.go:500-515` | Alarm events (current value, limit, status) |
| `uns_alarms_handler` | `modelUNS.go:519-531` | Alarm acknowledgment handlers |
| `uns_tag` | `modelUNS.go:688-699` | Topic-to-tag references |
| `uns_mount` | `modelUNS.go:593-611` | External mount configurations |
| `uns_mount_extend` | `modelUNS.go:574-588` | Extended mount metadata |
| `uns_history_delete_job` | `modelUNS.go:553-569` | Async cleanup jobs for deleted topics |
| `uns_webhook` / `uns_webhook_action` | `modelUNS.go:701-751` | Webhook integrations |
| `uns_export_record` | `modelUNS.go:537-548` | Export job tracking |
| `uns_person_config` | `modelUNS.go:616-627` | Per-user preferences (language) |
| `uns_dashboard` / `uns_dashboard_ref` | `uns.sql:52-91` | Dashboard definitions and UNS linkage |

### 1.6 DTOs

**File:** `backend/internal/common/dto/uns.go`

- **`UpdateUnsDto`** (lines 11-87): Main DTO for create/update operations. Carries model fields, parent references, data type, field definitions, calculation expressions, protocol settings, flags (addFlow, addDashboard, save2db), frequency, alarm rules, labels, and access level.
- **`SimpleUnsInstance`** (lines 191-236): Lightweight DTO for removal operations.
- **`SaveDataDto`** (lines 252-289): DTO for batch data persistence.
- **`BatchRemoveUnsDto`** (lines 309-316): DTO for batch alias-based deletion.

**File:** `backend/internal/common/dto/UnsSearchCondition.go`

- **`UnsSearchCondition`** (lines 5-35): Rich search filter supporting keyword, alias, path, template, labels, data type, time ranges, and extend fields.

---

## 2. MQTT Integration

### 2.1 MQTT Client Setup

**File:** `backend/internal/repo/event/subDev/mqtt.go` (lines 1-57)

- Uses EMQX shared subscription (`$share/uns/#`) at QoS 1
- Client wraps `clients.MqttClient` from the `unitedrhino` library
- Subscribes to all topics using wildcard `#`
- Shared subscription prefix: `$share/uns.rpc/` for distributed load balancing

### 2.2 Connection Management

**File:** `backend/internal/adapters/msg_consumer/UnsMessageConsumer.go` (lines 335-358)

Connection initialization occurs on `ContextRefreshedEvent`:

```go
func (u *UnsMessageConsumer) OnEventContextRefreshedEvent10(ev *event.ContextRefreshedEvent) {
    // Only if MQTT mode is configured
    if sv.Config.DevLink.Mode == "mqtt" {
        go func() {
            cli, er := subDev.NewMqttClient(&sv.Config.DevLink.Mqtt, u)
            if er != nil {
                // Exponential backoff retry: 5s, 10s, 20s, 40s... capped at 60s
                for i := int64(5); ; i <<= 1 {
                    if i < 0 { i = 60 }
                    time.Sleep(time.Duration(i) * time.Second)
                    cli, er = subDev.NewMqttClient(...)
                    if cli != nil && er == nil { break }
                }
            }
            cli.SubscribeAll()
        }()
    }
}
```

**Reconnection strategy:** Infinite retry with exponential backoff (5s → 10s → 20s → ... → capped at 60s). The bit-shift `i <<= 1` with overflow check `if i < 0` ensures it caps.

**Issue:** No maximum retry limit. If EMQX is permanently down, this goroutine loops forever.

### 2.3 Message Processing Pipeline

**File:** `backend/internal/adapters/msg_consumer/UnsMessageConsumer.go` (lines 42-88)

The `OnMsg` method processes each MQTT message:

1. **Topic Resolution** (lines 43-64):
   - If topic starts with a digit → try parse as ID → lookup by ID
   - If `UseAliasAsTopic` → lookup by alias, fallback to path
   - Otherwise → lookup by path, fallback to alias
   - Uses `UnsDefinitionService` cache (see §2.4)

2. **Payload Parsing** (line 75):
   - Trims leading whitespace
   - Parses JSON into `[]map[string]string` via `parseJsonList`

3. **Data Processing** (lines 81-84):
   - `procDataAndSendWs()` → field validation, WebSocket push, calculation trigger
   - `sendData()` → enqueue to disk queue for persistence

4. **Performance Tracking** (lines 85-87):
   - Logs slow messages (>500ms processing time)

### 2.4 UNS Definition Cache

**File:** `backend/internal/adapters/msg_consumer/UnsDefinitionServiceImpl.go` (lines 21-206)

- Uses `ccache` (LRU cache) with 1.1M max entries, 64 buckets
- Three key prefixes: `a:` (alias), `p:` (path), raw ID (base-36)
- Cache TTL: 10 minutes for definitions, 13 minutes for alias/path→ID mappings
- **Double-checked locking** with 1000 RWMutex shards for concurrent safety
- Negative caching: 1-2 minute TTL for missing entries
- Auto-extends TTL when item has <10s remaining (prevents mid-use expiry)
- Invalidated on: `BatchCreateTableEvent`, `RemoveTopicsEvent`, `UpdateInstanceEvent`

### 2.5 Message Field Validation

**File:** `backend/internal/adapters/msg_consumer/json_uns_filter.go` (lines 15-132)

`filterMsgByUns()` validates each field against the UNS schema:

| Field Type | Validation |
|-----------|-----------|
| INTEGER | Range: `MinInt32` to `MaxInt32` |
| LONG | Range: `MinInt64` to `MaxInt64` |
| FLOAT | Range: `-MaxFloat32` to `MaxFloat32` |
| DOUBLE | Range: `-MaxFloat64` to `MaxFloat64` |
| STRING | Max length check (`maxLen` from field definition) |
| DATETIME | Parses numeric timestamps (milliseconds) or date strings |
| BOOLEAN | Accepts `true`/`false` literals, coerces other values |

On validation failure, a **quality code (QoS)** is set:
- `0x400000000000000` = type error
- `0x80000000000000` = out-of-range

Single-field messages get auto-mapped if only one non-system field exists.

### 2.6 Disk Queue (Data Sink)

**File:** `backend/internal/adapters/msg_consumer/UnsQueueDataSinkService.go` (lines 22-130)

Messages are buffered through a **disk-backed queue** before persistence:

- Queue location: `{RootPath}/queue`
- Max message size: 4 MB
- Max file size: 64 MB
- Sync interval: 5 seconds
- **Batch processing**: Accumulates up to 10,000 rows or waits 1 second (whichever comes first)
- Messages encoded in **custom binary format** (not JSON) for efficiency

**Binary codec:** `backend/internal/adapters/msg_consumer/msg_codec.go`
- Format: `[count:4][{unsId:8, dataSrcId:2, dataCount:4, [{mapSize:2, [{keyLen:2, key, valLen:2, val}]}]}]`
- Uses big-endian encoding, zero-copy string conversion

**TODO (line 43):** Disk-full handling is unimplemented — `queue.Put()` error is silently ignored.

### 2.7 Realtime Calculation Service

**File:** `backend/internal/adapters/msg_consumer/UnsRealtimeCalcService.go` (lines 1-22)

**ENTIRELY UNIMPLEMENTED.** The `TryCalculate()` method:
- Checks if `def.RefUns` is populated (indicating a calculated field)
- Returns immediately with `nil` — **no actual calculation logic exists**
- Contains `//TODO 实时计算` marker

This means **calculated/derived UNS topics do not work at runtime**.

---

## 3. TimescaleDB Integration

### 3.1 Architecture: Single Hypertable + Views

**File:** `backend/internal/adapters/timescaledb/sql_table.go` (lines 11-73)

All time-series UNS data is stored in a **single physical hypertable** `uns_timeserial`:

```sql
CREATE TABLE IF NOT EXISTS public."uns_timeserial" (
    "tag" int8 NOT NULL,           -- UNS ID (partitioning key)
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),  -- Time column
    "_qos" int8 DEFAULT 0,         -- Quality code
    "double_1" float8 NULL,        -- Dynamic typed columns
    -- ... more columns added dynamically ...
    PRIMARY KEY ("tag", "timeStamp")
);
```

**Hypertable configuration:**
- Chunk interval: 2 hours
- Hash partitioning on `tag` column: 50 partitions
- Compression: enabled, segment by `tag`, ordered by `timeStamp DESC`
- Compression policy: after 1 hour, every 2 hours
- Retention policy: 2 years

### 3.2 View-per-Topic Strategy

Each UNS topic gets a **SQL VIEW** over the shared hypertable:

**File:** `backend/internal/adapters/timescaledb/sql_view.go` (lines 11-81)

```sql
CREATE OR REPLACE VIEW "my_topic_alias" AS
SELECT "timeStamp", "tag",
       "double_1"::real AS "temperature",
       "str_1" AS "status"
FROM uns_timeserial
WHERE "tag" = 12345;
```

### 3.3 Dynamic Column Mapping

**File:** `backend/internal/adapters/timescaledb/sql_field_mapping.go` (lines 17-192)

Fields are mapped to physical columns using a **typed prefix + sequential number** pattern:

| UNS Field Type | Physical Column Prefix | PostgreSQL Type |
|----------------|----------------------|-----------------|
| INTEGER | `long_N` | `int4` |
| LONG | `long_N` | `int8` |
| FLOAT | `double_N` | `float4` |
| DOUBLE | `double_N` | `float8` |
| BOOLEAN | `bool_N` | `boolean` |
| DATETIME | `date_N` | `timestamptz` |
| STRING | `str_N` | `varchar` |

New columns are added to the hypertable via `ALTER TABLE ADD COLUMN IF NOT EXISTS` when new field types are needed.

When a physical column is **reused** by a new topic (column previously existed), a `CASE WHEN` expression prevents returning stale data from before the reassignment:
```sql
CASE WHEN "timeStamp" < to_timestamp(1710500000) THEN NULL ELSE "double_3" END AS "pressure"
```

### 3.4 Data Persistence Pipeline

**File:** `backend/internal/adapters/timescaledb/tsdb_persistence.go` (lines 21-284)

1. **Preprocessing** (`tsdb_preprocess.go`):
   - Groups data by `[unsId, timestamp]` key
   - Splits into "conflict" (timestamp ≤ last known) and "normal" (new timestamps) buckets
   - Merges duplicate timestamps

2. **Normal path**: Direct `COPY FROM` into `uns_timeserial`
3. **Conflict path**: 
   - Creates temp table (`LIKE uns_timeserial EXCLUDING INDEXES`)
   - `COPY FROM` into temp table
   - `INSERT INTO ... SELECT * FROM temp ON CONFLICT DO UPDATE SET ...`
   - Transaction auto-drops temp table on commit

4. **Error recovery** (`shouldRetry`, line 215-228):
   - `42P01` (undefined table) → recreate table + retry
   - `42703` (undefined column) → fix table schema + retry
   - `08xxx` (connection errors) → retry without schema fix

5. **Connection timeout**: 15 minutes for persistence operations

### 3.5 View Sync: `GenerateSyncSQLs`

**File:** `backend/internal/adapters/timescaledb/sql_sync.go` (lines 9-52)

The sync process when creating/updating UNS topics:
1. Analyze required physical fields across all topics
2. If hypertable doesn't exist → `CREATE TABLE` + `create_hypertable()`
3. If new field types needed → `ALTER TABLE ADD COLUMN`
4. For each topic → `CREATE OR REPLACE VIEW`

**View conflict handling** (`TsdbPersistentService.go` lines 157-274):
- Error `42P16` (invalid view definition) → drop + recreate view
- Error `42809` (wrong object type, table exists where view expected) → rename existing table to `bk_<alias>`, create view, migrate data from backup, drop backup

---

## 4. UNS Business Logic

### 4.1 Add/Create Service

**File:** `backend/internal/logic/supos/uns/uns/service/UnsAddService.go` (lines 32-589)

`CreateModelAndInstancesInner()` is the core creation method:

1. **Auto-categorization** (line 66): Optionally appends category folders
2. **Parameter initialization** (`initParamsUns`): Resolves paths, validates structure
3. **Existing UNS lookup** (lines 76-114): Queries by alias and ID to detect updates vs inserts
4. **Dependency sorting** (lines 123-138): Sorts folders by parent→child dependency using topological sort; detects circular dependencies
5. **ID assignment** (`trySetId`): Generates snowflake IDs for new entries
6. **LayRec computation** (`setLayRecAndPath`): Computes materialized path strings
7. **Batch save** (`saveBatchAndSendEvent`, lines 344-450):
   - Begins transaction
   - Creates labels
   - Creates physical tables via `IPersistentService.Save()`
   - `MultiInsert` / `MultiUpdate` via GORM (batch size 1000)
   - Publishes `BatchCreateTableEvent`
   - Removes stale entries
   - Commits or rolls back

**Upsert strategy:** Uses `ON CONFLICT DO NOTHING` for inserts, `ON CONFLICT (id) DO UPDATE` for multi-updates.

### 4.2 Update Service

**File:** `backend/internal/logic/supos/uns/uns/service/UnsUpdateService.go` (lines 16-137)

- `UpdateDetail`: Loads existing record by ID or alias, merges flags (dashboard, flow, save2db, access level) using bitwise operations, delegates to `CreateModelInstance`
- `UpdateName`: Renames a node preserving all other fields
- `SubscribeModel`: **Not implemented** (empty method body, line 134-136)

### 4.3 Remove Service

**File:** `backend/internal/logic/supos/uns/uns/service/UnsRemoveService.go` and `uns_remove_helper.go`

Removal is **soft-delete** based (`LogicDeleteByIds` sets `status != 1`):

1. Separates input into folders, files, and templates
2. **Folders**: Exports subtree via CSV streaming, deletes in batches of 1000
3. **Files**: Direct batch deletion
4. **Templates**: Finds all instances using the template, removes them, then nullifies `model_id` on remaining folders
5. Publishes `RemoveTopicsEvent` to notify subscribers (TimescaleDB drops views/tables, SourceFlow deletes Node-RED flows, WebSocket cleans up subscriptions)

**Issue (line 179):** Reference checking before deletion is marked `//TODO 引用检查` — meaning UNS nodes referenced by other calculated/aggregate topics can be deleted without warning.

### 4.4 Tree Query

**File:** `backend/internal/repo/relationDB/unsNamespace_query_tree.go`

- Uses raw SQL with `nextIdLong()` function for next-level pagination
- Supports search types: 1=UNS (name+alias), 2=with labels, 3=with templates
- Keyword search uses `ILIKE` on path and `LIKE` on alias
- Filter by `subscribe_enable` flag
- Counts children by path type (folders vs files) using `STRING_AGG`
- System data excluded: `id > 1000`

### 4.5 Export/Import

**File:** `backend/internal/repo/relationDB/unsNamespace_export.go` (referenced)

Uses CSV streaming with `DoExportBatch()` pattern:
- Writes UNS records to `io.Writer` in batches
- Consumer callback processes each batch
- Used for both export and bulk delete operations

---

## 5. WebSocket / Real-time Data

### 5.1 Architecture

**File:** `backend/internal/logic/supos/uns/uns/service/WebsocketService.go` (lines 26-616)

`WebsocketService` manages four concurrent subscription maps:

| Map | Key | Value | Purpose |
|-----|-----|-------|---------|
| `sessions` | sessionId (string) | `*WsSubscription` | All active sessions |
| `idToSessionsMap` | unsId (int64) | `*sync.Map` of sessionIds | Subscribe by UNS ID |
| `topicToSessionsMap` | topic (string) | `*sync.Map` of sessionIds | Subscribe by topic path |
| `aliasToSessionsMap` | alias (string) | `*sync.Map` of sessionIds | Subscribe by alias (CMD_SUB) |

All maps use `sync.Map` for lock-free concurrent access.

### 5.2 Connection Lifecycle

1. **Connect** (`TryAddSession`, line 92): Checks session limit, creates `WsSubscription`
2. **Subscribe** (`HandleSessionConnected`, line 113):
   - URL params: `?id=123&topic=path/to/topic`
   - `?globalTopology=true` for topology subscriptions
   - Initial message pushed immediately on subscribe
3. **Message handling** (`HandleCmdMsg`, line 206):
   - `/send?t=alias&body=payload` → forward to MQTT consumer (TODO: not connected)
   - JSON `{"head":{"cmd":1},"data":{"sub_real_value":{...}}}` → alias subscriptions
4. **Disconnect** (`HandleSessionClosed`, line 277): Cleans all maps

### 5.3 Message Push Flow

When `UnsMessageConsumer` processes a message:
1. `sendToWebsocket()` → `WebsocketService.SendMessage()`
2. `SendMessage()` routes by UNS ID or topic path
3. `onWsMsg()` checks if field values need DB query (first message for a field)
4. `processWsMsg()` builds `TopicMessageInfo` JSON with typed values, timestamps, and error messages
5. Write lock per subscription prevents concurrent write corruption

### 5.4 Topology WebSocket

- `topologySessions` tracks sessions subscribing to network topology
- On `UnsTopologyChangeEvent` → broadcasts to all topology subscribers
- Initial topology data from `UnsTopologyService.GetLastMsg()`

### 5.5 Issues

- **Line 217:** `/send` command handler has `// TODO: Call topicMessageConsumer.onMessageByAlias` — WebSocket-to-MQTT write path is broken
- **Line 255:** `// TODO: Push real-time data` for alias subscriptions — initial data push after CMD_SUB is not implemented
- **Session count check (line 96):** Race condition in `TryAddSession` — counts sessions non-atomically before adding. Comment acknowledges this: "simplified approach"
- **No heartbeat/ping-pong:** No mechanism to detect stale WebSocket connections

---

## 6. EventFlow & SourceFlow

### 6.1 Architecture

EventFlow and SourceFlow are wrappers around Node-RED flow management. **EventFlow delegates entirely to SourceFlow** with a different flow type constant:

**File:** `backend/internal/logic/supos/eventflow/createEventFlowLogic.go` (lines 33-48)
```go
func (l *CreateEventFlowLogic) CreateEventFlow(req *types.EventFlowCreateReq) (string, error) {
    return sourceflow.NewCreateSourceFlowLogic(l.ctx, l.svcCtx).
        CreateFlowWithType(srcReq, constants.FlowTypeEVENTFLOW)
}
```

### 6.2 Flow Types

| Type | Constant | Node-RED Instance |
|------|----------|-------------------|
| SourceFlow | `FlowTypeNODERED` | `svcCtx.SourceNodeRed` |
| EventFlow | `FlowTypeEVENTFLOW` | `svcCtx.EventNodeRed` |

### 6.3 SourceFlow CRUD

**File:** `backend/internal/logic/supos/sourceflow/createSourceFlowLogic.go` (lines 37-82)

- **Create**: Validates name uniqueness per flow type, generates snowflake ID, inserts `NoderedSourceFlow` record with `FlowStatusDraft`
- **Deploy** (`deploySourceFlowLogic.go`): Marshals flow JSON, delegates to `flowcommon.DeployFlow()` which interacts with Node-RED API
- **Delete**: Removes Node-RED runtime flow via `DELETE /flow/{id}`, then deletes DB records

### 6.4 Auto-Provisioning on UNS Create

**File:** `backend/internal/logic/supos/sourceflow/event_sub.go` (lines 78-134)

`SourceFlowService.OnEventBatchCreateTableEvent()`:
1. Filters UNS files with `AddFlow == true`
2. Loads `relational-emqx.json.tpl` template
3. Renders template with `$uns_path`, `$model_alias`, `$payload` variables
4. `$payload` is built from field definitions with appropriate random generators per type
5. Creates `NoderedSourceFlow` record
6. Auto-deploys to Node-RED

### 6.5 Auto-Cleanup on UNS Delete

**File:** `backend/internal/logic/supos/sourceflow/event_sub.go` (lines 246-312)

`SourceFlowService.OnEventRemoveTopicsEvent()`:
1. Collects aliases from deleted topics
2. Finds associated flows by alias
3. Deletes Node-RED runtime flow via API
4. Removes DB records and top markers

### 6.6 Import/Export

**File:** `backend/internal/logic/supos/sourceflow/NodeRedFlowImport.go` (lines 22-340)

Import handles a JSON structure with sections:
- `flows` → `NoderedSourceFlow` records
- `nodes` → Direct Node-RED node definitions (deployed via `POST /flows`)
- `unsRefs` → Flow-to-UNS linkage records
- `tags` → Reference tags (sent to `/nodered-api/batchSave/tags`)

Uses streaming JSON decoder with `io.Pipe()` for memory-efficient processing.

### 6.7 Node-RED Proxy

**File:** `backend/internal/logic/supos/eventflow/proxyEventFlowsLogic.go` (lines 37-228)

- Loads flow from DB (draft JSON preferred over runtime)
- Falls back to Node-RED API: `GET /flow/{id}` → extracts `nodes` array
- Ensures tab node exists for the flow
- Appends global config nodes from `GET /flow/global`
- Returns combined nodes with version `rev` from Node-RED

---

## 7. Kong Integration

### 7.1 Status

**No `backend/internal/adapters/kong/` directory exists.** The glob search returned 0 files.

Kong-related logic stubs exist in `backend/internal/logic/supos/uns/kong/`:

| File | Status |
|------|--------|
| `pageListLogic.go` | `// todo: add your logic here and delete this line` |
| `confirmLogic.go` | `// todo: add your logic here and delete this line` |
| `updateLogic.go` | `// todo: add your logic here and delete this line` |
| `routeListLogic.go` | `// todo: add your logic here and delete this line` |

**Kong integration is entirely unimplemented.** Route registration, service discovery, and API gateway management are all placeholder stubs.

---

## 8. Skills / UNS Validation

### 8.1 UNS Structure Design

**File:** `skills/uns-structure/SKILL.md`

Defines the **three-tier terminal hierarchy**:
```
{Business Parent} / {Type Folder} / {Topic}
```

Type folders:
- `METRIC` → `TIME_SEQUENCE_TYPE` (numeric time-series)
- `STATE` → `JSONB_TYPE` (complex state payloads)
- `ACTION` → `JSONB_TYPE` (command triggers)

Alias generation: `Current_Alias = Parent_Alias + "_" + Current_Name`

### 8.2 Validation Rules

**File:** `skills/uns-validation/SKILL.md`

| Rule | Description | Implemented? |
|------|-------------|-------------|
| V1: No Orphan Topics | Every topic must be under a type folder | **NO** — no code enforces this |
| V2: METRIC type check | METRIC must use TIME_SEQUENCE_TYPE, no json field | **NO** — `UnsAddService` has `//TODO 计算，引用，聚合等类型的 校验和处理` |
| V3: STATE/ACTION type check | Must use JSONB_TYPE with json field | **NO** — same TODO |
| V4: Alias Path Inheritance | `node.alias == parent.alias + "_" + node.name` | **PARTIAL** — `PathUtil.GenerateAliasWithRandom()` exists but doesn't enforce inheritance |
| V5: No INFO type | dataType must not be INFO | **NO** — no validation code exists |
| Alias uniqueness | No duplicate aliases | **YES** — enforced by `UNIQUE INDEX idx_uns_spacex_alias` in PostgreSQL |

**Key finding:** The validation rules defined in the skills documents are almost entirely **unimplemented in code**. The only enforced constraint is alias uniqueness via a database unique index.

---

## 9. Complete Data Flow

### 9.1 Ingestion Pipeline (Source → UNS → Sink)

```
[MQTT/EMQX]
    │ $share/uns/# (QoS 1)
    ▼
[UnsMessageConsumer.OnMsg()]
    │ 1. Topic resolution (alias/path/ID lookup via cached definitions)
    │ 2. JSON parsing → []map[string]string
    │ 3. Field validation (type checking, range checking)
    │ 4. Timestamp merge (for time-series: dedup same-timestamp records)
    │ 5. Update field LastValue in-memory cache
    ▼
[procDataAndSendWs()]
    ├─→ [WebsocketService.SendMessage()] → Push to subscribed WebSocket clients
    ├─→ [UnsRealtimeCalcService.TryCalculate()] → (NOT IMPLEMENTED)
    ▼
[UnsQueueDataSinkService.Sink()]
    │ Binary encode → Disk queue (64MB files, 4MB max msg)
    ▼
[UnsQueueDataSinkService.fetchData()]  (background goroutine)
    │ Batch: max 10,000 rows or 1s timeout
    │ Route by DataSrcId to appropriate IPersistentService
    ▼
[TsdbPersistentService.Persistent()]  (or PostgreSQL equivalent)
    │ 1. Preprocess: group by [unsId, timestamp], split conflict/normal
    │ 2. Normal: COPY FROM → uns_timeserial
    │ 3. Conflict: temp table → COPY → INSERT ON CONFLICT UPDATE
    ▼
[TimescaleDB: uns_timeserial hypertable]
    │ Compressed after 1 hour
    │ Retained for 2 years
    ▼
[SQL Views: one per UNS topic]
    │ SELECT ... FROM uns_timeserial WHERE tag = {unsId}
```

### 9.2 UNS CRUD → Side Effects

```
[API: Create UNS Topic]
    ▼
[UnsAddService.CreateModelAndInstancesInner()]
    │ 1. Validate & resolve params
    │ 2. Sort by dependency
    │ 3. Assign IDs, compute layRec
    │ 4. BEGIN TRANSACTION
    │ 5. Create physical tables (IPersistentService.Save)
    │ 6. MultiInsert/MultiUpdate in PostgreSQL
    │ 7. Publish BatchCreateTableEvent
    │ 8. COMMIT
    ▼
[Event Subscribers]
    ├─→ [TsdbPersistentService] → CREATE VIEW, ALTER TABLE
    ├─→ [SourceFlowService] → Auto-create Node-RED flow (if AddFlow=true)
    ├─→ [UnsDefinitionService] → Invalidate cache
    └─→ [WebsocketService] → (no action on create)

[API: Delete UNS Topic]
    ▼
[UnsRemoveService.Remove()]
    │ 1. Soft-delete (status != 1)
    │ 2. Publish RemoveTopicsEvent
    ▼
[Event Subscribers]
    ├─→ [TsdbPersistentService] → DROP VIEW / DROP TABLE
    ├─→ [SourceFlowService] → Delete Node-RED flow + runtime
    ├─→ [UnsDefinitionService] → Invalidate cache
    └─→ [WebsocketService] → Clean up subscriptions
```

---

## 10. Missing Features & TODOs

| Location | TODO | Impact |
|----------|------|--------|
| `msg_consumer/UnsRealtimeCalcService.go:20` | `//TODO 实时计算` | **CRITICAL**: Calculated/derived topics produce no output |
| `msg_consumer/UnsQueueDataSinkService.go:43` | `//TODO: 磁盘满的处理` | **HIGH**: Disk full → silent data loss |
| `uns/service/WebsocketService.go:127` | `// TODO: Handle file import` | **MEDIUM**: File import via WebSocket not working |
| `uns/service/WebsocketService.go:217` | `// TODO: Call topicMessageConsumer.onMessageByAlias` | **HIGH**: WebSocket → MQTT write path broken |
| `uns/service/WebsocketService.go:255` | `// TODO: Push real-time data` | **HIGH**: CMD_SUB initial data push missing |
| `uns/service/UnsAddService.go:160` | `//TODO 计算，引用，聚合等类型的 校验和处理` | **HIGH**: No validation for calc/ref/aggregate types |
| `uns/service/uns_remove_helper.go:179` | `//TODO 引用检查` | **HIGH**: Can delete referenced UNS nodes |
| `modelUNS.go:114` | `return nil //TODO` on `GetCalculationType()` | **MEDIUM**: Calculation type metadata missing |
| `uns/kong/*.go` (4 files) | `// todo: add your logic here` | **HIGH**: Kong API gateway integration unimplemented |
| `uns/alarm/*.go` (4 files) | `// todo: add your logic here` | **HIGH**: Alarm CRUD fully unimplemented |
| `uns/dashboard/*.go` (4 files) | `// todo: add your logic here` | **HIGH**: Dashboard management partially unimplemented |
| `uns/external/treeLogic.go` | `// todo: add your logic here` | **MEDIUM**: External tree view unimplemented |
| `uns/file/batch*.go, blobLogic.go` (3 files) | `// todo: add your logic here` | **MEDIUM**: File batch operations unimplemented |
| `unsNamespace.go:56` | `//todo 添加过滤字段` | **LOW**: Filter struct is empty |
| `UnsUpdateService.go:134-136` | `SubscribeModel` returns nil | **MEDIUM**: Model subscription not working |

### Scaffolded-Only Logic Files (empty implementations)

These files contain only the goctl scaffold with `// todo: add your logic here`:

- `batchQueryFileLogic.go`
- `checkDuplicationNameLogic.go`
- `searchExternalTreeLogic.go`
- `parserTopicPayloadLogic.go`
- `externalTopicAddLogic.go`
- `batchQueryFileHistoryValueLogic.go`
- `clearExternalTreeLogic.go`
- `dataSrc2UnsFieldsLogic.go`

---

## 11. Data Integrity Issues

### 11.1 Race Condition in WebSocket Session Count

**File:** `WebsocketService.go:92-107`

`TryAddSession` counts sessions by iterating `sync.Map` then adds — not atomic. Under high concurrency, the limit can be exceeded.

### 11.2 Unsafe Pointer String Conversion

**File:** `UnsMessageConsumer.go:89-100`

```go
func b2s(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
```

Uses `unsafe.Pointer` for zero-copy byte-to-string conversion. If the underlying byte slice is reused (e.g., by the MQTT library), this creates **use-after-free** bugs. The `reflect.SliceHeader`/`reflect.StringHeader` approach in `s2b` is deprecated since Go 1.20.

### 11.3 Silent Data Loss on Disk Queue

**File:** `UnsQueueDataSinkService.go:43`

`queue.Put(binData)` error is ignored. If the disk queue rejects the message (size exceeded, disk full), data is silently lost with no metric or alert.

### 11.4 Cache Coherence Gap

**File:** `UnsDefinitionServiceImpl.go`

The definition cache has a 10-minute TTL. Between cache invalidation events and actual DB updates, there's a window where stale definitions can cause:
- Messages routed to wrong physical tables
- Field validation using outdated schemas
- Data written with incorrect column mappings

### 11.5 Transaction Scope in Persistence

**File:** `tsdb_persistence.go:29`

The persistence function acquires a single connection with a 15-minute timeout. If the TimescaleDB is under heavy load, this can cause connection pool exhaustion for other operations.

### 11.6 Soft Delete Without Cascade

UNS deletion is soft-delete (`status != 1`), but:
- `uns_label_ref` has no cascade — labels remain linked to "deleted" nodes
- `uns_tag` entries are not cleaned up
- `uns_attachment` records persist
- The `uns_history_delete_job` table exists for async cleanup, but no background job runner was found

### 11.7 No Foreign Key Constraints

The schema uses no foreign keys between tables. `parent_id`, `model_id`, `label_id`, `uns_id` are all unconstrained. This allows orphaned records and referential integrity violations.

### 11.8 Potential SQL Injection in Tree Queries

**File:** `unsNamespace_query_tree.go:106-136`

While `escapeSQL()` and `escapeLikePattern()` are used, the raw SQL construction with string concatenation is error-prone. Any missed escape path could lead to SQL injection.

---

## 12. Error Handling Analysis

### 12.1 MQTT Pipeline

| Stage | Error Handling | Quality |
|-------|---------------|---------|
| MQTT connection | Infinite retry with backoff | OK, but no max retry |
| Topic resolution | Returns nil → logs debug, sends raw to WebSocket | OK |
| JSON parsing | Returns error → sends error to WebSocket | OK |
| Field validation | Sets QoS quality code, replaces with default | Good |
| Disk queue write | **Error ignored** | BAD |
| Persistence | Logs error, retries on specific PG errors | Good |

### 12.2 CRUD Operations

| Operation | Error Handling | Quality |
|-----------|---------------|---------|
| Create | Transaction rollback on error | Good |
| Update | Returns error result | OK |
| Delete | Transaction rollback, event publishing failure logged | OK |
| Circular dependency | Detected in dependency sort, reported per-item | Good |

### 12.3 TimescaleDB

| Scenario | Handling | Quality |
|----------|---------|---------|
| Table doesn't exist | Auto-create + retry | Good |
| Column missing | ALTER TABLE ADD + retry | Good |
| View conflicts with table | Rename table, create view, migrate data | Good |
| Connection pool exhaustion | 15-min timeout, logged with pool stats | OK |
| Compression errors | Not specifically handled | Risk |

---

## Summary of Critical Findings

1. **Realtime calculation is completely unimplemented** — the most impactful TODO
2. **Kong API gateway integration is entirely scaffolded** — no route registration
3. **Alarm management is fully unimplemented** — all 4 logic files are empty
4. **UNS validation rules from skills docs are not enforced in code** — only DB unique constraint on alias
5. **Disk queue data loss on full disk** — error silently swallowed
6. **WebSocket `/send` command doesn't work** — can't write back to MQTT from frontend
7. **Reference integrity not checked on delete** — can break calculated topics
8. **No foreign keys** — soft-delete leaves orphaned label refs, tags, attachments
9. **Unsafe pointer operations** in hot path — potential memory safety issues
10. **View-per-topic approach scales well** for reads but `ALTER TABLE` on the shared hypertable under concurrent writes could cause lock contention
