#!/bin/bash

set -e

STOP_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

source "$STOP_SCRIPT_DIR/global/compose-context.sh"
prepare_compose_context

"${compose_args[@]}" stop
