# Nav Panel - 快速部署

基于 Docker Hub 公开镜像的一键部署方案。

## 🚀 快速开始

### 方式 1：自动部署脚本

```bash
# 下载部署文件
curl -O https://raw.githubusercontent.com/your-repo/nav-panel/main/deploy/docker-compose.yml
curl -O https://raw.githubusercontent.com/your-repo/nav-panel/main/deploy/deploy.sh

# 执行部署
chmod +x deploy.sh
./deploy.sh
```

### 方式 2：手动部署

```bash
# 下载配置文件
curl -O https://raw.githubusercontent.com/your-repo/nav-panel/main/deploy/docker-compose.yml

# 启动服务
docker compose up -d
```

### 方式 3：直接运行

```bash
# 克隆此目录
git clone https://github.com/your-repo/nav-panel.git
cd nav-panel/deploy

# 启动服务
docker compose up -d
```

## 📦 使用的镜像

- **前端**: `herio66/nav-panel-frontend:latest`
- **后端**: `herio66/nav-panel-backend:latest`

## 🔧 访问应用

部署完成后，访问：http://localhost:8080

## 🧪 测试部署

```bash
# 下载并运行测试脚本
curl -O https://raw.githubusercontent.com/your-repo/nav-panel/main/deploy/test.sh
chmod +x test.sh
./test.sh
```

测试脚本会验证：

- ✅ 服务容器状态
- ✅ 后端健康检查接口
- ✅ 前端页面访问

## 📋 管理命令

```bash
# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 更新镜像
docker compose pull
docker compose up -d

# 完全清理（包括数据）
docker compose down -v
```

## 💾 数据持久化

应用数据存储在 Docker 卷 `nav-panel-data` 中，即使删除容器，数据也会保留。

## 🆘 故障排除

### 端口冲突

如果 8080 端口被占用，修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - '8888:80' # 改为其他端口
```

### 服务无法启动

检查 Docker 服务是否运行：

```bash
docker --version
docker compose --version
```

### 查看详细日志

```bash
docker compose logs backend
docker compose logs frontend
```
