#!/bin/bash
# -------------------------------------------------------------
# 自动创建 Docker TLS 证书（脚本预置密码版）
# 修改说明：移除手动输入密码步骤，改为脚本内指定密码
# -------------------------------------------------------------

# ====== 配置区域（修改以下变量即可）====== #
HOST="$1"           # 服务器IP或域名
PASSWORD="Supos1304@" # 预置密码（特殊字符+数字+字母）
# -------------------------------------- #

# 其他配置（一般无需修改）
COUNTRY="CN"
STATE="ZheJiang"
CITY="HangZhou"
ORGANIZATION="LZ"
ORGANIZATIONAL_UNIT="Dev"
COMMON_NAME="$HOST"
EMAIL="admin@example.com"
NOW_PATH="/etc/docker/certs/"
mkdir -p "$NOW_PATH"
cd "$NOW_PATH" || exit

# ====== 证书生成流程 ====== #
# 1. 生成CA私钥（使用预置密码）
openssl genrsa -aes256 -passout "pass:$PASSWORD" -out "ca-key.pem" 4096

# 2. 生成CA根证书
openssl req -new -x509 -days 36500 -key "ca-key.pem" -sha256 -out "ca.pem" -passin "pass:$PASSWORD" \
  -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORGANIZATION/OU=$ORGANIZATIONAL_UNIT/CN=$COMMON_NAME/emailAddress=$EMAIL"

# 3. 生成服务端私钥
openssl genrsa -out "server-key.pem" 4096

# 4. 生成服务端证书签名请求
openssl req -subj "/CN=$COMMON_NAME" -sha256 -new -key "server-key.pem" -out server.csr

# 5. 配置证书扩展（含IP白名单）
cat > extfile.cnf <<EOF
subjectAltName = IP:127.0.0.1,DNS:host.docker.internal
extendedKeyUsage = serverAuth
EOF
openssl x509 -req -days 36500 -sha256 -in server.csr -passin "pass:$PASSWORD" -CA "ca.pem" -CAkey "ca-key.pem" \
  -CAcreateserial -out "server-cert.pem" -extfile extfile.cnf

# 6. 生成客户端证书
openssl genrsa -out "key.pem" 4096
openssl req -subj '/CN=client' -new -key "key.pem" -out client.csr
echo "extendedKeyUsage = clientAuth" > extfile-client.cnf
openssl x509 -req -days 36500 -sha256 -in client.csr -passin "pass:$PASSWORD" -CA "ca.pem" -CAkey "ca-key.pem" \
  -CAcreateserial -out "cert.pem" -extfile extfile-client.cnf

# ====== 清理与打包 ====== #
rm -f client.csr server.csr extfile*.cnf
chmod 0400 ca-key.pem key.pem server-key.pem     # 私钥严格权限
chmod 0444 ca.pem server-cert.pem cert.pem       # 证书只读权限
