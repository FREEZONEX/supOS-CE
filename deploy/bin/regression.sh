#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${DEPLOY_ROOT}/.env"
RUNTIME_ENV_FILE="${DEPLOY_ROOT}/.env.runtime"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"
if [[ -f "${ENV_FILE}" ]]; then
  source_env_pair "${ENV_FILE}" "${RUNTIME_ENV_FILE}"
fi

if [[ -z "${OPEN_SOURCE_REGRESSION_PASSWORD:-}" && -n "${ADMIN_INITIAL_PASSWORD:-}" ]]; then
  export OPEN_SOURCE_REGRESSION_PASSWORD="${ADMIN_INITIAL_PASSWORD}"
fi
if ! command -v node >/dev/null 2>&1; then
  echo "node is required for the runtime regression" >&2
  exit 2
fi
if ! command -v npx >/dev/null 2>&1; then
  echo "npx is required for the browser regression" >&2
  exit 2
fi

BASE_URL="${TIER0_REGRESSION_BASE_URL:-${ENTRANCE_URL:-}}"
if [[ -z "${BASE_URL}" ]]; then
  PROTOCOL="${ENTRANCE_PROTOCOL:-http}"
  HOST="${ENTRANCE_DOMAIN:-127.0.0.1}"
  [[ "${HOST}" == "0.0.0.0" || "${HOST}" == "localhost" ]] && HOST="127.0.0.1"
  PORT="${ENTRANCE_PORT:-8088}"
  [[ "${PROTOCOL}" == "https" ]] && PORT="${ENTRANCE_SSL_PORT:-8443}"
  BASE_URL="${PROTOCOL}://${HOST}:${PORT}"
fi
USERNAME="${TIER0_REGRESSION_USERNAME:-tier0}"
SCOPE="${TIER0_REGRESSION_SCOPE:-foundation}"
ARGS=(--base-url "${BASE_URL}" --username "${USERNAME}" --scope "${SCOPE}")
if [[ -n "${PRODUCT_VERSION:-}" ]]; then
  ARGS+=(--expected-version "${PRODUCT_VERSION}")
fi
case "${TIER0_REGRESSION_FAULT_INJECTION:-false}" in
  true|TRUE|1|yes|YES|on|ON)
    ARGS+=(--fault-injection true --deploy-root "${DEPLOY_ROOT}")
    ;;
esac

exec node "${DEPLOY_ROOT}/regression/runtime-regression.mjs" "${ARGS[@]}" "$@"
