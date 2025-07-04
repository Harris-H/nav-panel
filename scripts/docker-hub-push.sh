#!/bin/bash

# Nav Panel Docker Hub 推送脚本
# 使用方法: ./scripts/docker-hub-push.sh [username] [tag]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Nav Panel Docker Hub 推送脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [选项] [username] [tag]"
    echo ""
    echo "选项:"
    echo "  -u, --use-existing   使用现有镜像，跳过重新构建"
    echo "  -h, --help           显示此帮助信息"
    echo ""
    echo "配置文件:"
    echo "  创建 .docker-config 文件设置默认用户名和密码："
    echo "  DOCKER_USERNAME=your-username"
    echo "  DOCKER_PASSWORD=your-password"
    echo ""
    echo "示例:"
    echo "  $0 myusername v1.0.0"
    echo "  $0 --use-existing myusername latest"
    echo "  $0  # 使用配置文件中的用户名，默认tag为latest"
}

# 构建选项
USE_EXISTING="false"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --use-existing|-u)
            USE_EXISTING="true"
            shift
            ;;
        --help|-h|help)
            show_help
            exit 0
            ;;
        *)
            # 位置参数
            if [ -z "$DOCKER_USERNAME" ]; then
                DOCKER_USERNAME="$1"
            elif [ -z "$TAG" ]; then
                TAG="$1"
            fi
            shift
            ;;
    esac
done

# 加载配置文件
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/.docker-config"

if [ -f "$CONFIG_FILE" ]; then
    log_info "加载配置文件: $CONFIG_FILE"
    
    # 检查并修复Windows换行符问题
    if grep -q $'\r' "$CONFIG_FILE" 2>/dev/null; then
        log_warning "检测到Windows换行符，正在自动修复..."
        # 创建临时文件，去除\r
        sed 's/\r$//' "$CONFIG_FILE" > "${CONFIG_FILE}.tmp" && mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
        log_success "换行符已修复"
    fi
    
    source "$CONFIG_FILE"
fi

# 默认参数
TAG=${TAG:-"latest"}

# 检查必要参数
if [ -z "$DOCKER_USERNAME" ]; then
    log_error "未指定Docker Hub用户名"
    echo "请使用以下方式之一指定用户名："
    echo "1. 命令行参数: $0 username"
    echo "2. 配置文件: 在 $CONFIG_FILE 中设置 DOCKER_USERNAME=your-username"
    exit 1
fi

# 镜像名称
FRONTEND_IMAGE="nav-panel-frontend"
BACKEND_IMAGE="nav-panel-backend"
FRONTEND_TAG="${DOCKER_USERNAME}/${FRONTEND_IMAGE}:${TAG}"
BACKEND_TAG="${DOCKER_USERNAME}/${BACKEND_IMAGE}:${TAG}"

# 检查镜像是否存在
check_image_exists() {
    local image_name=$1
    if docker image inspect "${image_name}" >/dev/null 2>&1; then
        return 0  # 镜像存在
    else
        return 1  # 镜像不存在
    fi
}

# 验证本地镜像
validate_local_images() {
    log_info "验证本地镜像..."
    
    local missing_images=()
    
    if ! check_image_exists "${FRONTEND_IMAGE}:${TAG}"; then
        missing_images+=("${FRONTEND_IMAGE}:${TAG}")
    fi
    
    if ! check_image_exists "${BACKEND_IMAGE}:${TAG}"; then
        missing_images+=("${BACKEND_IMAGE}:${TAG}")
    fi
    
    if [ ${#missing_images[@]} -gt 0 ]; then
        log_error "以下镜像不存在:"
        for img in "${missing_images[@]}"; do
            echo "  ❌ $img"
        done
        echo ""
        log_error "使用现有镜像模式下无法继续，请先构建镜像或移除 --use-existing 选项"
        exit 1
    fi
    
    log_success "所有本地镜像验证通过"
    return 0  # 不需要构建
}

# 构建镜像
build_images() {
    log_info "开始构建镜像..."
    
    # 检查项目根目录
    if [ ! -f "package.json" ] || [ ! -d "backend" ]; then
        log_error "请在项目根目录下运行此脚本"
        exit 1
    fi
    
    # 构建前端镜像
    log_info "构建前端镜像: ${FRONTEND_IMAGE}:${TAG}"
    docker build -f Dockerfile.frontend -t ${FRONTEND_IMAGE}:${TAG} .
    
    # 构建后端镜像
    log_info "构建后端镜像: ${BACKEND_IMAGE}:${TAG}"
    docker build -f backend/Dockerfile -t ${BACKEND_IMAGE}:${TAG} ./backend
    
    log_success "镜像构建完成"
}

# Docker Hub 登录
docker_login() {
    log_info "登录 Docker Hub..."
    
    if [ -n "$DOCKER_PASSWORD" ]; then
        # 使用配置文件中的密码
        echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
    else
        # 交互式登录
        docker login -u "$DOCKER_USERNAME"
    fi
    
    log_success "Docker Hub 登录成功"
}

# 标记镜像
tag_images() {
    log_info "标记镜像..."
    
    # 标记前端镜像
    docker tag ${FRONTEND_IMAGE}:${TAG} ${FRONTEND_TAG}
    log_info "前端镜像标记为: ${FRONTEND_TAG}"
    
    # 标记后端镜像  
    docker tag ${BACKEND_IMAGE}:${TAG} ${BACKEND_TAG}
    log_info "后端镜像标记为: ${BACKEND_TAG}"
    
    log_success "镜像标记完成"
}

# 推送单个镜像（带重试）
push_image_with_retry() {
    local image_tag=$1
    local max_retries=3
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        log_info "推送镜像: ${image_tag} (尝试 $((retry_count + 1))/${max_retries})"
        
        if docker push "${image_tag}"; then
            log_success "镜像推送成功: ${image_tag}"
            return 0
        else
            retry_count=$((retry_count + 1))
            if [ $retry_count -lt $max_retries ]; then
                log_warning "推送失败，5秒后重试..."
                sleep 5
            else
                log_error "镜像推送失败: ${image_tag}"
                return 1
            fi
        fi
    done
}

# 推送镜像
push_images() {
    log_info "推送镜像到 Docker Hub..."
    
    # 推送前端镜像（带重试）
    if ! push_image_with_retry "${FRONTEND_TAG}"; then
        log_error "前端镜像推送失败"
        exit 1
    fi
    
    # 推送后端镜像（带重试）
    if ! push_image_with_retry "${BACKEND_TAG}"; then
        log_error "后端镜像推送失败"
        exit 1
    fi
    
    log_success "所有镜像推送完成"
}

# 更新镜像版本（跨平台兼容）
update_image_versions() {
    local compose_file=$1
    
    # 尝试使用 sed (Linux/macOS)
    if command -v sed >/dev/null 2>&1; then
        # 创建备份
        cp "$compose_file" "${compose_file}.bak"
        
        # 更新镜像版本
        sed -i.tmp "s|image: .*/nav-panel-frontend:.*|image: ${FRONTEND_TAG}|g" "$compose_file" 2>/dev/null && \
        sed -i.tmp "s|image: .*/nav-panel-backend:.*|image: ${BACKEND_TAG}|g" "$compose_file" 2>/dev/null
        
        # 清理临时文件
        rm -f "${compose_file}.tmp" 2>/dev/null
        
        return 0
    fi
    
    # 尝试使用 PowerShell (Windows)
    if command -v powershell >/dev/null 2>&1; then
        powershell -Command "
            \$content = Get-Content '$compose_file' -Raw
            \$content = \$content -replace 'image: .*/nav-panel-frontend:.*', 'image: $FRONTEND_TAG'
            \$content = \$content -replace 'image: .*/nav-panel-backend:.*', 'image: $BACKEND_TAG'
            Set-Content '$compose_file' -Value \$content -NoNewline
        " 2>/dev/null
        
        return $?
    fi
    
    # 尝试使用 Python (如果可用)
    if command -v python3 >/dev/null 2>&1 || command -v python >/dev/null 2>&1; then
        local python_cmd="python3"
        if ! command -v python3 >/dev/null 2>&1; then
            python_cmd="python"
        fi
        
        $python_cmd -c "
import re
import sys

try:
    with open('$compose_file', 'r') as f:
        content = f.read()
    
    content = re.sub(r'image: .*/nav-panel-frontend:.*', 'image: $FRONTEND_TAG', content)
    content = re.sub(r'image: .*/nav-panel-backend:.*', 'image: $BACKEND_TAG', content)
    
    with open('$compose_file', 'w') as f:
        f.write(content)
    
    sys.exit(0)
except Exception:
    sys.exit(1)
" 2>/dev/null
        
        return $?
    fi
    
    # 都不可用，返回失败
    return 1
}

# 更新或生成 docker-compose 文件
update_compose_file() {
    mkdir -p deploy
    
    local compose_file="deploy/docker-compose.yml"
    
    if [ -f "$compose_file" ]; then
        log_info "更新现有的 docker-compose.yml 文件..."
        
        # 更新镜像版本（跨平台兼容）
        if update_image_versions "$compose_file"; then
            log_success "docker-compose.yml 镜像版本已更新"
        else
            log_warning "更新失败，将重新生成文件"
            generate_new_compose_file "$compose_file"
        fi
    else
        log_info "生成新的 docker-compose.yml 文件..."
        generate_new_compose_file "$compose_file"
    fi
}

# 生成新的 docker-compose 文件
generate_new_compose_file() {
    local compose_file=$1
    
    cat > "$compose_file" << EOF
services:
  frontend:
    image: ${FRONTEND_TAG}
    container_name: nav-panel-frontend
    restart: unless-stopped
    ports:
      - "8080:80"
    depends_on:
      backend:
        condition: service_healthy
    networks:
      - nav-panel-network

  backend:
    image: ${BACKEND_TAG}
    container_name: nav-panel-backend
    restart: unless-stopped
    expose:
      - "8080"
    volumes:
      - nav-panel-data:/app/data
    networks:
      - nav-panel-network
    environment:
      - GIN_MODE=release
      - CGO_ENABLED=1
    healthcheck:
      test:
        [
          'CMD',
          'wget',
          '--no-verbose',
          '--tries=1',
          '-O',
          '/dev/null',
          'http://localhost:8080/api/ping',
        ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

volumes:
  nav-panel-data:
    driver: local

networks:
  nav-panel-network:
    driver: bridge
EOF
    
    log_success "docker-compose.yml 已生成到 deploy/ 目录"
}

# 显示完成信息
show_completion() {
    echo ""
    log_success "推送完成！"
    echo ""
    echo "🐳 Docker Hub 镜像:"
    echo "  前端: ${FRONTEND_TAG}"
    echo "  后端: ${BACKEND_TAG}"
    echo ""
    echo "📁 部署文件:"
    echo "  deploy/docker-compose.yml"
    echo ""
    echo "🚀 使用方式:"
    echo "  docker compose -f deploy/docker-compose.yml up -d"
    echo ""
    echo "🔗 Docker Hub 链接:"
    echo "  https://hub.docker.com/r/${DOCKER_USERNAME}/${FRONTEND_IMAGE}"
    echo "  https://hub.docker.com/r/${DOCKER_USERNAME}/${BACKEND_IMAGE}"
}

# 主函数
main() {
    echo "======================================"
    echo "   Nav Panel Docker Hub 推送"
    echo "======================================"
    echo ""
    echo "用户名: ${DOCKER_USERNAME}"
    echo "标签: ${TAG}"
    echo "使用现有镜像: ${USE_EXISTING}"
    echo ""
    
    # 检查是否需要构建
    if [ "$USE_EXISTING" = "true" ]; then
        log_info "使用现有镜像模式"
        validate_local_images
    else
        log_info "重新构建镜像模式"
        build_images
    fi
    
    docker_login
    tag_images
    push_images
    update_compose_file
    show_completion
}

# 执行主函数
main 