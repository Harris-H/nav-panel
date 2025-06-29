package model

// SearchEngine 搜索引擎模型
type SearchEngine struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	URL         string  `json:"url" db:"url"`
	Icon        *string `json:"icon,omitempty" db:"icon"`        // 直接存储DataURL格式
	Placeholder *string `json:"placeholder" db:"placeholder"`
	IsDefault   bool    `json:"isDefault" db:"is_default"`
}



// CreateSearchEngineRequest 创建搜索引擎请求
type CreateSearchEngineRequest struct {
	ID          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	URL         string  `json:"url" binding:"required,url"`
	Icon        *string `json:"icon,omitempty"`                  // DataURL格式图标
	Placeholder *string `json:"placeholder"`
	IsDefault   *bool   `json:"isDefault"`
}

// UpdateSearchEngineRequest 更新搜索引擎请求
type UpdateSearchEngineRequest struct {
	Name        *string `json:"name"`
	URL         *string `json:"url"`
	Icon        *string `json:"icon,omitempty"`          // DataURL格式图标
	Placeholder *string `json:"placeholder"`
	IsDefault   *bool   `json:"isDefault"`
} 