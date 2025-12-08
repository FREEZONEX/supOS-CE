CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
--
CREATE TABLE IF NOT EXISTS public."supos_timeserial_string" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" text NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_string',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hour'
);

ALTER TABLE supos_timeserial_string SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近30天的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_string', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_string',
    compress_after => INTERVAL '1 hour',
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);


CREATE TABLE IF NOT EXISTS public."supos_timeserial_integer" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" int4 NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_integer',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hour'
);

ALTER TABLE supos_timeserial_integer SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_integer', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_integer',
    compress_after => INTERVAL '1 hour', -- 压缩1小时后的数据
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);

CREATE TABLE IF NOT EXISTS public."supos_timeserial_long" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" int8 NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_long',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hour'
);

ALTER TABLE supos_timeserial_long SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_long', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_long',
    compress_after => INTERVAL '1 hour', -- 压缩1小时后的数据
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);


CREATE TABLE IF NOT EXISTS public."supos_timeserial_double" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" float8 NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_double',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hour'
);

ALTER TABLE supos_timeserial_double SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_double', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_double',
    compress_after => INTERVAL '1 hour', -- 压缩1小时后的数据
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);


CREATE TABLE IF NOT EXISTS public."supos_timeserial_float" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" float4 NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_float',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hour'
);

ALTER TABLE supos_timeserial_float SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_float', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_float',
    compress_after => INTERVAL '1 hour', -- 压缩1小时后的数据
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);

CREATE TABLE IF NOT EXISTS public."supos_timeserial_boolean" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" BOOLEAN  NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_boolean',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hours'
);

ALTER TABLE supos_timeserial_boolean SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_boolean', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_boolean',
    compress_after => INTERVAL '1 hour', -- 压缩1小时后的数据
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);

CREATE TABLE IF NOT EXISTS public."supos_timeserial_datetime" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" timestamptz(3) NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_datetime',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,        -- 按位号分50区
    chunk_time_interval => INTERVAL '2 hour'
);

ALTER TABLE supos_timeserial_datetime SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_datetime', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy(
    'supos_timeserial_datetime',
    compress_after => INTERVAL '1 hour', -- 压缩1小时后的数据
    schedule_interval => INTERVAL '2 hour'  -- 执行频率改为每2小时
);

CREATE TABLE IF NOT EXISTS public."supos_timeserial_blob" (
    "tag" int8 NOT NULL,
    "timeStamp" timestamptz(3) NOT NULL DEFAULT now(),
    "quality" int8 default 0,
    "value" varchar(512) NULL,
    PRIMARY KEY ("tag", "timeStamp")
);
SELECT create_hypertable (
    'supos_timeserial_blob',
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 10,        -- 按位号分10区
    chunk_time_interval => INTERVAL '1 day'
);

ALTER TABLE supos_timeserial_blob SET (
    timescaledb.compress,                          -- 启用压缩
    timescaledb.compress_segmentby = 'tag',   -- 按设备分组压缩
    timescaledb.compress_orderby = '"timeStamp" DESC'      -- 按时间降序排序
);

-- 保留最近2年的数据，自动删除旧分区
SELECT add_retention_policy('supos_timeserial_blob', INTERVAL '2 year');
-- 压缩每天的数据
SELECT add_compression_policy('supos_timeserial_blob', INTERVAL '24 hour');