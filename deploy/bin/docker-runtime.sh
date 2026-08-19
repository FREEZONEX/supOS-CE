#!/usr/bin/env bash

# Shared Docker runtime checks for Linux deployment entrypoints. This file only
# defines functions so it remains safe to source before --help handling.

tier0_resolve_docker_endpoint() {
  local explicit_endpoint="${1:-}"
  local endpoint=""

  if [[ -n "${explicit_endpoint}" ]]; then
    endpoint="${explicit_endpoint}"
  elif [[ -n "${DOCKER_HOST:-}" ]]; then
    endpoint="${DOCKER_HOST}"
  elif command -v docker >/dev/null 2>&1; then
    endpoint="$(docker context inspect --format '{{.Endpoints.docker.Host}}' 2>/dev/null | head -n 1 || true)"
  fi

  if [[ -z "${endpoint}" || "${endpoint}" == "<no value>" ]]; then
    case "$(uname -s)" in
      Linux) endpoint="unix:///var/run/docker.sock" ;;
      Darwin) endpoint="unix://${HOME}/.docker/run/docker.sock" ;;
    esac
  fi

  printf '%s\n' "${endpoint}"
}

tier0_safe_docker_endpoint() {
  local endpoint="${1:-<default Docker endpoint>}"

  case "${endpoint}" in
    unix://*|npipe://*) printf '%s\n' "${endpoint}" ;;
    *://*) printf '%s://<redacted>\n' "${endpoint%%://*}" ;;
    *) printf '<redacted Docker endpoint>\n' ;;
  esac
}

tier0_docker_endpoint_is_rootless() {
  local endpoint="${1:-}"
  local security_options="${2:-}"

  case "${endpoint}" in
    unix:///run/user/*|unix://run/user/*|/run/user/*) return 0 ;;
  esac

  [[ "${security_options,,}" == *rootless* ]]
}

tier0_require_rootful_docker() {
  local requested_endpoint="${1:-}"
  local endpoint display_endpoint security_options

  [[ "$(uname -s)" == "Linux" ]] || return 0

  if ! command -v docker >/dev/null 2>&1; then
    printf 'ERROR: Docker is required; Tier0 Edge on Linux supports only a rootful Docker daemon.\n' >&2
    return 1
  fi

  endpoint="$(tier0_resolve_docker_endpoint "${requested_endpoint}")"
  display_endpoint="$(tier0_safe_docker_endpoint "${endpoint}")"
  if tier0_docker_endpoint_is_rootless "${endpoint}"; then
    printf 'ERROR: Rootless Docker is not supported by Tier0 Edge on Linux.\n' >&2
    printf 'Detected Docker endpoint: %s\n' "${display_endpoint}" >&2
    printf 'Use a rootful Docker daemon, normally unix:///var/run/docker.sock.\n' >&2
    return 1
  fi

  if [[ -n "${requested_endpoint}" ]]; then
    if ! security_options="$(DOCKER_HOST="${endpoint}" docker info --format '{{json .SecurityOptions}}' 2>/dev/null)"; then
      printf 'ERROR: Unable to verify the Docker daemon at %s.\n' "${display_endpoint}" >&2
      printf 'Tier0 Edge on Linux requires a reachable rootful Docker daemon.\n' >&2
      return 1
    fi
  elif ! security_options="$(docker info --format '{{json .SecurityOptions}}' 2>/dev/null)"; then
    printf 'ERROR: Unable to verify the Docker daemon at %s.\n' "${display_endpoint}" >&2
    printf 'Tier0 Edge on Linux requires a reachable rootful Docker daemon.\n' >&2
    return 1
  fi

  if [[ -z "${security_options}" || "${security_options}" == "<no value>" ]]; then
    printf 'ERROR: Docker daemon security options could not be verified at %s.\n' "${display_endpoint}" >&2
    printf 'Tier0 Edge on Linux requires a verifiable rootful Docker daemon.\n' >&2
    return 1
  fi

  if tier0_docker_endpoint_is_rootless "${endpoint}" "${security_options}"; then
    printf 'ERROR: Rootless Docker is not supported by Tier0 Edge on Linux.\n' >&2
    printf 'Detected Docker endpoint: %s\n' "${display_endpoint}" >&2
    printf 'Use a rootful Docker daemon, normally unix:///var/run/docker.sock.\n' >&2
    return 1
  fi
}
