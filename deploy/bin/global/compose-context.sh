#!/bin/bash

COMPOSE_CONTEXT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
DEPLOY_ROOT="$(cd "$COMPOSE_CONTEXT_DIR/../.." && pwd)"
DOCKER_COMPOSE_FILE="$DEPLOY_ROOT/docker-compose.yml"

resolve_deploy_env_file() {
  ENV_FILE="$DEPLOY_ROOT/.env.default"
  if [ -f "$DEPLOY_ROOT/.env" ]; then
    ENV_FILE="$DEPLOY_ROOT/.env"
  fi
}

load_compose_profiles() {
  PROFILE_FILE="$VOLUMES_PATH/edge/system/active-services.txt"
  PROFILE_LINE=""
  if [ -f "$PROFILE_FILE" ]; then
    PROFILE_LINE="$(sed -n '2p' "$PROFILE_FILE")"
  fi

  COMPOSE_PROFILE_ARGS=()
  if [ -n "$PROFILE_LINE" ]; then
    read -r -a COMPOSE_PROFILE_ARGS <<< "$PROFILE_LINE"
  fi
}

load_all_compose_profiles() {
  ALL_COMPOSE_PROFILE_ARGS=()
  while IFS= read -r profile; do
    [ -n "$profile" ] || continue
    ALL_COMPOSE_PROFILE_ARGS+=(--profile "$profile")
  done < <(
    awk '
      /^[[:space:]]+profiles:[[:space:]]*$/ { in_profiles=1; next }
      in_profiles && /^[[:space:]]+-[[:space:]]+/ {
        line=$0
        sub(/^[[:space:]]*-[[:space:]]*/, "", line)
        if (!seen[line]++) {
          print line
        }
        next
      }
      in_profiles { in_profiles=0 }
    ' "$DOCKER_COMPOSE_FILE"
  )
}

build_compose_args() {
  compose_args=(
    docker compose
    --env-file "$ENV_FILE"
  )
  if [ -f "$DEPLOY_ROOT/.env.tmp" ]; then
    compose_args+=(--env-file "$DEPLOY_ROOT/.env.tmp")
  fi
  compose_args+=(--project-name tier0 "${COMPOSE_PROFILE_ARGS[@]}" -f "$DOCKER_COMPOSE_FILE")
}

prepare_compose_context() {
  resolve_deploy_env_file
  source "$ENV_FILE"
  load_compose_profiles
  build_compose_args
}

prepare_full_compose_context() {
  resolve_deploy_env_file
  source "$ENV_FILE"
  load_all_compose_profiles
  COMPOSE_PROFILE_ARGS=("${ALL_COMPOSE_PROFILE_ARGS[@]}")
  build_compose_args
}
