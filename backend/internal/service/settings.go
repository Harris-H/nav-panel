package service

import (
	"encoding/json"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/repository"
)

type SettingsService struct {
	settingsRepo     *repository.SettingsRepository
	websiteRepo      *repository.WebsiteRepository
	searchEngineRepo *repository.SearchEngineRepository
}

func NewSettingsService(settingsRepo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
	}
}

// SetRepositories 设置其他仓库依赖（用于导入导出功能）
func (s *SettingsService) SetRepositories(websiteRepo *repository.WebsiteRepository, searchEngineRepo *repository.SearchEngineRepository) {
	s.websiteRepo = websiteRepo
	s.searchEngineRepo = searchEngineRepo
}



func (s *SettingsService) Get() (*model.AppSettingsResponse, error) {
	// 获取设置数据
	settings, err := s.settingsRepo.Get()
	if err != nil {
		return nil, err
	}

	// 获取搜索引擎数据
	searchEngines, err := s.searchEngineRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 解析JSON配置
	var layout model.LayoutConfig
	if err := json.Unmarshal([]byte(settings.LayoutConfig), &layout); err != nil {
		// 使用默认值
		layout = model.LayoutConfig{
			Columns:    6,
			CardSize:   "medium",
			ShowLabels: true,
			Gap:        20,
		}
	}

	var background model.BackgroundConfig
	if err := json.Unmarshal([]byte(settings.BackgroundConfig), &background); err != nil {
		background = model.BackgroundConfig{
			Type:  "gradient",
			Value: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
		}
	}

	var cardStyle model.CardStyleConfig
	if err := json.Unmarshal([]byte(settings.CardStyleConfig), &cardStyle); err != nil {
		cardStyle = model.CardStyleConfig{
			BorderRadius: 12,
			Opacity:      0.9,
			Shadow:       true,
		}
	}

	var searchConfig model.SearchConfig
	if err := json.Unmarshal([]byte(settings.SearchConfig), &searchConfig); err != nil {
		searchConfig = model.SearchConfig{
			Enabled:         true,
			DefaultEngineId: "google",
			OpenInNewTab:    true,
		}
	}

	// 将搜索引擎数据添加到搜索配置中
	searchConfig.Engines = searchEngines

	return &model.AppSettingsResponse{
		Theme:      settings.Theme,
		Layout:     layout,
		Background: background,
		CardStyle:  cardStyle,
		Search:     searchConfig,
		UpdatedAt:  settings.UpdatedAt,
	}, nil
}

func (s *SettingsService) Update(req *model.UpdateSettingsRequest) (*model.AppSettingsResponse, error) {
	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Theme != nil {
		updates["theme"] = *req.Theme
	}

	if req.Layout != nil {
		layoutJSON, err := json.Marshal(*req.Layout)
		if err != nil {
			return nil, err
		}
		updates["layout_config"] = string(layoutJSON)
	}

	if req.Background != nil {
		backgroundJSON, err := json.Marshal(*req.Background)
		if err != nil {
			return nil, err
		}
		updates["background_config"] = string(backgroundJSON)
	}

	if req.CardStyle != nil {
		cardStyleJSON, err := json.Marshal(*req.CardStyle)
		if err != nil {
			return nil, err
		}
		updates["card_style_config"] = string(cardStyleJSON)
	}

	if req.Search != nil {
		// 处理搜索配置，但不包含engines数组（engines单独管理）
		searchMap, ok := (*req.Search).(map[string]interface{})
		if ok {
			// 移除engines字段，只保存其他搜索配置
			searchConfigOnly := make(map[string]interface{})
			for k, v := range searchMap {
				if k != "engines" {
					searchConfigOnly[k] = v
				}
			}
			
			searchJSON, err := json.Marshal(searchConfigOnly)
			if err != nil {
				return nil, err
			}
			updates["search_config"] = string(searchJSON)
		}
	}

	// 执行更新
	if len(updates) > 0 {
		err := s.settingsRepo.Update(updates)
		if err != nil {
			return nil, err
		}
	}

	// 返回更新后的设置
	return s.Get()
}

func (s *SettingsService) Export() (*model.ExportData, error) {
	// 获取所有数据
	websites, err := s.websiteRepo.GetAll()
	if err != nil {
		return nil, err
	}

	searchEngines, err := s.searchEngineRepo.GetAll()
	if err != nil {
		return nil, err
	}

	settings, err := s.settingsRepo.Get()
	if err != nil {
		return nil, err
	}

	return &model.ExportData{
		Websites:      websites,
		SearchEngines: searchEngines,
		Settings:      settings,
	}, nil
}

func (s *SettingsService) Import(data *model.ImportData) error {
	// 导入网站数据
	if data.Websites != nil {
		// 删除现有网站（可选：也可以选择合并）
		existing, _ := s.websiteRepo.GetAll()
		for _, site := range existing {
			s.websiteRepo.Delete(site.ID)
		}

		// 导入新网站
		for _, website := range data.Websites {
			s.websiteRepo.Create(&website)
		}
	}

	// 导入搜索引擎数据
	if data.SearchEngines != nil {
		// 删除现有搜索引擎
		existing, _ := s.searchEngineRepo.GetAll()
		for _, engine := range existing {
			s.searchEngineRepo.Delete(engine.ID)
		}

		// 导入新搜索引擎
		for _, engine := range data.SearchEngines {
			s.searchEngineRepo.Create(&engine)
		}
	}

	// 导入设置
	if data.Settings != nil {
		updates := map[string]interface{}{
			"theme":               data.Settings.Theme,
			"layout_config":       data.Settings.LayoutConfig,
			"background_config":   data.Settings.BackgroundConfig,
			"card_style_config":   data.Settings.CardStyleConfig,
			"search_config":       data.Settings.SearchConfig,
		}
		s.settingsRepo.Update(updates)
	}

	return nil
}



func (s *SettingsService) Reset() error {
	return s.settingsRepo.Reset()
} 