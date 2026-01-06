
CREATE TABLE if not exists uns_namespace (
	"id" bigint PRIMARY KEY NOT NULL,
	"lay_rec" text NOT NULL,
	"alias" varchar(128) NOT NULL,
	"parent_alias" varchar(128) NULL,
	"name" varchar(512) NOT NULL,
	"path" text NOT NULL,
	"path_type" int2 NOT NULL,
	"data_type" int2 NULL,
	"fields" json NULL,
	"create_at" timestamptz DEFAULT now() NULL,
	"status" smallint DEFAULT 1 NULL,
	"description" varchar(255),
	"update_at" timestamptz NULL,
	"protocol" varchar(2000) NULL,
	"data_path" varchar(128) NULL,
	 "with_flags" integer NULL default 0,
	 "data_src_id" int2 NULL,
	 "ref_uns" jsonb default '{}',
	 "refers" json NULL,
	 "expression" varchar(255) NULL,
	 "table_name" varchar(190) NULL,
	 "number_fields" int2 default NULL,
	 "parent_id" bigint default NULL,
	 "model_id" bigint default NULL,
	 "protocol_type" varchar(64) NULL,
	 "extend" jsonb DEFAULT '{}'
);
ALTER TABLE uns_namespace ALTER COLUMN fields TYPE json USING fields::json;
ALTER TABLE uns_namespace ALTER COLUMN refers TYPE json USING refers::json;
ALTER TABLE uns_namespace ALTER COLUMN extend TYPE jsonb USING extend::jsonb;
ALTER TABLE uns_namespace ADD IF NOT EXISTS "label_ids" jsonb default NULL;
ALTER TABLE uns_namespace ALTER COLUMN "label_ids" TYPE jsonb USING label_ids::jsonb;
ALTER TABLE uns_namespace ADD IF NOT EXISTS "subscribe_at" timestamptz NULL;
CREATE UNIQUE INDEX if not exists idx_uns_spacex_alias ON uns_namespace (alias);

ALTER TABLE uns_namespace ALTER COLUMN update_at SET DEFAULT now();
CREATE INDEX CONCURRENTLY if not exists idx_uns_namespace_update_at ON uns_namespace (update_at);
update uns_namespace set update_at=create_at where update_at is null;
CREATE INDEX if not exists idx_uns_namespace_parent_id ON uns_namespace(parent_id);
ALTER TABLE uns_namespace ADD IF NOT EXISTS pathash INTEGER;
CREATE INDEX if not exists idx_path_hash ON uns_namespace (pathash);
CREATE EXTENSION IF NOT EXISTS pg_trgm;
ALTER TABLE uns_namespace ADD COLUMN IF NOT EXISTS fields_text TEXT;
CREATE INDEX IF NOT EXISTS idx_namespace_fields_text ON uns_namespace USING GIN (fields_text gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idxgin_namespace_alias ON uns_namespace USING GIN (alias gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idxgin_namespace_path ON uns_namespace USING GIN (path gin_trgm_ops);

insert into uns_namespace("id","path_type","lay_rec","alias","name","path","description")values(0,1,'0','__templates__','tmplt','tmplt','模板顶级目录')ON CONFLICT (id) DO UPDATE set lay_rec=EXCLUDED.lay_rec,path=EXCLUDED.path,name=EXCLUDED.name;

CREATE TABLE if not exists uns_dashboard (
	id varchar(64) PRIMARY KEY NOT NULL,
	"name" varchar(255) NULL,
	description varchar(255) NULL,
    "json_content" text NULL,
	update_time timestamp(6) NULL,
	create_time timestamp(6) NULL
);

CREATE TABLE if not exists "uns_dashboard_ref" (
"dashboard_id" varchar(64) NOT NULL,
"uns_alias" varchar(255) NOT NULL,
"create_at" timestamptz(6) DEFAULT now(),
PRIMARY KEY (dashboard_id, uns_alias)
);
DELETE FROM "uns_dashboard_ref"
WHERE ctid IN (
    SELECT ctid
    FROM (
             SELECT
                 ctid,
                 ROW_NUMBER() OVER (
                PARTITION BY dashboard_id, uns_alias
                ORDER BY ctid
            ) AS row_num
             FROM "uns_dashboard_ref"
         ) t
    WHERE t.row_num > 1
);
ALTER TABLE uns_dashboard_ref ADD PRIMARY KEY (dashboard_id, uns_alias);

CREATE TABLE if not exists uns_dashboard_top_recodes (
    id varchar(64) NOT NULL,
	user_id varchar Not NULL,
	mark int2 default 1,
	mark_time timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	update_time timestamptz NULL
);
CREATE UNIQUE INDEX udx_dashboard_id_user ON uns_dashboard_top_recodes (id, user_id);
alter table uns_dashboard add if not exists "creator" varchar(128) NULL;

CREATE INDEX if not exists "idx_user_id" ON "supos_user_menu" USING btree ("user_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST);

CREATE TABLE if not exists "uns_alarms_data" (
"_id" BIGSERIAL PRIMARY KEY,
"uns" bigint NOT NULL,
"uns_path" varchar(255) NULL,
"current_value" float4,
"limit_value" float4,
"is_alarm" bool DEFAULT true,
"read_status" bool DEFAULT false,
"_ct" timestamptz(6) DEFAULT now()
);

ALTER TABLE "uns_alarms_data" ALTER COLUMN "read_status" SET DEFAULT false;
ALTER TABLE "uns_alarms_data" ALTER COLUMN "is_alarm" SET DEFAULT true;
ALTER TABLE "uns_alarms_data" ADD IF NOT EXISTS "uns_path" varchar(255) NULL;

CREATE index if not exists uns_alarms_data_uns_idx ON "uns_alarms_data" ("uns");

alter table "uns_dashboard"  add if not exists "type" int2 DEFAULT 1;

CREATE TABLE if not exists "uns_tag" (
"id" int8 NOT NULL,
"topic" varchar(255) COLLATE "pg_catalog"."default",
"tag_name" varchar(255) COLLATE "pg_catalog"."default",
"is_deleted" bool DEFAULT false,
"create_at" timestamptz(6) DEFAULT now(),
CONSTRAINT "uns_tag_pkey" PRIMARY KEY ("id")
);

CREATE TABLE if not exists uns_attachment (
    "id" bigint NOT NULL PRIMARY KEY,
	"uns_alias" varchar(128) NOT NULL,
	"original_name" varchar(255) NOT NULL,
	"attachment_name" varchar(255) NOT NULL,
	"attachment_path" varchar(255) NULL,
	"extension_name" varchar(20) NULL,
	"create_at" timestamptz DEFAULT now() NULL
);

CREATE TABLE if not exists "uns_label" (
"id" BIGSERIAL PRIMARY KEY,
"label_name" varchar(255) COLLATE "pg_catalog"."default",
"create_at" timestamptz(6) DEFAULT now()
);

alter table uns_label add if not exists "with_flags" integer NULL default 0;
alter table uns_label add if not exists "subscribe_frequency"  varchar(20) NULL;
alter table uns_label add if not exists "subscribe_at" timestamptz(6) NULL;
alter table uns_label add if not exists "update_at" timestamptz(6) DEFAULT now();

CREATE TABLE if not exists "uns_label_ref" (
"id" BIGSERIAL PRIMARY KEY,
"label_id" int8 NOT NULL,
"uns_id" bigint NOT NULL,
"create_at" timestamptz(6) DEFAULT now()
);

ALTER TABLE uns_label_ref DROP COLUMN id,
DROP CONSTRAINT if exists uns_label_ref_pkey,
ADD PRIMARY KEY (label_id, uns_id);

CREATE TABLE if not exists "supos_todo" (
"id" BIGSERIAL PRIMARY KEY,
"user_id" varchar(64) NOT NULL,
"username" varchar(64) NOT NULL,
"module_code" varchar(32),
"module_name" varchar(32),
"status" smallint DEFAULT 0 NULL,
"todo_msg" varchar(256) ,
"business_id" bigint,
"link" varchar(512),
"handler_user_id" varchar(64),
"handler_username" varchar(64),
"handler_time" timestamptz(6),
"create_at" timestamptz(6) DEFAULT now()
);

alter table supos_todo add if not exists "handler_time" timestamptz(6);
alter table supos_todo add if not exists "module_name" varchar(32);

COMMENT ON COLUMN "supos_todo"."user_id" IS '用户ID';
COMMENT ON COLUMN "supos_todo"."username" IS '用户名';
COMMENT ON COLUMN "supos_todo"."module_code" IS '模块编码';
COMMENT ON COLUMN "supos_todo"."status" IS '代办状态：0-未处理 1-已处理';
COMMENT ON COLUMN "supos_todo"."todo_msg" IS '事项信息';
COMMENT ON COLUMN "supos_todo"."business_id" IS '业务主键';
COMMENT ON COLUMN "supos_todo"."link" IS '链接';
COMMENT ON COLUMN "supos_todo"."handler_user_id" IS '处理人用户ID';
COMMENT ON COLUMN "supos_todo"."handler_username" IS '处理人用户名';
COMMENT ON COLUMN "supos_todo"."create_at" IS '创建时间';


CREATE TABLE if not exists "supos_example" (
"id" BIGSERIAL PRIMARY KEY,
"name" varchar(255) COLLATE "pg_catalog"."default",
"description" varchar(512) COLLATE "pg_catalog"."default",
"package_path" varchar(512) COLLATE "pg_catalog"."default",
"status" int2,
"type" int2,
"dashboard_type" int2,
"dashboard_id" varchar(64) COLLATE "pg_catalog"."default",
"dashboard_name" varchar(512) COLLATE "pg_catalog"."default",
"create_at" timestamptz(6) DEFAULT now());

COMMENT ON COLUMN "supos_example"."status" IS '安装状态：1-未安装，2-安装中，3已安装';
COMMENT ON COLUMN "supos_example"."type" IS '类型：1-OT 2-IT';

INSERT INTO "supos_example" ("id", "name", "description", "package_path", "status", "type", "dashboard_type", "dashboard_id", "dashboard_name", "create_at") VALUES (1, 'ot-demo', 'ot-demo', '/templates/example/ot.zip', 1, 1, NULL, NULL, NULL, '2025-02-25 08:31:06.039+00') ON CONFLICT (id) DO NOTHING;
INSERT INTO "supos_example" ("id", "name", "description", "package_path", "status", "type", "dashboard_type", "dashboard_id", "dashboard_name", "create_at") VALUES (2, 'it-demo', 'it-demo', '/templates/example/it.zip', 1, 2, NULL, NULL, NULL, '2025-02-25 08:31:06.039+00') ON CONFLICT (id) DO NOTHING;

CREATE TABLE if not exists "uns_alarms_handler" (
"id" BIGSERIAL PRIMARY KEY,
"uns_id" bigint,
"user_id" varchar(64),
"username" varchar(256),
"create_at" timestamptz(6) DEFAULT now());


alter table supos_todo add if not exists "process_id" int8 NULL;
alter table supos_todo add if not exists "process_instance_id" varchar(64) NULL;

CREATE DATABASE camunda;

CREATE TABLE if not exists "supos_workflow_process" (
"id" BIGSERIAL PRIMARY KEY,
"description" varchar(512),
"process_definition_id" varchar(64),
"process_definition_name" varchar(256),
"process_definition_key" varchar(256),
"status" int2 default 0,
"deploy_id" varchar(64),
"deploy_name" varchar(256),
"deploy_time" timestamptz(6),
"bpmn_xml" text,
"create_at" timestamptz(6) DEFAULT now());

alter table uns_namespace add if not exists "display_name" varchar(512) NULL;

CREATE TABLE if not exists "supos_app_key" (
"id" BIGSERIAL PRIMARY KEY,
"app_secret_key" varchar(200) NOT NULL,
"app_secret_value" varchar(200) NOT NULL,
"status" int2 default 1,
"create_time" timestamptz(6) DEFAULT now());

CREATE TABLE if not exists "global_export_record" (
	id int8 NOT NULL,
	user_id varchar(64) NULL,
	file_path varchar(2000) NULL,
	create_time timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	update_time timestamptz NULL,
	confirm bool NULL,
	CONSTRAINT global_export_record_pk PRIMARY KEY (id)
);

CREATE TABLE if not exists "uns_history_delete_job" (
    "id" BIGSERIAL PRIMARY KEY,
    "alias" varchar(128) NOT NULL,
    "name" varchar(512) NOT NULL,
    "table_name" varchar(128) NULL,
    "path" text NOT NULL,
    "path_type" int2 NOT NULL,
    "data_type" int2 NULL,
    "fields" json NULL,
    "status" smallint DEFAULT 1 NULL,
    "create_at" timestamptz DEFAULT now() NULL
);
CREATE INDEX if not exists idx_uns_delete_alias ON uns_history_delete_job (alias);

CREATE TABLE if not exists "uns_person_config" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" varchar(64) NOT NULL,
  "main_language" varchar(64) NULL,
  "create_at" timestamptz(6) DEFAULT now(),
  "update_at" timestamptz(6) DEFAULT now()
);

CREATE UNIQUE INDEX if not exists uq_idx_uns_person_config_uid ON uns_person_config (user_id);

CREATE TABLE if not exists "uns_export_record" (
	id int8 NOT NULL,
	user_id varchar(64) NULL,
	file_path varchar(2000) NULL,
	create_time timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	update_time timestamptz NULL,
	confirm bool NULL,
	CONSTRAINT uns_export_record_pk PRIMARY KEY (id)
);

alter table uns_namespace add if not exists "extend_field_flags" integer NULL default 0;

alter table uns_namespace add if not exists "mount_type" int2 DEFAULT 0 NULL;
alter table uns_namespace add if not exists "mount_source" varchar(256) NULL;

CREATE TABLE if not exists "uns_mount" (
"id" BIGSERIAL PRIMARY KEY,
"mount_seq" varchar(64),
"target_type" varchar(20),
"target_alias" varchar(64),
"mount_model" varchar(20),
"source_alias" varchar(1024),
"mount_status" int2,
"status" varchar(20),
"data_type" int2,
"with_flags" integer NULL default 0,
"version" varchar(20),
"next_version" varchar(20)
);

CREATE TABLE if not exists "uns_mount_extend" (
"id" BIGSERIAL PRIMARY KEY,
"source_sub_type" varchar(20),
"mount_seq" varchar(64),
"target_alias" varchar(64),
"first_source_alias" varchar(1024),
"second_source_alias" varchar(1024),
"source_name" varchar(1024),
"extend" text
);
alter table uns_namespace add if not exists "parent_data_type" int2 NULL;

alter table uns_dashboard add if not exists "need_init" bool DEFAULT false;
