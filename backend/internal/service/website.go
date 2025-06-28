package service

import (
	"fmt"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/repository"
	"time"
)

type WebsiteService struct {
	repo *repository.WebsiteRepository
}

func NewWebsiteService(repo *repository.WebsiteRepository) *WebsiteService {
	return &WebsiteService{repo: repo}
}

func (s *WebsiteService) GetAll() ([]model.Website, error) {
	return s.repo.GetAll()
}

func (s *WebsiteService) GetByID(id string) (*model.Website, error) {
	return s.repo.GetByID(id)
}

func (s *WebsiteService) Create(req *model.CreateWebsiteRequest) (*model.Website, error) {
	// 生成ID
	id := generateID()

	website := &model.Website{
		ID:          id,
		Name:        req.Name,
		URL:         req.URL,
		Icon:        req.Icon,
		Description: req.Description,
		Category:    req.Category,
	}

	err := s.repo.Create(website)
	if err != nil {
		return nil, err
	}

	return website, nil
}

func (s *WebsiteService) Update(id string, req *model.UpdateWebsiteRequest) (*model.Website, error) {
	// 检查网站是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("website not found: %w", err)
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}

	// 执行更新
	err = s.repo.Update(id, updates)
	if err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return s.repo.GetByID(id)
}

func (s *WebsiteService) Delete(id string) error {
	// 检查网站是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("website not found: %w", err)
	}

	return s.repo.Delete(id)
}

func (s *WebsiteService) Reorder(req *model.ReorderWebsitesRequest) ([]model.Website, error) {
	// 验证所有网站ID是否存在
	for _, id := range req.WebsiteIds {
		_, err := s.repo.GetByID(id)
		if err != nil {
			return nil, fmt.Errorf("website not found: %s", id)
		}
	}

	// 执行重新排序
	err := s.repo.Reorder(req.WebsiteIds)
	if err != nil {
		return nil, fmt.Errorf("failed to reorder websites: %w", err)
	}

	// 返回重新排序后的列表
	return s.repo.GetAll()
}

// generateID 生成唯一ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
} 