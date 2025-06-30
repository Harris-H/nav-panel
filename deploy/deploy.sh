#!/bin/bash

# Nav Panel Docker Hub 部署脚本

set -e

echo "======================================"
echo "  Nav Panel Docker Hub 部署"
echo "======================================"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    exit 1
fi

# 拉取最新镜像
echo "拉取最新镜像..."
docker compose pull

# 停止旧服务
echo "停止旧服务..."
docker compose down

# 启动新服务
echo "启动新服务..."
docker compose up -d

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
docker compose ps

echo ""
echo "部署完成！"
echo "镜像来源: Docker Hub (herio66)"
echo "前端: herio66/nav-panel-frontend:latest"
echo "后端: herio66/nav-panel-backend:latest"
echo "访问地址: http://localhost:8080" 