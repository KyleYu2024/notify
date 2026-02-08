#!/bin/bash

# Docker 部署测试脚本

set -e

echo "=== Docker 部署测试 ==="

# 检查 Docker 和 docker-compose 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: docker-compose 未安装"
    exit 1
fi

echo "✓ Docker 和 docker-compose 已安装"



# 停止并清理现有容器
echo "清理现有容器..."
docker-compose down --remove-orphans

# 构建并启动服务
echo "构建并启动服务..."
docker-compose up -d --build

# 等待服务启动
echo "等待服务启动..."
sleep 30

# 检查容器状态
echo "检查容器状态..."
if ! docker-compose ps | grep -q "Up"; then
    echo "错误: 容器未正常启动"
    docker-compose logs
    exit 1
fi

echo "✓ 容器已启动"
