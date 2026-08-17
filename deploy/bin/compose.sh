#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${DEPLOY_ROOT}/.env"
RUNTIME_ENV_FILE="${DEPLOY_ROOT}/.env.runtime"
EXPLICIT_BACKEND_IMAGE="${BACKEND_IMAGE:-}"
EXPLICIT_LOCAL_FRONTEND_DEV="${LOCAL_FRONTEND_DEV:-}"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"

# BACKEND_IMAGE is a deployment-time setting (controlled by .env or by the
# CI-provided environment variable). The generated .env.runtime must not
# override it; remove any stale entry so sudo/compose does not pick up an
# old manual value such as :branch-test.
if [[ -f "${RUNTIME_ENV_FILE}" ]]; then
  tmp_runtime_env="$(mktemp)"
  sed '/^[[:space:]]*\(export[[:space:]]\+[[:space:]]*\)\?BACKEND_IMAGE=/d' "${RUNTIME_ENV_FILE}" > "${tmp_runtime_env}"
  cat "${tmp_runtime_env}" > "${RUNTIME_ENV_FILE}"
  rm -f "${tmp_runtime_env}"
fi

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "deploy env not found: ${ENV_FILE}" >&2
  echo "run bash bin/install.sh --no-start first" >&2
  exit 2
fi

source_env_file "${ENV_FILE}"
# Capture the shell-expanded value from .env before .env.runtime is loaded.
# Release templates intentionally express image tags through PRODUCT_VERSION.
ENV_FILE_BACKEND_IMAGE="${BACKEND_IMAGE:-}"
if [[ -f "${RUNTIME_ENV_FILE}" ]]; then
  source_env_file "${RUNTIME_ENV_FILE}"
fi
if [[ -n "${EXPLICIT_BACKEND_IMAGE}" ]]; then
  BACKEND_IMAGE="${EXPLICIT_BACKEND_IMAGE}"
  export BACKEND_IMAGE
elif [[ -n "${ENV_FILE_BACKEND_IMAGE}" ]]; then
  # .env.runtime may still carry a stale BACKEND_IMAGE (e.g. from a manual
  # deployment). Prefer the value recorded in .env, which CI keeps current.
  BACKEND_IMAGE="${ENV_FILE_BACKEND_IMAGE}"
  export BACKEND_IMAGE
fi
if [[ -n "${EXPLICIT_LOCAL_FRONTEND_DEV}" ]]; then
  LOCAL_FRONTEND_DEV="${EXPLICIT_LOCAL_FRONTEND_DEV}"
  export LOCAL_FRONTEND_DEV
fi

normalize_windows_compose_paths() {
  local volumes_path="${VOLUMES_PATH:-}"
  local machine_id_path="${MACHINE_ID_HOST_PATH:-}"
  local machine_id_file
  local host_root
  local host_proc_path="${HOST_PROC_PATH:-}"
  local host_sys_path="${HOST_SYS_PATH:-}"

  if [[ -n "${volumes_path}" ]]; then
    export VOLUMES_PATH
    VOLUMES_PATH="$(to_compose_host_path "${volumes_path}")"
  fi

  if [[ -z "${machine_id_path}" || "${machine_id_path}" == "/etc" || "${machine_id_path}" == "/etc/machine-id" ]]; then
    if [[ -z "${volumes_path}" ]]; then
      volumes_path="${HOME}/volumes/tier0/data"
    fi
    machine_id_file="$(rewrite_windows_host_path "${volumes_path}")/machine-id"
    mkdir -p "$(dirname "${machine_id_file}")"
    if [[ ! -f "${machine_id_file}" ]]; then
      if [[ -f "/etc/machine-id" ]]; then
        cp "/etc/machine-id" "${machine_id_file}" 2>/dev/null || printf '%s\n' "$(hostname)-$(date +%s)" > "${machine_id_file}"
      else
        printf '%s\n' "$(hostname)-$(date +%s)" > "${machine_id_file}"
      fi
    fi
    export MACHINE_ID_HOST_PATH
    MACHINE_ID_HOST_PATH="$(to_compose_host_path "${machine_id_file}")"
  else
    export MACHINE_ID_HOST_PATH
    MACHINE_ID_HOST_PATH="$(to_compose_host_path "${machine_id_path}")"
  fi

  if [[ -z "${volumes_path}" ]]; then
    volumes_path="${HOME}/volumes/tier0/data"
  fi
  host_root="$(rewrite_windows_host_path "${volumes_path}")/host-metrics"
  if [[ -z "${host_proc_path}" || "${host_proc_path}" == "/proc" ]]; then
    mkdir -p "${host_root}/proc"
    HOST_PROC_PATH="$(to_compose_host_path "${host_root}/proc")"
  else
    HOST_PROC_PATH="$(to_compose_host_path "${host_proc_path}")"
  fi
  export HOST_PROC_PATH

  if [[ -z "${host_sys_path}" || "${host_sys_path}" == "/sys" ]]; then
    mkdir -p "${host_root}/sys"
    HOST_SYS_PATH="$(to_compose_host_path "${host_root}/sys")"
  else
    HOST_SYS_PATH="$(to_compose_host_path "${host_sys_path}")"
  fi
  export HOST_SYS_PATH
}

if is_windows_shell || is_wsl_shell; then
  normalize_windows_compose_paths
fi

if [[ -z "${COMPOSE_PROJECT_NAME:-}" ]]; then
  COMPOSE_PROJECT_NAME="edge"
  export COMPOSE_PROJECT_NAME
fi

cd "${DEPLOY_ROOT}"

# Stale docker-compose.override.yml files may carry a hardcoded backend image
# (e.g. from an earlier manual/branch deployment). Since the backend image must
# be controlled by .env / CI, strip any image line from the backend block of
# the override file before Compose merges it.
strip_backend_image_from_override() {
  local file="${1:-docker-compose.override.yml}"
  [[ -f "$file" ]] || return 0
  awk '
    /^  [a-zA-Z]/ { in_backend = ($0 ~ /^  backend:/) }
    in_backend && /^    image:/ { next }
    { print }
  ' "$file" > "${file}.tmp.$$" && mv "${file}.tmp.$$" "$file"
}
strip_backend_image_from_override

echo "[compose] BACKEND_IMAGE=${BACKEND_IMAGE:-<unset>}" >&2
compose_args=(docker compose --env-file "${ENV_FILE}")
if [[ -f "${RUNTIME_ENV_FILE}" ]]; then
  compose_args+=(--env-file "${RUNTIME_ENV_FILE}")
fi
compose_args+=(--project-name "${COMPOSE_PROJECT_NAME}")
"${compose_args[@]}" "$@"
