# -*- coding:utf-8 -*-
.PHONY: update

start:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker compose up -d


update:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@docker compose pull
	@docker compose down syssvr
	@docker compose down unssvr
	@docker compose down frontend
	@docker compose up -d
    @docker compose restart kong
