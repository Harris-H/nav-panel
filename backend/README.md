# Nav Panel Backend

基于 Go + Gin + SQLite 的导航面板后端 API 服务。

## 项目结构

```
backend/
├── cmd/
│   └── server/          # 应用程序入口
├── internal/
│   ├── database/        # 数据库连接和迁移
│   ├── handler/         # HTTP 处理器层
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
├── data/               # SQLite 数据库文件目录
├── go.mod              # Go 模块依赖
└── README.md           # 项目说明
```

## 技术栈

- **Go 1.21+** - 后端语言
- **Gin** - Web 框架
- **SQLite** - 数据库
- **CORS** - 跨域支持

## API 接口

### 网站管理

- `GET /api/websites` - 获取所有网站
- `POST /api/websites` - 创建网站
- `PUT /api/websites/:id` - 更新网站
- `DELETE /api/websites/:id` - 删除网站

### 搜索引擎管理

- `GET /api/search-engines` - 获取所有搜索引擎
- `POST /api/search-engines` - 创建搜索引擎
- `PUT /api/search-engines/:id` - 更新搜索引擎
- `DELETE /api/search-engines/:id` - 删除搜索引擎

### 应用设置

- `GET /api/settings` - 获取应用设置
- `PUT /api/settings` - 更新应用设置

### 数据导入导出

- `GET /api/export` - 导出所有数据
- `POST /api/import` - 导入数据

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行开发服务器

```bash
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

### 3. 构建生产版本

```bash
go build -o nav-panel-backend cmd/server/main.go
```

## 数据库

项目使用 SQLite 作为数据库，数据库文件会自动创建在 `data/nav-panel.db`。

数据库表结构：

- `websites` - 网站信息
- `search_engines` - 搜索引擎配置
- `app_settings` - 应用设置

## 配置

默认配置：

- 端口：8080
- 数据库：data/nav-panel.db
- CORS：支持 localhost:3000 和 localhost:5173

## 部署

1. 构建可执行文件
2. 复制可执行文件到服务器
3. 确保有写入权限创建 data 目录
4. 运行可执行文件

```bash
./nav-panel-backend
```

## 开发说明

项目采用分层架构：

- **Handler** - 处理 HTTP 请求和响应
- **Service** - 业务逻辑处理
- **Repository** - 数据访问操作
- **Model** - 数据模型定义

这种架构便于测试、维护和扩展。
