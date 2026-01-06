CREATE TABLE if not exists supos_node_flows (
                                                id bigint NOT NULL PRIMARY KEY,
                                                flow_id varchar NULL,
                                                flow_name varchar Not NULL,
                                                flow_status varchar Not NULL,
                                                flow_data TEXT NULL,
                                                description varchar NULL,
                                                template varchar null DEFAULT 'node-red',
                                                create_time timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
                                                update_time timestamptz NULL
);

CREATE TABLE if not exists supos_node_flow_models (
                                                      parent_id bigint NULL,
                                                      topic varchar NULL,
                                                      alias varchar NULL,
                                                      create_time timestamptz NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX if not exists  idx_flow_alias ON supos_node_flow_models (alias);

alter table supos_node_flow_models add if not exists "alias" varchar(128) NULL;

CREATE TABLE if not exists supos_node_flow_top_recodes (
                                                           id bigint NOT NULL,
                                                           user_id varchar Not NULL,
                                                           mark int2 default 1,
                                                           mark_time timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
                                                           update_time timestamptz NULL
);
CREATE UNIQUE INDEX idx_unique_id_user ON supos_node_flow_top_recodes (id, user_id);
alter table supos_node_flows add if not exists "creator" varchar(128) NULL;