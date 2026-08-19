#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
FROM=""
RESTORE_DB="true"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --from)
      FROM="${2:?missing backup file}"
      shift 2
      ;;
    --no-db)
      RESTORE_DB="false"
      shift
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2
      ;;
  esac
done

if [[ -z "${FROM}" ]]; then
  echo "--from is required" >&2
  exit 2
fi
if [[ ! -f "${FROM}" ]]; then
  echo "backup file not found: ${FROM}" >&2
  exit 1
fi

TMP="$(mktemp -d)"
trap 'rm -rf "${TMP}"' EXIT
tar -xzf "${FROM}" -C "${TMP}"

if [[ -d "${TMP}/deploy" ]]; then
  for file in .env .env.runtime docker-compose.yml docker-compose.override.yml; do
    if [[ -f "${TMP}/deploy/${file}" ]]; then
      cp "${TMP}/deploy/${file}" "${DEPLOY_ROOT}/${file}"
    fi
  done
fi

if [[ ! -f "${DEPLOY_ROOT}/.env" ]]; then
  echo "backup did not contain deploy/.env" >&2
  exit 1
fi

source_env_pair "${DEPLOY_ROOT}/.env" "${DEPLOY_ROOT}/.env.runtime"

UNS_DB_NAME="${UNS_DB_NAME:-postgres}"
UNS_DB_URL="${UNS_DB_URL:-}"

is_local_tsdb_url() {
  local url="$1"
  [[ "${url}" == *"@tsdb:"* || "${url}" == *"//tsdb:"* ]]
}

is_legacy_postgresql_url() {
  local url="$1"
  [[ "${url}" == *"@postgresql:"* || "${url}" == *"//postgresql:"* ]]
}

LOCAL_RESTORE_SERVICE=""
LOCAL_RESTORE_PASSWORD=""
if is_local_tsdb_url "${UNS_DB_URL}"; then
  LOCAL_RESTORE_SERVICE="tsdb"
  LOCAL_RESTORE_PASSWORD="${TSDB_PASSWORD:-}"
elif is_legacy_postgresql_url "${UNS_DB_URL}"; then
  LOCAL_RESTORE_SERVICE="postgresql"
  LOCAL_RESTORE_PASSWORD="${POSTGRES_PASSWORD:-}"
fi

DB_DUMP_FILE=""
if [[ -f "${TMP}/db/database.sql" ]]; then
  DB_DUMP_FILE="${TMP}/db/database.sql"
elif [[ -f "${TMP}/db/platform.sql" ]]; then
  DB_DUMP_FILE="${TMP}/db/platform.sql"
fi

if [[ "${RESTORE_DB}" == "true" && -n "${DB_DUMP_FILE}" ]]; then
  if [[ -z "${LOCAL_RESTORE_SERVICE}" ]]; then
    echo "cannot restore platform DB through external UNS_DB_URL" >&2
    exit 2
  fi
  bash "${DEPLOY_ROOT}/bin/compose.sh" up -d "${LOCAL_RESTORE_SERVICE}"
  exists="$(bash "${DEPLOY_ROOT}/bin/compose.sh" exec -T -e PGPASSWORD="${LOCAL_RESTORE_PASSWORD}" "${LOCAL_RESTORE_SERVICE}" \
    psql -U postgres -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${UNS_DB_NAME}'" | tr -d '[:space:]')"
  if [[ "${exists}" != "1" ]]; then
    bash "${DEPLOY_ROOT}/bin/compose.sh" exec -T -e PGPASSWORD="${LOCAL_RESTORE_PASSWORD}" "${LOCAL_RESTORE_SERVICE}" \
      psql -U postgres -d postgres -v ON_ERROR_STOP=1 -c "CREATE DATABASE \"${UNS_DB_NAME}\"" >/dev/null
  fi
  bash "${DEPLOY_ROOT}/bin/compose.sh" exec -T -e PGPASSWORD="${LOCAL_RESTORE_PASSWORD}" "${LOCAL_RESTORE_SERVICE}" \
    psql -U postgres -d "${UNS_DB_NAME}" < "${DB_DUMP_FILE}" >/dev/null

fi


echo "restored current deploy from ${FROM}"
