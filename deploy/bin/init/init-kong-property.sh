#!/usr/bin/env bash
# Render Kong's declarative configuration from .env and .env.tmp so both fresh
# installs and IP/profile changes converge on the same final gateway config.
set -e
ROOT_DIR=$1
ENV_FILE="$ROOT_DIR/.env.default"
if [ -f "$ROOT_DIR/.env" ]; then
  ENV_FILE="$ROOT_DIR/.env"
fi
# ---------------------------------------------------------------------------
# 0. Normalise .env line endings (Windows → Unix)
# ---------------------------------------------------------------------------
# Use the new variable name
sed -i 's/\r$//' "$ENV_FILE"

# ---------------------------------------------------------------------------
# 1. Load variables from .env
# ---------------------------------------------------------------------------
set -a
source "$ENV_FILE"
set +a

# Portainer OAuth uses built-in defaults unless callers explicitly override
# them in .env. Export them here so envsubst can render kong_config.yml.tpl
# without requiring these knobs to exist in .env.default.
export IAM_PORTAINER_CLIENT_ID="${IAM_PORTAINER_CLIENT_ID:-portainer}"
export IAM_PORTAINER_CLIENT_SECRET="${IAM_PORTAINER_CLIENT_SECRET:-Tier0PortainerSecret@1304}"

# ---------------------------------------------------------------------------
# 2. Build BASE_URL
# BASE_URL is reused by login redirects, Portainer OAuth, and any route rules
# that need the externally reachable entrance address.
# ---------------------------------------------------------------------------
REDIRECT_BASE_URL=""
if [ "$ENTRANCE_PROTOCOL" == "http" ]; then
  REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_PORT"
  if [[ "$ENTRANCE_PORT" == "80" ]]; then
    REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
  fi
fi
if [ "$ENTRANCE_PROTOCOL" == "https" ]; then
  REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN:$ENTRANCE_SSL_PORT"
  if [[ "$ENTRANCE_SSL_PORT" == "443" ]]; then
    REDIRECT_BASE_URL="$ENTRANCE_PROTOCOL://$ENTRANCE_DOMAIN"
  fi
fi

export BASE_URL="$REDIRECT_BASE_URL"

# ---------------------------------------------------------------------------
# 3. Authentication flag → KONG_AUTH_ENABLED
# Keep Kong auth aligned with OS_AUTH_ENABLE so localhost installs and normal
# Linux deployments render the same config from a single switch.
# ---------------------------------------------------------------------------
OS_AUTH_ENABLE=${OS_AUTH_ENABLE:-true}
if [[ "${OS_AUTH_ENABLE,,}" == "true" ]]; then
  export KONG_AUTH_ENABLED=true
else
  export KONG_AUTH_ENABLED=false
fi

# ---------------------------------------------------------------------------
# 4. Load .env.tmp overrides (if the file exists)
# .env.tmp is produced from the selected compose profiles. Load it after .env so
# profile-specific toggles win when templating kong_config.yml.
# ---------------------------------------------------------------------------
if [[ -f "$ROOT_DIR/.env.tmp" ]]; then
  set -a
  source "$ROOT_DIR/.env.tmp"
  set +a
fi
# ---------------------------------------------------------------------------
# 6. Render Kong configuration
# The generated file is later copied into the mounted data directory and then
# imported by Kong on startup via `kong config db_import`.
# ---------------------------------------------------------------------------
# Use the new variable name
TEMPLATE_FILE=$ROOT_DIR/mount/kong/kong_config.yml.tpl
OUTPUT_FILE=$ROOT_DIR/mount/kong/kong_config.yml

rm -f "$OUTPUT_FILE"
envsubst < "$TEMPLATE_FILE" > "$OUTPUT_FILE"

echo "Info: success to generate kong_config.yml (auth enabled = $KONG_AUTH_ENABLED)"
