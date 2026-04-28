#!/bin/bash

set -e

HIDE_NODERED_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
source "$HIDE_NODERED_DIR/../global/log.sh"

NODERED_CONTAINER="${NODERED_CONTAINER:-nodered}"
NODERED_PORT="${NODERED_PORT:-1880}"
HIDE_NODERED_HEALTH_WAIT_SECONDS="${HIDE_NODERED_HEALTH_WAIT_SECONDS:-360}"
HIDE_NODERED_CONFIG_WAIT_MS="${HIDE_NODERED_CONFIG_WAIT_MS:-300000}"
HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS="${HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS:-3000}"

info "start to init protocol nodes...."

wait_for_nodered_ready() {
    local container="$1"
    local port="$2"
    local timeout_seconds="$3"
    local elapsed=0

    while (( elapsed < timeout_seconds )); do
        local health_status
        health_status="$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$container" 2>/dev/null || true)"

        if [[ "$health_status" == "healthy" ]]; then
            info "$container is healthy."
            return 0
        fi

        if [[ "$health_status" == "none" ]] && lsof -i :"$port" > /dev/null 2>&1; then
            warn "$container has no healthcheck result, falling back to host port $port."
            return 0
        fi

        sleep 5
        elapsed=$((elapsed + 5))
    done

    return 1
}

dump_nodered_diagnostics() {
    local container="$1"
    warn "Dumping $container diagnostics..."
    docker ps --filter "name=^/${container}$" --format 'table {{.Names}}\t{{.Status}}' || true
    docker logs --tail 120 "$container" || true
    docker exec "$container" sh -c 'cd /data && npm ls node-red-contrib-modbus node-red-contrib-opcua --depth=0' || true
    docker exec "$container" sh -c 'if [ -f /data/.config.nodes.json ]; then grep -n "\"node-red-contrib-modbus\"\\|\"node-red-contrib-opcua\"" /data/.config.nodes.json || true; else echo ".config.nodes.json not found"; fi' || true
}

if ! wait_for_nodered_ready "$NODERED_CONTAINER" "$NODERED_PORT" "$HIDE_NODERED_HEALTH_WAIT_SECONDS"; then
    error "Timed out waiting for $NODERED_CONTAINER to become ready."
    dump_nodered_diagnostics "$NODERED_CONTAINER"
    exit 1
fi

if ! hide_output="$(
  docker exec \
    --user 0 \
    -e HIDE_NODERED_CONFIG_WAIT_MS="$HIDE_NODERED_CONFIG_WAIT_MS" \
    -e HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS="$HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS" \
    "$NODERED_CONTAINER" \
    bash -lc 'node - << "NODE"
const fs = require("fs");
const p = "/data/.config.nodes.json";

const MODBUS_KEY = "node-red-contrib-modbus";
const OPCUA_KEY  = "node-red-contrib-opcua";   // 若你装的是 iiot 版，改成 node-red-contrib-iiot-opcua

const KEEP_MODBUS = new Set(["Modbus-Read","Modbus-Client","Modbus-Server"]);
const KEEP_OPCUA  = new Set(["OpcUa-Item","OpcUa-Client","OpcUa-Server","OpcUa-Endpoint"]);

const waitMs = Number(process.env.HIDE_NODERED_CONFIG_WAIT_INTERVAL_MS || "3000");
const maxWait = Number(process.env.HIDE_NODERED_CONFIG_WAIT_MS || "300000");

function exists(pkg, j) { return j[pkg] && j[pkg].nodes; }
function sleep(ms) { return new Promise(r=>setTimeout(r,ms)); }

(async () => {
  const start = Date.now();
  let j;

  console.log("Waiting for .config.nodes.json packages with timeout(ms):", maxWait);

  while (true) {
    try { j = JSON.parse(fs.readFileSync(p, "utf8")); }
    catch (e) {
      if (Date.now() - start >= maxWait) {
        console.error("Timeout: cannot read", p);
        process.exit(1);
      }
      await sleep(waitMs);
      continue;
    }

    const hasModbus = exists(MODBUS_KEY, j);
    const hasOpcua  = exists(OPCUA_KEY , j);

    if (hasModbus && hasOpcua) break;

    if (Date.now() - start >= maxWait) {
      const missing = [
        hasModbus ? null : MODBUS_KEY,
        hasOpcua  ? null : OPCUA_KEY
      ].filter(Boolean).join(", ");
      console.error("Timeout waiting for packages in .config.nodes.json:", missing || "(unknown)");
      process.exit(1);
    }

    console.log("[wait] missing:",
      hasModbus ? "" : MODBUS_KEY,
      hasOpcua  ? "" : OPCUA_KEY,
      "— retry in 3s");
    await sleep(waitMs);
  }

  // 重新读一遍，防止等待期间被写入
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

  enableSubset(MODBUS_KEY, KEEP_MODBUS);
  enableSubset(OPCUA_KEY , KEEP_OPCUA);

  if (!changed) {
    console.log("__HIDE_NODERED_NO_CHANGES__");
    return;
  }

  fs.writeFileSync(p, JSON.stringify(j, null, 2));
  console.log("Saved:", p);
})().catch(e => {
  console.error("Error:", e.message);
  process.exit(1);
});
NODE'
)"; then
    error "Failed to update Node-RED palette visibility."
    dump_nodered_diagnostics "$NODERED_CONTAINER"
    exit 1
fi

printf '%s\n' "$hide_output"

if printf '%s\n' "$hide_output" | grep -q '__HIDE_NODERED_NO_CHANGES__$'; then
    info "Node-RED palette visibility is already up to date. Skipping restart."
    exit 0
fi

docker restart "$NODERED_CONTAINER" >/dev/null
