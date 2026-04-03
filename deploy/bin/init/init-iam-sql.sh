#!/bin/bash

set -e

INIT_IAM_SQL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
ENV_FILE="$INIT_IAM_SQL_DIR/../../.env.default"
if [ -f "$INIT_IAM_SQL_DIR/../../.env" ]; then
  ENV_FILE="$INIT_IAM_SQL_DIR/../../.env"
fi

source "$INIT_IAM_SQL_DIR/../global/log.sh"
source "$ENV_FILE"
source "$INIT_IAM_SQL_DIR/../../.env.tmp"

IAM_PORTAINER_CLIENT_ID="${IAM_PORTAINER_CLIENT_ID:-portainer}"
IAM_PORTAINER_CLIENT_SECRET="${IAM_PORTAINER_CLIENT_SECRET:-Tier0PortainerSecret@1304}"
IAM_PORTAINER_REDIRECT_URI="${BASE_URL}/portainer/home/"
SQL_FILE="$INIT_IAM_SQL_DIR/sql/seed-portainer-oauth.sql"

psql_exec() {
  # Keep stdin attached because the seed SQL is piped in via `-f - < file`.
  docker exec -i -e "PGPASSWORD=${POSTGRES_PASSWORD}" postgresql \
    psql -U postgres -d postgres "$@"
}

wait_for_postgres() {
  local retries=30
  # Wait for the database container itself first so repeated installs on Linux
  # or Windows do not race the container startup sequence.
  until psql_exec -tAc "SELECT 1;" >/dev/null 2>&1; do
    retries=$((retries - 1))
    if [ "$retries" -le 0 ]; then
      warn "PostgreSQL is not ready. Skip IAM OAuth SQL initialization."
      return 1
    fi
    sleep 2
  done
  return 0
}

wait_for_iam_table() {
  local retries=30
  # The IAM OAuth table is created by backend migration, not by PostgreSQL init
  # scripts, so wait for it explicitly before seeding the Portainer client.
  until psql_exec -tAc "SELECT to_regclass('supos.supos_oauth_client');" | grep -q "supos.supos_oauth_client"; do
    retries=$((retries - 1))
    if [ "$retries" -le 0 ]; then
      warn "supos.supos_oauth_client is not ready. Skip IAM OAuth SQL initialization."
      return 1
    fi
    sleep 2
  done
  return 0
}

wait_for_postgres || return 0 2>/dev/null || exit 0
wait_for_iam_table || return 0 2>/dev/null || exit 0

# Keep this idempotent so install, changeIp, and repeated local redeploys all
# converge on the same Portainer OAuth client definition. Pipe the SQL file
# explicitly instead of relying on shell input redirection for a shell
# function, because Git Bash / PowerShell combinations can otherwise hand an
# empty stdin to `docker exec`.
cat "$SQL_FILE" | psql_exec \
  -v ON_ERROR_STOP=1 \
  -v client_id="$IAM_PORTAINER_CLIENT_ID" \
  -v client_secret="$IAM_PORTAINER_CLIENT_SECRET" \
  -v redirect_uri="$IAM_PORTAINER_REDIRECT_URI" \
  -f - >/dev/null

info "Initialized IAM OAuth client for Portainer: ${IAM_PORTAINER_CLIENT_ID}"
