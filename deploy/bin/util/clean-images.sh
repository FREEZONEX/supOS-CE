#!/bin/bash
set -e

CLEAN_IMAGES=$1

if [ -n "$CLEAN_IMAGES" ]; then
    docker rmi -f $(docker images | grep "$CLEAN_IMAGES" | awk '{print $1 ":" $2}')
else 
    docker rmi -f $(docker images | awk '{print $1 ":" $2}')
fi