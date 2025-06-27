package model

// SearchEngine 搜索引擎模型
type SearchEngine struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	URL         string  `json:"url" db:"url"`
	Icon        *string `json:"icon,omitempty" db:"-"`                   // 临时字段，用于返回DataURL，不存储到数据库
	IconData    []byte  `json:"iconData,omitempty" db:"icon_data"`       // 图标二进制数据
	IconType    *string `json:"iconType,omitempty" db:"icon_type"`       // 图标MIME类型
	Placeholder *string `json:"placeholder" db:"placeholder"`
	IsDefault   bool    `json:"isDefault" db:"is_default"`
}

// GetIconDataURL 获取图标的Data URL格式（用于前端显示）
func (se *SearchEngine) GetIconDataURL() *string {
	if len(se.IconData) > 0 && se.IconType != nil {
		dataURL := "data:" + *se.IconType + ";base64," + string(se.IconData)
		return &dataURL
	}
	return nil
}

// CreateSearchEngineRequest 创建搜索引擎请求
type CreateSearchEngineRequest struct {
	ID          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	URL         string  `json:"url" binding:"required,url"`
	IconData    *string `json:"iconData,omitempty"`                     // Base64编码的图标数据
	IconType    *string `json:"iconType,omitempty"`                     // 图标MIME类型
	Placeholder *string `json:"placeholder"`
	IsDefault   *bool   `json:"isDefault"`
}

// UpdateSearchEngineRequest 更新搜索引擎请求
type UpdateSearchEngineRequest struct {
	Name        *string `json:"name"`
	URL         *string `json:"url"`
	IconData    *string `json:"iconData,omitempty"`         // Base64编码的图标数据
	IconType    *string `json:"iconType,omitempty"`         // 图标MIME类型
	Placeholder *string `json:"placeholder"`
	IsDefault   *bool   `json:"isDefault"`
} 