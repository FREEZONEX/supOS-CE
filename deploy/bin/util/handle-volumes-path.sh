# This script intelligently handles the VOLUMES_PATH.
# It checks if the path is empty, or if it's the Linux default on a Windows system,
# and applies the correct OS-specific default path if needed.
# It respects any user-defined custom paths.

info "Checking storage path (VOLUMES_PATH)..."
ENV_FILE="$SCRIPT_DIR/../.env.default"
if [ -f "$SCRIPT_DIR/../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../.env"
fi

write_volumes_path() {
    local updated_path="$1"
    sed -i "s|^VOLUMES_PATH=.*|VOLUMES_PATH=$updated_path|" "$ENV_FILE"
    source "$ENV_FILE"
}

if [[ "$platform" == MINGW64* ]]; then
    # On Git Bash for Windows, replace Linux/WSL-style paths with a Git Bash
    # drive path so mkdir/cp/chmod all operate on the Windows host filesystem.
    if [ -z "$VOLUMES_PATH" ] || [ "$VOLUMES_PATH" == "/volumes/tier0/data" ]; then
        info "Path is unset or is Linux default. Setting the correct path for Windows."
        if [ -d "/d" ]; then
            default_path="/d/tier0/data"
        else
            default_path="$HOME/volumes/tier0/data"
        fi
        info "Default storage path for Windows is set to: $default_path"
        write_volumes_path "$default_path"
    elif [[ "$VOLUMES_PATH" =~ ^/mnt/([A-Za-z])/(.*)$ ]]; then
        drive_letter="${BASH_REMATCH[1],,}"
        remainder_path="${BASH_REMATCH[2]}"
        normalized_path="/$drive_letter/$remainder_path"
        info "Normalizing WSL-style VOLUMES_PATH for Git Bash: $VOLUMES_PATH -> $normalized_path"
        write_volumes_path "$normalized_path"
    elif command -v cygpath >/dev/null 2>&1 && [[ "$VOLUMES_PATH" =~ ^[A-Za-z]:[\\/].*$ ]]; then
        normalized_path="$(cygpath -u "$VOLUMES_PATH")"
        info "Normalizing Windows-style VOLUMES_PATH for Git Bash: $VOLUMES_PATH -> $normalized_path"
        write_volumes_path "$normalized_path"
    else
        info "Using user-defined VOLUMES_PATH from .env: $VOLUMES_PATH"
    fi
else
    # On Linux/macOS, only set the default path if it's empty.
    if [ -z "$VOLUMES_PATH" ]; then
        info "VOLUMES_PATH is unset. Setting the default path for Linux."
        default_path="/volumes/tier0/data"
        info "Default storage path for Linux is set to: $default_path"
        sed -i "s|^VOLUMES_PATH=.*|VOLUMES_PATH=$default_path|" "$ENV_FILE"
        source "$ENV_FILE" # Reload .env for the current session
    else
        # WSL users sometimes keep a Git Bash style path like /d/tier0/data in
        # .env. Normalize it once so the later bash scripts and compose mounts
        # consistently see /mnt/d/... on Linux.
        if [[ "$VOLUMES_PATH" =~ ^/([A-Za-z])/(.*)$ ]] && [ ! -d "$VOLUMES_PATH" ]; then
            drive_letter="${BASH_REMATCH[1],,}"
            remainder_path="${BASH_REMATCH[2]}"
            normalized_path="/mnt/$drive_letter/$remainder_path"
            if [ -d "/mnt/$drive_letter" ]; then
                info "Normalizing Windows-style VOLUMES_PATH for WSL: $VOLUMES_PATH -> $normalized_path"
                write_volumes_path "$normalized_path"
            fi
        fi
        info "Using existing VOLUMES_PATH from .env: $VOLUMES_PATH"
    fi
fi
echo
