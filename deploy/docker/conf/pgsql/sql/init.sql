-- 创建 Kong 和 Konga 数据库
CREATE DATABASE kong;
CREATE DATABASE konga;


SET search_path TO public;


CREATE OR REPLACE FUNCTION sync_public_tables_to_schema(
    target_schema text,
    table_prefix text  -- 新增：表名前缀参数
)
    RETURNS void AS $$
DECLARE
tbl_name text;
    table_exists boolean;
    seq_name text;
    seq_owner text;
    fk_constraint record;
    seq_last_value bigint := 1; -- 强制默认值
    seq_increment_by integer := 1; -- 强制默认值
    seq_schema text;
    seq_short_name text;
    seq_found boolean;
    default_seq_name text; -- 用于生成默认序列名
BEGIN
    -- 确保目标schema存在
EXECUTE format('CREATE SCHEMA IF NOT EXISTS %I', target_schema);

-- 第一步：复制指定前缀的表结构
FOR tbl_name IN
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_type = 'BASE TABLE'
  AND table_name LIKE table_prefix || '%'  -- 仅匹配带前缀的表
    LOOP
SELECT EXISTS(
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = target_schema
      AND table_name = tbl_name
) INTO table_exists;

IF NOT table_exists THEN
                EXECUTE format(
                        'CREATE TABLE %I.%I (LIKE public.%I INCLUDING ALL)',
                        target_schema, tbl_name, tbl_name
                        );
                RAISE NOTICE '已创建表: %.%', target_schema, tbl_name;
END IF;
END LOOP;

    -- 第二步：复制序列并修正字段默认值（仅处理前缀匹配的表）
FOR tbl_name IN
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_type = 'BASE TABLE'
  AND table_name LIKE table_prefix || '%'  -- 仅匹配带前缀的表
    LOOP
            FOR seq_name, seq_owner IN
SELECT
    pg_get_serial_sequence(format('%I.%I', 'public', tbl_name), c.column_name)::regclass::text AS seq_name,
    c.column_name AS seq_owner
FROM information_schema.columns c
WHERE c.table_schema = 'public'
  AND c.table_name = tbl_name
  AND pg_get_serial_sequence(format('%I.%I', 'public', tbl_name), c.column_name) IS NOT NULL
    LOOP
DECLARE
target_seq_name text; -- 目标序列名（带schema）
BEGIN
                        -- 生成默认序列名（基于表名和列名）
                        default_seq_name := format('%s_%s_seq', tbl_name, seq_owner);

                        -- 解析原序列的schema和短名称
                        seq_schema := split_part(seq_name, '.', 1);
                        seq_short_name := split_part(seq_name, '.', 2);

                        -- 处理序列名异常
                        IF seq_short_name = '' OR seq_short_name IS NULL THEN
                            RAISE NOTICE '序列名称异常: %，使用默认名称: %', seq_name, default_seq_name;
                            seq_short_name := default_seq_name;
END IF;

                        -- 原序列schema默认值
                        IF seq_schema = '' OR seq_schema IS NULL THEN
                            seq_schema := 'public';
END IF;

                        -- 明确目标序列的schema和名称（确保在目标schema中）
                        target_seq_name := format('%I.%I', target_schema, seq_short_name);

                        -- 仅在目标序列不存在时创建
                        IF NOT EXISTS(SELECT 1 FROM pg_sequences WHERE schemaname = target_schema AND sequencename = seq_short_name) THEN
                            -- 重置参数
                            seq_last_value := 1;
                            seq_increment_by := 1;
                            seq_found := false;

                            -- 查询原序列属性
FOR seq_last_value, seq_increment_by, seq_found IN
SELECT
    COALESCE(last_value, 1),
    COALESCE(increment_by, 1),
    true
FROM pg_sequences
WHERE schemaname = seq_schema AND sequencename = seq_short_name
    LOOP
                                    EXIT;
END LOOP;

                            -- 强制参数有效
                            IF NOT seq_found THEN
                                RAISE NOTICE '未找到原序列 %，使用默认参数', seq_name;
                                seq_last_value := 1;
                                seq_increment_by := 1;
END IF;
                            IF seq_last_value IS NULL THEN seq_last_value := 1; END IF;
                            IF seq_increment_by IS NULL THEN seq_increment_by := 1; END IF;

                            -- 创建目标序列（明确在目标schema中）
EXECUTE format(
        'CREATE SEQUENCE %I.%I AS integer START WITH %s INCREMENT BY %s NO MAXVALUE NO MINVALUE CACHE 1',
        target_schema, seq_short_name,  -- 明确指定schema和短名称
        seq_last_value, seq_increment_by
        );

-- 绑定序列所有权（关键修复：明确序列的schema与表一致）
EXECUTE format(
        'ALTER SEQUENCE %I.%I OWNED BY %I.%I.%I',
        target_schema, seq_short_name,  -- 序列的schema（与表相同）
        target_schema, tbl_name, seq_owner  -- 表的schema
        );

-- 修正字段默认值
EXECUTE format(
        'ALTER TABLE %I.%I ALTER COLUMN %I SET DEFAULT nextval(%L::regclass)',
        target_schema, tbl_name, seq_owner,
        target_seq_name  -- 使用带schema的序列名
        );

RAISE NOTICE '已复制序列: % -> %', seq_name, target_seq_name;
END IF;
END;
END LOOP;
END LOOP;

    -- 第三步：修正外键约束（仅处理前缀匹配的表）
FOR fk_constraint IN
SELECT
    tc.constraint_name,
    tc.table_name AS fk_table,
    kcu.column_name AS fk_column,
    ccu.table_name AS ref_table,
    ccu.column_name AS ref_column
FROM
    information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu
             ON tc.constraint_name = kcu.constraint_name
                 AND tc.table_schema = kcu.table_schema
        JOIN information_schema.constraint_column_usage ccu
             ON tc.constraint_name = ccu.constraint_name
WHERE
    tc.table_schema = 'public'
  AND tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_name LIKE table_prefix || '%'  -- 仅匹配带前缀的外键表
  AND ccu.table_name LIKE table_prefix || '%'  -- 仅匹配带前缀的引用表
    LOOP
            IF EXISTS(
                SELECT 1 FROM information_schema.tables
                WHERE table_schema = target_schema
                  AND table_name = fk_constraint.fk_table
            ) AND EXISTS(
                SELECT 1 FROM information_schema.tables
                WHERE table_schema = target_schema
                  AND table_name = fk_constraint.ref_table
            ) THEN
                EXECUTE format('ALTER TABLE %I.%I DROP CONSTRAINT IF EXISTS %I',
                               target_schema, fk_constraint.fk_table, fk_constraint.constraint_name);

EXECUTE format('ALTER TABLE %I.%I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES %I.%I(%I)',
               target_schema, fk_constraint.fk_table, fk_constraint.constraint_name,
               fk_constraint.fk_column,
               target_schema, fk_constraint.ref_table, fk_constraint.ref_column);

RAISE NOTICE '已修正外键: %.% -> %.%',
                    target_schema, fk_constraint.fk_table,
                    target_schema, fk_constraint.ref_table;
END IF;
END LOOP;

    RAISE NOTICE '同步完成，目标schema: %，表前缀: %', target_schema, table_prefix;
END;
$$ LANGUAGE plpgsql;
