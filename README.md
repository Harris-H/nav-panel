# 导航面板 - 全栈导航站管理系统

一个基于 Vue3 + Go 的现代化全栈导航站管理系统，支持多搜索引擎集成和完整的网站管理功能。

## ✨ 特性

- 🏗️ **前后端分离架构** - Vue3 前端 + Go 后端，现代化技术栈
- 🔍 **多搜索引擎支持** - 集成百度、Google、必应等多个搜索引擎
- 🌐 **完整网站管理** - 支持网站的增删改查和分类管理
- 🎨 **优雅的界面设计** - 响应式卡片布局，适配所有设备
- 🚀 **一键部署** - Docker Compose 快速部署，开箱即用
- 💾 **数据持久化** - 后端数据库存储，数据安全可靠
- ⚡ **高性能** - Go 后端高并发处理，Vue3 前端快速响应

## 🚀 快速部署

### 使用 Docker Compose（推荐）

这是最简单快速的部署方式：

```bash
# 克隆项目
git clone <repository-url>
cd nav-panel

# 使用部署脚本快速启动
cd deploy
chmod +x deploy.sh
./deploy.sh
```

或者手动执行：

```bash
# 使用 Docker Compose 启动服务
docker compose up -d

# 查看服务状态
docker compose ps
```

访问 `http://localhost:8080` 即可使用。

### 开发环境部署

#### 后端开发

```bash
cd backend

# 安装 Go 依赖
go mod tidy

# 启动后端服务
go run cmd/server/main.go
```

#### 前端开发

```bash
# 安装前端依赖
npm install

# 启动开发服务器
npm run dev
```

## 📁 项目结构

```
nav-panel/
├── backend/                    # Go 后端
│   ├── cmd/server/            # 服务器启动入口
│   ├── internal/              # 内部业务逻辑
│   │   ├── handler/           # HTTP 处理器
│   │   ├── service/           # 业务服务层
│   │   ├── repository/        # 数据访问层
│   │   ├── model/             # 数据模型
│   │   └── database/          # 数据库配置
│   ├── migrations/            # 数据库迁移文件
│   └── go.mod                 # Go 模块文件
├── src/                       # Vue3 前端
│   ├── views/                 # 页面组件
│   ├── components/            # 通用组件
│   ├── stores/                # 状态管理
│   ├── utils/                 # 工具函数
│   └── types/                 # TypeScript 类型
├── deploy/                    # 部署相关
│   ├── deploy.sh              # 快速部署脚本
│   └── docker-compose.yml     # Docker 编排文件
├── scripts/                   # 构建脚本
└── nginx.conf                 # Nginx 配置
```

## 🎯 功能模块

### 网站管理

- **网站CRUD**: 完整的网站增删改查功能
- **分类管理**: 支持网站分类和标签
- **图标支持**: 自动获取网站图标或上传自定义图标
- **批量操作**: 支持批量导入导出网站数据

### 搜索引擎集成

- **多引擎支持**: 百度、Google、必应、搜狗等
- **快速切换**: 一键切换不同搜索引擎
- **自定义引擎**: 支持添加自定义搜索引擎
- **搜索历史**: 保存搜索历史记录

### 系统设置

- **主题定制**: 多种主题颜色和背景选择
- **布局配置**: 可调整卡片大小和排列方式
- **用户设置**: 个性化配置选项

## 🛠️ 技术栈

### 前端技术

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript 超集
- **Vite** - 极速前端构建工具
- **Pinia** - Vue3 状态管理库
- **Vue Router** - 官方路由管理器

### 后端技术

- **Go** - 高性能编程语言
- **Gin** - 轻量级 Web 框架
- **GORM** - Go ORM 库
- **SQLite/MySQL** - 数据库支持

### 部署技术

- **Docker** - 容器化部署
- **Docker Compose** - 服务编排
- **Nginx** - 反向代理服务器

## 🌟 项目优势

### 技术优势

- **现代化架构**: 采用前后端分离设计，技术栈新颖且主流
- **高性能**: Go 后端提供卓越的并发处理能力
- **类型安全**: TypeScript 确保前端代码质量
- **容器化**: Docker 确保环境一致性和易部署性

### 功能优势

- **一站式解决方案**: 集成搜索引擎和网站管理功能
- **用户体验优先**: 响应式设计，操作简洁流畅
- **数据安全**: 后端数据库存储，支持数据备份
- **扩展性强**: 模块化设计，易于功能扩展

### 部署优势

- **一键部署**: Docker Compose 配置文件，部署过程自动化
- **环境隔离**: 容器化确保不同环境的一致性
- **易于维护**: 统一的部署脚本和配置管理

## 📱 响应式支持

- **桌面端**: 最佳的浏览体验和功能完整性
- **平板端**: 优化的触控交互体验
- **移动端**: 完整功能的移动适配

## 🔧 配置说明

### 环境变量

复制 `env.example` 到 `.env` 并根据需要修改配置：

```bash
cp env.example .env
```

### 数据库配置

支持 SQLite（默认）和 MySQL，可在配置文件中切换。

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 📞 联系

如有问题或建议，请通过 GitHub Issues 联系。
