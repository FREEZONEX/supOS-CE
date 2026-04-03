#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

source "$SCRIPT_DIR/global/compose-context.sh"
prepare_compose_context

"${compose_args[@]}" start
