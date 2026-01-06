#!/bin/bash

# exit error
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

echo "Info: start loading images. it may take few minutes. please wait..."
for tar_file in $SCRIPT_DIR/../../images/*.tar; do
  docker load -i "$tar_file"
done
echo "Info: Loading completed."
