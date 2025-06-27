# 搜索引擎图标升级说明

## 🎯 升级概述

本次升级将搜索引擎的图标存储方式从**网址链接**改为**二进制数据**，解决因网络访问失败导致的图标显示问题。

## 🔧 主要变更

### 后端变更

#### 1. 数据库结构

- **新增字段**：
  - `icon_data` (BLOB): 存储图标的二进制数据
  - `icon_type` (TEXT): 存储图标的MIME类型
- **保留字段**：
  - `icon` (TEXT): 保留作为备选方案，支持向后兼容

#### 2. 模型层 (`backend/internal/model/search_engine.go`)

```go
type SearchEngine struct {
    ID          string  `json:"id" db:"id"`
    Name        string  `json:"name" db:"name"`
    URL         string  `json:"url" db:"url"`
    Icon        *string `json:"icon" db:"icon"`                     // 保留原有的URL字段作为备选
    IconData    []byte  `json:"iconData,omitempty" db:"icon_data"`  // 新增：图标二进制数据
    IconType    *string `json:"iconType,omitempty" db:"icon_type"`  // 新增：图标MIME类型
    Placeholder *string `json:"placeholder" db:"placeholder"`
    IsDefault   bool    `json:"isDefault" db:"is_default"`
}
```

#### 3. API 接口

- **新增接口**：
  - `POST /api/search-engines/with-icon` - 创建带图片上传的搜索引擎
  - `PUT /api/search-engines/:id/with-icon` - 更新带图片上传的搜索引擎
- **保留接口**：
  - `POST /api/search-engines` - 创建搜索引擎（URL方式）
  - `PUT /api/search-engines/:id` - 更新搜索引擎（URL方式）

### 前端变更

#### 1. API 客户端 (`src/utils/api.ts`)

```typescript
// 新增图片上传支持的方法
async createSearchEngineWithIcon(
  engine: { name: string; url: string; placeholder?: string; isDefault?: boolean },
  iconFile?: File,
  iconUrl?: string
): Promise<SearchEngine>

async updateSearchEngineWithIcon(
  id: string,
  updates: { name?: string; url?: string; placeholder?: string; isDefault?: boolean },
  iconFile?: File,
  iconUrl?: string
): Promise<SearchEngine>
```

#### 2. 用户界面

- **添加搜索引擎时**：提供选择上传本地文件或输入URL的选项
- **编辑搜索引擎时**：可选择是否更新图标，支持文件上传或URL
- **文件限制**：最大2MB，支持PNG、JPG、GIF、WebP、SVG格式

## 🚀 使用方法

### 1. 添加搜索引擎（本地图标）

1. 打开设置 → 搜索引擎设置
2. 点击"添加搜索引擎"
3. 填写名称和URL
4. 选择"上传本地图标文件"
5. 选择图片文件（最大2MB）
6. 确认添加

### 2. 添加搜索引擎（URL图标）

1. 打开设置 → 搜索引擎设置
2. 点击"添加搜索引擎"
3. 填写名称和URL
4. 选择"输入网址"
5. 输入图标URL
6. 确认添加

### 3. 更新现有搜索引擎图标

1. 点击搜索引擎的编辑按钮
2. 编辑基本信息
3. 选择"是否要更新图标"
4. 选择上传文件或输入URL
5. 确认更新

## 📋 技术特性

### 1. 数据处理流程

```
用户上传图片 → 前端转换为FormData → 后端解码base64 → 存储为二进制 → 返回时转换为Data URL
```

### 2. 向后兼容

- 现有的URL图标继续工作
- 自动优先使用二进制数据
- 如果二进制数据不可用，回退到URL

### 3. 性能优化

- 传输时自动转换为base64 Data URL
- 减少网络请求依赖
- 改善图标加载速度和可靠性

## 🔧 开发者说明

### 后端数据流

1. **接收**: Multipart form data（文件上传）
2. **验证**: 文件类型和大小
3. **转换**: 文件 → 二进制数据 + MIME类型
4. **存储**: 数据库BLOB字段
5. **返回**: 二进制 → base64 Data URL

### 前端数据流

1. **选择**: 文件选择器或URL输入
2. **验证**: 客户端文件大小检查
3. **发送**: FormData格式上传
4. **显示**: 接收Data URL直接使用

### 数据库迁移

```sql
-- 自动执行，向后兼容
ALTER TABLE search_engines ADD COLUMN icon_data BLOB;
ALTER TABLE search_engines ADD COLUMN icon_type TEXT;
```

## 🎯 优势

1. **可靠性**：不依赖外部图标服务
2. **性能**：减少额外的网络请求
3. **离线**：完全离线可用
4. **安全**：避免外部资源加载风险
5. **兼容**：保持向后兼容性

## 🧪 测试指南

### 1. 启动服务

```bash
cd backend
go run cmd/server/main.go
```

### 2. 前端开发服务器

```bash
npm run dev
```

### 3. 测试场景

- ✅ 上传PNG图标
- ✅ 上传JPG图标
- ✅ 上传SVG图标
- ✅ 使用URL图标
- ✅ 文件大小限制测试
- ✅ 格式验证测试
- ✅ 编辑现有搜索引擎
- ✅ 向后兼容性测试
