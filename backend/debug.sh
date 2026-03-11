#!/bin/bash

# --- 配置变量 ---
REMOTE_IP="192.168.235.152"
REMOTE_USER="root"               # 建议根据实际情况修改用户
REMOTE_PATH="/tmp"               # 远程服务器临时存放路径
CONTAINER_NAME="uns"             # 目标容器名称
CONTAINER_PATH="/app"            # 容器内目标路径
LOCAL_FILE="backend"

# 1. 编译 Go 程序
echo ">>> 正在编译 Go 程序..."
# 注意：在 Linux 环境下编译 Windows 风格路径需调整，这里使用标准路径
go build ./backend.go

if [ $? -ne 0 ]; then
    echo "❌ 编译失败，请检查代码！"
    exit 1
fi
echo "✅ 编译成功。"

# ssh-copy-id root@$REMOTE_IP

# 2. 将编译结果拷贝到远程服务器
echo ">>> 正在发送文件到 $REMOTE_IP..."
scp ./$LOCAL_FILE $REMOTE_USER@$REMOTE_IP:$REMOTE_PATH/

if [ $? -ne 0 ]; then
    echo "❌ 文件拷贝失败，请确认是否已配置 SSH 免密登录。"
    exit 1
fi

# 3, 4, 5. 远程执行命令 (拷贝到容器 -> 重启 -> 查看日志)
echo ">>> 正在远程部署并重启容器..."
ssh $REMOTE_USER@$REMOTE_IP << EOF
    # 将文件从宿主机拷贝到容器
    docker cp $REMOTE_PATH/$LOCAL_FILE $CONTAINER_NAME:$CONTAINER_PATH/

    # 修改权限（确保容器内可执行）
    docker exec $CONTAINER_NAME chmod +x $CONTAINER_PATH/$LOCAL_FILE

    # 重启容器
    echo ">>> 重启容器 $CONTAINER_NAME..."
    docker restart $CONTAINER_NAME

    # 查看日志（显示最后 50 行并持续跟踪）
    echo ">>> 正在查看日志 (Ctrl+C 退出)..."
    docker logs -f --tail 50 $CONTAINER_NAME |grep [app]
EOF