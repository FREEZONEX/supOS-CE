-- Backfill asset.file.download into pre-existing API keys.
--
-- apikey_scope.go added asset.file.download to the default read scope (and,
-- transitively, to the write/full default lists) when the OpenAPI asset file
-- endpoints were introduced. API keys created before that change keep their
-- old comma-separated resource_keys and would receive 403 from
-- /openapi/v1/assets/files/url and /openapi/v1/assets/files/download.
--
-- resource_keys is stored as a comma-separated TEXT list (see
-- APIKeyRepo.CreateAPIKey -> strings.Join(resourceKeys, ",")). Every
-- permission level (read_only/data_writer/full_access) includes the read
-- scope, so the key is backfilled on all live rows. The statement is
-- idempotent: rows already containing the key are untouched, and the
-- versioned migration ledger runs this file only once per database.
UPDATE sys_api_key
SET resource_keys = CASE
        WHEN BTRIM(resource_keys) = '' THEN 'asset.file.download'
        ELSE resource_keys || ',asset.file.download'
    END
WHERE deleted_time = 0
  AND (',' || resource_keys || ',') NOT LIKE '%,asset.file.download,%';
