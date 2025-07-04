package handler

import (
	"net/http"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// GetGroups 获取所有分组
func (h *GroupHandler) GetGroups(c *gin.Context) {
	groups, err := h.groupService.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// GetGroupsWithWebsites 获取所有分组及其网站
func (h *GroupHandler) GetGroupsWithWebsites(c *gin.Context) {
	groups, err := h.groupService.GetAllGroupsWithWebsites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// CreateGroup 创建分组
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req model.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.groupService.CreateGroup(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": group})
}

// UpdateGroup 更新分组
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req model.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedGroup, err := h.groupService.UpdateGroup(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedGroup})
}

// DeleteGroup 删除分组
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	err := h.groupService.DeleteGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// ReorderGroups 重新排序分组
func (h *GroupHandler) ReorderGroups(c *gin.Context) {
	var req model.ReorderGroupsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.groupService.ReorderGroups(req.GroupIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Groups reordered successfully"})
}

// MoveWebsiteToGroup 移动网站到分组
func (h *GroupHandler) MoveWebsiteToGroup(c *gin.Context) {
	var req model.MoveWebsiteToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.groupService.MoveWebsiteToGroup(req.WebsiteId, req.GroupId, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Website moved successfully"})
}