#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

ENV_FILE="$SCRIPT_DIR/../.env.default"
if [ -f "$SCRIPT_DIR/../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../.env"
fi

source $ENV_FILE
source $SCRIPT_DIR/global/log.sh
# Handle force flag
if [ "$1" = "-f" ]; then
  FORCE=true
else
  FORCE=false
fi

# Execute uninstall.sh with error handling
warn "Running uninstall script..."
if ! bash $SCRIPT_DIR/uninstall.sh; then
  error "Error: Uninstall script failed"
  exit 1
fi

# Confirmation for deletion
if [ "$FORCE" = false ]; then
  echo
  warn "This operation will remove all supOS data."
  echo
  read -p "Are you sure to delete supOS data directory \"$VOLUMES_PATH\"? (y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Deletion cancelled"
    exit 1
  fi
fi

# Safety check before deletion
if [[ -z "$VOLUMES_PATH" || "$VOLUMES_PATH" == "/" ]]; then
  error "VOLUMES_PATH is not set correctly. Aborting deletion."
  exit 1
fi

echo
warn "Removing all supOS data..."
echo

rm -rf "$VOLUMES_PATH"
warn "Cleanup completed"
