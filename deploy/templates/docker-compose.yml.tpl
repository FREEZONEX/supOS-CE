services:
  backend:
    image: ${BACKEND_IMAGE:?BACKEND_IMAGE is required}
    container_name: ${COMPOSE_PROJECT_NAME:-edge}-backend
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "${BACKEND_DOCKER_LOG_MAX_SIZE:-500m}"
        max-file: "${BACKEND_DOCKER_LOG_MAX_FILE:-7}"
    healthcheck:
      test: ["CMD-SHELL", "curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null"]
      interval: 5s
      timeout: 4s
      retries: 24
      start_period: 20s
    ports:
      - "${ENTRANCE_PORT:-8088}:8080"
    environment:
      UNS_DB_URL: ${UNS_DB_URL}
      SINK_DB_URL: ${SINK_DB_URL:-}
      TIMESERIES_RETENTION_YEARS: ${TIMESERIES_RETENTION_YEARS:-10}
      REDIS_ADDR: ${REDIS_ADDR:-redis:6379}
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      # set ZLM_STREAM_HOST to this host's LAN address for remote workers.
      ADMIN_INITIAL_PASSWORD: ${ADMIN_INITIAL_PASSWORD:-tier0}
      PRODUCT_VERSION: ${PRODUCT_VERSION:?PRODUCT_VERSION is required}
      LANGUAGE: ${LANGUAGE:-en-US}
      COMPOSE_PROJECT_NAME: ${COMPOSE_PROJECT_NAME:-edge}
      WEB_DIR: /app/web
      BACKEND_LOG_MODE: ${BACKEND_LOG_MODE:-file}
      BACKEND_LOG_ENCODING: ${BACKEND_LOG_ENCODING:-json}
      BACKEND_LOG_PATH: ${BACKEND_LOG_PATH:-/app/logs}
      BACKEND_LOG_LEVEL: ${BACKEND_LOG_LEVEL:-info}
      BACKEND_LOG_ROTATION: ${BACKEND_LOG_ROTATION:-daily}
      BACKEND_LOG_KEEP_DAYS: ${BACKEND_LOG_KEEP_DAYS:-14}
      BACKEND_LOG_COMPRESS: ${BACKEND_LOG_COMPRESS:-true}
      LOCAL_FRONTEND_DEV: ${LOCAL_FRONTEND_DEV:-false}
      FRONTEND_DEV_PROXY_URL: ${FRONTEND_DEV_PROXY_URL:-}
      FILESTORE_LOCAL_ROOT: /app/data/files
      SOURCEFLOW_URL: http://sourceflow:1880
      EVENTFLOW_URL: http://eventflow:1880
      EMQX_URL: http://emqx:18083
      TIER0_SDK_API_HOST: ${TIER0_SDK_API_HOST:-backend:8080}
      TIER0_SDK_MQTT_HOST: ${TIER0_SDK_MQTT_HOST:-}
      TIER0_SDK_MQTT_PORT: ${TIER0_SDK_MQTT_PORT:-}
      TIER0_API_KEY: ${TIER0_API_KEY:-}
      HOST_METRICS_ENABLED: ${HOST_METRICS_ENABLED:-false}
      DATAINGEST_ENABLED: ${DATAINGEST_ENABLED:-true}
      DATAINGEST_MQTT_BROKERS: ${DATAINGEST_MQTT_BROKERS:-tcp://emqx:1883}
      DATAINGEST_MQTT_CLIENT_ID: ${DATAINGEST_MQTT_CLIENT_ID:-backend-sink}
      DATAINGEST_MQTT_TOPIC: ${DATAINGEST_MQTT_TOPIC:-#}
      EMQX_API_KEY: ${EMQX_API_KEY:-}
      EMQX_API_SECRET: ${EMQX_API_SECRET:-}
      NODERED_INTERNAL_TOKEN: ${NODERED_INTERNAL_TOKEN:-}
      TIER0_INSTALLATION_ID: ${TIER0_INSTALLATION_ID:?TIER0_INSTALLATION_ID is required}
      DATAINGEST_QUEUE_SIZE: ${DATAINGEST_QUEUE_SIZE:-20000}
      DATAINGEST_BATCH_SIZE: ${DATAINGEST_BATCH_SIZE:-5000}
      DATAINGEST_FLUSH_INTERVAL_MS: ${DATAINGEST_FLUSH_INTERVAL_MS:-1000}
      EMQX_EXTERNAL_BROKER_HOST: ${EMQX_EXTERNAL_BROKER_HOST:-}
      ENTRANCE_DOMAIN: ${ENTRANCE_DOMAIN:-}
      OS_MQTT_TCP_PORT: ${OS_MQTT_TCP_PORT:-1883}
      OS_MQTT_WEBSOCKET_PORT: ${OS_MQTT_WEBSOCKET_PORT:-8083}
      OS_MQTT_WEBSOCKET_TSL_PORT: ${OS_MQTT_WEBSOCKET_TSL_PORT:-8084}
    depends_on:
      emqx:
        condition: service_healthy
    networks:
      - runtime
    volumes:
      - ${VOLUMES_PATH}/certs/tls:/app/certs:ro
      - ${VOLUMES_PATH}/backend/files:/app/data/files
      - ${VOLUMES_PATH}/backend/logs:/app/logs
      - ${HOST_PROC_PATH:-/proc}:/host/proc:ro
      - ${HOST_SYS_PATH:-/sys}:/host/sys:ro

  tsdb:
    image: harbor.tier0.dev/library/timescaledb:2.20.0-pg17-pgvector0.8.5
    container_name: ${COMPOSE_PROJECT_NAME:-edge}-tsdb
    profiles:
      - local-db
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "${TSDB_DOCKER_LOG_MAX_SIZE:-100m}"
        max-file: "${TSDB_DOCKER_LOG_MAX_FILE:-5}"
    ports:
      - "${TSDB_PORT:-5433}:5432"
    environment:
      TZ: UTC
      service_logo: postgresql-original.svg
      service_description: aboutus.postgresqlDescription
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${TSDB_PASSWORD}
      POSTGRES_DB: ${UNS_DB_NAME:-postgres}
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
    volumes:
      - ${VOLUMES_PATH}/tsdb/conf/postgresql.conf:/etc/postgresql/custom.conf
      - ${VOLUMES_PATH}/tsdb/data:/var/lib/postgresql/data
    command:
      - postgres
      - -c
      - config_file=/etc/postgresql/custom.conf
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 30s
    networks:
      - runtime

  redis:
    image: harbor.tier0.dev/library/redis:7-alpine
    container_name: ${COMPOSE_PROJECT_NAME:-edge}-redis
    profiles:
      - local-redis
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "${REDIS_DOCKER_LOG_MAX_SIZE:-100m}"
        max-file: "${REDIS_DOCKER_LOG_MAX_FILE:-5}"
    command: ["redis-server", "--appendonly", "yes", "--requirepass", "${REDIS_PASSWORD}"]
    volumes:
      - ${VOLUMES_PATH}/redis/data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 5s
      timeout: 3s
      retries: 20
    networks:
      - runtime

  emqx:
    image: harbor.tier0.dev/library/emqx:5.8
    container_name: ${COMPOSE_PROJECT_NAME:-edge}-emqx
    restart: ${EMQX_RESTART_POLICY:-unless-stopped}
    logging:
      driver: "json-file"
      options:
        max-size: "${EMQX_DOCKER_LOG_MAX_SIZE:-200m}"
        max-file: "${EMQX_DOCKER_LOG_MAX_FILE:-5}"
    ports:
      - "${OS_MQTT_TCP_PORT:-1883}:1883"
      - "${OS_MQTT_SSL_PORT:-8883}:8883"
      - "${OS_MQTT_WEBSOCKET_PORT:-8083}:8083"
      - "${OS_MQTT_WEBSOCKET_TSL_PORT:-8084}:8084"
    environment:
      EMQX_DASHBOARD__DEFAULT_USERNAME: ${EMQX_DASHBOARD_USERNAME:-admin}
      EMQX_DASHBOARD__DEFAULT_PASSWORD: ${EMQX_DASHBOARD_PASSWORD:?EMQX_DASHBOARD_PASSWORD is required}
    volumes:
      - ${VOLUMES_PATH}/certs/tls:/opt/emqx/etc/certs:ro
      - ${VOLUMES_PATH}/emqx/data:/opt/emqx/data
      - ${VOLUMES_PATH}/emqx/log:/opt/emqx/log
      - ${VOLUMES_PATH}/emqx/config/default_api_key.conf:/opt/emqx/etc/default_api_key.conf:ro
      - ./mount/emqx/emqx.conf:/opt/emqx/etc/emqx.conf:ro
      - ./mount/emqx/acl.conf:/opt/emqx/etc/acl.conf:ro
    healthcheck:
      test: ["CMD", "emqx", "ctl", "status"]
      interval: 10s
      timeout: 5s
      retries: 12
      start_period: 60s
    networks:
      - runtime

  sourceflow:
    image: harbor.tier0.dev/library/node-red:4.1.10-22
    container_name: ${COMPOSE_PROJECT_NAME:-edge}-sourceflow
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "${SOURCEFLOW_DOCKER_LOG_MAX_SIZE:-200m}"
        max-file: "${SOURCEFLOW_DOCKER_LOG_MAX_FILE:-5}"
    user: root
    environment:
      PORT: 1880
      FLOWS: /data/flows.json
      OS_LANG: ${LANGUAGE:-en-US}
      NODE_OPTIONS: --openssl-legacy-provider
      NODE_RED_SEED_NAME: sourceflow
      NPM_CONFIG_CACHE: /data/.npm
      NODERED_INTERNAL_TOKEN: ${NODERED_INTERNAL_TOKEN:-}
    entrypoint: ["/bin/sh", "-lc"]
    command:
      - sh /usr/local/bin/seed-node-red.sh --quick sourceflow && exec npm start -- --userDir /data
    healthcheck:
      test: ["CMD-SHELL", "node -e \"const token=process.env.NODERED_INTERNAL_TOKEN||'';fetch('http://127.0.0.1:1880/',{headers:{'X-Tier0-Internal-Token':token}}).then(r=>process.exit(r.ok?0:1)).catch(()=>process.exit(1))\""]
      interval: 10s
      timeout: 5s
      retries: 12
      start_period: 30s
    volumes:
      - ${VOLUMES_PATH}/sourceflow:/data
      - ./mount/sourceflow:/opt/tier0-mount/sourceflow:ro
      - ./mount/node-red-init/seed-node-red.sh:/usr/local/bin/seed-node-red.sh:ro
    networks:
      runtime:
        aliases:
          - nodered

  eventflow:
    image: harbor.tier0.dev/library/node-red:4.1.10-22
    container_name: ${COMPOSE_PROJECT_NAME:-edge}-eventflow
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "${EVENTFLOW_DOCKER_LOG_MAX_SIZE:-200m}"
        max-file: "${EVENTFLOW_DOCKER_LOG_MAX_FILE:-5}"
    user: root
    environment:
      PORT: 1880
      FLOWS: /data/flows.json
      OS_LANG: ${LANGUAGE:-en-US}
      NODE_OPTIONS: --openssl-legacy-provider
      NODE_RED_SEED_NAME: eventflow
      NODE_RED_SEED_VERSION: 2026.07.10-001
      NPM_CONFIG_CACHE: /data/.npm
      NODERED_INTERNAL_TOKEN: ${NODERED_INTERNAL_TOKEN:-}
    entrypoint: ["/bin/sh", "-lc"]
    command:
      - sh /usr/local/bin/seed-node-red.sh eventflow && exec npm start -- --userDir /data
    healthcheck:
      test: ["CMD-SHELL", "node -e \"const token=process.env.NODERED_INTERNAL_TOKEN||'';fetch('http://127.0.0.1:1880/',{headers:{'X-Tier0-Internal-Token':token}}).then(r=>process.exit(r.ok?0:1)).catch(()=>process.exit(1))\""]
      interval: 10s
      timeout: 5s
      retries: 12
      start_period: 30s
    volumes:
      - ${VOLUMES_PATH}/eventflow:/data
      - ./mount/eventflow:/opt/tier0-mount/eventflow:ro
      - ./mount/node-red-init/seed-node-red.sh:/usr/local/bin/seed-node-red.sh:ro
    networks:
      - runtime

networks:
  runtime:
    name: ${COMPOSE_PROJECT_NAME:-edge}_runtime
