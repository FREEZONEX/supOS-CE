#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
LOCAL_ROOT="${DEPLOY_ROOT}/local"

usage() {
  cat >&2 <<'USAGE'
usage:
  local-deploy.sh create <name> [--target PATH] [--volumes-path PATH] [--port-offset N] [--domain DOMAIN] [--install-nodered-packages] [--no-start]
  local-deploy.sh up <name> [name...]
  local-deploy.sh down <name> [name...]
  local-deploy.sh ps <name> [name...]

default local names:
  host   -> deploy/local/deploy-instance-edge-001, ../volumes/001, offset 100
  branch -> deploy/local/deploy-instance-edge-002, ../volumes/002, offset 200
USAGE
}

sanitize() {
  printf '%s' "$1" | tr -c 'A-Za-z0-9_' '_' | tr '[:upper:]' '[:lower:]'
}

default_target() {
  case "$1" in
    host) printf '%s\n' "${LOCAL_ROOT}/deploy-instance-edge-001" ;;
    branch) printf '%s\n' "${LOCAL_ROOT}/deploy-instance-edge-002" ;;
    *) printf '%s\n' "${LOCAL_ROOT}/deploy-instance-$(sanitize "$1")" ;;
  esac
}

default_volumes_path() {
  case "$1" in
    host) printf '%s\n' "../volumes/001" ;;
    branch) printf '%s\n' "../volumes/002" ;;
    *) printf '%s\n' "../volumes/$(sanitize "$1")" ;;
  esac
}

default_port_offset() {
  case "$1" in
    host) printf '%s\n' "100" ;;
    branch) printf '%s\n' "200" ;;
    *) printf '%s\n' "300" ;;
  esac
}

port_with_offset() {
  local base="$1"
  local offset="$2"
  echo $((base + offset))
}

copy_deploy_payload() {
  local target="$1"
  mkdir -p "${target}"
  rm -rf \
    "${target}/bin" \
    "${target}/templates" \
    "${target}/seed" \
    "${target}/mount" \
    "${target}/regression" \
    "${target}/sql" \
    "${target}/compose.d" \
    "${target}/docker-compose.yml" \
    "${target}/docker-compose.override.yml" \
    "${target}/.env.tmp" \
    "${target}/.install-state.json"

  cp -a "${DEPLOY_ROOT}/bin" "${target}/bin"
  cp -a "${DEPLOY_ROOT}/templates" "${target}/templates"
  cp -a "${DEPLOY_ROOT}/mount" "${target}/mount"
  if [[ -d "${DEPLOY_ROOT}/regression" ]]; then
    cp -a "${DEPLOY_ROOT}/regression" "${target}/regression"
  fi
  if [[ -d "${DEPLOY_ROOT}/sql" ]]; then
    cp -a "${DEPLOY_ROOT}/sql" "${target}/sql"
  fi
  if [[ -d "${DEPLOY_ROOT}/compose.d" ]]; then
    cp -a "${DEPLOY_ROOT}/compose.d" "${target}/compose.d"
  fi
  if [[ -f "${DEPLOY_ROOT}/.env.default" ]]; then
    cp "${DEPLOY_ROOT}/.env.default" "${target}/.env.default"
  fi
  if [[ -f "${DEPLOY_ROOT}/README.md" ]]; then
    cp "${DEPLOY_ROOT}/README.md" "${target}/README.md"
  fi
  if [[ -f "${DEPLOY_ROOT}/.dockerignore" ]]; then
    cp "${DEPLOY_ROOT}/.dockerignore" "${target}/.dockerignore"
  fi
  if [[ -f "${DEPLOY_ROOT}/checksums.sha256" ]]; then
    cp "${DEPLOY_ROOT}/checksums.sha256" "${target}/checksums.sha256"
  fi
}

set_env() {
  local file="$1"
  local key="$2"
  local value="$3"
  local tmp
  touch "${file}"
  tmp="${file}.tmp.$$"
  awk -v k="${key}" -v v="${value}" '
    BEGIN { done = 0 }
    $0 ~ "^" k "=" { print k "=" v; done = 1; next }
    { print }
    END { if (done == 0) print k "=" v }
  ' "${file}" > "${tmp}"
  mv "${tmp}" "${file}"
}

create_instance() {
  local name="$1"
  shift
  local target volumes_path port_offset no_start
  local domain skip_nodered_packages
  target="$(default_target "${name}")"
  volumes_path="$(default_volumes_path "${name}")"
  port_offset="$(default_port_offset "${name}")"
  domain="localhost"
  skip_nodered_packages="true"
  no_start="false"

  while [[ $# -gt 0 ]]; do
    case "$1" in
      --target)
        target="$2"
        shift 2
        ;;
      --volumes-path)
        volumes_path="$2"
        shift 2
        ;;
      --port-offset)
        port_offset="$2"
        shift 2
        ;;
      --domain)
        domain="$2"
        shift 2
        ;;
      --install-nodered-packages)
        skip_nodered_packages="false"
        shift
        ;;
      --no-start)
        no_start="true"
        shift
        ;;
      *)
        echo "unknown create argument: $1" >&2
        exit 2
        ;;
    esac
  done

  if [[ "$target" != /* && "$target" != [A-Za-z]:/* && "$target" != [A-Za-z]:\\* ]]; then
    target="${DEPLOY_ROOT}/../${target}"
  fi
  target="$(mkdir -p "${target}" && cd "${target}" && pwd -P)"
  copy_deploy_payload "${target}"

  local env_file="${target}/.env"
  if [[ ! -f "${env_file}" && -f "${target}/.env.default" ]]; then
    cp "${target}/.env.default" "${env_file}"
  elif [[ ! -f "${env_file}" ]]; then
    echo "deploy env template not found: ${target}/.env.default" >&2
    exit 1
  fi
  set_env "${env_file}" "OS_NAME" "${name}"
  set_env "${env_file}" "COMPOSE_PROJECT_NAME" "edge_$(sanitize "${name}")"
  set_env "${env_file}" "PORT_OFFSET" "${port_offset}"
  set_env "${env_file}" "VOLUMES_PATH" "${volumes_path}"
  set_env "${env_file}" "UNS_DB_URL" ""
  set_env "${env_file}" "SINK_DB_URL" ""
  set_env "${env_file}" "ENTRANCE_PROTOCOL" "http"
  set_env "${env_file}" "ENTRANCE_DOMAIN" "${domain}"

  local -a install_args
  install_args=(--no-start)
  if [[ "${skip_nodered_packages}" == "true" ]]; then
    install_args+=(--skip-nodered-packages)
  fi
  (cd "${target}" && bash bin/install.sh "${install_args[@]}")
  if [[ "${no_start}" != "true" ]]; then
    install_args=()
    if [[ "${skip_nodered_packages}" == "true" ]]; then
      install_args+=(--skip-nodered-packages)
    fi
    (cd "${target}" && bash bin/install.sh "${install_args[@]}")
  fi

  echo "${name}: ${target}"
}

target_for_name() {
  local name="$1"
  local target
  target="$(default_target "${name}")"
  if [[ ! -d "${target}" ]]; then
    echo "local deploy not found for ${name}: ${target}" >&2
    echo "run local-deploy.sh create ${name} first" >&2
    exit 1
  fi
  printf '%s\n' "${target}"
}

run_for_instances() {
  local action="$1"
  shift
  if [[ $# -eq 0 ]]; then
    usage
    exit 2
  fi
  local name target
  for name in "$@"; do
    target="$(target_for_name "${name}")"
    case "${action}" in
      up) (cd "${target}" && bash bin/install.sh) ;;
      down) (cd "${target}" && bash bin/compose.sh down) ;;
      ps) (cd "${target}" && bash bin/compose.sh ps) ;;
    esac
  done
}

if [[ $# -lt 1 ]]; then
  usage
  exit 2
fi

cmd="$1"
shift
case "${cmd}" in
  create)
    if [[ $# -lt 1 ]]; then
      usage
      exit 2
    fi
    name="$1"
    shift
    create_instance "${name}" "$@"
    ;;
  up|down|ps)
    run_for_instances "${cmd}" "$@"
    ;;
  *)
    usage
    exit 2
    ;;
esac
