.PHONY: up down restart logs ps clean kong-init kong-migrate help

# 启动所有服务
up:
	docker-compose up -d
	@echo "所有服务已启动！"
	@echo "访问地址："
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - EMQX: http://localhost:18083"
	@echo "  - Kong Proxy: http://localhost:8000"
	@echo "  - Kong Admin API: http://localhost:8001"
	@echo "  - Kong Manager: http://localhost:8002"
	@echo "  - Konga UI: http://localhost:1337"

# 停止所有服务
down:
	docker-compose down
	@echo "所有服务已停止！"

# 重启所有服务
restart:
	docker-compose restart
	@echo "所有服务已重启！"

# 查看服务日志
logs:
	docker-compose logs -f

# 查看特定服务日志
logs-kong:
	docker-compose logs -f kong

logs-konga:
	docker-compose logs -f konga

logs-postgres:
	docker-compose logs -f postgres

logs-emqx:
	docker-compose logs -f emqx

# 查看服务状态
ps:
	docker-compose ps

# 清理所有容器和数据卷
clean:
	docker-compose down -v
	@echo "已清理所有容器和数据卷！"

# 初始化 Kong 数据库
kong-migrate:
	docker-compose run --rm kong kong migrations bootstrap
	@echo "Kong 数据库已初始化！"

# 升级 Kong 数据库
kong-migrate-up:
	docker-compose run --rm kong kong migrations up
	@echo "Kong 数据库已升级！"

# 重置 Kong 数据库
kong-reset:
	docker-compose run --rm kong kong migrations reset --yes
	@echo "Kong 数据库已重置！"

# 检查 Kong 配置
kong-check:
	docker exec kong kong check

# 重新加载 Kong 配置
kong-reload:
	docker exec kong kong reload
	@echo "Kong 配置已重新加载！"

# 帮助信息
help:
	@echo "可用命令："
	@echo "  make up              - 启动所有服务"
	@echo "  make down            - 停止所有服务"
	@echo "  make restart         - 重启所有服务"
	@echo "  make logs            - 查看所有服务日志"
	@echo "  make logs-kong       - 查看 Kong 日志"
	@echo "  make logs-konga      - 查看 Konga 日志"
	@echo "  make logs-postgres   - 查看 PostgreSQL 日志"
	@echo "  make logs-emqx       - 查看 EMQX 日志"
	@echo "  make ps              - 查看服务状态"
	@echo "  make clean           - 清理所有容器和数据卷"
	@echo "  make kong-migrate    - 初始化 Kong 数据库"
	@echo "  make kong-migrate-up - 升级 Kong 数据库"
	@echo "  make kong-reset      - 重置 Kong 数据库"
	@echo "  make kong-check      - 检查 Kong 配置"
	@echo "  make kong-reload     - 重新加载 Kong 配置"
	@echo "  make help            - 显示此帮助信息"

