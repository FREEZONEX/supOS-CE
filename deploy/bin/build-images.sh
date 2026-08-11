#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
ROOT_DIR="$(cd "${DEPLOY_DIR}/.." && pwd)"
source "${SCRIPT_DIR}/env-loader.sh"
source_env_file "${DEPLOY_DIR}/.env.default"
if [[ -f "${DEPLOY_DIR}/.env" ]]; then
  source_env_pair "${DEPLOY_DIR}/.env" "${DEPLOY_DIR}/.env.runtime"
fi
: "${PRODUCT_VERSION:?PRODUCT_VERSION is required}"
image_ref="${BACKEND_IMAGE:-harbor.tier0.dev/tier0/tier0-edge-backend:${PRODUCT_VERSION}}"
LOCAL_FRONTEND_DEV="${LOCAL_FRONTEND_DEV:-false}" "${SCRIPT_DIR}/prepare-web-artifacts.sh"
build_args=(--pull=false)
if [[ -n "${DOCKER_BUILD_NETWORK:-}" ]]; then
  build_args+=(--network "${DOCKER_BUILD_NETWORK}")
fi
for proxy_name in HTTP_PROXY HTTPS_PROXY NO_PROXY http_proxy https_proxy no_proxy; do
  if [[ -n "${!proxy_name:-}" ]]; then
    build_args+=(--build-arg "${proxy_name}=${!proxy_name}")
  fi
done
docker build "${build_args[@]}" -t "${image_ref}" -f "${ROOT_DIR}/backend/Dockerfile" "${ROOT_DIR}/backend"
echo "built ${image_ref}"
