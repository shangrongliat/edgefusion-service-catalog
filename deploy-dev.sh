#!/bin/bash

# 设置编译参数，如果需要的话
export GOOS=linux
#export GOARCH=amd64
export GOARCH=arm64
REMOTE_USER="root"
REMOTE_IPS=("172.16.100.111" "172.16.100.112" "172.16.100.91" "172.16.100.92" "172.16.100.93")
REMOTE_PATH="/home/work/catalog/"
PROJECT_NAME="catalog"
# 项目根目录
PROJECT_DIR=$(pwd)

echo "rm local file......."
rm -rf $PROJECT_DIR/build/edgefusion*
sleep 1
# 编译Go项目
echo "build $PROJECT_NAME file......."
go build -o build/$PROJECT_NAME main.go

# 检查编译是否成功
if [ $? -eq 0 ]; then
    echo "Compilation successful."
else
    echo "Compilation failed."
    exit 1
fi


# 打包可执行文件和配置文件
#echo "Packaging binary and config files..."
#tar -czf "$PROJECT_DIR/build/$PROJECT_NAME.tar.gz" -C "$PROJECT_DIR/build" $PROJECT_NAME -C "$PROJECT_DIR" etc/* -C "$PROJECT_DIR/" start.sh -C "$PROJECT_DIR/" stop.sh
# 遍历所有远程IP地址
for REMOTE_IP in "${REMOTE_IPS[@]}"; do
  # 删除远程服务器上的旧包（需要先登录远程服务器）
  ssh  $REMOTE_USER@$REMOTE_IP "rm -rf $REMOTE_PATH$PROJECT_NAME"

  # 上传并解压到远程服务器
  echo "Uploading and extracting package to remote server..."
  scp "${PROJECT_DIR}/build/$PROJECT_NAME" $REMOTE_USER@$REMOTE_IP:$REMOTE_PATH
done

echo "Deployment completed successfully."