#!/bin/bash

set -e

INIT_EVENTFLOW_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
source "$INIT_EVENTFLOW_DIR/../global/log.sh"

info "start to init eventflow modules ..."

ensure_container_running() {
    local container="$1"
    local status=""
    status="$(docker inspect -f '{{.State.Status}}' "$container" 2>/dev/null || echo missing)"
    if [ "$status" != "running" ]; then
        warn "$container is not running yet (status=$status). Continuing package initialization."
    fi
    return 0
}

ensure_container_running eventflow
modules_changed=0

read_installed_package_version() {
    local container="$1"
    local package="$2"
    docker exec "$container" sh -c "node -e \"try { process.stdout.write(require('/data/node_modules/${package}/package.json').version) } catch (e) { process.exit(1) }\"" 2>/dev/null
}

ensure_registry_package() {
    local container="$1"
    local package="$2"
    local version="$3"
    local label="${4:-$package}"
    local current_version
    current_version="$(read_installed_package_version "$container" "$package" || true)"
    if [ "$current_version" = "$version" ]; then
        info "$container package ${label}@${version} already installed, skipping."
        return 0
    fi
    info "Installing $container package ${label}@${version}..."
    if ! docker exec "$container" sh -c "cd /data && npm install --no-audit --offline ${package}@${version}"; then
        error "node-red install ${label} failed!"
        return 1
    fi
    modules_changed=1
}

ensure_optional_registry_package() {
    local container="$1"
    local package="$2"
    local version="$3"
    local label="${4:-$package}"
    local current_version
    current_version="$(read_installed_package_version "$container" "$package" || true)"
    if [ "$current_version" = "$version" ]; then
        info "$container package ${label}@${version} already installed, skipping."
        return 0
    fi
    warn "Skipping optional $container package ${label}@${version}; it is not installed and offline cache may be incomplete."
    return 0
}

# --verbose
#docker exec eventflow sh -c "cd /data && npm install --no-audit --offline @supcon-international/node-red-dev-copilot@1.7.5" \
#|| error "node-red install node-red-dev-copilot failed!"
ensure_registry_package nodered "@supcon-international/node-red-mcp-server" "1.3.4" "node-red-mcp-server"

ensure_optional_registry_package eventflow "@flowfuse/node-red-dashboard" "1.26.0" "node-red-dashboard"

ensure_registry_package eventflow "factory-agent-actions" "1.1.0"
ensure_registry_package eventflow "factory-agent-deepseek" "1.1.1"
ensure_registry_package eventflow "factory-agent-gemini" "1.0.6"
ensure_registry_package eventflow "factory-agent-states" "1.1.8"

#
#docker exec eventflow sh -c "cd /data && npm install  --no-audit --offline node-red-contrib-modbus@5.43.0" \
#|| error "node-red install modbus failed!"
#
#docker exec eventflow sh -c "cd /data && npm install --no-audit --offline node-red-contrib-opcua@0.2.339" \
#|| error "node-red install opcua failed!"
#
#docker exec eventflow sh -c "cd /data && npm install --no-audit --offline node-red-contrib-opcda-client@0.0.7" \
#|| error "node-red install opcda failed!"
#
#docker exec eventflow sh -c "cd /data && npm install  --no-audit --offline node-red-contrib-buffer-parser@3.2.2" \
#|| error "node-red install buffer-parser failed!"

# license: GPL-3.0-or-later 默认不安装，用户可以自主安装
#docker exec $2 sh -c "cd /data && npm install $3 --no-audit --offline node-red-contrib-s7@3.1.0" \
#|| error "node-red install Siemens s7 failed!"

#docker exec eventflow sh -c "cd /data && npm install --no-audit --offline node-red-contrib-mcprotocol@1.2.1" \
#|| error "node-red install MITSUBISHI mcprotocol failed!"
#
#docker exec eventflow sh -c "cd /data && npm install --no-audit --offline node-red-contrib-omron-fins@0.5.0" \
#|| error "node-red install OMRON fins failed!"

#docker exec eventflow sh -c "cd /data && npm install --unsafe-perm /data/offline_modules/modules/node-xlsx-0.24.0.tgz"
#docker exec eventflow sh -c "cd /data && npm install --unsafe-perm /data/offline_modules/modules/formidable-3.5.4.tgz"

ensure_registry_package eventflow "node-red-contrib-postgresql" "0.14.2" "postgresql"

if [ "$modules_changed" -eq 1 ]; then
    info "EventFlow packages changed. Restarting eventflow..."
    docker restart eventflow >/dev/null
else
    info "EventFlow packages already up to date. Skipping eventflow restart."
fi
