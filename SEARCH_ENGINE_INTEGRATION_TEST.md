# 搜索引擎后端集成测试

## 测试步骤

### 1. 启动后端服务

```bash
cd backend
go run cmd/server/main.go
```

### 2. 启动前端服务

```bash
npm run dev
```

### 3. 测试搜索引擎API集成

#### a) 数据加载测试

- 打开应用，查看浏览器控制台
- 应该看到 "Loading data from API..." 日志
- 如果后端没有搜索引擎数据，会看到 "No search engines found in backend, creating defaults..."
- 应该显示 "Loaded search engines: X" 其中X是搜索引擎数量

#### b) 添加搜索引擎测试

1. 点击设置按钮打开设置模态框
2. 在搜索引擎设置部分，点击"+ 添加搜索引擎"
3. 填写信息：
   - 名称：DuckDuckGo
   - URL：https://duckduckgo.com/?q={}
   - 图标：https://duckduckgo.com/favicon.ico
   - 占位符：使用 DuckDuckGo 搜索...
4. 确认添加，应该看到"搜索引擎添加成功！"提示
5. 新搜索引擎应该出现在列表中

#### c) 编辑搜索引擎测试

1. 点击某个搜索引擎的编辑按钮（✏️）
2. 修改名称或其他信息
3. 确认修改，应该看到"搜索引擎更新成功！"提示
4. 修改应该立即反映在列表中

#### d) 删除搜索引擎测试

1. 点击某个搜索引擎的删除按钮（🗑️）
2. 确认删除，应该看到"搜索引擎删除成功！"提示
3. 搜索引擎应该从列表中消失
4. 注意：如果只有一个搜索引擎，删除按钮应该被禁用

#### e) 搜索引擎选择测试

1. 在搜索框左侧点击搜索引擎图标
2. 应该显示所有可用的搜索引擎
3. 选择不同的搜索引擎
4. 搜索框的占位符文本应该更新
5. 进行搜索测试，确保使用正确的搜索引擎

### 4. 数据持久化测试

1. 刷新页面
2. 添加的搜索引擎应该仍然存在
3. 设置的默认搜索引擎应该保持不变

### 5. 错误处理测试

1. 停止后端服务
2. 尝试添加/编辑/删除搜索引擎
3. 应该看到错误提示，不会崩溃
4. 重新启动后端服务，功能应该恢复正常

## 预期结果

✅ **成功标准**：

- 所有搜索引擎CRUD操作都通过API与后端同步
- 数据持久化到数据库
- 错误处理优雅，有用户友好的提示
- 界面响应流畅，有加载状态指示
- 搜索功能正常工作

❌ **失败指标**：

- 数据不持久化
- 添加/编辑/删除操作失败
- 界面崩溃或冻结
- 错误信息不友好
- 搜索功能异常

## 日志监控

在浏览器控制台中监控以下日志：

- `Loading data from API...`
- `Loaded search engines: X`
- `Search engine added successfully: {...}`
- `Search engine updated successfully: {...}`
- `Search engine deleted successfully: ID`
- 任何错误信息

## API调用验证

使用浏览器开发者工具的网络标签，验证以下API调用：

- `GET /api/search-engines` - 获取搜索引擎列表
- `POST /api/search-engines` - 创建搜索引擎
- `PUT /api/search-engines/:id` - 更新搜索引擎
- `DELETE /api/search-engines/:id` - 删除搜索引擎
- `PUT /api/settings` - 更新设置（包含默认搜索引擎）
