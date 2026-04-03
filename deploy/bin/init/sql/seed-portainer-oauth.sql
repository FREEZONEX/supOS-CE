\set ON_ERROR_STOP on

INSERT INTO supos.supos_oauth_client (
  client_id,
  client_secret,
  client_name,
  redirect_uris,
  scopes,
  grant_types,
  enabled,
  trusted,
  created_at,
  updated_at
) VALUES (
  :'client_id',
  :'client_secret',
  'portainer',
  :'redirect_uri',
  'openid',
  'authorization_code',
  true,
  true,
  NOW(),
  NOW()
)
ON CONFLICT (client_id) DO UPDATE SET
  client_secret = EXCLUDED.client_secret,
  client_name = EXCLUDED.client_name,
  redirect_uris = EXCLUDED.redirect_uris,
  scopes = EXCLUDED.scopes,
  grant_types = EXCLUDED.grant_types,
  enabled = EXCLUDED.enabled,
  trusted = EXCLUDED.trusted,
  updated_at = NOW();
