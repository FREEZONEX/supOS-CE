CREATE TABLE IF NOT EXISTS uns_cloud_sync_log (
    sync_type VARCHAR(32) NOT NULL DEFAULT 'uns',
    details JSONB NOT NULL DEFAULT '{}'::JSONB
);
