#!/bin/bash

ROOT_DIR="$1"
shift

PROFILE_ARGS="$*"
ENV_TMP_FILE="$ROOT_DIR/.env.tmp"

{
  # 设置.env.tmp临时环境变量
  if [ "$LANGUAGE" == "zh-CN" ]; then
    echo "GRAFANA_LANG=zh-Hans"
    echo "FUXA_LANG=zh-cn"
  else
    echo "GRAFANA_LANG=en-US"
    echo "FUXA_LANG=en"
  fi

  # 判断是否启用了ELK
  if echo "$PROFILE_ARGS" | grep -q "elk"; then
     echo "ENABLE_ELK=true"
     echo "ENABLE_ELK_MENU=menu"
  else
     echo "ENABLE_ELK=false"
     echo "ENABLE_ELK_MENU=none"
  fi

  # 判断MQTT
  if echo "$PROFILE_ARGS" | grep -q "gmqtt"; then
     echo "MQTT_PLUG=gmqtt"
  else
     echo "MQTT_PLUG=emqx"
  fi

  if echo "$PROFILE_ARGS" | grep -q "gitea"; then
     echo "ENABLE_GITEA_MENU=menu"
  else
     echo "ENABLE_GITEA_MENU=none"
  fi

  if echo "$PROFILE_ARGS" | grep -q "mcpclient"; then
     echo "ENABLE_MCP=menu"
  else
     echo "ENABLE_MCP=none"
  fi

  REDIRECT_BASE_URL=""
  if [ "$ENTRANCE_PROTOCOL" == "http" ]; then
    REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_PORT"
    if [[ "$ENTRANCE_PORT" == "80" ]]; then
      REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
    fi
  fi
  if [ "$ENTRANCE_PROTOCOL" == "https" ]; then
    REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_SSL_PORT"
    if [[ "$ENTRANCE_SSL_PORT" == "443" ]]; then
      REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
    fi
  fi

  echo "BASE_URL=$REDIRECT_BASE_URL"
  echo "ENABLE_PORTAINER=menu"

  echo "UNS_DB_URL=postgres://postgres:${POSTGRES_PASSWORD}@postgresql:5432/postgres?search_path=supos"
  echo "SINK_PG_URL=postgres://postgres:${POSTGRES_PASSWORD}@tsdb:5432/postgres"
  echo "SINK_TSDB_URL=postgres://postgres:${TSDB_PASSWORD}@tsdb:5432/postgres"
  echo "OS_LLM_TYPE=${LLM_TYPE}"
} > "$ENV_TMP_FILE"
