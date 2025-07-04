package repository

import (
	"database/sql"
	"nav-panel-backend/internal/model"
	"time"
)

type WebsiteRepository struct {
	db *sql.DB
}

func NewWebsiteRepository(db *sql.DB) *WebsiteRepository {
	return &WebsiteRepository{db: db}
}

func (r *WebsiteRepository) GetAll() ([]model.Website, error) {
	query := `
		SELECT id, name, url, icon, description, category, group_id, created_at, updated_at, COALESCE(sort_order, 0) as sort_order
		FROM websites 
		ORDER BY sort_order ASC, created_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var websites []model.Website
	for rows.Next() {
		var w model.Website
		err := rows.Scan(&w.ID, &w.Name, &w.URL, &w.Icon, &w.Description, &w.Category, &w.GroupId, &w.CreatedAt, &w.UpdatedAt, &w.SortOrder)
		if err != nil {
			return nil, err
		}
		websites = append(websites, w)
	}

	return websites, nil
}

func (r *WebsiteRepository) GetByID(id string) (*model.Website, error) {
	query := `
		SELECT id, name, url, icon, description, category, group_id, created_at, updated_at, COALESCE(sort_order, 0) as sort_order
		FROM websites 
		WHERE id = ?
	`
	
	var w model.Website
	err := r.db.QueryRow(query, id).Scan(&w.ID, &w.Name, &w.URL, &w.Icon, &w.Description, &w.Category, &w.GroupId, &w.CreatedAt, &w.UpdatedAt, &w.SortOrder)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *WebsiteRepository) Create(website *model.Website) error {
	query := `
		INSERT INTO websites (id, name, url, icon, description, category, group_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	website.CreatedAt = now
	website.UpdatedAt = now

	_, err := r.db.Exec(query, website.ID, website.Name, website.URL, website.Icon, website.Description, website.Category, website.GroupId, website.CreatedAt, website.UpdatedAt)
	return err
}

func (r *WebsiteRepository) Update(id string, updates map[string]interface{}) error {
	// 动态构建更新查询
	query := "UPDATE websites SET updated_at = ?"
	args := []interface{}{time.Now()}

	for key, value := range updates {
		query += ", " + key + " = ?"
		args = append(args, value)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *WebsiteRepository) Delete(id string) error {
	query := "DELETE FROM websites WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *WebsiteRepository) Reorder(websiteIds []string) error {
	// 开始事务
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 更新每个网站的排序
	for i, id := range websiteIds {
		query := "UPDATE websites SET sort_order = ?, updated_at = ? WHERE id = ?"
		_, err := tx.Exec(query, i, time.Now(), id)
		if err != nil {
			return err
		}
	}

	// 提交事务
	return tx.Commit()
} 