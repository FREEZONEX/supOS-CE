#!/bin/bash

set -e

INIT_PORTAINER_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
ENV_FILE="$INIT_PORTAINER_DIR/../../.env.default"
if [ -f "$INIT_PORTAINER_DIR/../../.env" ]; then
  ENV_FILE="$INIT_PORTAINER_DIR/../../.env"
fi

source "$INIT_PORTAINER_DIR/../global/log.sh"
source "$ENV_FILE"
source "$INIT_PORTAINER_DIR/../../.env.tmp"

PORTAINER_ADMIN_PASSWORD="${PORTAINER_ADMIN_PASSWORD:-adminpassword}"
IAM_PORTAINER_CLIENT_ID="${IAM_PORTAINER_CLIENT_ID:-portainer}"
IAM_PORTAINER_CLIENT_SECRET="${IAM_PORTAINER_CLIENT_SECRET:-Tier0PortainerSecret@1304}"
IAM_BOOTSTRAP_USERNAME="${IAM_BOOTSTRAP_USERNAME:-tier0}"
IAM_BOOTSTRAP_PASSWORD="${IAM_BOOTSTRAP_PASSWORD:-tier0}"
IAM_BOOTSTRAP_EMAIL="${IAM_BOOTSTRAP_EMAIL-}"
PORTAINER_REDIRECT_URI="${BASE_URL}/portainer/home/"
PORTAINER_AUTHORIZATION_URI="${BASE_URL}/inter-api/iam/oauth2/authorize"
PORTAINER_ACCESS_TOKEN_URI="http://kong:8000/inter-api/iam/oauth2/token"
PORTAINER_RESOURCE_URI="http://kong:8000/inter-api/iam/oauth2/userinfo"

json_escape() {
  printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

portainer_api() {
  local method=$1 path=$2
  shift 2
  docker exec nodered curl -sk -X "$method" "https://portainer:9443${path}" \
    -H "Authorization: Bearer $PORTAINER_JWT" "$@" 2>/dev/null
}

portainer_setting_field() {
  local field="$1"
  local settings_json
  settings_json="$(portainer_api GET /api/settings)"
  printf '%s' "$settings_json" | sed -n "s/.*\"${field}\":\"\\([^\"]*\\)\".*/\\1/p" | head -n 1
}

find_portainer_user_id() {
  local username="$1"
  portainer_api GET /api/users \
  | tr '{' '\n' \
  | grep "\"Username\":\"$(json_escape "$username")\"" \
  | sed -n 's/.*"Id":\([0-9]\+\).*/\1/p' \
  | head -n 1
}

ensure_admin_user() {
  local username="$1" email="$2"
  local payload="{\"Username\":\"$(json_escape "$username")\",\"Role\":1}"
  if [ -n "$email" ]; then
    payload="{\"Username\":\"$(json_escape "$username")\",\"Role\":1,\"Email\":\"$(json_escape "$email")\"}"
  fi

  portainer_api POST /api/users \
    -H "Content-Type: application/json" \
    -d "$payload" > /dev/null 2>&1 || true

  local uid
  uid="$(find_portainer_user_id "$username")"
  if [ -z "$uid" ]; then
    warn "Portainer user '$username' not found — cannot promote to admin"
    return 1
  fi

  portainer_api PUT "/api/users/${uid}" \
    -H "Content-Type: application/json" \
    -d '{"Role":1}' > /dev/null 2>&1
  info "Ensured Portainer user '$username' (id=$uid) is administrator"
}

PORTAINER_JWT=""
for _ in $(seq 1 60); do
  set +e
  PORTAINER_JWT=$(
    docker exec nodered curl -skX POST https://portainer:9443/api/auth \
      -H "Content-Type: application/json" \
      -d "{\"username\":\"admin\",\"password\":\"$(json_escape "$PORTAINER_ADMIN_PASSWORD")\"}" 2>/dev/null \
    | awk -F'"' '/jwt/ {print $4}'
  )
  set -e
  [ -n "$PORTAINER_JWT" ] && break
  sleep 2
done

if [ -z "$PORTAINER_JWT" ]; then
  error "Failed to obtain Portainer JWT after waiting for Portainer readiness."
  error "Set PORTAINER_ADMIN_PASSWORD in .env to the plaintext Portainer admin password if you changed the default hashed password."
  return 1 2>/dev/null || exit 1
fi
info "Obtained Portainer admin JWT"

docker exec nodered curl -s -X POST http://portainer:9000/api/endpoints \
  -H "Authorization: Bearer $PORTAINER_JWT" \
  -F "Name=local" \
  -F "EndpointCreationType=1" \
  -F "ContainerEngine=docker" > /dev/null 2>&1 \
|| if [ "$1" == "--verbose" ]; then warn "Failed to create Portainer local endpoint (it may already exist)"; fi

portainer_api PUT /api/settings \
  -H "Content-Type: application/json" \
  -d "{
    \"authenticationMethod\": 3,
    \"oauthSettings\": {
      \"AccessTokenURI\": \"${PORTAINER_ACCESS_TOKEN_URI}\",
      \"AuthStyle\": 0,
      \"AuthorizationURI\": \"${PORTAINER_AUTHORIZATION_URI}\",
      \"ClientID\": \"${IAM_PORTAINER_CLIENT_ID}\",
      \"ClientSecret\": \"${IAM_PORTAINER_CLIENT_SECRET}\",
      \"OAuthAutoCreateUsers\": true,
      \"RedirectURI\": \"${PORTAINER_REDIRECT_URI}\",
      \"ResourceURI\": \"${PORTAINER_RESOURCE_URI}\",
      \"SSO\": true,
      \"UserIdentifier\":\"preferred_username\",
      \"Scopes\":\"openid\"
    },
    \"userSessionTimeout\": \"1h\"
  }" > /dev/null 2>&1 \
&& info "Configured Portainer OAuth against platform IAM" \
|| warn "Failed to configure Portainer OAuth"

CURRENT_AUTHORIZATION_URI="$(portainer_setting_field AuthorizationURI)"
CURRENT_REDIRECT_URI="$(portainer_setting_field RedirectURI)"

if [ "$CURRENT_AUTHORIZATION_URI" != "$PORTAINER_AUTHORIZATION_URI" ] || [ "$CURRENT_REDIRECT_URI" != "$PORTAINER_REDIRECT_URI" ]; then
  error "Portainer OAuth settings mismatch after update."
  error "Expected AuthorizationURI=$PORTAINER_AUTHORIZATION_URI, actual=$CURRENT_AUTHORIZATION_URI"
  error "Expected RedirectURI=$PORTAINER_REDIRECT_URI, actual=$CURRENT_REDIRECT_URI"
  return 1 2>/dev/null || exit 1
fi

info "Verified Portainer OAuth settings: RedirectURI=$CURRENT_REDIRECT_URI"

ensure_admin_user "$IAM_BOOTSTRAP_USERNAME" "$IAM_BOOTSTRAP_EMAIL"

info "Portainer initialization complete"
