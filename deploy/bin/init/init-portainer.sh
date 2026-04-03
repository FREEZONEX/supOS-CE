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
# AuthorizationURI is reached by the end user's browser, so it must use the
# externally accessible BASE_URL. Token and userinfo are called by the
# Portainer container itself, so they stay on the internal Docker network.
PORTAINER_AUTHORIZATION_URI="${BASE_URL}/inter-api/iam/oauth2/authorize"
PORTAINER_ACCESS_TOKEN_URI="http://kong:8000/inter-api/iam/oauth2/token"
PORTAINER_RESOURCE_URI="http://kong:8000/inter-api/iam/oauth2/userinfo"

json_escape() {
  printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

find_portainer_user_id() {
  docker exec nodered curl -sk "https://portainer:9443/api/users" \
    -H "Authorization: Bearer $PORTAINER_JWT" \
  | tr '{' '\n' \
  | grep "\"Username\":\"$(json_escape "$IAM_BOOTSTRAP_USERNAME")\"" \
  | sed -n 's/.*"Id":\([0-9]\+\).*/\1/p' \
  | head -n 1
}

PORTAINER_JWT=""
PORTAINER_USER_PAYLOAD="{\"Username\":\"$(json_escape "$IAM_BOOTSTRAP_USERNAME")\",\"Role\":1}"
if [ -n "$IAM_BOOTSTRAP_EMAIL" ]; then
  PORTAINER_USER_PAYLOAD="{\"Username\":\"$(json_escape "$IAM_BOOTSTRAP_USERNAME")\",\"Role\":1,\"Email\":\"$(json_escape "$IAM_BOOTSTRAP_EMAIL")\"}"
fi

# Use the Node-RED container as the curl runner because it ships with curl and
# shares the same compose network as Portainer and Kong on both Linux and WSL.
for _ in $(seq 1 60); do
  set +e
  PORTAINER_JWT=$(docker exec nodered curl -skX POST https://portainer:9443/api/auth \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"admin\",\"password\":\"$(json_escape "$PORTAINER_ADMIN_PASSWORD")\"}" 2>/dev/null | awk -F'"' '/jwt/ {print $4}')
  set -e
  if [ -n "$PORTAINER_JWT" ]; then
    break
  fi
  sleep 2
done

if [ -z "$PORTAINER_JWT" ]; then
  error "Failed to obtain Portainer JWT after waiting for Portainer readiness."
  error "Set PORTAINER_ADMIN_PASSWORD in .env to the plaintext for Portainer admin if you changed the default hashed password."
  return 1 2>/dev/null || exit 1
fi

info "Successfully obtained Portainer JWT"

# The local Docker endpoint may already exist on repeated installs or IP
# changes, so treat create failures as non-fatal here.
docker exec nodered curl -s -X POST http://portainer:9000/api/endpoints \
  -H "Authorization: Bearer $PORTAINER_JWT" \
  -F "Name=local" \
  -F "EndpointCreationType=1" \
  -F "ContainerEngine=docker" > /dev/null 2>&1 \
|| if [ "$1" == "--verbose" ]; then warn "Failed to create Portainer local endpoint (it may already exist)"; fi

docker exec nodered curl -skX POST "https://portainer:9443/api/users" \
  -H "Authorization: Bearer $PORTAINER_JWT" \
  -H "Content-Type: application/json" \
  -d "$PORTAINER_USER_PAYLOAD" \
  > /dev/null 2>&1 \
&& info "Successfully created Portainer administrator '$IAM_BOOTSTRAP_USERNAME'" \
|| if [ "$1" == "--verbose" ]; then warn "Failed to create Portainer administrator '$IAM_BOOTSTRAP_USERNAME'"; fi

PORTAINER_BOOTSTRAP_USER_ID="$(find_portainer_user_id)"
if [ -n "$PORTAINER_BOOTSTRAP_USER_ID" ]; then
  docker exec nodered curl -skX PUT "https://portainer:9443/api/users/${PORTAINER_BOOTSTRAP_USER_ID}" \
    -H "Authorization: Bearer $PORTAINER_JWT" \
    -H "Content-Type: application/json" \
    -d "{\"Username\":\"$(json_escape "$IAM_BOOTSTRAP_USERNAME")\",\"Role\":1}" \
    > /dev/null 2>&1 \
  && info "Ensured Portainer user '$IAM_BOOTSTRAP_USERNAME' is an administrator" \
  || { if [ "$1" == "--verbose" ]; then warn "Failed to promote Portainer user '$IAM_BOOTSTRAP_USERNAME' to administrator"; fi; }
else
  if [ "$1" == "--verbose" ]; then warn "Portainer user '$IAM_BOOTSTRAP_USERNAME' was not found after create attempt"; fi
fi

docker exec nodered curl -skX PUT "https://portainer:9443/api/settings" \
  -H "Authorization: Bearer $PORTAINER_JWT" \
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
&& info "Successfully configured Portainer OAuth against platform IAM" \
|| { if [ "$1" == "--verbose" ]; then warn "Failed to configure Portainer OAuth"; fi; }

info "Finished setting Portainer OAuth"
