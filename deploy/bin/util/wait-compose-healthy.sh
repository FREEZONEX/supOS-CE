#!/usr/bin/env bash
# wait_compose_healthy — wait until every active compose service is running and
# (when Docker healthchecks exist) healthy.
#
# Caller must set: ENV_FILE, DOCKER_COMPOSE_FILE, COMPOSE_PROFILE_ARGS (array)
# log.sh must already be sourced (info / error).
#
# Do not use SCRIPT_DIR here: sourcing deploy/bin/init/*.sh overwrites SCRIPT_DIR and would
# point .env.tmp at deploy/bin/.env.tmp instead of deploy/.env.tmp.

_WAIT_COMPOSE_DEPLOY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
_WAIT_COMPOSE_SET_TEMP_ENV="${_WAIT_COMPOSE_DEPLOY_ROOT}/bin/util/set-temp-env.sh"

wait_compose_healthy() {
  local max_wait="${1:-900}"
  local interval="${2:-5}"
  local env_tmp="${_WAIT_COMPOSE_DEPLOY_ROOT}/.env.tmp"

  if [[ ! -f "$env_tmp" ]] && [[ -f "$_WAIT_COMPOSE_SET_TEMP_ENV" ]]; then
    warn "wait_compose_healthy: missing $env_tmp; recreating via set-temp-env.sh"
    # shellcheck disable=SC1090
    source "$_WAIT_COMPOSE_SET_TEMP_ENV" "${_WAIT_COMPOSE_DEPLOY_ROOT}" "${COMPOSE_PROFILE_ARGS[*]}"
  fi
  local elapsed=0
  local last_info=-30
  local compose=(
    docker compose
    --env-file "$ENV_FILE"
    --env-file "$env_tmp"
    --project-name tier0
    "${COMPOSE_PROFILE_ARGS[@]}"
    -f "$DOCKER_COMPOSE_FILE"
  )

  if [[ ! -f "$env_tmp" ]]; then
    error "wait_compose_healthy: missing env tmp file: $env_tmp"
    return 1
  fi

  while (( elapsed < max_wait )); do
    local failed=()
    local svc_list

    if ! svc_list="$("${compose[@]}" config --services 2>/dev/null)"; then
      error "wait_compose_healthy: docker compose config --services failed"
      return 1
    fi

    if [[ -z "${svc_list// }" ]]; then
      error "wait_compose_healthy: no services returned from compose config"
      return 1
    fi

    while IFS= read -r svc; do
      [[ -z "$svc" ]] && continue
      local cid
      cid="$("${compose[@]}" ps -q "$svc" 2>/dev/null | head -n 1)"
      if [[ -z "$cid" ]]; then
        failed+=("$svc:not_running")
        continue
      fi

      local status health
      status="$(docker inspect -f '{{.State.Status}}' "$cid" 2>/dev/null || echo unknown)"
      health="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$cid" 2>/dev/null || echo unknown)"

      if [[ "$status" != "running" ]]; then
        failed+=("$svc:$status")
        continue
      fi
      if [[ "$health" != "none" && "$health" != "healthy" ]]; then
        failed+=("$svc:health=$health")
      fi
    done <<< "$svc_list"

    if (( ${#failed[@]} == 0 )); then
      return 0
    fi

    if (( elapsed - last_info >= 30 )); then
      last_info=$elapsed
      info "Waiting for containers (${elapsed}s / ${max_wait}s): ${failed[*]}"
    fi

    sleep "$interval"
    elapsed=$((elapsed + interval))
  done

  error "Timeout after ${max_wait}s; containers not ready: ${failed[*]}"
  return 1
}
