#!/bin/sh
set -e

if kong migrations bootstrap; then
  kong migrations up
  kong migrations finish
  kong config db_import /etc/kong/kong_config.yml
else
  kong migrations up
  kong migrations finish
fi

exec kong start
