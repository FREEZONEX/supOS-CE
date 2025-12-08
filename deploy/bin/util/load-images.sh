#!/bin/bash

info "start loading images. it may take few minutes. please wait..."
for tar_file in $SCRIPT_DIR/../images/*.tar*; do
  docker load -i "$tar_file" &
done
wait
info "Loading completed."
