DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'timescaledb extension is not available, use plain PostgreSQL tables: %', SQLERRM;
END
$$;

DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'pg_stat_statements extension is not available: %', SQLERRM;
END
$$;

CREATE SCHEMA IF NOT EXISTS uns;

CREATE TABLE IF NOT EXISTS uns."uns_timeserial" (
    "_id" int8 NOT NULL,
    "_timestamp" timestamptz(3) NOT NULL DEFAULT now(),
    "_quality" int8 DEFAULT 0,
    "double_1" float8 NULL,
    PRIMARY KEY ("_id", "_timestamp")
);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'uns'
          AND table_name = 'uns_timeserial'
          AND column_name = 'timeStamp'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'uns'
          AND table_name = 'uns_timeserial'
          AND column_name = '_timestamp'
    ) THEN
        ALTER TABLE uns."uns_timeserial" RENAME COLUMN "timeStamp" TO "_timestamp";
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'uns'
          AND table_name = 'uns_timeserial'
          AND column_name = 'quality'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'uns'
          AND table_name = 'uns_timeserial'
          AND column_name = '_quality'
    ) THEN
        ALTER TABLE uns."uns_timeserial" RENAME COLUMN "quality" TO "_quality";
    END IF;
END
$$;

DO $$
DECLARE
    qualified_table text;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb') THEN
        RETURN;
    END IF;

    qualified_table := format('%I.%I', 'uns', 'uns_timeserial');

    EXECUTE format(
        'SELECT create_hypertable(%L, %L, partitioning_column => %L, number_partitions => 4, chunk_time_interval => %L::interval, if_not_exists => TRUE)',
        qualified_table,
        '_timestamp',
        '_id',
        '7 days'
    );

    EXECUTE format(
        'ALTER TABLE %s SET (timescaledb.compress, timescaledb.compress_segmentby = %L, timescaledb.compress_orderby = %L)',
        qualified_table,
        '_id',
        '"_timestamp" DESC'
    );

    IF NOT EXISTS (
        SELECT 1
        FROM timescaledb_information.jobs
        WHERE hypertable_schema = 'uns'
          AND hypertable_name = 'uns_timeserial'
          AND proc_name = 'policy_retention'
    ) THEN
        PERFORM add_retention_policy(qualified_table, INTERVAL '10 years');
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM timescaledb_information.jobs
        WHERE hypertable_schema = 'uns'
          AND hypertable_name = 'uns_timeserial'
          AND proc_name = 'policy_compression'
    ) THEN
        PERFORM add_compression_policy(
            qualified_table,
            compress_after => INTERVAL '7 days',
            schedule_interval => INTERVAL '1 day'
        );
    END IF;
END
$$;
