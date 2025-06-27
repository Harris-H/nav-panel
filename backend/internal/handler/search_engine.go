package handler

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchEngineHandler struct {
	service *service.SearchEngineService
}

func NewSearchEngineHandler(service *service.SearchEngineService) *SearchEngineHandler {
	return &SearchEngineHandler{service: service}
}

// GetAll 获取所有搜索引擎
func (h *SearchEngineHandler) GetAll(c *gin.Context) {
	engines, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": engines})
}

// Create 创建搜索引擎
func (h *SearchEngineHandler) Create(c *gin.Context) {
	var req model.CreateSearchEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": engine})
}

// CreateWithIcon 创建带图片上传的搜索引擎
func (h *SearchEngineHandler) CreateWithIcon(c *gin.Context) {
	// 解析multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 限制10MB
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// 获取基本字段
	req := model.CreateSearchEngineRequest{
		ID:          c.PostForm("id"),
		Name:        c.PostForm("name"),
		URL:         c.PostForm("url"),
		Placeholder: stringPtr(c.PostForm("placeholder")),
	}

	// 验证必填字段
	if req.ID == "" || req.Name == "" || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID, name and URL are required"})
		return
	}

	// 处理图标上传
	file, header, err := c.Request.FormFile("icon")
	if err == nil {
		defer file.Close()

		// 检查文件类型
		contentType := header.Header.Get("Content-Type")
		if !isValidImageType(contentType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type. Only PNG, JPG, GIF, WebP and SVG are allowed"})
			return
		}

		// 检查文件大小（限制2MB）
		if header.Size > 2<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image too large. Maximum size is 2MB"})
			return
		}

		// 读取文件内容
		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image file"})
			return
		}

		// 转换为base64
		base64Data := base64.StdEncoding.EncodeToString(data)
		dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
		req.IconData = &dataURL
	}

	// 处理isDefault字段
	if isDefaultStr := c.PostForm("isDefault"); isDefaultStr == "true" {
		isDefault := true
		req.IsDefault = &isDefault
	}

	engine, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": engine})
}

// Update 更新搜索引擎
func (h *SearchEngineHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var req model.UpdateSearchEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine, err := h.service.Update(id, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Search engine not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": engine})
}

// UpdateWithIcon 更新带图片上传的搜索引擎
func (h *SearchEngineHandler) UpdateWithIcon(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// 解析multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 限制10MB
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// 构建更新请求
	req := model.UpdateSearchEngineRequest{}

	if name := c.PostForm("name"); name != "" {
		req.Name = &name
	}
	if url := c.PostForm("url"); url != "" {
		req.URL = &url
	}
	if placeholder := c.PostForm("placeholder"); placeholder != "" {
		req.Placeholder = &placeholder
	}

	// 处理图标上传
	file, header, err := c.Request.FormFile("icon")
	if err == nil {
		defer file.Close()

		// 检查文件类型
		contentType := header.Header.Get("Content-Type")
		if !isValidImageType(contentType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type. Only PNG, JPG, GIF, WebP and SVG are allowed"})
			return
		}

		// 检查文件大小（限制2MB）
		if header.Size > 2<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image too large. Maximum size is 2MB"})
			return
		}

		// 读取文件内容
		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image file"})
			return
		}

		// 转换为base64
		base64Data := base64.StdEncoding.EncodeToString(data)
		dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
		req.IconData = &dataURL
	}

	// 处理isDefault字段
	if isDefaultStr := c.PostForm("isDefault"); isDefaultStr != "" {
		isDefault := isDefaultStr == "true"
		req.IsDefault = &isDefault
	}

	engine, err := h.service.Update(id, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Search engine not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": engine})
}

// Delete 删除搜索引擎
func (h *SearchEngineHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Search engine not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Search engine deleted successfully"})
}

// 辅助函数
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/png",
		"image/jpeg",
		"image/jpg", 
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}
	
	for _, validType := range validTypes {
		if strings.HasPrefix(contentType, validType) {
			return true
		}
	}
	return false
} 