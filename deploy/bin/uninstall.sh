#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
ENV_FILE="$SCRIPT_DIR/../.env.default"
if [ -f "$SCRIPT_DIR/../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../.env"
fi

# 删除kong数据库
docker rm -f kong > /dev/null 2>&1
docker exec -it postgresql bash -c 'PGPASSWORD="postgres" psql -U postgres -c "DROP DATABASE IF EXISTS kong;"'
docker exec -it postgresql bash -c 'PGPASSWORD="postgres" psql -U postgres -c "CREATE DATABASE kong;"'
docker exec -it postgresql bash -c 'PGPASSWORD="postgres" psql -U postgres -d keycloak -c "TRUNCATE TABLE \"public\".\"offline_user_session\";"'
docker exec -it postgresql bash -c 'PGPASSWORD="postgres" psql -U postgres -d keycloak -c "TRUNCATE TABLE \"public\".\"offline_client_session\";"'

source $ENV_FILE

DOCKER_COMPOSE_FILE=$SCRIPT_DIR/../docker-compose.yml

# 卸载所有服务
command="--profile fuxa --profile grafana --profile minio --profile eventflow --profile konga"

if [ -f $SCRIPT_DIR/../.env.tmp ]; then 
  docker compose --env-file $ENV_FILE --env-file $SCRIPT_DIR/../.env.tmp --project-name tier0 $command -f $DOCKER_COMPOSE_FILE down && rm -f $VOLUMES_PATH/edge/system/active-services.txt
else 
  docker compose --env-file $ENV_FILE --project-name tier0 $command -f $DOCKER_COMPOSE_FILE down && rm -f $VOLUMES_PATH/edge/system/active-services.txt
fi

# 删除所有容器
docker ps -a -q --filter "network=tier0_edge_network" | xargs --no-run-if-empty docker rm -f > /dev/null 2>&1

rm -f $SCRIPT_DIR/../.env.tmp > /dev/null 2>&1

