package model

import "time"

// AppSettings 应用设置模型（数据库存储格式）
type AppSettings struct {
	ID            int       `json:"id" db:"id"`
	Theme         string    `json:"theme" db:"theme"`
	LayoutConfig  string    `json:"layout" db:"layout_config"`
	BackgroundConfig string `json:"background" db:"background_config"`
	CardStyleConfig string  `json:"cardStyle" db:"card_style_config"`
	SearchConfig  string    `json:"search" db:"search_config"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// LayoutConfig 布局配置
type LayoutConfig struct {
	Columns   int    `json:"columns"`
	CardSize  string `json:"cardSize"`
	ShowLabels bool  `json:"showLabels"`
	Gap       int    `json:"gap"`
}

// BackgroundConfig 背景配置
type BackgroundConfig struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// CardStyleConfig 卡片样式配置
type CardStyleConfig struct {
	BorderRadius int     `json:"borderRadius"`
	Opacity      float64 `json:"opacity"`
	Shadow       bool    `json:"shadow"`
}

// SearchConfig 搜索配置
type SearchConfig struct {
	Enabled         bool            `json:"enabled"`
	DefaultEngineId string          `json:"defaultEngineId"`
	OpenInNewTab    bool            `json:"openInNewTab"`
	Engines         []SearchEngine  `json:"engines"`
}

// AppSettingsResponse 应用设置响应格式（给前端的格式）
type AppSettingsResponse struct {
	Theme      string            `json:"theme"`
	Layout     LayoutConfig      `json:"layout"`
	Background BackgroundConfig  `json:"background"`
	CardStyle  CardStyleConfig   `json:"cardStyle"`
	Search     SearchConfig      `json:"search"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	Theme        *string `json:"theme"`
	Layout       *interface{} `json:"layout"`
	Background   *interface{} `json:"background"`
	CardStyle    *interface{} `json:"cardStyle"`
	Search       *interface{} `json:"search"`
}

// ExportData 导出数据结构
type ExportData struct {
	Websites      []Website      `json:"websites"`
	SearchEngines []SearchEngine `json:"searchEngines"`
	Settings      *AppSettings   `json:"settings"`
}

// ImportData 导入数据结构
type ImportData struct {
	Websites      []Website      `json:"websites"`
	SearchEngines []SearchEngine `json:"searchEngines"`
	Settings      *AppSettings   `json:"settings"`
} 