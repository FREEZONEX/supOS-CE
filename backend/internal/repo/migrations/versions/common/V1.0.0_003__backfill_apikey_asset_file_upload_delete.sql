-- Backfill asset.file.upload / asset.file.delete into pre-existing API keys.
--
-- apikey_scope.go added asset.file.upload to the default write scope (and,
-- transitively, to the full default list) and asset.file.delete to the
-- default full scope when the OpenAPI asset file endpoints were introduced.
-- V1.0.0_002 only backfilled asset.file.download, so API keys created before
-- this change keep their old comma-separated resource_keys and would still
-- receive 403 from POST /openapi/v1/assets/files (upload) and
-- /openapi/v1/assets/files/delete.
--
-- resource_keys is stored as a comma-separated TEXT list (see
-- APIKeyRepo.CreateAPIKey -> strings.Join(resourceKeys, ",")). Permission
-- levels compose read < write < full, so upload goes to data_writer and
-- full_access rows, delete only to full_access rows. Both statements are
-- idempotent: rows already containing the key are untouched, and the
-- versioned migration ledger runs this file only once per database.
UPDATE sys_api_key
SET resource_keys = CASE
        WHEN BTRIM(resource_keys) = '' THEN 'asset.file.upload'
        ELSE resource_keys || ',asset.file.upload'
    END
WHERE deleted_time = 0
  AND permission IN ('data_writer', 'full_access')
  AND (',' || resource_keys || ',') NOT LIKE '%,asset.file.upload,%';

UPDATE sys_api_key
SET resource_keys = CASE
        WHEN BTRIM(resource_keys) = '' THEN 'asset.file.delete'
        ELSE resource_keys || ',asset.file.delete'
    END
WHERE deleted_time = 0
  AND permission = 'full_access'
  AND (',' || resource_keys || ',') NOT LIKE '%,asset.file.delete,%';
