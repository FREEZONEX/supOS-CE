#!/bin/bash

set -e

HANDLE_VOLUMES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
DEPLOY_BIN="$(cd "$HANDLE_VOLUMES_DIR/.." && pwd)"

source "$DEPLOY_BIN/global/log.sh"
source "$DEPLOY_BIN/util/normalize-volume-shells.sh"

echo "Start creating volumes"

if [ "${SKIP_VOLUME_SYNC:-false}" = "true" ]; then
  info "Skipping volume synchronization because --skip-volumes-sync was provided."
  if [ ! -d "$VOLUMES_PATH" ] || [ -z "$(ls -A "$VOLUMES_PATH" 2>/dev/null)" ]; then
    error "Cannot skip volume synchronization because $VOLUMES_PATH is empty."
    error "Remove --skip-volumes-sync or initialize volumes once before using this shortcut."
    exit 1
  fi
else
  # Check for a specific sub-directory to reliably detect an existing installation.
  if [ -d "$VOLUMES_PATH/postgresql" ]; then
    info "Existing installation detected. Stopping services and updating volumes..."
    source "$DEPLOY_BIN/stop.sh"
    source "$DEPLOY_BIN/init/update-volumes.sh"
  else
    info "New installation detected. Initializing volumes..."
    source "$DEPLOY_BIN/init/init-volumes.sh"
  fi
fi

# After volumes are created, copy the service config file to its final destination.
SOURCE_CONFIG_FILE="$DEPLOY_BIN/global/active-services.txt"
FINAL_CONFIG_FILE="$VOLUMES_PATH/edge/system/active-services.txt"
if [ -f "$SOURCE_CONFIG_FILE" ]; then
  info "Activating selected service profile..."
  mkdir -p "$(dirname "$FINAL_CONFIG_FILE")"
  cp "$SOURCE_CONFIG_FILE" "$FINAL_CONFIG_FILE"
fi

# Windows checkouts often leave CRLF in mounted shell scripts. Normalize the
# final volume tree before docker compose starts any container that bind-mounts
# scripts from VOLUMES_PATH (notably Kong's start.sh).
normalize_volume_shell_scripts "$VOLUMES_PATH"
