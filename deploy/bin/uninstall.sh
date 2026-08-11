#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
if [[ $# -ne 0 ]]; then
  echo "usage: bash bin/uninstall.sh" >&2
  exit 2
fi
if [[ ! -f "${DEPLOY_ROOT}/.env" ]]; then
  echo "missing ${DEPLOY_ROOT}/.env; run bin/install.sh first" >&2
  exit 1
fi
source "${SCRIPT_DIR}/env-loader.sh"
source_env_pair "${DEPLOY_ROOT}/.env" "${DEPLOY_ROOT}/.env.runtime"
COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-edge}"
echo "Stopping Tier0 Edge deployment: ${COMPOSE_PROJECT_NAME}"
(cd "${DEPLOY_ROOT}" && bash "${SCRIPT_DIR}/compose.sh" down --remove-orphans)
rm -f "${DEPLOY_ROOT}/.env.tmp"
