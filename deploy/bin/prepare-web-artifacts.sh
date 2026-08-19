#!/usr/bin/env bash
if [ "${BASH##*/}" != "bash" ]; then
  exec bash "$0" "$@"
fi
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
ROOT_DIR="$(cd "${DEPLOY_DIR}/.." && pwd)"
FRONTEND_DIR="${ROOT_DIR}/frontend"
WEB_ARTIFACT_DIR="${ROOT_DIR}/backend/.build/web"
FRONTEND_BUILD_NODE_OPTIONS="${FRONTEND_BUILD_NODE_OPTIONS:-${NODE_OPTIONS:---max-old-space-size=8192}}"

is_true() {
  case "${1:-}" in
    true|TRUE|1|yes|YES|on|ON) return 0 ;;
    *) return 1 ;;
  esac
}

build_main_web=true
if is_true "${LOCAL_FRONTEND_DEV:-false}"; then
  build_main_web=false
  echo "LOCAL_FRONTEND_DEV=true: skipping main web production build"
fi

frontend_windows_path() {
  if command -v wslpath >/dev/null 2>&1; then
    wslpath -w "${FRONTEND_DIR}"
  elif (cd "${FRONTEND_DIR}" && pwd -W >/dev/null 2>&1); then
    (cd "${FRONTEND_DIR}" && pwd -W)
  else
    echo "cannot convert frontend path for Windows pnpm" >&2
    return 1
  fi
}

cd "${FRONTEND_DIR}"
if command -v node >/dev/null 2>&1 && command -v pnpm >/dev/null 2>&1; then
  export NODE_OPTIONS="${FRONTEND_BUILD_NODE_OPTIONS}"
  pnpm install --frozen-lockfile --ignore-scripts
  if [[ "${build_main_web}" == "true" ]]; then
    pnpm build:scripts
  fi
  if [[ "${build_main_web}" == "true" ]]; then
    # Git Bash/MSYS otherwise rewrites VITE_ASSET_PREFIX="/" to a Windows path.
    (cd apps/web && \
      MSYS2_ENV_CONV_EXCL="VITE_ASSET_PREFIX${MSYS2_ENV_CONV_EXCL:+;${MSYS2_ENV_CONV_EXCL}}" \
      VITE_ASSET_PREFIX="${VITE_ASSET_PREFIX:-/}" \
      VITE_SKIP_LICENSE_CHECKS="${VITE_SKIP_LICENSE_CHECKS:-true}" \
      pnpm build)
  fi
elif command -v cmd.exe >/dev/null 2>&1; then
  FRONTEND_WIN="$(frontend_windows_path)"
  cmd.exe /c "cd /d ${FRONTEND_WIN} && set \"NODE_OPTIONS=${FRONTEND_BUILD_NODE_OPTIONS}\" && pnpm install --frozen-lockfile --ignore-scripts"
  if [[ "${build_main_web}" == "true" ]]; then
    cmd.exe /c "cd /d ${FRONTEND_WIN} && set \"NODE_OPTIONS=${FRONTEND_BUILD_NODE_OPTIONS}\" && pnpm build:scripts"
  fi
  if [[ "${build_main_web}" == "true" ]]; then
    cmd.exe /c "cd /d ${FRONTEND_WIN}\\apps\\web && set \"NODE_OPTIONS=${FRONTEND_BUILD_NODE_OPTIONS}\" && if not defined VITE_ASSET_PREFIX set VITE_ASSET_PREFIX=/ && if not defined VITE_SKIP_LICENSE_CHECKS set VITE_SKIP_LICENSE_CHECKS=true && pnpm build"
  fi
  cmd.exe /c "cd /d ${FRONTEND_WIN}\\apps\\anchor && set \"NODE_OPTIONS=${FRONTEND_BUILD_NODE_OPTIONS}\" && pnpm build"
else
  echo "node/pnpm not found" >&2
  exit 1
fi

rm -rf "${WEB_ARTIFACT_DIR}"
mkdir -p "${WEB_ARTIFACT_DIR}"
if [[ "${build_main_web}" == "true" ]]; then
  cp -R "${FRONTEND_DIR}/apps/web/dist/." "${WEB_ARTIFACT_DIR}/"
  test -f "${WEB_ARTIFACT_DIR}/index.html"
fi

echo "web artifacts prepared: ${WEB_ARTIFACT_DIR}"
