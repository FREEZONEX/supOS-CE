#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

if ! declare -f sed_i > /dev/null 2>&1; then
  source "$SCRIPT_DIR/util/platform-compat.sh"
fi

source "$SCRIPT_DIR/global/compose-context.sh"
prepare_compose_context
source "$SCRIPT_DIR/global/log.sh"

normalize_volumes_path() {
  local target="$1"
  local platform_name
  platform_name="$(uname -s)"

  if command -v cygpath >/dev/null 2>&1 && [[ "$target" =~ ^[A-Za-z]:[\\/].*$ ]]; then
    target="$(cygpath -u "$target")"
  fi

  if [[ "$platform_name" == MINGW* || "$platform_name" == MSYS* ]]; then
    # Git Bash works with /d/... paths. Normalize any WSL-style /mnt/d/... input
    # back to the Windows host drive path before attempting deletion.
    if [[ "$target" =~ ^/mnt/([A-Za-z])/(.*)$ ]]; then
      local drive_letter
      drive_letter="$(echo "${BASH_REMATCH[1]}" | tr '[:upper:]' '[:lower:]')"
      local remainder_path="${BASH_REMATCH[2]}"
      if [ -d "/$drive_letter" ]; then
        target="/$drive_letter/$remainder_path"
      fi
    fi
  else
    # WSL/Linux prefers /mnt/d/... mounts. Normalize Git Bash style /d/... input
    # so cleanup works consistently when scripts are executed from Linux shells.
    if [[ "$target" =~ ^/([A-Za-z])/(.*)$ ]] && [ ! -d "$target" ]; then
      local drive_letter
      drive_letter="$(echo "${BASH_REMATCH[1]}" | tr '[:upper:]' '[:lower:]')"
      local remainder_path="${BASH_REMATCH[2]}"
      if [ -d "/mnt/$drive_letter" ]; then
        target="/mnt/$drive_letter/$remainder_path"
      fi
    fi
  fi

  printf '%s' "$target"
}

purge_volumes_path() {
  local target
  target="$(normalize_volumes_path "$1")"
  if [ -z "$target" ]; then
    error "VOLUMES_PATH is empty. Aborting cleanup."
    return 1
  fi

  local resolved
  # realpath -m normalises the path without requiring the target to exist.
  # macOS's realpath does not support -m; fall back to python3 (always present
  # on macOS) or return the path as-is when neither is available.
  resolved="$(realpath -m "$target" 2>/dev/null \
    || python3 -c "import os,sys; print(os.path.normpath(os.path.abspath(sys.argv[1])))" "$target" 2>/dev/null \
    || echo "$target")"
  if [ "$resolved" = "/" ] || [ "$resolved" = "/mnt" ] || [ "$resolved" = "/volumes" ] || [ "$resolved" = "/home" ]; then
    error "Refusing to purge unsafe VOLUMES_PATH: $resolved"
    return 1
  fi

  if [ ! -d "$resolved" ]; then
    info "No persisted deployment data found at: $resolved"
    return 0
  fi

  if [ ! -f "$resolved/edge/system/docker-compose.yml" ] && [ ! -d "$resolved/postgresql" ] && [ ! -d "$resolved/kong" ]; then
    error "'$resolved' does not look like a Tier0 deploy data directory. Aborting cleanup."
    return 1
  fi

  # Keep the mount root itself and remove only the deployment contents so the
  # next install can reuse a host-precreated directory with the same ownership.
  info "Removing persisted deployment data under: $resolved"
  find "$resolved" -mindepth 1 -maxdepth 1 -exec rm -rf {} +
}

# Handle force flag
if [ "$1" = "-f" ]; then
  FORCE=true
else
  FORCE=false
fi

# Confirmation for deletion
if [ "$FORCE" = false ]; then
  echo
  warn "This operation will remove all supOS data."
  echo
  read -p "Are you sure to stop the deployment and delete data under \"$VOLUMES_PATH\"? (y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cleanup cancelled"
    exit 1
  fi
fi

# Execute uninstall.sh with error handling after the confirmation prompt so a
# declined cleanup does not unexpectedly stop the running stack.
warn "Running uninstall script..."
if ! bash "$SCRIPT_DIR/uninstall.sh"; then
  error "Error: Uninstall script failed"
  exit 1
fi

echo
warn "Removing all supOS data..."
echo

purge_volumes_path "$VOLUMES_PATH"
warn "Cleanup completed"
