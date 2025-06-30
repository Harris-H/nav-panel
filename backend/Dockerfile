# 第一阶段：构建阶段
FROM golang:1.21-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 第二阶段：运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates sqlite wget

# 创建非 root 用户
RUN adduser -D -g '' appuser

# 设置工作目录
WORKDIR /app

# 创建数据目录
RUN mkdir -p /app/data && chown appuser:appuser /app/data

# 复制构建结果
COPY --from=builder /app/main .

# 复制迁移文件（如果存在）
COPY --from=builder /app/migrations ./migrations

# 设置文件权限
RUN chown appuser:appuser /app/main

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider -O /dev/null http://localhost:8080/api/ping || exit 1

# 启动应用
CMD ["./main"] 