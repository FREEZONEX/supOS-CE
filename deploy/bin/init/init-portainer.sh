#!/bin/bash

# exit error
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
ENV_FILE="$SCRIPT_DIR/../../.env.default"
if [ -f "$SCRIPT_DIR/../../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../../.env"
fi
source $SCRIPT_DIR/../global/log.sh
source $ENV_FILE
source $SCRIPT_DIR/../../.env.tmp

PORTAINER_ADMIN_PASSWORD="${PORTAINER_ADMIN_PASSWORD:-adminpassword}"

#info "start to init portainer OAuth ..."
set +e
PORTAINER_JWT=$(docker exec nodered curl -skX POST https://portainer:9443/api/auth \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"admin\", \"password\": \"${PORTAINER_ADMIN_PASSWORD}\"}" 2>/dev/null | awk -F'"' '/jwt/ {print $4}')
set -e

if [ -n "$PORTAINER_JWT" ]; then
  info "Successfully obtained Portainer JWT"
else
  warn "Failed to obtain Portainer JWT (admin login). Portainer OAuth/RedirectURI will not be updated."
  warn "Set PORTAINER_ADMIN_PASSWORD in .env to the plaintext for Portainer admin (must match docker-compose --admin-password hash; default: adminpassword)."
  return 0 2>/dev/null || exit 0
fi

#init endpoint
docker exec nodered curl -s -X POST http://portainer:9000/api/endpoints \
  -H "Authorization: Bearer $PORTAINER_JWT" \
  -F "Name=local" \
  -F "EndpointCreationType=1" \
  -F "ContainerEngine=docker"


docker exec nodered curl -skX POST "https://portainer:9443/api/users"   -H "Authorization: Bearer $PORTAINER_JWT"   -H "Content-Type: application/json"   -d '{    "Username": "tier0", "Password": "tier0@tier1304", "Email": "tier0@supos.com","Role": 1 }' > /dev/null 2>&1 && echo "Successfully created administrator 'tier0'"\
|| if [ "$1" == "--verbose" ]; then warn "Failed to create administrator 'tier0'"; fi

if [ "$ENTRANCE_DOMAIN" = "127.0.0.1" ] || [ "$ENTRANCE_DOMAIN" = "localhost" ]; then
  info "ENTRANCE_DOMAIN=$ENTRANCE_DOMAIN, disabling Portainer OAuth ..."
  docker exec nodered curl -skX PUT "https://portainer:9443/api/settings" \
    -H "Authorization: Bearer $PORTAINER_JWT" \
    -H "Content-Type: application/json" \
    -d '{
      "authenticationMethod": 1,
      "userSessionTimeout": "1h"
    }' > /dev/null 2>&1 \
  && info "Successfully disabled Portainer OAuth" \
  || { if [ "$1" == "--verbose" ]; then warn "Failed to disable Portainer OAuth"; fi; }
else
  docker exec nodered curl -skX PUT "https://portainer:9443/api/settings" \
    -H "Authorization: Bearer $PORTAINER_JWT" \
    -H "Content-Type: application/json" \
    -d "{
      \"authenticationMethod\": 3,
      \"oauthSettings\": {
        \"AccessTokenURI\": \"http://kong:8000/keycloak/home/auth/realms/tier0/protocol/openid-connect/token\",
        \"AuthStyle\": 0,
        \"AuthorizationURI\": \"$BASE_URL/keycloak/home/auth/realms/tier0/protocol/openid-connect/auth\",
        \"ClientID\": \"$OAUTH_CLIENT_ID\",
        \"ClientSecret\": \"$OAUTH_CLIENT_SECRET\",
        \"OAuthAutoCreateUsers\": true,
        \"RedirectURI\": \"$BASE_URL/portainer/home/\",
        \"ResourceURI\": \"http://kong:8000/keycloak/home/auth/realms/tier0/protocol/openid-connect/userinfo\",
        \"SSO\": true,
        \"UserIdentifier\":\"preferred_username\",
        \"Scopes\":\"openid\"
      },
      \"userSessionTimeout\": \"1h\"
    }" > /dev/null 2>&1 \
  && info "Successfully set Portainer OAuth" \
  || { if [ "$1" == "--verbose" ]; then warn "Failed to set Portainer OAuth"; fi; }
fi

info "Finished setting Portainer OAuth"
