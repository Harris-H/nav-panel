package model

import "time"

// Group 分组模型
type Group struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Color       *string   `json:"color" db:"color"`
	Icon        *string   `json:"icon" db:"icon"`
	SortOrder   int       `json:"sortOrder" db:"sort_order"`
	IsCollapsed bool      `json:"isCollapsed" db:"is_collapsed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateGroupRequest 创建分组请求
type CreateGroupRequest struct {
	Name  string  `json:"name" binding:"required"`
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
}

// UpdateGroupRequest 更新分组请求
type UpdateGroupRequest struct {
	Name        *string `json:"name"`
	Color       *string `json:"color"`
	Icon        *string `json:"icon"`
	IsCollapsed *bool   `json:"isCollapsed"`
}

// ReorderGroupsRequest 重新排序分组请求
type ReorderGroupsRequest struct {
	GroupIds []string `json:"groupIds" binding:"required"`
}

// MoveWebsiteToGroupRequest 移动网站到分组请求
type MoveWebsiteToGroupRequest struct {
	WebsiteId string  `json:"websiteId" binding:"required"`
	GroupId   *string `json:"groupId"` // null表示移出分组
	Position  *int    `json:"position"` // 在分组中的位置
}

// GroupWithWebsites 带网站的分组
type GroupWithWebsites struct {
	*Group
	Websites []Website `json:"websites"`
}