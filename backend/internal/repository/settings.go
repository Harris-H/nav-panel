package repository

import (
	"database/sql"
	"nav-panel-backend/internal/model"
	"time"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get() (*model.AppSettings, error) {
	query := `
		SELECT id, theme, layout_config, background_config, card_style_config, search_config, updated_at 
		FROM app_settings 
		WHERE id = 1
	`
	
	var s model.AppSettings
	err := r.db.QueryRow(query).Scan(&s.ID, &s.Theme, &s.LayoutConfig, &s.BackgroundConfig, &s.CardStyleConfig, &s.SearchConfig, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SettingsRepository) Update(updates map[string]interface{}) error {
	// 动态构建更新查询
	query := "UPDATE app_settings SET updated_at = ?"
	args := []interface{}{time.Now()}

	for key, value := range updates {
		query += ", " + key + " = ?"
		args = append(args, value)
	}

	query += " WHERE id = 1"

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *SettingsRepository) Reset() error {
	query := `
		UPDATE app_settings SET 
			theme = 'light',
			layout_config = '{"columns":6,"cardSize":"medium","showLabels":true,"gap":20}',
			background_config = '{"type":"gradient","value":"linear-gradient(135deg, #667eea 0%, #764ba2 100%)"}',
			card_style_config = '{"borderRadius":12,"opacity":0.9,"shadow":true}',
			search_config = '{"enabled":true,"defaultEngineId":"google","openInNewTab":true}',
			updated_at = ?
		WHERE id = 1
	`
	
	_, err := r.db.Exec(query, time.Now())
	return err
} 