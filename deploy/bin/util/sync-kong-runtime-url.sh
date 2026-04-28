#!/bin/bash

set -e

SYNC_KONG_URL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
DEPLOY_ROOT="$(cd "$SYNC_KONG_URL_DIR/../.." && pwd)"
ENV_FILE="$DEPLOY_ROOT/.env.default"
if [ -f "$DEPLOY_ROOT/.env" ]; then
  ENV_FILE="$DEPLOY_ROOT/.env"
fi

source "$SYNC_KONG_URL_DIR/../global/log.sh"

set -a
source "$ENV_FILE"
set +a

if [ -f "$DEPLOY_ROOT/.env.tmp" ]; then
  set -a
  source "$DEPLOY_ROOT/.env.tmp"
  set +a
fi

POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-postgres}"

# These plugin IDs come from deploy/mount/kong/kong_config.yml.tpl and are
# reserved for system-managed entrance URL settings only.
AUTH_CHECKER_PLUGIN_ID="1845ee75-d704-40e1-a8b0-aa2baaf9d71b"
PORTAINER_REDIRECT_PLUGIN_ID="46bda5cf-63ea-401f-9f06-b9e024aa5597"
PORTAINER_ORIGIN_PLUGIN_ID="476cc68a-e2be-4cf1-831b-334b61f97ae4"

trim() {
  printf '%s' "$1" | tr -d '[:space:]'
}

sql_escape() {
  printf '%s' "$1" | sed "s/'/''/g"
}

build_base_url() {
  local base_url=""

  if [ "$ENTRANCE_PROTOCOL" == "http" ]; then
    base_url="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_PORT"
    if [[ "$ENTRANCE_PORT" == "80" ]]; then
      base_url="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
    fi
  fi

  if [ "$ENTRANCE_PROTOCOL" == "https" ]; then
    base_url="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_SSL_PORT"
    if [[ "$ENTRANCE_SSL_PORT" == "443" ]]; then
      base_url="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
    fi
  fi

  printf '%s' "$base_url"
}

psql_scalar() {
  docker exec -i -e "PGPASSWORD=${POSTGRES_PASSWORD}" postgresql \
    psql -U postgres -d kong -v ON_ERROR_STOP=1 -At "$@"
}

plugin_exists() {
  local plugin_id="$1"
  local exists
  exists="$(trim "$(psql_scalar -c "SELECT 1 FROM plugins WHERE id = '$plugin_id'::uuid;")")"
  if [ "$exists" == "1" ]; then
    return 0
  fi
  return 1
}

BASE_URL="$(build_base_url)"
LOGIN_URL="${BASE_URL}/tier0-login"
PORTAINER_REDIRECT_QUERY="redirect_uri:${BASE_URL}/portainer/home/"
PORTAINER_ORIGIN_HEADER="Origin:${BASE_URL}"
PORTAINER_REFERER_HEADER="Referer:${BASE_URL}/portainer/home/"

LOGIN_URL_SQL="$(sql_escape "$LOGIN_URL")"
PORTAINER_REDIRECT_QUERY_SQL="$(sql_escape "$PORTAINER_REDIRECT_QUERY")"
PORTAINER_ORIGIN_HEADER_SQL="$(sql_escape "$PORTAINER_ORIGIN_HEADER")"
PORTAINER_REFERER_HEADER_SQL="$(sql_escape "$PORTAINER_REFERER_HEADER")"

auth_changed="$(trim "$(psql_scalar \
  -c "WITH updated AS (
        UPDATE plugins
           SET config = jsonb_set(
             config,
             '{login_url}',
             to_jsonb('$LOGIN_URL_SQL'::text),
             false
           )
         WHERE id = '$AUTH_CHECKER_PLUGIN_ID'::uuid
           AND config #>> '{login_url}' IS DISTINCT FROM '$LOGIN_URL_SQL'
         RETURNING 1
      )
      SELECT count(*) FROM updated;")")"

redirect_changed="$(trim "$(psql_scalar \
  -c "WITH updated AS (
        UPDATE plugins
           SET config = jsonb_set(
             config,
             '{append,querystring,3}',
             to_jsonb('$PORTAINER_REDIRECT_QUERY_SQL'::text),
             false
           )
         WHERE id = '$PORTAINER_REDIRECT_PLUGIN_ID'::uuid
           AND config #>> '{append,querystring,3}' IS DISTINCT FROM '$PORTAINER_REDIRECT_QUERY_SQL'
         RETURNING 1
      )
      SELECT count(*) FROM updated;")")"

origin_changed="$(trim "$(psql_scalar \
  -c "WITH updated AS (
        UPDATE plugins
           SET config = jsonb_set(
             jsonb_set(
               config,
               '{replace,headers,0}',
               to_jsonb('$PORTAINER_ORIGIN_HEADER_SQL'::text),
               false
             ),
             '{replace,headers,1}',
             to_jsonb('$PORTAINER_REFERER_HEADER_SQL'::text),
             false
           )
         WHERE id = '$PORTAINER_ORIGIN_PLUGIN_ID'::uuid
           AND (
             config #>> '{replace,headers,0}' IS DISTINCT FROM '$PORTAINER_ORIGIN_HEADER_SQL'
             OR config #>> '{replace,headers,1}' IS DISTINCT FROM '$PORTAINER_REFERER_HEADER_SQL'
           )
         RETURNING 1
      )
      SELECT count(*) FROM updated;")")"

changed=0

if [ "$auth_changed" -gt 0 ]; then
  info "Updated Kong login URL to ${LOGIN_URL}"
  changed=1
elif plugin_exists "$AUTH_CHECKER_PLUGIN_ID"; then
  info "Kong login URL already matches ${LOGIN_URL}"
else
  warn "Kong auth-checker plugin (${AUTH_CHECKER_PLUGIN_ID}) not found; skipped login URL sync"
fi

if [ "$redirect_changed" -gt 0 ]; then
  info "Updated Kong Portainer redirect querystring to ${PORTAINER_REDIRECT_QUERY}"
  changed=1
elif plugin_exists "$PORTAINER_REDIRECT_PLUGIN_ID"; then
  info "Kong Portainer redirect querystring already matches ${PORTAINER_REDIRECT_QUERY}"
else
  warn "Kong Portainer redirect plugin (${PORTAINER_REDIRECT_PLUGIN_ID}) not found; skipped redirect sync"
fi

if [ "$origin_changed" -gt 0 ]; then
  info "Updated Kong Portainer Origin/Referer headers to ${BASE_URL}"
  changed=1
elif plugin_exists "$PORTAINER_ORIGIN_PLUGIN_ID"; then
  info "Kong Portainer Origin/Referer headers already match ${BASE_URL}"
else
  warn "Kong Portainer Origin plugin (${PORTAINER_ORIGIN_PLUGIN_ID}) not found; skipped header sync"
fi

echo "SYNC_KONG_RUNTIME_URL_CHANGED=${changed}"
