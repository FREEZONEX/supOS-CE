#!/bin/bash

# exit error
set -e

docker ps -a -q --filter "network=tier0_edge_network" | xargs --no-run-if-empty docker stop \
&& echo "stopped" || echo "failed"
