#!/bin/sh
set -e

if kong migrations bootstrap; then
  kong migrations up
  kong migrations finish
  kong config db_import /etc/kong/kong_config.yml
else
  kong migrations up || true
  kong migrations finish || true
  kong config db_import /etc/kong/kong_config.yml || true
fi

exec kong start
