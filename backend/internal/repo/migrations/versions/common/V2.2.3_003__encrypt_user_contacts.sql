ALTER TABLE sys_user_info ALTER COLUMN email TYPE TEXT;
ALTER TABLE sys_user_info ALTER COLUMN phone TYPE TEXT;
DROP INDEX IF EXISTS idx_sys_user_info_tc_email;

ALTER TABLE sys_user_oper_log DROP COLUMN IF EXISTS operator_email;

CREATE OR REPLACE FUNCTION pg_temp.tier0_strip_user_contact_pii(value JSONB)
RETURNS JSONB
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
    stripped JSONB;
BEGIN
    CASE jsonb_typeof(value)
        WHEN 'object' THEN
            SELECT COALESCE(
                jsonb_object_agg(item.key, pg_temp.tier0_strip_user_contact_pii(item.value)),
                '{}'::JSONB
            )
            INTO stripped
            FROM jsonb_each(value) AS item
            WHERE lower(item.key) NOT IN ('email', 'phone');
            RETURN stripped;
        WHEN 'array' THEN
            SELECT COALESCE(
                jsonb_agg(pg_temp.tier0_strip_user_contact_pii(item.value)),
                '[]'::JSONB
            )
            INTO stripped
            FROM jsonb_array_elements(value) AS item;
            RETURN stripped;
        ELSE
            RETURN value;
    END CASE;
END;
$$;

CREATE OR REPLACE FUNCTION pg_temp.tier0_strip_user_contact_pii_text(value TEXT)
RETURNS TEXT
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
    IF value IS NULL OR btrim(value) = '' THEN
        RETURN value;
    END IF;
    RETURN pg_temp.tier0_strip_user_contact_pii(value::JSONB)::TEXT;
EXCEPTION
    WHEN others THEN
        RETURN '{}'::TEXT;
END;
$$;

UPDATE sys_user_oper_log
SET detail = pg_temp.tier0_strip_user_contact_pii_text(detail)
WHERE detail IS NOT NULL AND btrim(detail) <> '';

UPDATE uns_cloud_sync_log
SET details = pg_temp.tier0_strip_user_contact_pii(details)
WHERE sync_type = 'user';

UPDATE sys_outbox_event
SET payload = pg_temp.tier0_strip_user_contact_pii(payload)
WHERE event_type = 'fleet.user.sync.result.ready';
