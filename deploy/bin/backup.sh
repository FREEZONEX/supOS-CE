#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${DEPLOY_ROOT}/.env"
RUNTIME_ENV_FILE="${DEPLOY_ROOT}/.env.runtime"
OUT_DIR="${DEPLOY_ROOT}/backups"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --out)
      OUT_DIR="${2:?missing output directory}"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2
      ;;
  esac
done

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "deploy env not found: ${ENV_FILE}" >&2
  exit 2
fi

source_env_pair "${ENV_FILE}" "${RUNTIME_ENV_FILE}"

BACKUP_NAME="${OS_NAME:-${COMPOSE_PROJECT_NAME:-edge}}"
BACKUP_NAME="$(printf '%s' "${BACKUP_NAME}" | tr -c 'A-Za-z0-9_.-' '_')"
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

LOCAL_DUMP_SERVICE=""
LOCAL_DUMP_PASSWORD=""
if is_local_tsdb_url "${UNS_DB_URL}"; then
  LOCAL_DUMP_SERVICE="tsdb"
  LOCAL_DUMP_PASSWORD="${TSDB_PASSWORD:-}"
elif is_legacy_postgresql_url "${UNS_DB_URL}"; then
  LOCAL_DUMP_SERVICE="postgresql"
  LOCAL_DUMP_PASSWORD="${POSTGRES_PASSWORD:-}"
fi

TS="$(date +%Y%m%d%H%M%S)"
OUT="${OUT_DIR}/backup-${BACKUP_NAME}-${TS}.tar.gz"
WORK="$(mktemp -d)"
trap 'rm -rf "${WORK}"' EXIT

mkdir -p "${OUT_DIR}" "${WORK}/deploy" "${WORK}/db" "${WORK}/meta"
for file in .env .env.runtime docker-compose.yml docker-compose.override.yml; do
  if [[ -f "${DEPLOY_ROOT}/${file}" ]]; then
    cp "${DEPLOY_ROOT}/${file}" "${WORK}/deploy/${file}"
  fi
done

VOLUMES_ABS=""
if [[ -n "${VOLUMES_PATH:-}" ]]; then
  if [[ "${VOLUMES_PATH}" == /* || "${VOLUMES_PATH}" == [A-Za-z]:/* || "${VOLUMES_PATH}" == [A-Za-z]:\\* ]]; then
    VOLUMES_ABS="${VOLUMES_PATH}"
  else
    VOLUMES_ABS="${DEPLOY_ROOT}/${VOLUMES_PATH}"
  fi
  if [[ -f "${VOLUMES_ABS}/.install-state.json" ]]; then
    cp "${VOLUMES_ABS}/.install-state.json" "${WORK}/deploy/.install-state.json"
  fi
fi

cat > "${WORK}/meta/backup.json" <<JSON
{"osName":"${OS_NAME:-}","composeProjectName":"${COMPOSE_PROJECT_NAME:-}","createdAt":"$(date -u +%Y-%m-%dT%H:%M:%SZ)","format":"deploy-backup-v2"}
JSON

if [[ -n "${LOCAL_DUMP_SERVICE}" ]]; then
  if bash "${DEPLOY_ROOT}/bin/compose.sh" ps --status running "${LOCAL_DUMP_SERVICE}" >/dev/null 2>&1; then
    bash "${DEPLOY_ROOT}/bin/compose.sh" exec -T -e PGPASSWORD="${LOCAL_DUMP_PASSWORD}" "${LOCAL_DUMP_SERVICE}" \
      pg_dump -U postgres -d "${UNS_DB_NAME}" --clean --if-exists --no-owner --no-privileges > "${WORK}/db/database.sql"

  fi
else
  echo "warning: skipping platform DB dump for external UNS_DB_URL" >&2
fi

tar -czf "${OUT}" -C "${WORK}" .
echo "${OUT}"
