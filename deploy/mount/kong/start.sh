#!/bin/bash
# kong config  db_export  /tmp/kong_config.yml  && docker cp kong:/tmp/kong_config.yml .

kong migrations bootstrap && kong config db_import /etc/kong/kong_config.yml

kong start