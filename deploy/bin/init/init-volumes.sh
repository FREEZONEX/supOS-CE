#!/bin/bash

set -e
INIT_VOLUMES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
DEPLOY_ROOT="$(cd "$INIT_VOLUMES_DIR/../.." && pwd)"

is_windows_bind_path() {
    [[ "$1" =~ ^/mnt/[A-Za-z]/ ]] || [[ "$1" =~ ^/[A-Za-z]/ ]]
}


# load npm cache
tar -xf "$DEPLOY_ROOT/mount/node-red/npmCache.tar.xz" -C "$DEPLOY_ROOT/mount/node-red/"
tar -xf "$DEPLOY_ROOT/mount/node-red/npmCache.tar.xz" -C "$DEPLOY_ROOT/mount/eventflow/"

info "loading npm cache complete."
find "$DEPLOY_ROOT/mount/grafana/data/plugins/" -type f -name "*.tar.gz" -exec tar -xzf {} -C "$DEPLOY_ROOT/mount/grafana/data/plugins/" \;

# 创建volumes目录
mkdir -p "$VOLUMES_PATH"
for source_dir in "$DEPLOY_ROOT"/mount/*; do
    cp -r "$source_dir" "$VOLUMES_PATH"
done
# Linux-native mounts keep the original ownership model. Windows-backed mounts
# under /mnt/<drive> or /d/... do not support the same uid/gid semantics, so
# skip the ownership mutation there and rely on Docker Desktop's bind mounts.
if is_windows_bind_path "$VOLUMES_PATH" || is_macos; then
    warn "Skipping chmod/chown for this platform: $VOLUMES_PATH"
else
    chown 999:0 -R "$VOLUMES_PATH/postgresql"
    chmod 644 "$VOLUMES_PATH"/postgresql/conf/*.conf
    chown 1000:1000 -R "$VOLUMES_PATH/emqx"
    chown 755:0 -R "$VOLUMES_PATH/grafana"
fi

cp "$DEPLOY_ROOT/docker-compose.yml" "$VOLUMES_PATH/edge/system/"
# 设置.sh文件为可执行文件
if is_windows_bind_path "$VOLUMES_PATH" || is_macos; then
    warn "Skipping chmod scan for this platform: $VOLUMES_PATH"
else
    find "$VOLUMES_PATH" -name "*.sh" -exec chmod +x {} \;
fi

info "success to create folder: $VOLUMES_PATH"
