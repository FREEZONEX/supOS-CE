#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

SERVICE="sourceflow"
WAIT_SECONDS="${HIDE_NODERED_HEALTH_WAIT_SECONDS:-360}"
CONFIG_WAIT_MS="${HIDE_NODERED_CONFIG_WAIT_MS:-60000}"
CONFIG_WAIT_INTERVAL_MS="${HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS:-3000}"

usage() {
  cat >&2 <<'USAGE'
usage: hide-nodered.sh [--wait-seconds N]

Hides nodes that should not appear after sourceflow installs managed packages.
EventFlow does not install these packages and is not handled here.
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
    --wait-seconds)
      WAIT_SECONDS="${2:?missing wait seconds}"
      shift 2
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

compose() {
  (cd "${DEPLOY_ROOT}" && bash "${DEPLOY_ROOT}/bin/compose.sh" "$@")
}

container_id_for_service() {
  compose ps -q "${SERVICE}" 2>/dev/null | grep -E '^[0-9a-fA-F]{12,64}$' | head -n 1
}

service_http_ready() {
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

wait_for_service_ready() {
  local elapsed=0
  local cid health running

  while (( elapsed < WAIT_SECONDS )); do
    cid="$(container_id_for_service || true)"
    if [[ -n "${cid}" ]]; then
      running="$(docker inspect --format '{{.State.Running}}' "${cid}" 2>/dev/null || true)"
      health="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "${cid}" 2>/dev/null || true)"
      if [[ "${running}" == "true" && "${health}" == "healthy" ]]; then
        return 0
      fi
      if [[ "${running}" == "true" && "${health}" == "none" ]] && service_http_ready; then
        return 0
      fi
    fi
    sleep 5
    elapsed=$((elapsed + 5))
  done

  return 1
}

dump_diagnostics() {
  warn "dumping ${SERVICE} diagnostics..."
  compose ps "${SERVICE}" || true
  compose logs --tail 120 "${SERVICE}" || true
  compose exec -T "${SERVICE}" sh -c 'cd /data && npm ls node-red-contrib-modbus node-red-contrib-opcua node-red-node-supmodel --depth=0' || true
  compose exec -T "${SERVICE}" sh -c 'if [ -f /data/.config.nodes.json ]; then grep -n "\"node-red-contrib-modbus\"\\|\"node-red-contrib-opcua\"\\|\"node-red-node-supmodel\"" /data/.config.nodes.json || true; else echo ".config.nodes.json not found"; fi' || true
}

if [[ ! -f "${DEPLOY_ROOT}/.env" ]]; then
  error "missing ${DEPLOY_ROOT}/.env; run bin/install.sh first"
  exit 1
fi

info "waiting for sourceflow Node-RED runtime..."
if ! wait_for_service_ready; then
  error "timed out waiting for sourceflow to become ready"
  dump_diagnostics
  exit 1
fi

if ! hide_output="$(
  compose exec -T \
    --user 0 \
    -e HIDE_NODERED_CONFIG_WAIT_MS="${CONFIG_WAIT_MS}" \
    -e HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS="${CONFIG_WAIT_INTERVAL_MS}" \
    "${SERVICE}" \
    sh -lc 'node - << "NODE"
const fs = require("fs");
const p = "/data/.config.nodes.json";

const MODBUS_KEY = "node-red-contrib-modbus";
const OPCUA_KEY = "node-red-contrib-opcua";
const OPCUA_OPEN62541_KEY = "node-red-contrib-opcua-open62541";
const SUPMODEL_KEY = "node-red-node-supmodel";
const KEEP_MODBUS = new Set(["Modbus-Read", "Modbus-Client", "Modbus-Server"]);
const KEEP_OPCUA_LEGACY = new Set(["OpcUa-Item", "OpcUa-Client", "OpcUa-Endpoint"]);
const KEEP_OPCUA_OPEN62541 = new Set(["opcua-connection", "opcua-read", "opcua-write", "opcua-subscribe"]);

const waitMs = Number(process.env.HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS || "3000");
const maxWait = Number(process.env.HIDE_NODERED_CONFIG_WAIT_MS || "300000");

function exists(pkg, j) {
  return j[pkg] && j[pkg].nodes;
}
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
function installed(pkg) {
  try {
    require.resolve(`${pkg}/package.json`, { paths: ["/data"] });
    return true;
  } catch {
    return false;
  }
}

(async () => {
  const newOpcuaInstalled = installed(OPCUA_OPEN62541_KEY);
  // When the new open62541 client-only package is present, hide all legacy
  // node-red-contrib-opcua nodes to avoid palette duplication.
  const keepOpcua = newOpcuaInstalled ? new Set() : KEEP_OPCUA_LEGACY;

  const requiredPackages = [
    [MODBUS_KEY, KEEP_MODBUS],
    [OPCUA_KEY, keepOpcua],
    [SUPMODEL_KEY, new Set()],
    ...(newOpcuaInstalled ? [[OPCUA_OPEN62541_KEY, KEEP_OPCUA_OPEN62541]] : []),
  ].filter(([pkg]) => installed(pkg));

  if (requiredPackages.length === 0) {
    console.log("__HIDE_NODERED_SKIPPED__: managed packages are not installed");
    return;
  }

  const start = Date.now();
  let j;

  while (true) {
    try {
      j = JSON.parse(fs.readFileSync(p, "utf8"));
    } catch (e) {
      if (Date.now() - start >= maxWait) {
        console.error("Timeout: cannot read", p);
        process.exit(1);
      }
      await sleep(waitMs);
      continue;
    }

    const missingPackages = requiredPackages
      .map(([pkg]) => (exists(pkg, j) ? null : pkg))
      .filter(Boolean);
    if (missingPackages.length === 0) break;

    if (Date.now() - start >= maxWait) {
      const missing = missingPackages.join(", ");
      console.error("Timeout waiting for packages in .config.nodes.json:", missing || "(unknown)");
      process.exit(1);
    }

    await sleep(waitMs);
  }

  j = JSON.parse(fs.readFileSync(p, "utf8"));
  let changed = false;

  function enableSubset(pkgKey, keepSet) {
    const nodes = j[pkgKey].nodes || {};
    for (const [name, meta] of Object.entries(nodes)) {
      const enabled = keepSet.has(name);
      if (meta.enabled !== enabled) {
        meta.enabled = enabled;
        changed = true;
      }
    }
    console.log("Updated", pkgKey, "=> enabled:", [...keepSet].join(", "));
  }

  for (const [pkg, keepSet] of requiredPackages) {
    enableSubset(pkg, keepSet);
  }

  if (!changed) {
    console.log("__HIDE_NODERED_NO_CHANGES__");
    return;
  }

  fs.writeFileSync(p, JSON.stringify(j, null, 2));
  console.log("Saved:", p);
})().catch((e) => {
  console.error("Error:", e.message);
  process.exit(1);
});
NODE'
)"; then
  error "failed to update ${SERVICE} Node-RED palette visibility"
  dump_diagnostics
  exit 1
fi

printf '%s\n' "${hide_output}"

if printf '%s\n' "${hide_output}" | grep -q '__HIDE_NODERED_NO_CHANGES__$'; then
  info "${SERVICE} Node-RED palette visibility is already up to date"
  exit 0
fi
if printf '%s\n' "${hide_output}" | grep -q '^__HIDE_NODERED_SKIPPED__'; then
  info "${SERVICE} managed packages are not installed; skipping palette visibility update"
  exit 0
fi

compose restart "${SERVICE}" >/dev/null
info "waiting for ${SERVICE} after palette visibility restart..."
if ! wait_for_service_ready; then
  error "timed out waiting for ${SERVICE} after palette visibility restart"
  dump_diagnostics
  exit 1
fi
info "updated ${SERVICE} Node-RED palette visibility"
