#!/bin/bash

# exit error
set -e

docker ps -a -q --filter "network=supos_default_network" | xargs --no-run-if-empty docker stop \
&& echo "stopped" || echo "failed"
