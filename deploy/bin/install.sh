#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

NO_START="false"
ADOPT_VOLUMES="false"
TLS_MODE=""
VOLUMES_PATH_ARG=""
SKIP_NODERED_PACKAGES=""
LOCAL_FRONTEND_DEV_ARG=""


usage() {
  printf 'Tier0 Edge installer\n\n'
  cat <<'EOF'
Usage: bash bin/install.sh [options]

Options:
  --no-start                 Write configuration without starting services.
  --volumes-path <path>      Persistent data path for this member.
  --adopt-volumes            Adopt an existing volume path.
  --tls <mode>               TLS mode: off, self-signed, or provided.
  --skip-nodered-packages    Skip Node-RED package installation.
  --local-frontend-dev       Enable local frontend development proxying.
  --no-local-frontend-dev    Disable local frontend development proxying.
  -h, --help                 Show this help.

EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --no-start)
      NO_START="true"
      shift
      ;;
    --volumes-path)
      VOLUMES_PATH_ARG="${2:?missing volumes path}"
      shift 2
      ;;
    --adopt-volumes)
      ADOPT_VOLUMES="true"
      shift
      ;;
    --tls)
      TLS_MODE="${2:?missing tls mode}"
      shift 2
      ;;
    --skip-nodered-packages)
      SKIP_NODERED_PACKAGES="true"
      shift
      ;;
    --local-frontend-dev)
      LOCAL_FRONTEND_DEV_ARG="true"
      shift
      ;;
    --no-local-frontend-dev)
      LOCAL_FRONTEND_DEV_ARG="false"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2
      ;;
  esac
done

# shellcheck source=common.sh
source "${SCRIPT_DIR}/common.sh"

tier0_require_rootful_docker || exit 1
run_deploy_workflow
