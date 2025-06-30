#!/bin/bash

# Nav Panel 部署测试脚本

echo "======================================"
echo "     Nav Panel 部署测试"
echo "======================================"

# 检查服务状态
echo "检查服务状态..."
docker compose ps

echo ""
echo "等待服务完全启动..."
sleep 5

# 测试健康检查接口
echo "测试后端健康检查..."
if curl -f -s http://localhost:8080/api/ping > /dev/null; then
    echo "✅ 后端健康检查 - 正常"
else
    echo "❌ 后端健康检查 - 失败"
    echo "后端日志："
    docker compose logs --tail=10 backend
fi

# 测试前端
echo "测试前端访问..."
if curl -f -s http://localhost:8080 > /dev/null; then
    echo "✅ 前端访问 - 正常"
else
    echo "❌ 前端访问 - 失败"
    echo "前端日志："
    docker compose logs --tail=10 frontend
fi

echo ""
echo "======================================"
echo "测试完成！"
echo ""
echo "如果所有测试都通过，你可以访问："
echo "🌐 http://localhost:8080"
echo ""
echo "管理命令："
echo "  查看状态: docker compose ps"
echo "  查看日志: docker compose logs -f"
echo "  停止服务: docker compose down"
echo "======================================" 