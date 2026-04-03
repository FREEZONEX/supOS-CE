#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

source "$SCRIPT_DIR/global/compose-context.sh"
source "$SCRIPT_DIR/global/log.sh"
prepare_full_compose_context

# Uninstall only tears down the compose stack. Persistent data cleanup stays in
# clean-all.sh so operators can stop/remove the deployment without wiping host data.
"${compose_args[@]}" down --remove-orphans

rm -f "$DEPLOY_ROOT/.env.tmp" > /dev/null 2>&1
