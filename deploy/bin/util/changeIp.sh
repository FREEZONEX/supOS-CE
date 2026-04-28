#!/bin/bash
set -e

CHANGE_IP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
ENV_FILE="$CHANGE_IP_DIR/../../.env.default"
if [ -f "$CHANGE_IP_DIR/../../.env" ]; then
  ENV_FILE="$CHANGE_IP_DIR/../../.env"
fi
# 设置.env环境变量
source "$ENV_FILE"
source "$CHANGE_IP_DIR/../global/log.sh"
source "$CHANGE_IP_DIR/wait-compose-healthy.sh"

BASE_URL=""
if [ "$ENTRANCE_PROTOCOL" == "http" ]; then
  BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_PORT"
  if [[ "$ENTRANCE_PORT" == "80" ]]; then
    BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
  fi
fi
if [ "$ENTRANCE_PROTOCOL" == "https" ]; then
  BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_SSL_PORT"
  if [[ "$ENTRANCE_SSL_PORT" == "443" ]]; then
    BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
  fi
fi

command=$(sed -n '2p' "$VOLUMES_PATH/edge/system/active-services.txt")
read -r -a COMPOSE_PROFILE_ARGS <<< "$command"
source "$CHANGE_IP_DIR/set-temp-env.sh" "$CHANGE_IP_DIR/../.." "$command"
# The entrance address is baked into Kong's declarative config and Portainer's
# OAuth redirect URI, so changing IP requires re-rendering Kong and re-running
# the IAM/Portainer bootstrap steps after restart.
source "$CHANGE_IP_DIR/../init/init-kong-property.sh" "$CHANGE_IP_DIR/../../" && cp "$CHANGE_IP_DIR/../../mount/kong/kong_config.yml" "$CHANGE_IP_DIR/../../mount/kong/start.sh" "$VOLUMES_PATH/kong/"
sync_kong_output="$(bash "$CHANGE_IP_DIR/sync-kong-runtime-url.sh")"
printf '%s\n' "$sync_kong_output"

info "IP修改成功, 正在重启服务..."

docker rm -f uns > /dev/null 2>&1 || true
docker rm -f kong > /dev/null 2>&1 || true
docker rm -f portainer > /dev/null 2>&1 || true

DOCKER_COMPOSE_FILE="$CHANGE_IP_DIR/../../docker-compose.yml"
# Portainer also needs to restart here so its runtime settings and OAuth flow
# pick up the regenerated entrance address.
docker compose --env-file "$ENV_FILE" --env-file "$CHANGE_IP_DIR/../../.env.tmp" --project-name tier0 "${COMPOSE_PROFILE_ARGS[@]}" -f "$DOCKER_COMPOSE_FILE" up -d --remove-orphans uns kong portainer
wait_compose_healthy 600 5
source "$CHANGE_IP_DIR/../init/init-iam-sql.sh"
source "$CHANGE_IP_DIR/../init/init-portainer.sh"

info "重启完成, 新的访问地址：$BASE_URL"
