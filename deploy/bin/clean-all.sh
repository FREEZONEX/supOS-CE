#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${DEPLOY_ROOT}/.env"
RUNTIME_ENV_FILE="${DEPLOY_ROOT}/.env.runtime"

FORCE="false"
SKIP_UNINSTALL="false"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"

usage() {
  cat >&2 <<'USAGE'
usage: clean-all.sh [-f|--force] [--skip-uninstall]

Stops the current deploy compose stack and removes data under VOLUMES_PATH.
It also removes Docker named volumes declared by the current compose project.
After a full cleanup, it removes the generated .env.runtime for a fresh reinstall.
It does not remove Docker itself and does not touch other deploy roots.
USAGE
}

info() {
  printf '[INFO] %s\n' "$*"
}

warn() {
  printf '[WARN] %s\n' "$*" >&2
}

error() {
  printf '[ERROR] %s\n' "$*" >&2
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -f|--force)
      FORCE="true"
      shift
      ;;
    --skip-uninstall)
      SKIP_UNINSTALL="true"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      exit 2
      ;;
  esac
done

if [[ ! -f "${ENV_FILE}" ]]; then
  error "missing ${ENV_FILE}; refusing to infer a default data path"
  exit 1
fi

source_env_pair "${ENV_FILE}" "${RUNTIME_ENV_FILE}"

if [[ -z "${VOLUMES_PATH:-}" ]]; then
  error "VOLUMES_PATH is empty"
  exit 1
fi

is_absolute_path() {
  case "$1" in
    /*|[A-Za-z]:/*|[A-Za-z]:\\*) return 0 ;;
    *) return 1 ;;
  esac
}

is_windows_shell() {
  case "$(uname -s 2>/dev/null || true)" in
    MINGW*|MSYS*|CYGWIN*) return 0 ;;
  esac
  case "${OSTYPE:-}" in
    msys*|cygwin*) return 0 ;;
  esac
  return 1
}

is_wsl_shell() {
  [[ -r /proc/version ]] && grep -qiE 'microsoft|wsl' /proc/version
}

normalize_host_path() {
  local path="$1"
  local drive rest
  path="${path//$'\r'/}"

  if command -v cygpath >/dev/null 2>&1 && [[ "${path}" =~ ^[A-Za-z]:[\\/].*$ ]]; then
    cygpath -u "${path}"
    return
  fi

  if is_windows_shell && [[ "${path}" =~ ^/mnt/([A-Za-z])/(.*)$ ]]; then
    drive="$(printf '%s' "${BASH_REMATCH[1]}" | tr '[:upper:]' '[:lower:]')"
    rest="${BASH_REMATCH[2]}"
    if [[ -e "/${drive}/${rest}" ]]; then
      printf '/%s/%s\n' "${drive}" "${rest}"
      return
    fi
  fi

  if is_wsl_shell && [[ "${path}" =~ ^/([A-Za-z])/(.*)$ ]]; then
    drive="$(printf '%s' "${BASH_REMATCH[1]}" | tr '[:upper:]' '[:lower:]')"
    rest="${BASH_REMATCH[2]}"
    if [[ -e "/mnt/${drive}/${rest}" ]]; then
      printf '/mnt/%s/%s\n' "${drive}" "${rest}"
      return
    fi
  fi

  printf '%s\n' "${path}"
}

docker_mount_path() {
  local path="$1"
  local docker_bin

  if command -v cygpath >/dev/null 2>&1; then
    cygpath -am "${path}"
    return
  fi

  docker_bin="$(command -v docker 2>/dev/null || true)"
  if is_wsl_shell && [[ "${docker_bin}" == *".exe" ]] && command -v wslpath >/dev/null 2>&1; then
    wslpath -w "${path}"
    return
  fi

  printf '%s\n' "${path}"
}

resolve_path() {
  local path="$1"
  local dir base
  path="$(normalize_host_path "${path}")"
  if ! is_absolute_path "${path}"; then
    path="${DEPLOY_ROOT}/${path}"
  fi

  if [[ -d "${path}" ]]; then
    (cd "${path}" && pwd -P)
    return
  fi

  dir="$(dirname "${path}")"
  base="$(basename "${path}")"
  if [[ -d "${dir}" ]]; then
    printf '%s/%s\n' "$(cd "${dir}" && pwd -P)" "${base}"
    return
  fi

  printf '%s\n' "${path}"
}

has_known_layout() {
  local root="$1"
  local name
  for name in backend tsdb redis emqx sourceflow eventflow certs; do
    if [[ -e "${root}/${name}" ]]; then
      return 0
    fi
  done
  return 1
}

has_only_ignorable_contents() {
  local root="$1"
  [[ -z "$(find "${root}" -mindepth 1 -maxdepth 1 \
    ! -name '.DS_Store' \
    ! -name '.localized' \
    ! -name '.fseventsd' \
    ! -name '.Spotlight-V100' \
    ! -name '.TemporaryItems' \
    ! -name '.Trashes' \
    -print -quit 2>/dev/null || true)" ]]
}

volume_has_contents() {
  local root="$1"
  [[ -n "$(find "${root}" -mindepth 1 -maxdepth 1 -print -quit 2>/dev/null || true)" ]]
}

cleanup_helper_image() {
  local image
  for image in \
    "redis:7-alpine" \
    "timescale/timescaledb:2.20.0-pg17"; do
    if docker image inspect "${image}" >/dev/null 2>&1; then
      printf '%s\n' "${image}"
      return 0
    fi
  done
  return 1
}

remove_volume_contents_with_host() {
  local root="$1"
  chmod -R u+rwX "${root}" 2>/dev/null || true
  find "${root}" -mindepth 1 -maxdepth 1 -exec rm -rf {} +
}

remove_volume_contents_with_docker() {
  local root="$1"
  local image mount_path

  if ! command -v docker >/dev/null 2>&1; then
    return 1
  fi
  if ! image="$(cleanup_helper_image)"; then
    return 1
  fi

  mount_path="$(docker_mount_path "${root}")"
  docker run --rm --user 0:0 \
    -v "${mount_path}:/target" \
    --entrypoint sh \
    "${image}" \
    -c 'set -eu; find /target -mindepth 1 -maxdepth 1 -exec rm -rf {} +'
}

remove_volume_contents() {
  local root="$1"
  info "removing persisted deployment data under: ${root}"

  if ! remove_volume_contents_with_host "${root}" || volume_has_contents "${root}"; then
    warn "host cleanup did not remove all files; retrying with a Docker root helper"
    if ! remove_volume_contents_with_docker "${root}"; then
      error "cleanup incomplete; remaining files may be owned by Docker. Try running from the shell that owns VOLUMES_PATH, or remove the directory with elevated permissions: ${root}"
      return 1
    fi
  fi

  if volume_has_contents "${root}"; then
    error "cleanup incomplete; VOLUMES_PATH still contains data: ${root}"
    find "${root}" -mindepth 1 -maxdepth 1 -print | head -n 20 >&2 || true
    return 1
  fi
}

remove_runtime_state() {
  if [[ "${SKIP_UNINSTALL}" == "true" ]]; then
    warn "keeping ${RUNTIME_ENV_FILE} because --skip-uninstall leaves compose containers and named volumes in place"
    return
  fi

  if [[ -f "${RUNTIME_ENV_FILE}" ]]; then
    rm -f "${RUNTIME_ENV_FILE}"
    info "removed generated runtime env: ${RUNTIME_ENV_FILE}"
  fi
  rm -f "${DEPLOY_ROOT}/.env.tmp"
}

assert_safe_volume_root() {
  local root="$1"
  local repo_root
  repo_root="$(cd "${DEPLOY_ROOT}/.." && pwd -P)"

  case "${root}" in
    ""|"/"|"/mnt"|"/volumes"|"/Volumes"|"/home"|"/Users"|"/tmp"|"/var"|"/usr"|"/opt")
      error "refusing to purge unsafe VOLUMES_PATH: ${root}"
      exit 1
      ;;
  esac

  if [[ "${root}" == "${DEPLOY_ROOT}" || "${root}" == "${repo_root}" ]]; then
    error "refusing to purge workspace path: ${root}"
    exit 1
  fi

  if [[ ! -d "${root}" ]]; then
    info "no persisted deployment data found at: ${root}"
    exit 0
  fi

  local lock_file="${root}/.deploy-lock"
  if [[ -f "${lock_file}" ]]; then
    local locked_project current_project
    locked_project="$(sed -n 's/^COMPOSE_PROJECT_NAME=//p' "${lock_file}" | head -n 1)"
    locked_project="${locked_project//$'\r'/}"
    current_project="${COMPOSE_PROJECT_NAME:-}"
    current_project="${current_project//$'\r'/}"
    if [[ -n "${locked_project}" && -n "${current_project}" && "${locked_project}" != "${current_project}" ]]; then
      error "VOLUMES_PATH is locked by '${locked_project}', not '${current_project}'"
      exit 1
    fi
    return 0
  fi

  if ! has_known_layout "${root}"; then
    if has_only_ignorable_contents "${root}"; then
      info "only empty/macOS metadata content found; nothing to clean"
      exit 0
    fi
    error "'${root}' does not look like a deploy data directory"
    error "expected ${root}/.deploy-lock or one of: backend tsdb redis emqx sourceflow eventflow certs"
    error "current contents:"
    find "${root}" -mindepth 1 -maxdepth 1 -print 2>/dev/null | head -n 20 >&2 || true
    exit 1
  fi
}

VOLUMES_ABS="$(resolve_path "${VOLUMES_PATH}")"
assert_safe_volume_root "${VOLUMES_ABS}"

if [[ "${FORCE}" != "true" ]]; then
  echo
  if [[ "${SKIP_UNINSTALL}" == "true" ]]; then
    warn "This operation will delete data under:"
    warn "--skip-uninstall is set; compose containers and named volumes will not be removed."
  else
    warn "This operation will stop the current deploy, remove compose named volumes, and delete data under:"
  fi
  warn "${VOLUMES_ABS}"
  echo
  read -r -p "Continue? (y/N) " reply
  case "${reply}" in
    y|Y|yes|YES) ;;
    *)
      echo "cleanup cancelled"
      exit 1
      ;;
  esac
fi

if [[ "${SKIP_UNINSTALL}" != "true" ]]; then
  warn "stopping current deploy and removing compose named volumes..."
  (cd "${DEPLOY_ROOT}" && bash "${SCRIPT_DIR}/compose.sh" down --volumes --remove-orphans)
  rm -f "${DEPLOY_ROOT}/.env.tmp"
fi

remove_volume_contents "${VOLUMES_ABS}"
remove_runtime_state
warn "cleanup completed"
