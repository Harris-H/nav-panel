package service

import (
	"fmt"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/repository"
)

type SearchEngineService struct {
	repo *repository.SearchEngineRepository
}

func NewSearchEngineService(repo *repository.SearchEngineRepository) *SearchEngineService {
	return &SearchEngineService{repo: repo}
}

func (s *SearchEngineService) GetAll() ([]model.SearchEngine, error) {
	return s.repo.GetAll()
}

func (s *SearchEngineService) GetByID(id string) (*model.SearchEngine, error) {
	return s.repo.GetByID(id)
}



func (s *SearchEngineService) Create(req *model.CreateSearchEngineRequest) (*model.SearchEngine, error) {
	// 检查是否已存在相同ID的搜索引擎
	existing, _ := s.repo.GetByID(req.ID)
	if existing != nil {
		return nil, fmt.Errorf("search engine with id '%s' already exists", req.ID)
	}

	engine := &model.SearchEngine{
		ID:          req.ID,
		Name:        req.Name,
		URL:         req.URL,
		Icon:        req.Icon,
		Placeholder: req.Placeholder,
		IsDefault:   req.IsDefault != nil && *req.IsDefault,
	}

	err := s.repo.Create(engine)
	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (s *SearchEngineService) Update(id string, req *model.UpdateSearchEngineRequest) (*model.SearchEngine, error) {
	// 检查搜索引擎是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("search engine not found: %w", err)
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
		if *req.Icon == "" {
			// 空字符串表示清除图标，设置为 NULL
			updates["icon"] = nil
		} else {
			updates["icon"] = *req.Icon
		}
	}
	if req.Placeholder != nil {
		updates["placeholder"] = *req.Placeholder
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}

	// 执行更新
	err = s.repo.Update(id, updates)
	if err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return s.GetByID(id)
}



func (s *SearchEngineService) Delete(id string) error {
	// 检查搜索引擎是否存在
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("search engine not found: %w", err)
	}

	// 如果删除的是默认搜索引擎，需要警告或设置新的默认引擎
	if existing.IsDefault {
		// 这里可以选择自动设置第一个为默认，或者返回错误要求用户先设置其他为默认
		engines, err := s.repo.GetAll()
		if err == nil && len(engines) > 1 {
			// 找到第一个不是当前删除的引擎，设为默认
			for _, engine := range engines {
				if engine.ID != id {
					s.repo.Update(engine.ID, map[string]interface{}{"is_default": true})
					break
				}
			}
		}
	}

	return s.repo.Delete(id)
} 