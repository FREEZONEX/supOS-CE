#!/usr/bin/env bash

# Runtime env files contain credentials. Disable inherited xtrace before any
# runtime state is loaded so install diagnostics cannot echo secrets.
if [[ "$-" == *x* ]]; then
  set +x
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${DEPLOY_ROOT}/.env"
RUNTIME_ENV_FILE="${DEPLOY_ROOT}/.env.runtime"
DEFAULT_ENV_FILE="${DEPLOY_ROOT}/.env.default"
DEFAULT_PRODUCT_VERSION="$(
  awk -F= '
    /^[[:space:]]*PRODUCT_VERSION=/ {
      value = substr($0, index($0, "=") + 1)
      sub(/\r$/, "", value)
      print value
      exit
    }
  ' "${DEFAULT_ENV_FILE}"
)"
if [[ ! "${DEFAULT_PRODUCT_VERSION}" =~ ^[0-9A-Za-z][0-9A-Za-z._+-]*$ ]]; then
  printf 'error: invalid PRODUCT_VERSION in %s: %s\n' "${DEFAULT_ENV_FILE}" "${DEFAULT_PRODUCT_VERSION}" >&2
  return 1 2>/dev/null || exit 1
fi
CURRENT_INSTALL_SCHEMA_VERSION=3
INSTALL_SCHEMA_VERSION="${CURRENT_INSTALL_SCHEMA_VERSION}"
BACKEND_IMAGE_REPOSITORY="harbor.tier0.dev/tier0/tier0-edge-backend"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"

# Script-maintainer switch only. Keep this out of .env to avoid making package
# install timing part of the user-facing deployment contract.
NODERED_PACKAGE_INSTALL_MODE_HARDCODED="async"
NODERED_PACKAGE_INSTALL_MODE="${NODERED_PACKAGE_INSTALL_MODE:-${NODERED_PACKAGE_INSTALL_MODE_HARDCODED}}"

NO_START="${NO_START:-false}"
ADOPT_VOLUMES="${ADOPT_VOLUMES:-false}"
TLS_MODE="${TLS_MODE:-}"
VOLUMES_PATH_ARG="${VOLUMES_PATH_ARG:-}"
SKIP_NODERED_PACKAGES="${SKIP_NODERED_PACKAGES:-}"
LOCAL_FRONTEND_DEV_ARG="${LOCAL_FRONTEND_DEV_ARG:-}"


random_hex() {
  local bytes="${1:-16}"
  local chars=$((bytes * 2))
  local seed hash
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex "${bytes}"
  elif command -v sha256sum >/dev/null 2>&1; then
    seed="${RANDOM}${RANDOM}$(date +%s)$$"
    hash="$(printf '%s' "${seed}" | sha256sum | awk '{print $1}')"
    printf '%s\n' "${hash:0:${chars}}"
  elif command -v shasum >/dev/null 2>&1; then
    seed="${RANDOM}${RANDOM}$(date +%s)$$"
    hash="$(printf '%s' "${seed}" | shasum -a 256 | awk '{print $1}')"
    printf '%s\n' "${hash:0:${chars}}"
  else
    hash=""
    while [[ "${#hash}" -lt "${chars}" ]]; do
      hash="${hash}${RANDOM}"
    done
    printf '%s\n' "${hash:0:${chars}}"
  fi
}

is_windows_host_shell() {
  [[ "${OSTYPE:-}" == msys* || "${OSTYPE:-}" == cygwin* ]] || command -v cmd.exe >/dev/null 2>&1
}

is_non_linux_host_shell() {
  is_windows_host_shell || [[ "${OSTYPE:-}" == darwin* ]]
}

default_volumes_path() {
  if is_windows_host_shell; then
    printf '%s\n' "${HOME}/volumes/tier0/data"
  elif [[ "${OSTYPE:-}" == darwin* ]]; then
    printf '%s\n' "${HOME}/volumes/tier0/data"
  else
    printf '%s\n' "/volumes/tier0/data"
  fi
}

is_absolute_path() {
  case "$1" in
    /*|[A-Za-z]:/*|[A-Za-z]:\\*) return 0 ;;
    *) return 1 ;;
  esac
}

absolute_path_for() {
  local path="$1"
  if is_windows_shell || is_wsl_shell; then
    path="$(rewrite_windows_host_path "${path}")"
  fi
  if ! is_absolute_path "${path}"; then
    path="${DEPLOY_ROOT}/${path}"
  fi
  mkdir -p "${path}"
  (cd "${path}" && pwd -P)
}

set_env() {
  local key="$1"
  local value="$2"
  local file="${3:-${ENV_FILE}}"
  local tmp
  mkdir -p "$(dirname "${file}")"
  touch "${file}"
  tmp="${file}.tmp.$$"
  awk -v k="${key}" -v v="${value}" '
    BEGIN { done = 0 }
    $0 ~ "^" k "=" { print k "=" v; done = 1; next }
    { print }
    END { if (done == 0) print k "=" v }
  ' "${file}" > "${tmp}"
  mv "${tmp}" "${file}"
}


ensure_env_template() {
  if [[ -f "${ENV_FILE}" ]]; then
    return
  fi
  if [[ -f "${DEPLOY_ROOT}/.env.default" ]]; then
    cp "${DEPLOY_ROOT}/.env.default" "${ENV_FILE}"
  else
    error "deploy env template not found: ${DEPLOY_ROOT}/.env.default"
  fi
}

load_env() {
  source_env_file "${ENV_FILE}"
}

is_local_tsdb_url() {
  local url="$1"
  [[ "${url}" == *"@tsdb:"* || "${url}" == *"//tsdb:"* ]]
}

is_legacy_postgresql_url() {
  local url="$1"
  [[ "${url}" == *"@postgresql:"* || "${url}" == *"//postgresql:"* ]]
}

is_generated_local_db_url() {
  local url="${1:-}"
  [[ -z "${url}" ]] || is_local_tsdb_url "${url}" || is_legacy_postgresql_url "${url}"
}

runtime_env_value() {
  local key="$1"
  [[ -f "${RUNTIME_ENV_FILE}" ]] || return 1
  awk -v k="${key}" '
    index($0, k "=") == 1 {
      sub("^[^=]*=", "")
      sub(/\r$/, "")
      if ((substr($0, 1, 1) == "\047" && substr($0, length($0), 1) == "\047") ||
          (substr($0, 1, 1) == "\042" && substr($0, length($0), 1) == "\042")) {
        $0 = substr($0, 2, length($0) - 2)
      }
      print
      found = 1
      exit
    }
    END { exit found ? 0 : 1 }
  ' "${RUNTIME_ENV_FILE}"
}


env_file_value() {
  local file="$1"
  local key="$2"
  [[ -f "${file}" ]] || return 1
  awk -v k="${key}" '
    index($0, k "=") == 1 {
      sub("^[^=]*=", "")
      sub(/\r$/, "")
      print
      found = 1
      exit
    }
    END { exit found ? 0 : 1 }
  ' "${file}"
}


load_runtime_value() {
  local key="$1"
  local value
  if value="$(runtime_env_value "${key}")" && [[ -n "${value}" ]]; then
    printf -v "${key}" '%s' "${value}"
    export "${key}"
  fi
}

load_runtime_default() {
  local key="$1"
  local value
  if [[ -n "${!key:-}" ]]; then
    return
  fi
  if value="$(runtime_env_value "${key}")" && [[ -n "${value}" ]]; then
    printf -v "${key}" '%s' "${value}"
    export "${key}"
  fi
}

load_runtime_generated_db_url_default() {
  local key="$1"
  local value
  if ! is_generated_local_db_url "${!key:-}"; then
    return
  fi
  if value="$(runtime_env_value "${key}")" && [[ -n "${value}" ]] && is_generated_local_db_url "${value}"; then
    printf -v "${key}" '%s' "${value}"
    export "${key}"
  fi
}

load_runtime_defaults() {
  load_runtime_value TSDB_PASSWORD
  load_runtime_value REDIS_PASSWORD
  load_runtime_value JWT_SECRET
  load_runtime_default ADMIN_INITIAL_PASSWORD
  load_runtime_default EMQX_API_KEY
  load_runtime_default EMQX_API_SECRET
  load_runtime_default EMQX_DASHBOARD_USERNAME
  load_runtime_default EMQX_DASHBOARD_PASSWORD
  load_runtime_default NODERED_INTERNAL_TOKEN
  load_runtime_default TIER0_API_KEY
  load_runtime_default TIER0_INSTALLATION_ID
  load_runtime_default COMPOSE_PROJECT_NAME
  load_runtime_generated_db_url_default UNS_DB_URL
  load_runtime_generated_db_url_default SINK_DB_URL
}

resolve_tier0_installation_id() {
  if [[ -z "${TIER0_INSTALLATION_ID:-}" ]]; then
    TIER0_INSTALLATION_ID="$(random_hex 16)"
  fi
  if [[ ! "${TIER0_INSTALLATION_ID}" =~ ^[A-Za-z0-9][A-Za-z0-9_.-]{7,127}$ ]]; then
    error "TIER0_INSTALLATION_ID must be 8-128 characters using only letters, numbers, dot, underscore, or hyphen"
  fi
  export TIER0_INSTALLATION_ID
}

restore_runtime_env_from_volume() {
  local backup_file="${VOLUMES_ABS}/.env.runtime"
  local backup_project
  if [[ -f "${RUNTIME_ENV_FILE}" || ! -f "${backup_file}" ]]; then
    return
  fi
  backup_project="$(env_file_value "${backup_file}" COMPOSE_PROJECT_NAME || true)"
  if [[ -n "${backup_project}" && "${backup_project}" != "${COMPOSE_PROJECT_NAME}" ]]; then
    warn "ignoring runtime environment from another Compose project: ${backup_project}"
    return
  fi
  cp "${backup_file}" "${RUNTIME_ENV_FILE}"
  chmod 600 "${RUNTIME_ENV_FILE}" 2>/dev/null || true
  install_detail "Restored runtime environment: ${backup_file}"
  load_runtime_defaults
}

error() {
  printf '[ERROR] %s\n' "$*" >&2
  exit 1
}

warn() {
  printf '[WARN] %s\n' "$*" >&2
}

install_step() {
  local index="$1"
  local total="$2"
  local message="$3"
  printf '\n==> [%s/%s] %s\n' "${index}" "${total}" "${message}"
}

install_detail() {
  printf '    - %s\n' "$*"
}

is_true() {
  case "${1:-}" in
    true|TRUE|1|yes|YES|on|ON) return 0 ;;
    *) return 1 ;;
  esac
}

append_compose_profile() {
  local profile="$1"
  if [[ -n "${COMPOSE_PROFILES:-}" ]]; then
    COMPOSE_PROFILES="${COMPOSE_PROFILES},${profile}"
  else
    COMPOSE_PROFILES="${profile}"
  fi
}


frontend_dev_bash_command() {
  printf 'cd frontend/apps/web && API_PROXY_URL=%s VITE_ASSET_PREFIX=%s VITE_SKIP_LICENSE_CHECKS=true pnpm dev --host 0.0.0.0 --port %s' \
    "${ENTRANCE_URL}" "${ENTRANCE_URL}" "${FRONTEND_DEV_PORT}"
}

frontend_dev_powershell_command() {
  printf "cd frontend/apps/web; \$env:API_PROXY_URL='%s'; \$env:VITE_ASSET_PREFIX='%s'; \$env:VITE_SKIP_LICENSE_CHECKS='true'; pnpm dev --host 0.0.0.0 --port %s" \
    "${ENTRANCE_URL}" "${ENTRANCE_URL}" "${FRONTEND_DEV_PORT}"
}

platform_url() {
  printf '%s/uns\n' "${ENTRANCE_URL%/}"
}


backend_ready_url() {
  printf '%s/readyz\n' "${ENTRANCE_URL%/}"
}

http_ready() {
  local url="$1"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSk --max-time 5 "${url}" >/dev/null 2>&1
    return
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -q --no-check-certificate -T 5 -O /dev/null "${url}" >/dev/null 2>&1
    return
  fi
  echo "curl or wget is required for backend health checks" >&2
  return 2
}

backend_container_http_ready() {
  local path="$1"
  local url="http://127.0.0.1:8080${path}"
  bash "${DEPLOY_ROOT}/bin/compose.sh" exec -T backend sh -c '
url="$1"
if command -v curl >/dev/null 2>&1; then
  curl -fsS --max-time 5 "$url" >/dev/null
elif command -v wget >/dev/null 2>&1; then
  wget -q -T 5 -O /dev/null "$url"
else
  exit 127
fi
' sh "${url}" >/dev/null 2>&1
}

backend_ready() {
  local external_url="$1"
  backend_container_http_ready "/readyz" || http_ready "${external_url}"
}


# 加载 build-images.sh 打包的 Node 运行时镜像，供内网/离线环境应用导入使用。


verify_bundled_file_checksum() {
  local source="$1"
  local checksum_file="${DEPLOY_ROOT}/checksums.sha256"
  local relative expected actual
  [[ -f "${checksum_file}" ]] || return 0
  relative="${source#${DEPLOY_ROOT}/}"
  [[ "${relative}" != "${source}" ]] || error "bundled file is outside the deployment root: ${source}"
  expected="$(awk -v path="${relative}" '{ gsub(/\r$/, "", $2); if ($2 == path) { print $1; exit } }' "${checksum_file}")"
  [[ -n "${expected}" ]] || error "checksum is missing for bundled file: ${relative}"
  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "${source}" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    actual="$(shasum -a 256 "${source}" | awk '{print $1}')"
  else
    error "sha256sum or shasum is required to verify bundled files"
  fi
  [[ "${actual}" == "${expected}" ]] || error "checksum verification failed for bundled file: ${relative}"
}

load_bundled_image_archive() {
  local image="$1"
  local archive_path="$2"

  install_detail "Loading bundled image: ${image}"
  verify_bundled_file_checksum "${archive_path}"
  case "${archive_path}" in
    *.tar.gz)
      gzip -dc "${archive_path}" | docker load
      ;;
    *.tar)
      docker load < "${archive_path}"
      ;;
    *)
      error "unsupported bundled image archive: ${archive_path}"
      ;;
  esac
  docker image inspect "${image}" >/dev/null 2>&1 ||
    error "bundled archive did not provide the expected image: ${image}"
  install_detail "Loaded bundled image: ${image}"
}

# A complete offline package includes images/manifest.txt with one
# "<image><TAB><archive>" entry per image. Archives are loaded before Compose
# starts so Docker never needs registry access on the target machine.
load_bundled_images() {
  local images_dir="${DEPLOY_ROOT}/images"
  local manifest="${images_dir}/manifest.txt"
  local image archive_name archive_path loaded_archives="|"
  [[ -f "${manifest}" ]] || return 0
  command -v docker >/dev/null 2>&1 || error "Docker is required to load bundled image archives"
  verify_bundled_file_checksum "${manifest}"
  while IFS=$'\t' read -r image archive_name || [[ -n "${image}${archive_name}" ]]; do
    image="${image%$'\r'}"
    archive_name="${archive_name%$'\r'}"
    [[ -n "${image}" ]] || continue
    [[ "${image}" != \#* ]] || continue
    [[ -n "${archive_name}" ]] || error "missing archive name for bundled image: ${image}"
    [[ "${archive_name}" == "$(basename "${archive_name}")" ]] || error "invalid bundled image archive path: ${archive_name}"
    archive_path="${images_dir}/${archive_name}"
    [[ -f "${archive_path}" ]] || error "bundled image archive is missing: ${archive_path}"
    if [[ "${loaded_archives}" == *"|${archive_path}|"* ]]; then
      docker image inspect "${image}" >/dev/null 2>&1 || error "bundled archive did not provide image: ${image}"
      continue
    fi
    load_bundled_image_archive "${image}" "${archive_path}"
    loaded_archives="${loaded_archives}${archive_path}|"
  done < "${manifest}"
}

print_backend_diagnostics() {
  echo "backend health check diagnostics:" >&2
  bash "${DEPLOY_ROOT}/bin/compose.sh" ps backend >&2 || true
  bash "${DEPLOY_ROOT}/bin/compose.sh" logs --tail=80 backend >&2 || true
}

wait_backend_ready() {
  local url elapsed interval max_wait
  url="$(backend_ready_url)"
  elapsed=0
  interval=3
  max_wait=180

  install_detail "Backend readiness: ${url}"
  while (( elapsed <= max_wait )); do
    if backend_ready "${url}"; then
      install_detail "Backend readiness: ok"
      return 0
    fi
    if (( elapsed > 0 && elapsed % 15 == 0 )); then
      install_detail "Waiting for backend readiness (${elapsed}s / ${max_wait}s)"
    fi
    sleep "${interval}"
    elapsed=$((elapsed + interval))
  done

  print_backend_diagnostics
  error "backend readiness check failed after ${max_wait}s: ${url}"
}

print_install_summary() {
  local package_log="${VOLUMES_ABS}/sourceflow/package-install.log"
  printf '\n============================================================\n'
  printf 'Tier0 Edge deployment %s.\n\n' "$([[ "${NO_START}" == "true" ]] && printf 'prepared' || printf 'started')"
  printf 'Open the platform: %s\n' "$(platform_url)"
  printf 'Default username: tier0\n'
  printf 'Default password: tier0\n'
  if is_true "${LOCAL_FRONTEND_DEV}"; then
    printf 'Frontend dev server: %s\n' "${FRONTEND_DEV_URL}"
    printf 'Frontend command (PowerShell): %s\n' "$(frontend_dev_powershell_command)"
    printf 'Frontend command (Git Bash): %s\n' "$(frontend_dev_bash_command)"
  fi
  if [[ "${NO_START}" != "true" && "${SKIP_NODERED_PACKAGES}" != "true" ]]; then
    printf 'SourceFlow package log: %s\n' "${package_log}"
  fi
  printf 'Directory: %s\nVolumes: %s\n' "${DEPLOY_ROOT}" "${VOLUMES_ABS}"
  printf '============================================================\n'
}

port_with_offset() {
  local base="$1"
  local offset="$2"
  echo $((base + offset))
}

resolve_port() {
  local key="$1"
  local base="$2"
  local offset="$3"
  local current="${!key:-}"
  if [[ -z "${current}" || ("${offset}" != "0" && "${current}" == "${base}") ]]; then
    port_with_offset "${base}" "${offset}"
  else
    printf '%s\n' "${current}"
  fi
}


build_entrance_base_url() {
  local protocol="$1"
  local domain="$2"
  local port="$3"
  local ssl_port="$4"
  local effective_port
  effective_port="${port}"
  if [[ "${protocol}" == "https" ]]; then
    effective_port="${ssl_port}"
  fi
  if [[ ("${protocol}" == "http" && "${effective_port}" == "80") || ("${protocol}" == "https" && "${effective_port}" == "443") ]]; then
    printf '%s://%s\n' "${protocol}" "${domain}"
  else
    printf '%s://%s:%s\n' "${protocol}" "${domain}" "${effective_port}"
  fi
}

ensure_self_signed_cert() {
  local dir="$1"
  local cfg
  local host
  local san_lines
  mkdir -p "${dir}"
  if ! command -v openssl >/dev/null 2>&1; then
    echo "openssl is required to generate TLS certs; provide ${dir}/cert.pem and ${dir}/key.pem instead" >&2
    exit 3
  fi
  host="${ENTRANCE_DOMAIN:-127.0.0.1}"
  san_lines=$(
    cat <<SAN
DNS.1 = localhost
DNS.2 = emqx
IP.1 = 127.0.0.1
SAN
  )
  if [[ -n "${host}" && "${host}" != "localhost" && "${host}" != "127.0.0.1" ]]; then
    if [[ "${host}" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      san_lines="${san_lines}"$'\n'"IP.2 = ${host}"
    else
      san_lines="${san_lines}"$'\n'"DNS.3 = ${host}"
    fi
  fi
  if [[ -f "${dir}/cert.pem" && -f "${dir}/key.pem" ]]; then
    if openssl x509 -in "${dir}/cert.pem" -noout -ext subjectAltName 2>/dev/null | grep -Fq "${host}"; then
      return
    fi
    install_detail "Regenerating TLS cert for ${host}"
  fi
  cfg="${dir}/openssl-san.cnf.$$"
  cat > "${cfg}" <<CONF
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
CN = localhost

[v3_req]
subjectAltName = @alt_names

[alt_names]
${san_lines}
CONF
  if ! openssl req -x509 -nodes -newkey rsa:2048 \
    -keyout "${dir}/key.pem" \
    -out "${dir}/cert.pem" \
    -days 3650 \
    -config "${cfg}" \
    -extensions v3_req >/dev/null 2>&1; then
    rm -f "${cfg}"
    echo "failed to generate TLS certs; provide ${dir}/cert.pem and ${dir}/key.pem instead" >&2
    exit 3
  fi
  rm -f "${cfg}"
  chmod 600 "${dir}/key.pem" 2>/dev/null || true
}

ensure_provided_cert() {
  local dir="$1"
  if [[ ! -f "${dir}/cert.pem" || ! -f "${dir}/key.pem" ]]; then
    echo "provided TLS mode requires ${dir}/cert.pem and ${dir}/key.pem" >&2
    exit 3
  fi
}

write_deploy_lock() {
  local volumes_abs="$1"
  local lock_file="${volumes_abs}/.deploy-lock"
  if [[ -f "${lock_file}" ]]; then
    local locked_project
    locked_project="$(sed -n 's/^COMPOSE_PROJECT_NAME=//p' "${lock_file}" | head -n 1)"
    locked_project="${locked_project//$'\r'/}"
    if [[ -n "${locked_project}" && "${locked_project}" != "${COMPOSE_PROJECT_NAME}" && "${ADOPT_VOLUMES}" != "true" ]]; then
      echo "volumes path is locked by compose project '${locked_project}': ${volumes_abs}" >&2
      echo "rerun with --adopt-volumes only if this deploy should take over that data directory" >&2
      exit 4
    fi
  fi
  cat > "${lock_file}" <<LOCK
COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME}
OS_NAME=${OS_NAME}
TIER0_INSTALLATION_ID=${TIER0_INSTALLATION_ID}
LOCK
}

prepare_volume_dirs() {
  local root="$1"
  mkdir -p \
    "${root}/backend/files" "${root}/backend/logs" \
    "${root}/certs/tls" "${root}/tsdb/conf" "${root}/tsdb/data" \
    "${root}/redis/data" "${root}/emqx/config" "${root}/emqx/data" "${root}/emqx/log" \
    "${root}/sourceflow" "${root}/eventflow" \
    "${root}/host-metrics/proc" "${root}/host-metrics/sys"
  if [[ "$(uname -s)" == "Linux" ]]; then
    chown -R 1000:1000 "${root}/emqx/data" "${root}/emqx/log" 2>/dev/null || \
      chmod -R a+rwX "${root}/emqx/data" "${root}/emqx/log" || error "failed to prepare EMQX data directories"
  fi
  local tsdb_conf="${root}/tsdb/conf/postgresql.conf"
  [[ ! -d "${tsdb_conf}" ]] || error "expected TSDB config file but found directory: ${tsdb_conf}"
  if [[ -f "${DEPLOY_ROOT}/mount/tsdb/conf/postgresql.conf" && ! -f "${tsdb_conf}" ]]; then
    cp "${DEPLOY_ROOT}/mount/tsdb/conf/postgresql.conf" "${tsdb_conf}"
  fi
  [[ ! -f "${tsdb_conf}" ]] || chmod a+r "${tsdb_conf}" || error "failed to make TSDB config readable"
}

ensure_host_metrics_paths() {
  HOST_METRICS_ENABLED="${HOST_METRICS_ENABLED:-false}"

  if is_windows_shell || is_wsl_shell || [[ "${OSTYPE:-}" == darwin* ]]; then
    # On Windows/WSL/macOS Docker runs inside a Linux VM, so /proc and /sys can
    # never expose real host metrics. Mount empty placeholder directories to keep
    # the compose file valid, and only honor HOST_METRICS_ENABLED=true when the
    # user supplied custom paths pointing at real metrics data.
    local has_custom_paths="true"
    if [[ -z "${HOST_PROC_PATH:-}" || "${HOST_PROC_PATH}" == "/proc" ]]; then
      HOST_PROC_PATH="${VOLUMES_ABS}/host-metrics/proc"
      has_custom_paths="false"
    fi
    if [[ -z "${HOST_SYS_PATH:-}" || "${HOST_SYS_PATH}" == "/sys" ]]; then
      HOST_SYS_PATH="${VOLUMES_ABS}/host-metrics/sys"
      has_custom_paths="false"
    fi
    if [[ "${HOST_METRICS_ENABLED}" == "true" && "${has_custom_paths}" != "true" ]]; then
      warn "HOST_METRICS_ENABLED=true has no effect on this platform: /proc and /sys only exist inside the Docker Linux VM. Disabling host metrics; set HOST_PROC_PATH/HOST_SYS_PATH to valid metrics paths to force-enable."
      HOST_METRICS_ENABLED="false"
    fi
    mkdir -p "${HOST_PROC_PATH}" "${HOST_SYS_PATH}"
  else
    HOST_PROC_PATH="${HOST_PROC_PATH:-/proc}"
    HOST_SYS_PATH="${HOST_SYS_PATH:-/sys}"
    # /proc and /sys already exist on Linux; only create custom override paths.
    if [[ ! -d "${HOST_PROC_PATH}" ]]; then
      mkdir -p "${HOST_PROC_PATH}"
    fi
    if [[ ! -d "${HOST_SYS_PATH}" ]]; then
      mkdir -p "${HOST_SYS_PATH}"
    fi
  fi

  export HOST_METRICS_ENABLED HOST_PROC_PATH HOST_SYS_PATH
}

configure_sourceflow_package_install() {
  local root="$1"
  local marker="${root}/sourceflow/.skip-package-install"
  if [[ "${SKIP_NODERED_PACKAGES:-}" == "true" ]]; then
    printf 'created by install.sh --skip-nodered-packages\n' > "${marker}"
  else
    rm -f "${marker}"
  fi
}

sourceflow_package_install_skipped() {
  local root="$1"
  [[ -f "${root}/sourceflow/.skip-package-install" ]]
}

validate_nodered_package_install_mode() {
  case "${NODERED_PACKAGE_INSTALL_MODE}" in
    async|sync) ;;
    *) error "unsupported NODERED_PACKAGE_INSTALL_MODE: ${NODERED_PACKAGE_INSTALL_MODE}" ;;
  esac
}

reset_nodered_package_install_mode() {
  NODERED_PACKAGE_INSTALL_MODE="${NODERED_PACKAGE_INSTALL_MODE_HARDCODED}"
}

start_sourceflow_package_install_async() {
  local log_file="${VOLUMES_ABS}/sourceflow/package-install.log"
  mkdir -p "${VOLUMES_ABS}/sourceflow"
  install_detail "Sourceflow package install: async (${log_file})"
  if command -v nohup >/dev/null 2>&1; then
    (cd "${DEPLOY_ROOT}" && nohup bash "${DEPLOY_ROOT}/bin/nodered-packages.sh" --background >/dev/null 2>&1 &)
  else
    (cd "${DEPLOY_ROOT}" && bash "${DEPLOY_ROOT}/bin/nodered-packages.sh" --background >/dev/null 2>&1 &)
  fi
}

run_sourceflow_package_install_sync() {
  install_detail "Sourceflow package install: sync"
  if ! bash "${DEPLOY_ROOT}/bin/nodered-packages.sh"; then
    echo "warning: failed to install sourceflow Node-RED packages; run bin/nodered-packages.sh to retry" >&2
  fi
}

run_post_start_initialization() {
  if sourceflow_package_install_skipped "${VOLUMES_ABS}"; then
    install_detail "Skipping SourceFlow package install"
    return
  fi
  case "${NODERED_PACKAGE_INSTALL_MODE}" in
    async) start_sourceflow_package_install_async ;;
    sync) run_sourceflow_package_install_sync ;;
  esac
}


render_compose_files() {
  cp "${DEPLOY_ROOT}/templates/docker-compose.yml.tpl" "${DEPLOY_ROOT}/docker-compose.yml"
  cp "${DEPLOY_ROOT}/templates/docker-compose.override.yml.tpl" "${DEPLOY_ROOT}/docker-compose.override.yml"
  if [[ "${LOCAL_DB_ENABLED:-true}" != "true" && "${LOCAL_REDIS_ENABLED:-true}" != "true" ]]; then
    return
  fi
  {
    echo "services:"
    echo "  backend:"
    echo "    depends_on:"
    if [[ "${LOCAL_DB_ENABLED:-true}" == "true" ]]; then
      printf '      tsdb:\n        condition: service_healthy\n'
    fi
    if [[ "${LOCAL_REDIS_ENABLED:-true}" == "true" ]]; then
      printf '      redis:\n        condition: service_healthy\n'
    fi
  } > "${DEPLOY_ROOT}/docker-compose.override.yml"
}


resolve_runtime_secret() {
  local key="$1"
  local bytes="${2:-24}"
  local value="${!key:-}"
  if [[ -z "${value}" ]]; then
    case "${key}" in
      EMQX_API_KEY) value="tier0-$(random_hex 4)" ;;
      *) value="$(random_hex "${bytes}")" ;;
    esac
    printf -v "${key}" '%s' "${value}"
  fi
  export "${key}"
}

resolve_internal_runtime_secrets() {
  resolve_runtime_secret EMQX_API_KEY 4
  resolve_runtime_secret EMQX_API_SECRET 24
  EMQX_DASHBOARD_USERNAME="${EMQX_DASHBOARD_USERNAME:-admin}"
  export EMQX_DASHBOARD_USERNAME
  resolve_runtime_secret EMQX_DASHBOARD_PASSWORD 24
  resolve_runtime_secret NODERED_INTERNAL_TOKEN 24
}


ensure_runtime_file_target() {
  local path="$1"
  local backup
  local stamp
  local suffix
  if [[ -d "${path}" && ! -L "${path}" ]]; then
    stamp="$(date +%Y%m%d%H%M%S)"
    backup="${path}.directory-backup.${stamp}"
    suffix=1
    while [[ -e "${backup}" ]]; do
      backup="${path}.directory-backup.${stamp}.${suffix}"
      suffix=$((suffix + 1))
    done
    mv "${path}" "${backup}" || error "failed to move directory at runtime file path: ${path}"
    warn "moved directory at runtime file path to backup: ${path} -> ${backup}"
  fi
  if [[ -e "${path}" && ! -f "${path}" ]]; then
    error "runtime file path is not a regular file: ${path}"
  fi
}

write_emqx_runtime_files() {
  local dir="${VOLUMES_ABS}/emqx/config"
  local api_key_file="${dir}/default_api_key.conf"
  mkdir -p "${dir}"
  ensure_runtime_file_target "${api_key_file}"
  printf '%s:%s:administrator\n' "${EMQX_API_KEY}" "${EMQX_API_SECRET}" > "${api_key_file}"
  chmod 600 "${api_key_file}" || error "failed to secure EMQX API credential file: ${api_key_file}"
}


write_runtime_env() {
  cat > "${RUNTIME_ENV_FILE}" <<ENV
# Tier0 Edge runtime config generated from .env. Edit .env and rerun bin/install.sh.
COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME}
TIER0_INSTALLATION_ID=${TIER0_INSTALLATION_ID}
LOCAL_FRONTEND_DEV=${LOCAL_FRONTEND_DEV}
FRONTEND_DEV_PORT=${FRONTEND_DEV_PORT}
FRONTEND_DEV_URL=${FRONTEND_DEV_URL}
FRONTEND_DEV_PROXY_URL=${FRONTEND_DEV_PROXY_URL}
VOLUMES_PATH="${VOLUMES_PATH}"
ENTRANCE_PROTOCOL=${ENTRANCE_PROTOCOL}
ENTRANCE_DOMAIN=${ENTRANCE_DOMAIN}
ENTRANCE_PORT=${ENTRANCE_PORT}
ENTRANCE_SSL_PORT=${ENTRANCE_SSL_PORT}
HOST_METRICS_ENABLED=${HOST_METRICS_ENABLED}
HOST_PROC_PATH="${HOST_PROC_PATH}"
HOST_SYS_PATH="${HOST_SYS_PATH}"
PRODUCT_VERSION=${PRODUCT_VERSION}
UNS_DB_URL=${UNS_DB_URL}
SINK_DB_URL=${SINK_DB_URL}
TIMESERIES_RETENTION_YEARS=${TIMESERIES_RETENTION_YEARS:-10}
COMPOSE_PROFILES=${COMPOSE_PROFILES}
TSDB_PORT=${TSDB_PORT}
OS_MQTT_TCP_PORT=${OS_MQTT_TCP_PORT}
OS_MQTT_SSL_PORT=${OS_MQTT_SSL_PORT}
OS_MQTT_WEBSOCKET_PORT=${OS_MQTT_WEBSOCKET_PORT}
OS_MQTT_WEBSOCKET_TSL_PORT=${OS_MQTT_WEBSOCKET_TSL_PORT}
TSDB_PASSWORD=${TSDB_PASSWORD}
REDIS_PASSWORD=${REDIS_PASSWORD}
REDIS_ADDR=${REDIS_ADDR}
JWT_SECRET=${JWT_SECRET}
TIER0_API_KEY=${TIER0_API_KEY}
ADMIN_INITIAL_PASSWORD=${ADMIN_INITIAL_PASSWORD}
EMQX_API_KEY=${EMQX_API_KEY}
EMQX_API_SECRET=${EMQX_API_SECRET}
EMQX_DASHBOARD_USERNAME=${EMQX_DASHBOARD_USERNAME}
EMQX_DASHBOARD_PASSWORD=${EMQX_DASHBOARD_PASSWORD}
NODERED_INTERNAL_TOKEN=${NODERED_INTERNAL_TOKEN}
ENV
  cp "${RUNTIME_ENV_FILE}" "${VOLUMES_ABS}/.env.runtime"
  chmod 600 "${RUNTIME_ENV_FILE}" "${VOLUMES_ABS}/.env.runtime" 2>/dev/null || true
}

write_install_state() {
  local services_started="$1"
  cat > "${VOLUMES_ABS}/.install-state.json" <<STATE
{
  "productVersion": "${PRODUCT_VERSION}",
  "installSchemaVersion": ${INSTALL_SCHEMA_VERSION},
  "backendImage": "${BACKEND_IMAGE}",
  "osName": "${OS_NAME}",
  "composeProjectName": "${COMPOSE_PROJECT_NAME}",
  "installationID": "${TIER0_INSTALLATION_ID}",
  "volumesPath": "${VOLUMES_PATH}",
  "entranceProtocol": "${ENTRANCE_PROTOCOL}",
  "entrancePort": ${ENTRANCE_PORT},
  "entranceSslPort": ${ENTRANCE_SSL_PORT},
  "servicesStarted": ${services_started}
}
STATE
  touch "${VOLUMES_ABS}/.install_complete"
  install_detail "State file: ${VOLUMES_ABS}/.install-state.json"
}

run_deploy_workflow() {
  install_step 1 7 "Loading Tier0 Edge deployment configuration"
  ensure_env_template
  load_env
  reset_nodered_package_install_mode
  validate_nodered_package_install_mode
  load_runtime_defaults
  install_detail "Config file: ${ENV_FILE}"

  install_step 2 7 "Resolving instance settings"
  OS_NAME="${OS_NAME:-Tier0 Edge}"
  PRODUCT_VERSION="${PRODUCT_VERSION:-${DEFAULT_PRODUCT_VERSION}}"
  LANGUAGE="${LANGUAGE:-en-US}"
  COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-edge}"
  BACKEND_IMAGE="${BACKEND_IMAGE:-${BACKEND_IMAGE_REPOSITORY}:${PRODUCT_VERSION}}"
  export PRODUCT_VERSION COMPOSE_PROJECT_NAME BACKEND_IMAGE
  if [[ -z "$(env_file_value "${ENV_FILE}" BACKEND_IMAGE || true)" ]]; then
    set_env BACKEND_IMAGE "${BACKEND_IMAGE}" "${ENV_FILE}"
  fi

  PORT_OFFSET="${PORT_OFFSET:-0}"
  DEFAULT_VOLUMES_PATH="$(default_volumes_path)"
  if [[ -n "${VOLUMES_PATH_ARG}" ]]; then
    VOLUMES_PATH="${VOLUMES_PATH_ARG}"
  elif [[ "${VOLUMES_PATH:-}" == "/volumes/tier0/data" ]] && is_non_linux_host_shell; then
    VOLUMES_PATH="${DEFAULT_VOLUMES_PATH}"
  else
    VOLUMES_PATH="${VOLUMES_PATH:-${DEFAULT_VOLUMES_PATH}}"
  fi
  ENTRANCE_PORT="$(resolve_port ENTRANCE_PORT 8088 "${PORT_OFFSET}")"
  ENTRANCE_SSL_PORT="$(resolve_port ENTRANCE_SSL_PORT 8443 "${PORT_OFFSET}")"
  TSDB_PORT="$(resolve_port TSDB_PORT 5433 "${PORT_OFFSET}")"
  OS_MQTT_TCP_PORT="$(resolve_port OS_MQTT_TCP_PORT 1883 "${PORT_OFFSET}")"
  OS_MQTT_SSL_PORT="$(resolve_port OS_MQTT_SSL_PORT 8883 "${PORT_OFFSET}")"
  OS_MQTT_WEBSOCKET_PORT="$(resolve_port OS_MQTT_WEBSOCKET_PORT 8083 "${PORT_OFFSET}")"
  OS_MQTT_WEBSOCKET_TSL_PORT="$(resolve_port OS_MQTT_WEBSOCKET_TSL_PORT 8084 "${PORT_OFFSET}")"
  ENTRANCE_PROTOCOL="${ENTRANCE_PROTOCOL:-http}"
  if [[ -n "${TLS_MODE}" && "${TLS_MODE}" != "off" ]]; then ENTRANCE_PROTOCOL="https"; fi
  case "${ENTRANCE_PROTOCOL}" in
    http) TLS_MODE="off" ;;
    https) TLS_MODE="${TLS_MODE:-self-signed}" ;;
    *) error "unsupported entrance protocol: ${ENTRANCE_PROTOCOL}" ;;
  esac
  case "${TLS_MODE}" in off|self-signed|provided) ;; *) error "unsupported TLS mode: ${TLS_MODE}" ;; esac
  ENTRANCE_DOMAIN="${ENTRANCE_DOMAIN:-127.0.0.1}"
  ENTRANCE_URL="$(build_entrance_base_url "${ENTRANCE_PROTOCOL}" "${ENTRANCE_DOMAIN}" "${ENTRANCE_PORT}" "${ENTRANCE_SSL_PORT}")"
  LOCAL_FRONTEND_DEV="${LOCAL_FRONTEND_DEV_ARG:-${LOCAL_FRONTEND_DEV:-false}}"
  FRONTEND_DEV_HOST="${FRONTEND_DEV_HOST:-127.0.0.1}"
  FRONTEND_DEV_PORT="${FRONTEND_DEV_PORT:-$(port_with_offset 5173 "${PORT_OFFSET}")}"
  FRONTEND_DEV_URL="${FRONTEND_DEV_URL:-http://${FRONTEND_DEV_HOST}:${FRONTEND_DEV_PORT}}"
  FRONTEND_DEV_PROXY_URL="${FRONTEND_DEV_PROXY_URL:-http://host.docker.internal:${FRONTEND_DEV_PORT}}"

  TSDB_PASSWORD_CONFIGURED="false"
  [[ -z "${TSDB_PASSWORD:-}" ]] || TSDB_PASSWORD_CONFIGURED="true"
  TSDB_PASSWORD="${TSDB_PASSWORD:-$(random_hex 12)}"
  REDIS_PASSWORD="${REDIS_PASSWORD:-$(random_hex 12)}"
  JWT_SECRET="${JWT_SECRET:-$(random_hex 32)}"
  TIER0_API_KEY="${TIER0_API_KEY:-$(random_hex 24)}"
  ADMIN_INITIAL_PASSWORD="${ADMIN_INITIAL_PASSWORD:-tier0}"
  UNS_DB_NAME="${UNS_DB_NAME:-postgres}"
  DEFAULT_UNS_DB_URL="postgres://postgres:${TSDB_PASSWORD}@tsdb:5432/${UNS_DB_NAME}?sslmode=disable"
  if [[ -z "${UNS_DB_URL:-}" ]] || is_legacy_postgresql_url "${UNS_DB_URL:-}"; then UNS_DB_URL="${DEFAULT_UNS_DB_URL}"; fi
  SINK_DB_URL="${SINK_DB_URL:-}"
  [[ ! "${SINK_DB_URL}" == *"@tsdb:"* && ! "${SINK_DB_URL}" == *"//tsdb:"* ]] || SINK_DB_URL=""
  if is_local_tsdb_url "${UNS_DB_URL}"; then LOCAL_DB_ENABLED="true"; else LOCAL_DB_ENABLED="false"; fi
  if [[ -n "${REDIS_ADDR:-}" && "${REDIS_ADDR}" != "redis:6379" ]]; then
    LOCAL_REDIS_ENABLED="false"
  else
    LOCAL_REDIS_ENABLED="true"
    REDIS_ADDR="redis:6379"
  fi
  COMPOSE_PROFILES=""
  [[ "${LOCAL_DB_ENABLED}" != "true" ]] || append_compose_profile local-db
  [[ "${LOCAL_REDIS_ENABLED}" != "true" ]] || append_compose_profile local-redis
  DATAINGEST_MQTT_BROKERS="${DATAINGEST_MQTT_BROKERS:-tcp://emqx:1883}"
  DATAINGEST_MQTT_TOPIC="${DATAINGEST_MQTT_TOPIC:-#}"
  install_detail "Product: Tier0 Edge ${PRODUCT_VERSION}"
  install_detail "Access URL: ${ENTRANCE_URL}"

  install_step 3 7 "Preparing volumes and TLS assets"
  VOLUMES_ABS="$(absolute_path_for "${VOLUMES_PATH}")"
  restore_runtime_env_from_volume
  resolve_tier0_installation_id
  resolve_internal_runtime_secrets
  if [[ "${LOCAL_DB_ENABLED}" == "true" && "${TSDB_PASSWORD_CONFIGURED}" != "true" && ! -f "${RUNTIME_ENV_FILE}" && -f "${VOLUMES_ABS}/tsdb/data/PG_VERSION" ]]; then
    error "existing database found without its runtime password; restore .env.runtime or set TSDB_PASSWORD"
  fi
  prepare_volume_dirs "${VOLUMES_ABS}"
  configure_sourceflow_package_install "${VOLUMES_ABS}"
  ensure_host_metrics_paths
  TLS_CERT_DIR="${VOLUMES_ABS}/certs/tls"
  if [[ "${TLS_MODE}" == "provided" ]]; then ensure_provided_cert "${TLS_CERT_DIR}"; else ensure_self_signed_cert "${TLS_CERT_DIR}"; fi
  write_deploy_lock "${VOLUMES_ABS}"

  install_step 4 7 "Writing runtime environment and Compose files"
  write_emqx_runtime_files
  write_runtime_env
  render_compose_files
  install_detail "Compose file: ${DEPLOY_ROOT}/docker-compose.yml"

  install_step 5 7 "Loading bundled images"
  load_bundled_images
  if [[ "${NO_START}" != "true" ]]; then
    install_step 6 7 "Starting Tier0 Edge services"
    (cd "${DEPLOY_ROOT}" && bash "${DEPLOY_ROOT}/bin/compose.sh" up -d)
    install_step 7 7 "Waiting for readiness"
    wait_backend_ready
    run_post_start_initialization
    write_install_state true
  else
    install_step 6 7 "Skipping service startup"
    install_step 7 7 "Finalizing deployment output"
    write_install_state false
  fi
  print_install_summary
}
