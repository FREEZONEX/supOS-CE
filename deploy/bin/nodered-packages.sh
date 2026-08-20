#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${DEPLOY_ROOT}/.env"
RUNTIME_ENV_FILE="${DEPLOY_ROOT}/.env.runtime"
SERVICE="sourceflow"
BACKGROUND="false"
DATA_WAIT_SECONDS="${NODERED_PACKAGE_DATA_WAIT_SECONDS:-60}"
RUNTIME_WAIT_SECONDS="${NODERED_PACKAGE_RUNTIME_WAIT_SECONDS:-360}"

# shellcheck source=env-loader.sh
source "${SCRIPT_DIR}/env-loader.sh"

usage() {
  cat >&2 <<'USAGE'
usage: nodered-packages.sh [--background]

Installs managed SourceFlow Node-RED packages, hides noisy protocol nodes,
and restarts sourceflow when package installation changed the runtime.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --background)
      BACKGROUND="true"
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
  echo "deploy env not found: ${ENV_FILE}" >&2
  exit 2
fi

source_env_pair "${ENV_FILE}" "${RUNTIME_ENV_FILE}"

is_absolute_path() {
  [[ "$1" == /* || "$1" == [A-Za-z]:/* || "$1" == [A-Za-z]:\\* ]]
}

sourceflow_data_dir() {
  local root="${VOLUMES_PATH:-/volumes/tier0/data}"
  if is_absolute_path "${root}"; then
    printf '%s/sourceflow\n' "${root}"
  else
    printf '%s/%s/sourceflow\n' "${DEPLOY_ROOT}" "${root}"
  fi
}

LOG_DIR="$(sourceflow_data_dir)"
mkdir -p "${LOG_DIR}"
LOG_FILE="${LOG_DIR}/package-install.log"
RUNTIME_STATUS_FILE="${LOG_DIR}/.seed-packages-runtime-status"

if [[ "${BACKGROUND}" == "true" ]]; then
  exec >> "${LOG_FILE}" 2>&1
else
  exec > >(tee -a "${LOG_FILE}") 2>&1
fi

info() {
  printf '[INFO] %s\n' "$*"
}

warn() {
  printf '[WARN] %s\n' "$*" >&2
}

error() {
  printf '[ERROR] %s\n' "$*" >&2
}

RUNTIME_STATUS="running"

write_runtime_status() {
  local status="$1"
  local detail="${2:-}"
  printf '%s %s%s\n' \
    "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    "${status}" \
    "${detail:+ ${detail}}" > "${RUNTIME_STATUS_FILE}"
}

record_runtime_failure() {
  local code="$?"
  if [[ "${RUNTIME_STATUS}" == "running" ]]; then
    write_runtime_status "failed" "exit=${code}"
  fi
  exit "${code}"
}

trap record_runtime_failure EXIT
write_runtime_status "running"

compose() {
  (cd "${DEPLOY_ROOT}" && bash "${DEPLOY_ROOT}/bin/compose.sh" "$@")
}

container_id_for_service() {
  compose ps -q "${SERVICE}" 2>/dev/null | grep -E '^[0-9a-fA-F]{12,64}$' | head -n 1
}

sourceflow_data_ready() {
  compose exec -T --user 0 "${SERVICE}" sh -lc \
    'test -d /data && test -w /data && mkdir -p /data/.npm && test -d /data/.npm' \
    >/dev/null 2>&1
}

wait_for_sourceflow_data() {
  local elapsed=0
  while (( elapsed < DATA_WAIT_SECONDS )); do
    if sourceflow_data_ready; then
      return 0
    fi
    sleep 2
    elapsed=$((elapsed + 2))
  done
  return 1
}

dump_sourceflow_diagnostics() {
  local cid
  warn "dumping ${SERVICE} diagnostics..."
  compose ps "${SERVICE}" || true
  compose logs --tail 100 "${SERVICE}" || true
  cid="$(container_id_for_service || true)"
  if [[ -n "${cid}" ]]; then
    docker inspect "${cid}" \
      --format '{{range .Mounts}}{{println .Type .Source "->" .Destination .RW}}{{end}}' || true
  fi
}

sourceflow_package_skipped() {
  compose exec -T "${SERVICE}" sh -lc 'test -f /data/.skip-package-install' >/dev/null 2>&1
}

sourceflow_packages_current() {
  compose exec -T "${SERVICE}" sh -lc '
seed_name="${NODE_RED_SEED_NAME:-sourceflow}"
seed_dir="${NODE_RED_SEED_DIR:-/opt/tier0-mount/${seed_name}}"
version="${NODE_RED_SEED_VERSION:-}"
if [ -z "$version" ] && [ -f "${seed_dir}/.seed-version" ]; then
  version="$(cat "${seed_dir}/.seed-version")"
fi
if [ -z "$version" ]; then
  version="dev"
fi
test -f /data/.seed-packages-version &&
  test "$(cat /data/.seed-packages-version)" = "$version" &&
  test -d /data/node_modules
' >/dev/null 2>&1
}

ensure_sourceflow_running() {
  compose up -d "${SERVICE}" >/dev/null
  if wait_for_sourceflow_data; then
    return 0
  fi

  warn "${SERVICE} /data mount did not become writable; recreating only ${SERVICE} and retrying"
  compose up -d --force-recreate "${SERVICE}" >/dev/null
  if wait_for_sourceflow_data; then
    return 0
  fi

  error "${SERVICE} /data mount is unavailable after recreate"
  dump_sourceflow_diagnostics
  return 1
}

install_sourceflow_packages() {
  compose exec -T --user 0 "${SERVICE}" sh -lc 'sh /usr/local/bin/seed-node-red.sh --packages sourceflow'
}

restart_sourceflow() {
  compose restart "${SERVICE}" >/dev/null
}

sourceflow_http_ready() {
  compose exec -T "${SERVICE}" node -e '
const http = require("http");
const req = http.get({ host: "127.0.0.1", port: 1880, path: "/" }, (res) => {
  process.exit(res.statusCode < 500 ? 0 : 1);
});
req.on("error", () => process.exit(1));
req.setTimeout(2000, () => {
  req.destroy();
  process.exit(1);
});
' >/dev/null 2>&1
}

sourceflow_runtime_ready() {
  local cid health running
  cid="$(container_id_for_service || true)"
  [[ -n "${cid}" ]] || return 1
  running="$(docker inspect --format '{{.State.Running}}' "${cid}" 2>/dev/null || true)"
  health="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "${cid}" 2>/dev/null || true)"
  [[ "${running}" == "true" ]] || return 1
  [[ "${health}" == "healthy" ]] || { [[ "${health}" == "none" ]] && sourceflow_http_ready; }
}

sourceflow_expected_types_registered() {
  compose exec -T "${SERVICE}" node -e '
const fs = require("fs");
const seedName = process.env.NODE_RED_SEED_NAME || "sourceflow";
const seedDir = process.env.NODE_RED_SEED_DIR || `/opt/tier0-mount/${seedName}`;
const expectedFile = `${seedDir}/expected-node-types.txt`;
if (!fs.existsSync(expectedFile)) {
  console.error(`expected node type manifest is missing: ${expectedFile}`);
  process.exit(2);
}
const expected = fs.readFileSync(expectedFile, "utf8")
  .split(/\r?\n/)
  .map((line) => line.replace(/#.*/, "").trim())
  .filter(Boolean);
const token = process.env.NODERED_INTERNAL_TOKEN || "";
fetch("http://127.0.0.1:1880/nodes", {
  headers: { "X-Tier0-Internal-Token": token }
}).then(async (response) => {
  const body = await response.text();
  if (!response.ok) {
    console.error(`Node-RED /nodes returned HTTP ${response.status}`);
    process.exit(1);
  }
  const registered = new Set(
    [...body.matchAll(/RED\.nodes\.registerType\(\s*["\x27]([^"\x27]+)["\x27]/g)]
      .map((match) => match[1])
  );
  const missing = expected.filter((type) => !registered.has(type));
  if (missing.length > 0) {
    console.error(`Node-RED runtime is missing expected node types: ${missing.join(", ")}`);
    process.exit(1);
  }
  console.log(`verified ${expected.length} expected Node-RED node types`);
}).catch((error) => {
  console.error(`failed to query Node-RED /nodes: ${error.message}`);
  process.exit(1);
});
'
}

wait_for_current_runtime_state() {
  local elapsed=0
  while (( elapsed < RUNTIME_WAIT_SECONDS )); do
    if sourceflow_expected_types_registered >/dev/null 2>&1; then
      return 0
    fi
    if sourceflow_runtime_ready; then
      return 1
    fi
    sleep 5
    elapsed=$((elapsed + 5))
  done
  return 1
}

wait_for_expected_node_types() {
  local max_wait="${1:-${RUNTIME_WAIT_SECONDS}}"
  local elapsed=0
  while (( elapsed < max_wait )); do
    if sourceflow_expected_types_registered >/dev/null 2>&1; then
      return 0
    fi
    sleep 5
    elapsed=$((elapsed + 5))
  done
  return 1
}

info "sourceflow package installation started at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
info "log file: ${LOG_FILE}"

ensure_sourceflow_running

if sourceflow_package_skipped; then
  info "sourceflow package installation skipped: /data/.skip-package-install exists"
  RUNTIME_STATUS="skipped"
  write_runtime_status "skipped"
  exit 0
fi

if sourceflow_packages_current; then
  info "sourceflow package files are already installed for the current seed version"
  if ! wait_for_current_runtime_state; then
    warn "installed package files are not active in Node-RED; restarting ${SERVICE}"
    restart_sourceflow
  fi
else
  install_sourceflow_packages

  info "restarting sourceflow to load installed packages"
  restart_sourceflow
fi

if ! bash "${DEPLOY_ROOT}/bin/hide-nodered.sh"; then
  error "failed to reconcile sourceflow palette visibility"
  exit 1
fi

info "verifying sourceflow Node-RED runtime node registration"
if ! wait_for_expected_node_types; then
  sourceflow_expected_types_registered || true
  error "sourceflow packages exist on disk but expected node types are not registered"
  dump_sourceflow_diagnostics
  exit 1
fi

RUNTIME_STATUS="success"
write_runtime_status "success"
info "sourceflow package installation completed at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
