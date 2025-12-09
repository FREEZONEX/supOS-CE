# This script intelligently handles the VOLUMES_PATH.
# It checks if the path is empty, or if it's the Linux default on a Windows system,
# and applies the correct OS-specific default path if needed.
# It respects any user-defined custom paths.

info "Checking storage path (VOLUMES_PATH)..."
ENV_FILE="$SCRIPT_DIR/../.env.default"
if [ -f "$SCRIPT_DIR/../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../.env"
fi

if [[ "$platform" == MINGW64* ]]; then
    # On Windows, correct the path if it's empty or set to the Linux default.
    if [ -z "$VOLUMES_PATH" ] || [ "$VOLUMES_PATH" == "/volumes/supos/data" ]; then
        info "Path is unset or is Linux default. Setting the correct path for Windows."
        default_path="$HOME/volumes/supos/data"
        info "Default storage path for Windows is set to: $default_path"
        sed -i "s|^VOLUMES_PATH=.*|VOLUMES_PATH=$default_path|" "$ENV_FILE"
        source "$ENV_FILE" # Reload .env for the current session
    else
        info "Using user-defined VOLUMES_PATH from .env: $VOLUMES_PATH"
    fi
else
    # On Linux/macOS, only set the default path if it's empty.
    if [ -z "$VOLUMES_PATH" ]; then
        info "VOLUMES_PATH is unset. Setting the default path for Linux."
        default_path="/volumes/supos/data"
        info "Default storage path for Linux is set to: $default_path"
        sed -i "s|^VOLUMES_PATH=.*|VOLUMES_PATH=$default_path|" "$ENV_FILE"
        source "$ENV_FILE" # Reload .env for the current session
    else
        info "Using existing VOLUMES_PATH from .env: $VOLUMES_PATH"
    fi
fi
echo