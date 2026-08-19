CREATE TABLE IF NOT EXISTS sys_user_info (
    user_id BIGSERIAL PRIMARY KEY,
    user_name VARCHAR(60) NOT NULL,
    nick_name VARCHAR(60) NOT NULL DEFAULT '',
    password VARCHAR(128) NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    github_pid VARCHAR(128) NOT NULL DEFAULT '',
    google_pid VARCHAR(128) NOT NULL DEFAULT '',
    last_ip VARCHAR(128) NOT NULL DEFAULT '',
    reg_ip VARCHAR(128) NOT NULL DEFAULT '',
    avatar VARCHAR(256) NOT NULL DEFAULT '',
    pf_role_code VARCHAR(60) NOT NULL DEFAULT 'member',
    status BIGINT NOT NULL DEFAULT 1,
    last_login BIGINT NOT NULL DEFAULT 0,
    first_login BIGINT NOT NULL DEFAULT 0,
    is_random_pwd BOOLEAN NOT NULL DEFAULT false,
    statistic_workspace_joined_count BIGINT NOT NULL DEFAULT 1,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_user_info_tc_un ON sys_user_info(user_name, deleted_time);
ALTER TABLE sys_user_info ADD COLUMN IF NOT EXISTS phone TEXT NOT NULL DEFAULT '';
DROP INDEX IF EXISTS idx_sys_user_info_tc_email;

CREATE TABLE IF NOT EXISTS sys_user_config (
    user_id BIGINT PRIMARY KEY,
    home_page VARCHAR(255) NOT NULL DEFAULT '/home',
    main_language VARCHAR(32) NOT NULL DEFAULT '',
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sys_role_info (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    "desc" VARCHAR(200) NOT NULL DEFAULT '',
    type BIGINT NOT NULL DEFAULT 1,
    status BIGINT NOT NULL DEFAULT 1,
    default_home_page VARCHAR(255) NOT NULL DEFAULT '/home',
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE sys_role_info ADD COLUMN IF NOT EXISTS default_home_page VARCHAR(255) NOT NULL DEFAULT '/home';
ALTER TABLE sys_role_info ADD COLUMN IF NOT EXISTS "desc" VARCHAR(200) NOT NULL DEFAULT '';
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'sys_role_info' AND column_name = 'description'
    ) THEN
        EXECUTE 'UPDATE sys_role_info SET "desc" = description WHERE "desc" = '''' AND description <> ''''';
    END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_role_info_name ON sys_role_info(name, deleted_time);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_role_info_code ON sys_role_info(code, deleted_time);

CREATE TABLE IF NOT EXISTS sys_workspace_user (
    id BIGSERIAL PRIMARY KEY,
    workspace_id BIGINT NOT NULL DEFAULT 1,
    user_id BIGINT NOT NULL,
    role_code VARCHAR(100) NOT NULL DEFAULT '',
    role_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_workspace_user_uid ON sys_workspace_user(workspace_id, user_id, deleted_time);

CREATE TABLE IF NOT EXISTS sys_resource_info (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL DEFAULT 0,
    resource_key VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    type BIGINT NOT NULL DEFAULT 5,
    route_path VARCHAR(255) NOT NULL DEFAULT '',
    url_type BIGINT NOT NULL DEFAULT 1,
    open_type BIGINT NOT NULL DEFAULT 0,
    icon VARCHAR(100) NOT NULL DEFAULT '',
    sort BIGINT NOT NULL DEFAULT 0,
    enabled BIGINT NOT NULL DEFAULT 1,
    system_generated BIGINT NOT NULL DEFAULT 1,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE sys_resource_info ADD COLUMN IF NOT EXISTS url_type BIGINT NOT NULL DEFAULT 1;
ALTER TABLE sys_resource_info ADD COLUMN IF NOT EXISTS open_type BIGINT NOT NULL DEFAULT 0;
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_resource_info_key ON sys_resource_info(resource_key, deleted_time);

CREATE TABLE IF NOT EXISTS sys_resource_action (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    action_type VARCHAR(20) NOT NULL,
    action_value VARCHAR(200) NOT NULL,
    methods VARCHAR(50) NOT NULL DEFAULT '',
    enabled BIGINT NOT NULL DEFAULT 1,
    system_generated BIGINT NOT NULL DEFAULT 1,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_resource_action_entry ON sys_resource_action(action_type, action_value, methods, deleted_time);

CREATE TABLE IF NOT EXISTS sys_role_resource (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_role_resource_pair ON sys_role_resource(role_id, resource_id, deleted_time);

CREATE TABLE IF NOT EXISTS sys_api_key (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(128) NOT NULL,
    key_prefix VARCHAR(32) NOT NULL DEFAULT '',
    key_suffix VARCHAR(16) NOT NULL DEFAULT '',
    owner_id BIGINT NOT NULL DEFAULT 0,
    owner_type VARCHAR(32) NOT NULL DEFAULT 'personal',
    usage_type VARCHAR(32) NOT NULL DEFAULT 'external',
    permission VARCHAR(32) NOT NULL DEFAULT 'read_only',
    resource_keys TEXT NOT NULL DEFAULT '',
    status BIGINT NOT NULL DEFAULT 1,
    last_used_time BIGINT NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE sys_api_key DROP COLUMN IF EXISTS key_type;
ALTER TABLE sys_api_key ALTER COLUMN permission SET DEFAULT 'read_only';
ALTER TABLE sys_api_key ALTER COLUMN owner_id TYPE BIGINT USING CASE WHEN owner_id::text ~ '^[0-9]+$' THEN owner_id::text::BIGINT ELSE 0 END;
ALTER TABLE sys_api_key ALTER COLUMN owner_id SET DEFAULT 0;
DO $$
DECLARE
    max_api_key_id BIGINT;
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'sys_api_key'
          AND column_name = 'id'
          AND column_default IS NULL
    ) THEN
        CREATE SEQUENCE IF NOT EXISTS sys_api_key_id_seq OWNED BY sys_api_key.id;
        SELECT COALESCE(MAX(id), 0) INTO max_api_key_id FROM sys_api_key;
        PERFORM setval('sys_api_key_id_seq', GREATEST(max_api_key_id, 1), max_api_key_id > 0);
        ALTER TABLE sys_api_key ALTER COLUMN id SET DEFAULT nextval('sys_api_key_id_seq');
    END IF;
END $$;
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND indexname = 'idx_sys_api_key_hash'
          AND indexdef LIKE '%key_hash, deleted_time%'
    ) THEN
        DROP INDEX idx_sys_api_key_hash;
    END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_api_key_hash ON sys_api_key(key_hash);
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND indexname = 'idx_sys_api_key_name_owner'
          AND indexdef NOT LIKE '%deleted_time%'
    ) THEN
        DROP INDEX idx_sys_api_key_name_owner;
    END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_api_key_name_owner ON sys_api_key(name, owner_id, owner_type, deleted_time);
CREATE INDEX IF NOT EXISTS idx_sys_api_key_owner_id ON sys_api_key(owner_id);
CREATE INDEX IF NOT EXISTS idx_sys_api_key_status ON sys_api_key(status);
CREATE INDEX IF NOT EXISTS idx_sys_api_key_created_time ON sys_api_key(created_time DESC);

CREATE TABLE IF NOT EXISTS sys_gateway_route (
    id BIGSERIAL PRIMARY KEY,
    route_key VARCHAR(96) NOT NULL,
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL DEFAULT '',
    target_type VARCHAR(32) NOT NULL,
    match_type VARCHAR(32) NOT NULL,
    path_pattern VARCHAR(256) NOT NULL,
    methods JSONB NOT NULL DEFAULT '[]',
    target_url VARCHAR(512) NOT NULL DEFAULT '',
    target_path VARCHAR(256) NOT NULL DEFAULT '',
    strip_prefix BOOLEAN NOT NULL DEFAULT false,
    rewrite_path VARCHAR(256) NOT NULL DEFAULT '',
    auth_policy VARCHAR(32) NOT NULL,
    resource_key VARCHAR(128) NOT NULL DEFAULT '',
    timeout_ms INTEGER NOT NULL DEFAULT 10000,
    priority INTEGER NOT NULL DEFAULT 100,
    enabled BOOLEAN NOT NULL DEFAULT true,
    system_builtin BOOLEAN NOT NULL DEFAULT false,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_gateway_route_key ON sys_gateway_route(route_key, deleted_time);

CREATE TABLE IF NOT EXISTS sys_outbox_event (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(64) NOT NULL,
    event_type VARCHAR(128) NOT NULL,
    aggregate_type VARCHAR(64) NOT NULL DEFAULT '',
    aggregate_id VARCHAR(128) NOT NULL DEFAULT '',
    payload JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    next_retry_time BIGINT NOT NULL DEFAULT 0,
    last_error TEXT NOT NULL DEFAULT '',
    locked_by VARCHAR(128) NOT NULL DEFAULT '',
    locked_until_time BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_outbox_event_event_id ON sys_outbox_event(event_id);

CREATE TABLE IF NOT EXISTS sys_async_job (
    id BIGSERIAL PRIMARY KEY,
    job_key VARCHAR(64) NOT NULL,
    job_type VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    progress INTEGER NOT NULL DEFAULT 0,
    request_json JSONB NOT NULL DEFAULT '{}',
    result_json JSONB NOT NULL DEFAULT '{}',
    error_message TEXT NOT NULL DEFAULT '',
    created_by BIGINT NOT NULL DEFAULT 0,
    started_time BIGINT NOT NULL DEFAULT 0,
    finished_time BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_async_job_job_key ON sys_async_job(job_key);

CREATE TABLE IF NOT EXISTS sys_user_oper_log (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    res_type VARCHAR(50) NOT NULL DEFAULT '',
    res_id VARCHAR(50) NOT NULL DEFAULT '',
    res_name VARCHAR(256) NOT NULL DEFAULT '',
    business_type BIGINT NOT NULL,
    detail VARCHAR(500) NOT NULL DEFAULT '',
    code BIGINT NOT NULL DEFAULT 200,
    is_show_in_recent BIGINT NOT NULL DEFAULT 1,
    operator_name VARCHAR(100) NOT NULL DEFAULT '',
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_sys_user_oper_log_created_time ON sys_user_oper_log(created_time DESC);
CREATE INDEX IF NOT EXISTS idx_sys_user_oper_log_user_recent ON sys_user_oper_log(user_id, is_show_in_recent, created_time DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_sys_user_oper_log_resource_time ON sys_user_oper_log(res_type, res_id, created_time DESC, id DESC);

CREATE TABLE IF NOT EXISTS uns_namespace_node_info (
    id BIGINT PRIMARY KEY,
    id_path VARCHAR(1024) NOT NULL DEFAULT '',
    parent_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(128) NOT NULL,
    display_name VARCHAR(256) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    namespace VARCHAR(512) NOT NULL,
    alias VARCHAR(128) NOT NULL DEFAULT '',
    type SMALLINT NOT NULL DEFAULT 1,
    topic_type SMALLINT NOT NULL DEFAULT 0,
    sort_key BIGINT NOT NULL DEFAULT 0,
    schema JSONB NOT NULL DEFAULT '[]'::jsonb,
    extendproperties JSONB NOT NULL DEFAULT '{}',
    enable_history SMALLINT NOT NULL DEFAULT 2,
    mock_data SMALLINT NOT NULL DEFAULT 2,
    template_id BIGINT DEFAULT NULL,
    is_favorite SMALLINT NOT NULL DEFAULT 2,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    recycle_is_del SMALLINT NOT NULL DEFAULT 0,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    deleted_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS recycle_is_del SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS id_path VARCHAR(1024) NOT NULL DEFAULT '';
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS sort_key BIGINT NOT NULL DEFAULT 0;
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS extendproperties JSONB NOT NULL DEFAULT '{}';
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS enable_history SMALLINT NOT NULL DEFAULT 2;
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS mock_data SMALLINT NOT NULL DEFAULT 2;
ALTER TABLE uns_namespace_node_info ADD COLUMN IF NOT EXISTS is_favorite SMALLINT NOT NULL DEFAULT 2;
ALTER TABLE uns_namespace_node_info ALTER COLUMN schema SET DEFAULT '[]'::jsonb;
UPDATE uns_namespace_node_info SET schema = '[]'::jsonb WHERE schema IS NULL OR jsonb_typeof(schema) = 'null';
UPDATE uns_namespace_node_info SET schema = COALESCE(schema->'fields', '[]'::jsonb) WHERE jsonb_typeof(schema) = 'object';
ALTER TABLE uns_namespace_node_info ALTER COLUMN schema SET NOT NULL;
ALTER TABLE uns_namespace_node_info ALTER COLUMN template_id DROP DEFAULT;
ALTER TABLE uns_namespace_node_info ALTER COLUMN template_id DROP NOT NULL;
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'uns_namespace_node_info' AND column_name = 'extend_properties'
    ) THEN
        EXECUTE 'UPDATE uns_namespace_node_info SET extendproperties = extend_properties WHERE extendproperties = ''{}''::jsonb';
    END IF;
END $$;
UPDATE uns_namespace_node_info SET template_id = NULL WHERE template_id = 0;
WITH RECURSIVE tree AS (
    SELECT id, parent_id, id::text AS path FROM uns_namespace_node_info WHERE parent_id = 0
    UNION ALL
    SELECT child.id, child.parent_id, tree.path || '/' || child.id::text
    FROM uns_namespace_node_info child
    JOIN tree ON child.parent_id = tree.id
)
UPDATE uns_namespace_node_info node SET id_path = tree.path FROM tree WHERE node.id = tree.id AND (node.id_path = '' OR node.id_path IS NULL);
UPDATE uns_namespace_node_info SET id_path = id::text WHERE id_path = '' OR id_path IS NULL;
ALTER TABLE uns_namespace_node_info DROP COLUMN IF EXISTS extend_properties;
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND tablename = 'uns_namespace_node_info'
          AND indexname = 'idx_uns_node_namespace'
          AND indexdef NOT ILIKE '%WHERE%deleted_time = 0%'
    ) THEN
        DROP INDEX idx_uns_node_namespace;
    END IF;

    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND tablename = 'uns_namespace_node_info'
          AND indexname = 'idx_uns_node_alias'
          AND (indexdef NOT ILIKE '%WHERE%deleted_time = 0%' OR indexdef NOT ILIKE '%alias%<>%')
    ) THEN
        DROP INDEX idx_uns_node_alias;
    END IF;

    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND tablename = 'uns_namespace_node_info'
          AND indexname = 'idx_uns_node_parent_name'
          AND indexdef NOT ILIKE '%WHERE%deleted_time = 0%'
    ) THEN
        DROP INDEX idx_uns_node_parent_name;
    END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_node_namespace ON uns_namespace_node_info(namespace) WHERE deleted_time = 0;
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_node_alias ON uns_namespace_node_info(alias) WHERE deleted_time = 0 AND alias <> '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_node_parent_name ON uns_namespace_node_info(parent_id, name) WHERE deleted_time = 0;
CREATE INDEX IF NOT EXISTS idx_uns_node_id_path ON uns_namespace_node_info(id_path);

CREATE TABLE IF NOT EXISTS uns_namespace_node_version (
    id BIGSERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 0,
    etag VARCHAR(64) NOT NULL DEFAULT '',
    checksum VARCHAR(64) NOT NULL DEFAULT '',
    unstree JSONB NOT NULL DEFAULT 'null'::jsonb,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE uns_namespace_node_version ADD COLUMN IF NOT EXISTS etag VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE uns_namespace_node_version ADD COLUMN IF NOT EXISTS checksum VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE uns_namespace_node_version ADD COLUMN IF NOT EXISTS unstree JSONB NOT NULL DEFAULT 'null'::jsonb;
ALTER TABLE uns_namespace_node_version ADD COLUMN IF NOT EXISTS created_by BIGINT NOT NULL DEFAULT 0;
ALTER TABLE uns_namespace_node_version ADD COLUMN IF NOT EXISTS updated_by BIGINT NOT NULL DEFAULT 0;
DO $$
DECLARE
    max_uns_namespace_node_version_id BIGINT;
BEGIN
    CREATE SEQUENCE IF NOT EXISTS uns_namespace_node_version_id_seq OWNED BY uns_namespace_node_version.id;
    SELECT COALESCE(MAX(id), 0) INTO max_uns_namespace_node_version_id FROM uns_namespace_node_version;
    PERFORM setval(
        'uns_namespace_node_version_id_seq',
        GREATEST(max_uns_namespace_node_version_id, 1),
        max_uns_namespace_node_version_id > 0
    );
    ALTER TABLE uns_namespace_node_version ALTER COLUMN id SET DEFAULT nextval('uns_namespace_node_version_id_seq');
END $$;
CREATE INDEX IF NOT EXISTS idx_uns_namespace_node_version_version ON uns_namespace_node_version(version);
CREATE INDEX IF NOT EXISTS idx_uns_namespace_node_version_etag ON uns_namespace_node_version(etag);

CREATE TABLE IF NOT EXISTS uns_namespace_template_info (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    topic_type SMALLINT NOT NULL,
    schema JSONB NOT NULL DEFAULT '[]'::jsonb,
    extend_properties JSONB NOT NULL DEFAULT '{}',
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    deleted_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE uns_namespace_template_info ALTER COLUMN schema SET DEFAULT '[]'::jsonb;
UPDATE uns_namespace_template_info SET schema = '[]'::jsonb WHERE schema IS NULL OR jsonb_typeof(schema) = 'null';
UPDATE uns_namespace_template_info SET schema = COALESCE(schema->'fields', '[]'::jsonb) WHERE jsonb_typeof(schema) = 'object';
ALTER TABLE uns_namespace_template_info ALTER COLUMN schema SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_template_name ON uns_namespace_template_info(name) WHERE deleted_time = 0;

CREATE TABLE IF NOT EXISTS uns_namespace_label_info (
    id BIGINT PRIMARY KEY,
    label_key VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    color VARCHAR(32) NOT NULL DEFAULT '',
    description VARCHAR(512) NOT NULL DEFAULT '',
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_label_key ON uns_namespace_label_info(label_key) WHERE deleted_time = 0;

CREATE TABLE IF NOT EXISTS uns_namespace_label_nodeid (
    label_id BIGINT NOT NULL,
    node_id BIGINT NOT NULL,
    created_time BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY(label_id, node_id)
);

CREATE TABLE IF NOT EXISTS uns_nodered_flow (
    id BIGINT PRIMARY KEY,
    flow_id VARCHAR(64) NOT NULL DEFAULT '',
    flow_name VARCHAR(512) NOT NULL,
    flow_data TEXT NOT NULL DEFAULT '',
    flow_status VARCHAR(32) NOT NULL DEFAULT 'draft',
    template VARCHAR(64) NOT NULL DEFAULT '',
    description VARCHAR(512) NOT NULL DEFAULT '',
    flow_type INTEGER NOT NULL DEFAULT 1,
    node_type INTEGER NOT NULL DEFAULT 1,
    parent_id BIGINT NOT NULL DEFAULT 0,
    sort_key BIGINT NOT NULL DEFAULT 0,
    is_favorite SMALLINT NOT NULL DEFAULT 2,
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
DROP INDEX IF EXISTS idx_uns_flow_name;
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND tablename = 'uns_nodered_flow'
          AND indexname = 'idx_uns_flow_parent_name_del'
          AND (
              indexdef NOT ILIKE '%WHERE%deleted_time = 0%'
              OR indexdef NOT ILIKE '%flow_type%'
          )
    ) THEN
        DROP INDEX idx_uns_flow_parent_name_del;
    END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_flow_parent_name_del ON uns_nodered_flow(flow_type, parent_id, flow_name) WHERE deleted_time = 0;
CREATE INDEX IF NOT EXISTS idx_uns_nodered_flow_flow_id ON uns_nodered_flow(flow_id);
CREATE INDEX IF NOT EXISTS idx_uns_flow_parent ON uns_nodered_flow(parent_id);
CREATE INDEX IF NOT EXISTS idx_uns_flow_sort ON uns_nodered_flow(sort_key);

CREATE TABLE IF NOT EXISTS uns_nodered_flow_node (
    parent_id BIGINT NOT NULL,
    node_id BIGINT NOT NULL,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(node_id)
);
CREATE INDEX IF NOT EXISTS idx_uns_flow_node_parent ON uns_nodered_flow_node(parent_id);
CREATE INDEX IF NOT EXISTS idx_uns_flow_node_node_id ON uns_nodered_flow_node(node_id);

CREATE TABLE IF NOT EXISTS uns_dashboard (
    id VARCHAR(128) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type SMALLINT NOT NULL DEFAULT 1,
    group_id BIGINT NOT NULL DEFAULT 0,
    node_type INTEGER NOT NULL DEFAULT 2,
    parent_id VARCHAR(128) NOT NULL DEFAULT '',
    need_init BOOLEAN NOT NULL DEFAULT false,
    description TEXT NOT NULL DEFAULT '',
    creator VARCHAR(128) NOT NULL DEFAULT '',
    create_time BIGINT NOT NULL DEFAULT 0,
    update_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE uns_dashboard ADD COLUMN IF NOT EXISTS node_type INTEGER NOT NULL DEFAULT 2;
ALTER TABLE uns_dashboard ADD COLUMN IF NOT EXISTS parent_id VARCHAR(128) NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_uns_dashboard_parent ON uns_dashboard(node_type, parent_id, create_time);

CREATE TABLE IF NOT EXISTS uns_dashboard_ref (
    dashboard_id VARCHAR(128) NOT NULL,
    uns_alias VARCHAR(128) NOT NULL,
    create_time BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY(dashboard_id, uns_alias)
);

CREATE TABLE IF NOT EXISTS uns_dashboard_top_recodes (
    id VARCHAR(128) NOT NULL,
    user_id VARCHAR(128) NOT NULL,
    mark SMALLINT NOT NULL DEFAULT 0,
    mark_time BIGINT NOT NULL DEFAULT 0,
    update_time BIGINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_dashboard_top_record ON uns_dashboard_top_recodes(id, user_id);

CREATE TABLE IF NOT EXISTS uns_folder (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id BIGINT NOT NULL DEFAULT 0,
    path VARCHAR(1024) NOT NULL DEFAULT '',
    owner VARCHAR(128) NOT NULL DEFAULT '',
    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND indexname = 'idx_uns_folder_name_parent'
          AND indexdef NOT LIKE '%deleted_time%'
    ) THEN
        DROP INDEX idx_uns_folder_name_parent;
    END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS idx_uns_folder_name_parent ON uns_folder(name, parent_id, deleted_time);
CREATE INDEX IF NOT EXISTS idx_uns_folder_parent ON uns_folder(parent_id);
CREATE INDEX IF NOT EXISTS idx_uns_folder_owner ON uns_folder(owner);

DO $$
DECLARE
    tbl text;
BEGIN
    FOREACH tbl IN ARRAY ARRAY[
        'sys_user_info',
        'sys_user_config',
        'user_config',
        'sys_role_info',
        'sys_workspace_user',
        'sys_resource_info',
        'sys_resource_action',
        'sys_role_resource',
        'sys_api_key',
        'sys_gateway_route',
        'sys_outbox_event',
        'sys_async_job',
        'uns_namespace_node_info',
        'uns_namespace_node_version',
        'uns_namespace_template_info',
        'uns_namespace_label_info',
        'uns_namespace_label_nodeid',
        'uns_nodered_flow',
        'uns_nodered_flow_node',
        'uns_folder'
    ]
    LOOP
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = current_schema()
              AND table_name = tbl
              AND column_name = 'created_time'
              AND data_type = 'bigint'
        ) THEN
            EXECUTE format('ALTER TABLE %I ALTER COLUMN created_time DROP DEFAULT', tbl);
            EXECUTE format(
                'ALTER TABLE %I ALTER COLUMN created_time TYPE TIMESTAMPTZ USING CASE WHEN created_time > 0 THEN to_timestamp(created_time / 1000.0) ELSE CURRENT_TIMESTAMP END',
                tbl
            );
            EXECUTE format('ALTER TABLE %I ALTER COLUMN created_time SET DEFAULT CURRENT_TIMESTAMP', tbl);
            EXECUTE format('ALTER TABLE %I ALTER COLUMN created_time SET NOT NULL', tbl);
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = current_schema()
              AND table_name = tbl
              AND column_name = 'updated_time'
              AND data_type = 'bigint'
        ) THEN
            EXECUTE format('ALTER TABLE %I ALTER COLUMN updated_time DROP DEFAULT', tbl);
            EXECUTE format(
                'ALTER TABLE %I ALTER COLUMN updated_time TYPE TIMESTAMPTZ USING CASE WHEN updated_time > 0 THEN to_timestamp(updated_time / 1000.0) ELSE CURRENT_TIMESTAMP END',
                tbl
            );
            EXECUTE format('ALTER TABLE %I ALTER COLUMN updated_time SET DEFAULT CURRENT_TIMESTAMP', tbl);
            EXECUTE format('ALTER TABLE %I ALTER COLUMN updated_time SET NOT NULL', tbl);
        END IF;
    END LOOP;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu
          ON tc.constraint_name = kcu.constraint_name
         AND tc.table_schema = kcu.table_schema
         AND tc.table_name = kcu.table_name
        WHERE tc.table_schema = current_schema()
          AND tc.table_name = 'uns_nodered_flow_node'
          AND tc.constraint_type = 'PRIMARY KEY'
          AND kcu.column_name = 'parent_id'
    ) THEN
        ALTER TABLE uns_nodered_flow_node DROP CONSTRAINT IF EXISTS uns_nodered_flow_node_pkey;
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu
          ON tc.constraint_name = kcu.constraint_name
         AND tc.table_schema = kcu.table_schema
         AND tc.table_name = kcu.table_name
        WHERE tc.table_schema = current_schema()
          AND tc.table_name = 'uns_nodered_flow_node'
          AND tc.constraint_type = 'PRIMARY KEY'
          AND kcu.column_name = 'node_id'
    ) THEN
        ALTER TABLE uns_nodered_flow_node ADD CONSTRAINT uns_nodered_flow_node_pkey PRIMARY KEY(node_id);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS asset_file (
    id BIGSERIAL PRIMARY KEY,
    file_key VARCHAR(512) NOT NULL,
    original_name VARCHAR(256) NOT NULL,
    content_type VARCHAR(128) NOT NULL DEFAULT '',
    size_bytes BIGINT NOT NULL DEFAULT 0,
    sha256 VARCHAR(64) NOT NULL DEFAULT '',
    storage_driver VARCHAR(32) NOT NULL DEFAULT 'local',
    storage_key VARCHAR(512) NOT NULL DEFAULT '',
    visibility VARCHAR(32) NOT NULL DEFAULT 'private',
    status VARCHAR(32) NOT NULL DEFAULT 'temp',
    project_id BIGINT NOT NULL DEFAULT 0,
    app_instance_id VARCHAR(64) NOT NULL DEFAULT '',
    session_id VARCHAR(64) NOT NULL DEFAULT '',
    created_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    updated_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);
ALTER TABLE asset_file ADD COLUMN IF NOT EXISTS project_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE asset_file ADD COLUMN IF NOT EXISTS app_instance_id VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE asset_file ADD COLUMN IF NOT EXISTS session_id VARCHAR(64) NOT NULL DEFAULT '';
-- bucket records the bucket chosen at upload time so later config changes
-- the wrong bucket; rows with '' fall back to visibility-derived resolution.
ALTER TABLE asset_file ADD COLUMN IF NOT EXISTS bucket VARCHAR(128) NOT NULL DEFAULT '';
-- 两列均可空，存量行不受影响（仅分片上传的 temp 记录使用）。
-- file_key mirrors the project-scoped storage key ({projectID}/assets/...), which
-- exceeds the legacy 64-char limit once appInstanceId/sessionId segments are present.
ALTER TABLE asset_file ALTER COLUMN file_key TYPE VARCHAR(512);
CREATE INDEX IF NOT EXISTS idx_asset_file_project_id ON asset_file(project_id);
CREATE INDEX IF NOT EXISTS idx_asset_file_app_instance_id ON asset_file(app_instance_id);
CREATE INDEX IF NOT EXISTS idx_asset_file_storage_key ON asset_file(storage_key);
CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_file_key ON asset_file(file_key) WHERE deleted_time = 0;

CREATE TABLE IF NOT EXISTS asset_binding (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    owner_type VARCHAR(32) NOT NULL,
    owner_id BIGINT NOT NULL DEFAULT 0,
    usage VARCHAR(32) NOT NULL DEFAULT 'attachment',
    sort_key INTEGER NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_by BIGINT NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL DEFAULT 0,
    deleted_time BIGINT NOT NULL DEFAULT 0
);





DO $$
DECLARE
    tbl text;
BEGIN
    FOREACH tbl IN ARRAY ARRAY[
        'asset_file',
        'asset_binding',
        'project',
        'uns_connect_node_info'
    ]
    LOOP
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = current_schema()
              AND table_name = tbl
              AND column_name = 'created_time'
              AND data_type = 'bigint'
        ) THEN
            EXECUTE format('ALTER TABLE %I ALTER COLUMN created_time DROP DEFAULT', tbl);
            EXECUTE format(
                'ALTER TABLE %I ALTER COLUMN created_time TYPE TIMESTAMPTZ USING CASE WHEN created_time > 0 THEN to_timestamp(created_time / 1000.0) ELSE CURRENT_TIMESTAMP END',
                tbl
            );
            EXECUTE format('ALTER TABLE %I ALTER COLUMN created_time SET DEFAULT CURRENT_TIMESTAMP', tbl);
            EXECUTE format('ALTER TABLE %I ALTER COLUMN created_time SET NOT NULL', tbl);
        END IF;

        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = current_schema()
              AND table_name = tbl
              AND column_name = 'updated_time'
              AND data_type = 'bigint'
        ) THEN
            EXECUTE format('ALTER TABLE %I ALTER COLUMN updated_time DROP DEFAULT', tbl);
            EXECUTE format(
                'ALTER TABLE %I ALTER COLUMN updated_time TYPE TIMESTAMPTZ USING CASE WHEN updated_time > 0 THEN to_timestamp(updated_time / 1000.0) ELSE CURRENT_TIMESTAMP END',
                tbl
            );
            EXECUTE format('ALTER TABLE %I ALTER COLUMN updated_time SET DEFAULT CURRENT_TIMESTAMP', tbl);
            EXECUTE format('ALTER TABLE %I ALTER COLUMN updated_time SET NOT NULL', tbl);
        END IF;
    END LOOP;
END $$;
