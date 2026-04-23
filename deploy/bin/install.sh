#!/bin/bash

set -e

# --- 1. Initialization ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

ENV_FILE="$SCRIPT_DIR/../.env.default"
if [ -f "$SCRIPT_DIR/../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../.env"
fi


sed -i 's/\r$//' "$ENV_FILE" # Clean .env file
source "$ENV_FILE"          # Load initial environment variables
source "$SCRIPT_DIR/global/log.sh"
source "$SCRIPT_DIR/global/choose-profile-command.sh"

SKIP_VOLUME_SYNC=false
for arg in "$@"; do
  case "$arg" in
    --skip-volumes-sync|--skip-volumes)
      SKIP_VOLUME_SYNC=true
      ;;
    -h|--help)
      cat <<'EOF'
Usage: bash install.sh [--skip-volumes-sync]

  --skip-volumes-sync, --skip-volumes
      Skip copying and updating files under VOLUMES_PATH. Use this only when
      the volume directory has already been initialized.
EOF
      exit 0
      ;;
  esac
done
export SKIP_VOLUME_SYNC

platform=$(uname -s)
info "Starting installation on platform: $platform"
echo

# --- 2. Configuration Setup (sourcing from /util) ---
source "$SCRIPT_DIR/util/handle-volumes-path.sh"
source "$SCRIPT_DIR/util/select-ip-address.sh"

# --- 3. Dependency Installation ---
bash "$SCRIPT_DIR/deb/install-docker.sh"

# --- 4. Service Profile Selection ---
# This script will set the 'command' variable for docker-compose
source "$SCRIPT_DIR/util/select-service-profile.sh"

# --- 5. Pre-run Initialization ---
# Render temporary env overrides from the selected compose profiles first so
# Kong and the init scripts all see the same service matrix.
source "$SCRIPT_DIR/util/set-temp-env.sh" "$SCRIPT_DIR/../" "${COMPOSE_PROFILE_ARGS[@]}"
# Kong is configured declaratively. Re-render it before each install so login
# redirects and enabled routes match the current entrance address and profiles.
bash "$SCRIPT_DIR/init/init-kong-property.sh" "$SCRIPT_DIR/.."
source "$SCRIPT_DIR/util/wait-compose-healthy.sh"

DOCKER_COMPOSE_FILE="$SCRIPT_DIR/../docker-compose.yml"

# --- 6. Volume and Image Management ---
source "$SCRIPT_DIR/util/handle-volumes.sh"

if [ -d "$SCRIPT_DIR/../images/" ] && [ "$(ls -A "$SCRIPT_DIR/../images/")" ]; then
  bash "$SCRIPT_DIR/util/load-images.sh"
fi

# --- 7. Main Execution: Start services and run post-init scripts ---
info "Starting Docker containers in detached mode..."
if ! docker compose --env-file "$ENV_FILE" --env-file "$SCRIPT_DIR/../.env.tmp" --project-name tier0 "${COMPOSE_PROFILE_ARGS[@]}" -f "$DOCKER_COMPOSE_FILE" up -d --build --remove-orphans; then
    error "Failed to start Docker containers. Please check the logs above."
    exit 1
fi
info "Containers started successfully. Waiting for core services to become healthy..."
echo
if ! WAIT_COMPOSE_HEALTH_ALLOW_RUNNING_SERVICES="nodered eventflow" wait_compose_healthy 300 5; then
    error "Containers did not become healthy in time."
    exit 1
fi
info "All required containers are healthy."

info "Synchronizing Kong runtime URL settings..."
sync_kong_output="$(bash "$SCRIPT_DIR/util/sync-kong-runtime-url.sh")"
printf '%s\n' "$sync_kong_output"
if printf '%s\n' "$sync_kong_output" | grep -q '^SYNC_KONG_RUNTIME_URL_CHANGED=1$'; then
    KONG_RUNTIME_RESTART_WAIT_SECONDS="${KONG_RUNTIME_RESTART_WAIT_SECONDS:-20}"
    info "Kong runtime URL settings changed. Restarting Kong to apply the updated database config..."
    docker restart kong >/dev/null
    if ! wait_compose_healthy "$KONG_RUNTIME_RESTART_WAIT_SECONDS" 5; then
        error "Kong did not become healthy after applying runtime URL changes."
        exit 1
    fi
fi

# Backend migration creates the IAM OAuth tables, so seed the built-in
# Portainer OAuth client only after the stack is healthy.
info "Seeding IAM OAuth client..."
bash "$SCRIPT_DIR/init/init-iam-sql.sh"

# Run each initialization script individually for clearer error reporting
info "Initializing Node-RED modules..."
bash "$SCRIPT_DIR/init/init-nodered.sh" || {
    error "Failed to initialize Node-RED."
    exit 1
}

info "Initializing EventFlow modules..."
bash "$SCRIPT_DIR/init/init-eventflow.sh" || {
    error "Failed to initialize EventFlow."
    exit 1
}

info "Hiding Node-RED entrypoints..."
bash "$SCRIPT_DIR/init/hide-nodered.sh" || {
    error "Failed to hide Node-RED entrypoints."
    exit 1
}

# Portainer is initialized last because it depends on Kong, IAM OAuth, and
# the final externally reachable BASE_URL.
info "Initializing Portainer OAuth..."
bash "$SCRIPT_DIR/init/init-portainer.sh" || {
    error "Failed to initialize Portainer OAuth."
    exit 1
}

# MinIO is optional. Run its init script in a subshell so its internal exits
# never terminate this installer when the minio profile is disabled.
bash "$SCRIPT_DIR/init/init-minio.sh" "$1" || {
    error "Failed to initialize MinIO."
    exit 1
}


# --- 8. Success ---
echo -e "\n============================================================"
echo -e "🎉  All services are up and running!"
echo -e "👉  Open the platform in your browser:\n"

if [[ "$ENTRANCE_PORT" == "80" || "$ENTRANCE_PORT" == "443" ]]; then
  PLATFORM_URL="${ENTRANCE_PROTOCOL}://${ENTRANCE_DOMAIN}/uns"
else
  PLATFORM_URL="${BASE_URL}/uns"
fi

IAM_BOOTSTRAP_USERNAME="${IAM_BOOTSTRAP_USERNAME:-tier0}"
IAM_BOOTSTRAP_PASSWORD="${IAM_BOOTSTRAP_PASSWORD:-tier0}"
IAM_BOOTSTRAP_EMAIL="${IAM_BOOTSTRAP_EMAIL-}"

echo -e "      $PLATFORM_URL\n"
echo -e "    Default username: ${IAM_BOOTSTRAP_USERNAME}\n"
echo -e "            password: ${IAM_BOOTSTRAP_PASSWORD}\n"
if [[ -n "${IAM_BOOTSTRAP_EMAIL}" ]]; then
  echo -e "               email: ${IAM_BOOTSTRAP_EMAIL}\n"
fi
echo -e "============================================================"
