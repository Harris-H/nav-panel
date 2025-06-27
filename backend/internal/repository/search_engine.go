package repository

import (
	"database/sql"
	"nav-panel-backend/internal/model"
)

type SearchEngineRepository struct {
	db *sql.DB
}

func NewSearchEngineRepository(db *sql.DB) *SearchEngineRepository {
	return &SearchEngineRepository{db: db}
}

func (r *SearchEngineRepository) GetAll() ([]model.SearchEngine, error) {
	query := `
		SELECT id, name, url, icon_data, icon_type, placeholder, is_default 
		FROM search_engines 
		ORDER BY is_default DESC, name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var engines []model.SearchEngine
	for rows.Next() {
		var e model.SearchEngine
		err := rows.Scan(&e.ID, &e.Name, &e.URL, &e.IconData, &e.IconType, &e.Placeholder, &e.IsDefault)
		if err != nil {
			return nil, err
		}
		engines = append(engines, e)
	}

	return engines, nil
}

func (r *SearchEngineRepository) GetByID(id string) (*model.SearchEngine, error) {
	query := `
		SELECT id, name, url, icon_data, icon_type, placeholder, is_default 
		FROM search_engines 
		WHERE id = ?
	`
	
	var e model.SearchEngine
	err := r.db.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.URL, &e.IconData, &e.IconType, &e.Placeholder, &e.IsDefault)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *SearchEngineRepository) Create(engine *model.SearchEngine) error {
	// 如果设置为默认，先取消其他默认引擎
	if engine.IsDefault {
		if err := r.clearDefaultFlags(); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO search_engines (id, name, url, icon_data, icon_type, placeholder, is_default)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query, engine.ID, engine.Name, engine.URL, engine.IconData, engine.IconType, engine.Placeholder, engine.IsDefault)
	return err
}

func (r *SearchEngineRepository) Update(id string, updates map[string]interface{}) error {
	// 如果更新默认状态为true，先清除其他默认状态
	if isDefault, exists := updates["is_default"]; exists && isDefault.(bool) {
		if err := r.clearDefaultFlags(); err != nil {
			return err
		}
	}

	// 动态构建更新查询
	query := "UPDATE search_engines SET "
	args := []interface{}{}
	first := true

	for key, value := range updates {
		if !first {
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		first = false
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *SearchEngineRepository) Delete(id string) error {
	query := "DELETE FROM search_engines WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SearchEngineRepository) clearDefaultFlags() error {
	query := "UPDATE search_engines SET is_default = FALSE"
	_, err := r.db.Exec(query)
	return err
} 