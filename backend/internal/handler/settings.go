package handler

import (
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	service *service.SettingsService
}

func NewSettingsHandler(service *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// Get 获取应用设置
func (h *SettingsHandler) Get(c *gin.Context) {
	settings, err := h.service.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// Update 更新应用设置
func (h *SettingsHandler) Update(c *gin.Context) {
	var req model.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.service.Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// Export 导出所有数据
func (h *SettingsHandler) Export(c *gin.Context) {
	data, err := h.service.Export()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=nav-panel-backup.json")
	c.JSON(http.StatusOK, data)
}

// Import 导入数据
func (h *SettingsHandler) Import(c *gin.Context) {
	var req model.ImportData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Import(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})
}

// Reset 重置设置为默认值
func (h *SettingsHandler) Reset(c *gin.Context) {
	err := h.service.Reset()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings reset to default"})
} 