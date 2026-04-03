#!/bin/bash

LOAD_IMAGES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

info "start loading images. it may take few minutes. please wait..."
for tar_file in "$LOAD_IMAGES_DIR"/../../images/*.tar*; do
  docker load -i "$tar_file" &
done
wait
info "Loading completed."
