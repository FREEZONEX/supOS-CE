# This script handles service profile selection for the first installation.

info "Checking for existing service configuration..."

if [ ! -f "$VOLUMES_PATH/edge/system/active-services.txt" ]; then
    # First-time installation
    echo "------------------------------------------------------------"
    echo "This is a first-time installation."
    echo "The default profile will be installed."
    echo -n "You can press [Enter] within 3 seconds to choose a custom profile..."

    if read -t 3; then
        # MANUAL SELECTION: Call the function in interactive mode (no arguments)
        echo; info "Manual profile selection activated."
        if [ "$OS_RESOURCE_SPEC" == "1" ]; then
            command=$(chooseProfile1)
        else
            command=$(chooseProfile2)
        fi
    else
        # AUTOMATIC DEFAULT (TIMEOUT): Call the function in non-interactive "default" mode
        echo; info "Applying default profile..."
        if [ "$OS_RESOURCE_SPEC" == "1" ]; then
            command=$(chooseProfile1 "default")
        else
            command=$(chooseProfile2 "default")
        fi
    fi
    echo "------------------------------------------------------------"
else
    # This is for subsequent runs
    info "Found existing service configuration, using it."
    command=$(sed -n '2p' "$VOLUMES_PATH/edge/system/active-services.txt")
fi

declare -g COMPOSE_PROFILE_ARGS_STR="$command"

declare -ag COMPOSE_PROFILE_ARGS=($command)
