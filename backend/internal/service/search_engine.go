package service

import (
	"encoding/base64"
	"fmt"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/repository"
	"strings"
)

type SearchEngineService struct {
	repo *repository.SearchEngineRepository
}

func NewSearchEngineService(repo *repository.SearchEngineRepository) *SearchEngineService {
	return &SearchEngineService{repo: repo}
}

func (s *SearchEngineService) GetAll() ([]model.SearchEngine, error) {
	engines, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// 为每个搜索引擎处理图标数据
	for i := range engines {
		s.processIconData(&engines[i])
	}

	return engines, nil
}

func (s *SearchEngineService) GetByID(id string) (*model.SearchEngine, error) {
	engine, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.processIconData(engine)
	return engine, nil
}

// processIconData 处理图标数据，将二进制数据转换为前端可用的格式
func (s *SearchEngineService) processIconData(engine *model.SearchEngine) {
	if len(engine.IconData) > 0 && engine.IconType != nil {
		// 将二进制数据转换为base64字符串
		base64Data := base64.StdEncoding.EncodeToString(engine.IconData)
		dataURL := "data:" + *engine.IconType + ";base64," + base64Data
		engine.Icon = &dataURL
		// 清空二进制数据，减少传输大小
		engine.IconData = nil
		engine.IconType = nil
	}
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
		Placeholder: req.Placeholder,
		IsDefault:   req.IsDefault != nil && *req.IsDefault,
	}

	// 处理base64图标数据
	if req.IconData != nil && *req.IconData != "" {
		iconData, iconType, err := s.parseBase64Icon(*req.IconData)
		if err != nil {
			return nil, fmt.Errorf("invalid icon data: %w", err)
		}
		engine.IconData = iconData
		engine.IconType = iconType
	}

	err := s.repo.Create(engine)
	if err != nil {
		return nil, err
	}

	// 返回时处理图标数据
	s.processIconData(engine)
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
	if req.IconData != nil && *req.IconData != "" {
		iconData, iconType, err := s.parseBase64Icon(*req.IconData)
		if err != nil {
			return nil, fmt.Errorf("invalid icon data: %w", err)
		}
		updates["icon_data"] = iconData
		updates["icon_type"] = iconType
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

// parseBase64Icon 解析base64图标数据
func (s *SearchEngineService) parseBase64Icon(base64Data string) ([]byte, *string, error) {
	// 检查是否是data URL格式
	if strings.HasPrefix(base64Data, "data:") {
		parts := strings.Split(base64Data, ",")
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid data URL format")
		}
		
		// 提取MIME类型
		header := parts[0]
		mimeType := strings.TrimPrefix(header, "data:")
		mimeType = strings.TrimSuffix(mimeType, ";base64")
		
		// 解码base64数据
		data, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode base64 data: %w", err)
		}
		
		return data, &mimeType, nil
	} else {
		// 直接是base64数据，假设为PNG格式
		data, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode base64 data: %w", err)
		}
		
		mimeType := "image/png"
		return data, &mimeType, nil
	}
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