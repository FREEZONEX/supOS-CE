#!/usr/bin/env bash

ensure_deploy_command_path() {
  local candidate new_path
  new_path="${PATH}"
  for candidate in \
    /usr/local/bin \
    /opt/homebrew/bin \
    /Applications/Docker.app/Contents/Resources/bin; do
    if [[ -d "${candidate}" ]]; then
      case ":${new_path}:" in
        *:"${candidate}":*) ;;
        *) new_path="${candidate}:${new_path}" ;;
      esac
    fi
  done
  export PATH="${new_path}"
}

ensure_deploy_command_path

is_windows_shell() {
  case "$(uname -s 2>/dev/null || true)" in
    MINGW*|MSYS*|CYGWIN*) return 0 ;;
  esac
  case "${OSTYPE:-}" in
    msys*|cygwin*) return 0 ;;
  esac
  return 1
}

is_wsl_shell() {
  [[ -r /proc/version ]] && grep -qiE 'microsoft|wsl' /proc/version
}

to_wsl_host_path() {
  local path="$1"
  local drive
  local rest

  case "${path}" in
    /[A-Za-z]/*)
      drive="$(printf '%s' "${path:1:1}" | tr '[:upper:]' '[:lower:]')"
      rest="${path:3}"
      printf '/mnt/%s/%s\n' "${drive}" "${rest}"
      return
      ;;
    [A-Za-z]:/*|[A-Za-z]:\\*)
      if command -v wslpath >/dev/null 2>&1; then
        wslpath -u "${path}"
      else
        printf '%s\n' "${path}"
      fi
      return
      ;;
  esac

  printf '%s\n' "${path}"
}

rewrite_windows_host_path() {
  local path="$1"
  local rest

  case "${path}" in
    /[A-Za-z]/*|[A-Za-z]:/*|[A-Za-z]:\\*)
      printf '%s\n' "${path}"
      return
      ;;
    /home/*/*)
      rest="${path#/home/}"
      rest="${rest#*/}"
      printf '%s\n' "${HOME}/${rest}"
      return
      ;;
    /Users/*/*)
      rest="${path#/Users/}"
      rest="${rest#*/}"
      printf '%s\n' "${HOME}/${rest}"
      return
      ;;
    /volumes/*)
      printf '%s\n' "${HOME}${path}"
      return
      ;;
  esac

  printf '%s\n' "${path}"
}

to_compose_host_path() {
  local path="$1"
  path="$(rewrite_windows_host_path "${path}")"

  if is_wsl_shell; then
    to_wsl_host_path "${path}"
  elif command -v cygpath >/dev/null 2>&1; then
    cygpath -am "${path}"
  else
    printf '%s\n' "${path}"
  fi
}

source_env_file() {
  local file="$1"
  local tmp
  [[ -f "${file}" ]] || return 1
  tmp="$(mktemp "${TMPDIR:-/tmp}/tier0-env.XXXXXX")"
  sed 's/\r$//' "${file}" > "${tmp}"

  set -a
  # shellcheck disable=SC1090
  source "${tmp}"
  set +a
  rm -f "${tmp}"
}

source_env_pair() {
  local env_file="$1"
  local runtime_env_file="${2:-}"

  source_env_file "${env_file}"
  if [[ -n "${runtime_env_file}" && -f "${runtime_env_file}" ]]; then
    source_env_file "${runtime_env_file}"
  fi
}
