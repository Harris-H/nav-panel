# 前后端集成说明

## 概述

本项目已将前端从使用本地 localStorage 存储改为与后端 API 进行集成。现在所有数据（网站、设置、搜索引擎）都通过 API 与后端数据库进行同步。

## 主要变更

### 1. 移除的文件

- `src/utils/storage.ts` - 本地存储工具文件已删除

### 2. 新增的文件

- `src/utils/api.ts` - API 客户端，包含所有与后端通信的方法
- `src/components/ErrorMessage.vue` - 错误消息显示组件
- `src/components/LoadingSpinner.vue` - 加载状态指示器组件

### 3. 修改的文件

- `src/stores/app.ts` - 完全重写，使用 API 替代 localStorage
- `src/App.vue` - 添加错误和加载状态组件，应用启动时从 API 加载数据

## API 端点

### 网站管理

- `GET /api/websites` - 获取所有网站
- `POST /api/websites` - 创建新网站
- `PUT /api/websites/:id` - 更新网站
- `DELETE /api/websites/:id` - 删除网站

### 搜索引擎管理

- `GET /api/search-engines` - 获取所有搜索引擎
- `POST /api/search-engines` - 创建新搜索引擎
- `PUT /api/search-engines/:id` - 更新搜索引擎
- `DELETE /api/search-engines/:id` - 删除搜索引擎

### 设置管理

- `GET /api/settings` - 获取应用设置
- `PUT /api/settings` - 更新应用设置

### 数据导入导出

- `GET /api/export` - 导出所有数据
- `POST /api/import` - 导入数据

## 启动说明

### 1. 启动后端服务器

```bash
cd backend
go run cmd/server/main.go
```

后端服务器将在 `http://localhost:8080` 启动

### 2. 启动前端开发服务器

```bash
npm install
npm run dev
```

前端开发服务器将在 `http://localhost:5173` 启动

## 功能说明

### 自动数据加载

- 应用启动时会自动从 API 加载所有数据
- 如果某个 API 调用失败，会使用默认值并显示警告

### 错误处理

- 所有 API 错误都会显示在右上角的错误提示中
- 用户可以点击关闭按钮消除错误提示

### 加载状态

- 在进行 API 调用时会显示全局加载指示器
- 确保用户了解正在进行数据操作

### 数据同步

- 所有数据修改都会立即同步到后端
- 不再依赖本地存储，确保数据的一致性

## 开发注意事项

### API 配置

- API 基础 URL 在 `src/utils/api.ts` 中配置
- 默认后端地址为 `http://localhost:8080/api`
- 如需修改，请更新 `API_BASE_URL` 常量

### 错误处理

- 所有 API 调用都包含错误处理
- 网络错误会自动重试（可在 api.ts 中配置）
- 确保用户体验不会因网络问题中断

### 类型安全

- 所有 API 响应都有完整的 TypeScript 类型定义
- 确保前后端数据结构的一致性

## 未来改进

1. **离线支持**: 可以考虑添加 Service Worker 实现离线缓存
2. **实时同步**: 可以使用 WebSocket 实现多客户端间的实时数据同步
3. **分页加载**: 对于大量数据，可以实现分页或无限滚动
4. **缓存策略**: 实现智能缓存减少不必要的 API 调用

## 故障排除

### 常见问题

1. **无法连接到后端**

   - 确认后端服务器正在运行
   - 检查 CORS 配置是否正确
   - 验证 API_BASE_URL 是否正确

2. **数据加载失败**

   - 检查后端数据库连接
   - 查看浏览器控制台的错误信息
   - 确认后端 API 响应格式正确

3. **CORS 错误**
   - 确认后端 CORS 配置包含前端地址
   - 检查前端开发服务器地址是否在允许列表中
