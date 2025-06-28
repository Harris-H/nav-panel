package model

import "time"

// Website 网站模型
type Website struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	URL         string    `json:"url" db:"url"`
	Icon        *string   `json:"icon" db:"icon"`
	Description *string   `json:"description" db:"description"`
	Category    *string   `json:"category" db:"category"`
	SortOrder   int       `json:"sortOrder" db:"sort_order"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateWebsiteRequest 创建网站请求
type CreateWebsiteRequest struct {
	Name        string  `json:"name" binding:"required"`
	URL         string  `json:"url" binding:"required,url"`
	Icon        *string `json:"icon"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
}

// UpdateWebsiteRequest 更新网站请求
type UpdateWebsiteRequest struct {
	Name        *string `json:"name"`
	URL         *string `json:"url"`
	Icon        *string `json:"icon"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
}

// ReorderWebsitesRequest 重新排序网站请求
type ReorderWebsitesRequest struct {
	WebsiteIds []string `json:"websiteIds" binding:"required"`
} 