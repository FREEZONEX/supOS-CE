#!/bin/bash

APPLY_NODERED_DEBUG_PATCH_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
source "$APPLY_NODERED_DEBUG_PATCH_DIR/../global/log.sh"
DEPLOY_ROOT="$(cd "$APPLY_NODERED_DEBUG_PATCH_DIR/../.." && pwd)"

apply_nodered_debug_patch() {
    local container="$1"
    local patch_file="$DEPLOY_ROOT/mount/patches/21-debug.html"
    local target_file="/usr/src/node-red/node_modules/@node-red/nodes/core/common/21-debug.html"

    if [ -z "$container" ]; then
        error "node-red debug patch failed: container name is empty"
        return 1
    fi

    if [ ! -f "$patch_file" ]; then
        error "node-red debug patch failed: patch file not found at $patch_file"
        return 1
    fi

    if ! docker inspect "$container" >/dev/null 2>&1; then
        error "node-red debug patch failed: container '$container' not found"
        return 1
    fi

    docker cp "$patch_file" "${container}:${target_file}" || {
        error "node-red debug patch failed: unable to copy patch into $container"
        return 1
    }

    info "patched $container:$target_file"
}
