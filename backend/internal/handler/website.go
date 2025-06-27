package handler

import (
	"database/sql"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebsiteHandler struct {
	service *service.WebsiteService
}

func NewWebsiteHandler(service *service.WebsiteService) *WebsiteHandler {
	return &WebsiteHandler{service: service}
}

// GetAll 获取所有网站
func (h *WebsiteHandler) GetAll(c *gin.Context) {
	websites, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": websites})
}

// Create 创建网站
func (h *WebsiteHandler) Create(c *gin.Context) {
	var req model.CreateWebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	website, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": website})
}

// Update 更新网站
func (h *WebsiteHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req model.UpdateWebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	website, err := h.service.Update(id, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": website})
}

// Delete 删除网站
func (h *WebsiteHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Website deleted successfully"})
} 