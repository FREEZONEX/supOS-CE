#!/bin/bash

set -e
UPDATE_VOLUMES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
DEPLOY_ROOT="$(cd "$UPDATE_VOLUMES_DIR/../.." && pwd)"

is_windows_bind_path() {
    [[ "$1" =~ ^/mnt/[A-Za-z]/ ]] || [[ "$1" =~ ^/[A-Za-z]/ ]]
}

# is_macos is exported by platform-compat.sh; provide a local fallback so
# update-volumes.sh stays safe if invoked outside the normal install flow.
if ! declare -f is_macos > /dev/null 2>&1; then
  is_macos() { [[ "${PLATFORM:-$(uname -s)}" == "Darwin" ]]; }
fi

# load npm cache
rm -rf "$DEPLOY_ROOT/mount/node-red/.npm/"
tar -xf "$DEPLOY_ROOT/mount/node-red/npmCache.tar.xz" -C "$DEPLOY_ROOT/mount/node-red/" > /dev/null 2>&1

info "loading npm cache complete."
find "$DEPLOY_ROOT/mount/grafana/data/plugins/" -type f -name "*.tar.gz" -exec tar -xzf {} -C "$DEPLOY_ROOT/mount/grafana/data/plugins/" \;

# 更新volumes目录
cp -r "$DEPLOY_ROOT/mount/grafana/data/plugins/"* "$VOLUMES_PATH/grafana/data/plugins/"
cp -r "$DEPLOY_ROOT/mount/kong/"* "$VOLUMES_PATH/kong/"
cp "$DEPLOY_ROOT/mount/emqx/config/"* "$VOLUMES_PATH/emqx/config/"
mkdir -p "$VOLUMES_PATH/certs"
mkdir -p "$VOLUMES_PATH/certs/tls"
if [ ! -f "$VOLUMES_PATH/certs/tls/cert.pem" ]; then
    if [ -f "$VOLUMES_PATH/kong/certificationfile/fullchain.cer" ]; then
        cp "$VOLUMES_PATH/kong/certificationfile/fullchain.cer" "$VOLUMES_PATH/certs/tls/cert.pem"
    elif [ -f "$DEPLOY_ROOT/mount/certs/tls/cert.pem" ]; then
        cp "$DEPLOY_ROOT/mount/certs/tls/cert.pem" "$VOLUMES_PATH/certs/tls/cert.pem"
    fi
fi
if [ ! -f "$VOLUMES_PATH/certs/tls/key.pem" ]; then
    if [ -f "$VOLUMES_PATH/kong/certificationfile/private.key" ]; then
        cp "$VOLUMES_PATH/kong/certificationfile/private.key" "$VOLUMES_PATH/certs/tls/key.pem"
    elif [ -f "$DEPLOY_ROOT/mount/certs/tls/key.pem" ]; then
        cp "$DEPLOY_ROOT/mount/certs/tls/key.pem" "$VOLUMES_PATH/certs/tls/key.pem"
    fi
fi

mkdir -p "$VOLUMES_PATH/node-red"
mkdir -p "$VOLUMES_PATH/node-red/offline_modules"
mkdir -p "$VOLUMES_PATH/node-red/themes"
mkdir -p "$VOLUMES_PATH/node-red/override"
rm -rf "$VOLUMES_PATH/node-red/.npm" && cp -r "$DEPLOY_ROOT/mount/node-red/.npm" "$VOLUMES_PATH/node-red/"
find "$DEPLOY_ROOT/mount/node-red" -maxdepth 1 -type f -exec cp {} "$VOLUMES_PATH/node-red/" \;
cp -r "$DEPLOY_ROOT/mount/node-red/offline_modules/"* "$VOLUMES_PATH/node-red/offline_modules/"
cp -r "$DEPLOY_ROOT/mount/node-red/themes/"* "$VOLUMES_PATH/node-red/themes/"
cp -r "$DEPLOY_ROOT/mount/node-red/template" "$VOLUMES_PATH/node-red/"

mkdir -p "$VOLUMES_PATH/eventflow/" && cp -r "$DEPLOY_ROOT/mount/eventflow/"* "$VOLUMES_PATH/eventflow/"
cp -r "$DEPLOY_ROOT/mount/postgresql/"* "$VOLUMES_PATH/postgresql/"

# Keep Linux ownership fixes for native volumes, but avoid failing on WSL /
# Docker Desktop bind mounts where chmod/chown is either ignored or rejected.
if is_windows_bind_path "$VOLUMES_PATH" || is_macos; then
    warn "Skipping chown for this platform: $VOLUMES_PATH"
else
    chown 999:0 -R "$VOLUMES_PATH/postgresql"
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

info "success to update folder: $VOLUMES_PATH"
