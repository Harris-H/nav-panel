# 前后端搜索引擎同步问题修复

## 🐛 问题描述

用户反映：在执行 `PUT /api/settings` 请求成功后，前端搜索引擎仍然没有更新，而且默认应该有搜索引擎。

## 🔍 问题分析

### 原始问题：

1. **后端数据格式问题**：

   - 后端将搜索配置存储为JSON字符串
   - 前端期望对象格式
   - 搜索引擎数据和搜索设置分离存储，但前端混合处理

2. **前端同步问题**：

   - `updateSettings` 方法没有正确处理搜索引擎数据
   - 缺少对搜索引擎数据的重新加载逻辑

3. **数据结构不匹配**：
   - 前端 `AppSettings.search.engines` 期望包含搜索引擎数组
   - 后端设置表只存储搜索配置，不包含搜索引擎数据

## ✅ 解决方案

### 1. 后端修复

#### a) 新增数据结构 (`backend/internal/model/settings.go`)

```go
// 新增响应格式，将JSON字符串转换为对象
type AppSettingsResponse struct {
    Theme      string            `json:"theme"`
    Layout     LayoutConfig      `json:"layout"`
    Background BackgroundConfig  `json:"background"`
    CardStyle  CardStyleConfig   `json:"cardStyle"`
    Search     SearchConfig      `json:"search"`  // 包含搜索引擎数组
    UpdatedAt  time.Time         `json:"updatedAt"`
}

type SearchConfig struct {
    Enabled         bool            `json:"enabled"`
    DefaultEngineId string          `json:"defaultEngineId"`
    OpenInNewTab    bool            `json:"openInNewTab"`
    Engines         []SearchEngine  `json:"engines"`  // 重要：包含搜索引擎数据
}
```

#### b) 修改设置服务 (`backend/internal/service/settings.go`)

```go
func (s *SettingsService) Get() (*model.AppSettingsResponse, error) {
    // 获取设置数据
    settings, err := s.settingsRepo.Get()
    if err != nil {
        return nil, err
    }

    // 获取搜索引擎数据
    searchEngines, err := s.searchEngineRepo.GetAll()
    if err != nil {
        return nil, err
    }

    // 解析JSON配置并合并搜索引擎数据
    var searchConfig model.SearchConfig
    json.Unmarshal([]byte(settings.SearchConfig), &searchConfig)
    searchConfig.Engines = searchEngines  // 关键：合并搜索引擎数据

    return &model.AppSettingsResponse{
        // ... 其他字段
        Search: searchConfig,
    }, nil
}
```

### 2. 前端修复

#### a) 修改 `updateSettings` 方法 (`src/stores/app.ts`)

```typescript
const updateSettings = async (newSettings: Partial<AppSettings>) => {
  try {
    const updatedSettings = await api.updateSettings(newSettings)

    // 更新本地设置，确保搜索引擎数据正确同步
    settings.value = {
      ...settings.value,
      ...updatedSettings,
    }

    // 如果后端返回的搜索引擎数据无效，从API重新获取
    if (!settings.value.search?.engines || !Array.isArray(settings.value.search.engines)) {
      const searchEngines = await api.getSearchEngines()
      settings.value.search = {
        ...settings.value.search,
        engines: searchEngines.length > 0 ? searchEngines : defaultSearchEngines,
      }
    }

    // 确保默认搜索引擎 ID 有效
    const validEngineIds = settings.value.search.engines.map((e) => e.id)
    if (!validEngineIds.includes(settings.value.search.defaultEngineId)) {
      settings.value.search.defaultEngineId = settings.value.search.engines[0]?.id || 'google'
    }

    // 更新当前搜索引擎
    currentSearchEngine.value = defaultSearchEngine.value
  } catch (err) {
    handleError(err, 'Error updating settings')
  }
}
```

## 🚀 修复效果

### 现在的数据流：

1. **设置获取** (`GET /api/settings`):

   ```json
   {
     "theme": "light",
     "layout": { "columns": 6, "cardSize": "medium", ... },
     "search": {
       "enabled": true,
       "defaultEngineId": "google",
       "openInNewTab": true,
       "engines": [  // 重要：现在包含搜索引擎数据
         { "id": "google", "name": "Google", ... },
         { "id": "baidu", "name": "百度", ... }
       ]
     }
   }
   ```

2. **设置更新** (`PUT /api/settings`):

   - 前端发送设置更新请求
   - 后端保存设置（不包含engines）
   - 后端响应包含完整的设置数据（包含当前搜索引擎）
   - 前端正确同步搜索引擎数据

3. **数据一致性**：
   - 搜索引擎通过独立的API管理 (`/api/search-engines`)
   - 设置API返回时自动包含最新的搜索引擎数据
   - 前端始终获得完整的、一致的数据

## 🧪 测试步骤

1. **启动后端**：

   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **启动前端**：

   ```bash
   npm run dev
   ```

3. **验证修复**：
   - 打开浏览器开发者工具
   - 查看 `GET /api/settings` 响应，确认包含 `search.engines` 数组
   - 修改设置并保存
   - 查看 `PUT /api/settings` 响应，确认搜索引擎数据正确返回
   - 验证前端搜索引擎选择器正常工作

## 🎯 预期结果

✅ **修复后**：

- 设置API返回完整的搜索引擎数据
- 前端正确同步和显示搜索引擎
- 设置更新后搜索引擎立即可用
- 默认搜索引擎正确加载和显示
- 搜索功能正常工作

❌ **修复前的问题**：

- 设置API不包含搜索引擎数据
- 前端搜索引擎列表为空或过时
- 设置更新后搜索功能异常
- 默认搜索引擎丢失
