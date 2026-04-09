#!/bin/bash

set -e

INIT_MINIO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
source "$INIT_MINIO_DIR/../global/log.sh"

# MinIO is optional. When its compose profile is disabled, skip the bootstrap
# work so later init steps (for example Portainer) still run.
if ! docker ps --format '{{.Names}}' | grep -qx 'minio'; then
  info "MinIO container is not running. Skipping MinIO initialization."
  return 0 2>/dev/null || exit 0
fi

info "Starting MinIO initialization..."

docker exec minio sh -c "mc alias set myminio http://minio:9000 admin adminpassword" >/dev/null

docker exec minio sh -c "mc admin policy info myminio public-delete-policy >/dev/null 2>&1 || mc admin policy create myminio public-delete-policy /data/public-delete-policy.json" >/dev/null

#初始化minio supos bucket
#docker exec minio sh -c "mc mb myminio/supos" \
#|| if [ "$1" == "--verbose" ]; then warn "failed to make bucket 'supos', perhaps already exited?"; fi

#supos 变成公共访问 bucket
#docker exec minio sh -c "mc anonymous set public myminio/supos" \
#|| if [ "$1" == "--verbose" ]; then error "failed to set bucket 'supos' as public access"; fi

#复制 /data/supos 目录下的文件到 supos bucket
#docker exec minio sh -c "mc mirror --overwrite /data/supos myminio/supos > /dev/null && echo 'Copy frontend .svg completed successfully.' \
#|| warn 'Copy frontend .svg failed.'" > /dev/null 2>&1
#docker restart minio && echo '<< minio restarted successfully' || error "failed to restart minio"

info "MinIO initialization complete"
