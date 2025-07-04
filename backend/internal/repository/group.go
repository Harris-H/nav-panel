package repository

import (
	"database/sql"
	"fmt"
	"nav-panel-backend/internal/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

// GetAll 获取所有分组
func (r *GroupRepository) GetAll() ([]model.Group, error) {
	query := `SELECT id, name, color, icon, sort_order, is_collapsed, created_at, updated_at 
	          FROM groups ORDER BY sort_order ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Name, &g.Color, &g.Icon, &g.SortOrder, &g.IsCollapsed, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

// GetAllWithWebsites 获取所有分组及其网站
func (r *GroupRepository) GetAllWithWebsites() ([]model.GroupWithWebsites, error) {
	// 首先获取所有分组
	groups, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var result []model.GroupWithWebsites
	for _, group := range groups {
		groupWithWebsites := model.GroupWithWebsites{
			Group: &group,
		}

		// 获取该分组下的所有网站
		query := `SELECT id, name, url, description, icon, sort_order, group_id, created_at, updated_at 
		          FROM websites WHERE group_id = ? ORDER BY sort_order ASC`
		
		rows, err := r.db.Query(query, group.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var websites []model.Website
		for rows.Next() {
			var w model.Website
			err := rows.Scan(&w.ID, &w.Name, &w.URL, &w.Description, &w.Icon, &w.SortOrder, &w.GroupId, &w.CreatedAt, &w.UpdatedAt)
			if err != nil {
				return nil, err
			}
			websites = append(websites, w)
		}

		groupWithWebsites.Websites = websites
		result = append(result, groupWithWebsites)
	}

	return result, nil
}

// Create 创建分组
func (r *GroupRepository) Create(req model.CreateGroupRequest) (*model.Group, error) {
	// 获取最大排序号
	var maxSortOrder sql.NullInt64
	err := r.db.QueryRow("SELECT MAX(sort_order) FROM groups").Scan(&maxSortOrder)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	sortOrder := 0
	if maxSortOrder.Valid {
		sortOrder = int(maxSortOrder.Int64) + 1
	}

	group := &model.Group{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Color:       req.Color,
		Icon:        req.Icon,
		SortOrder:   sortOrder,
		IsCollapsed: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO groups (id, name, color, icon, sort_order, is_collapsed, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err = r.db.Exec(query, group.ID, group.Name, group.Color, group.Icon, group.SortOrder, group.IsCollapsed, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// Update 更新分组
func (r *GroupRepository) Update(id string, req model.UpdateGroupRequest) (*model.Group, error) {
	setParts := []string{}
	args := []interface{}{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *req.Color)
	}
	if req.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *req.Icon)
	}
	if req.IsCollapsed != nil {
		setParts = append(setParts, "is_collapsed = ?")
		args = append(args, *req.IsCollapsed)
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id) // 最后添加WHERE条件的参数

	query := fmt.Sprintf("UPDATE groups SET %s WHERE id = ?", strings.Join(setParts, ", "))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// 返回更新后的分组
	return r.GetByID(id)
}

// GetByID 根据ID获取分组
func (r *GroupRepository) GetByID(id string) (*model.Group, error) {
	query := `SELECT id, name, color, icon, sort_order, is_collapsed, created_at, updated_at 
	          FROM groups WHERE id = ?`
	
	row := r.db.QueryRow(query, id)
	
	var group model.Group
	err := row.Scan(&group.ID, &group.Name, &group.Color, &group.Icon, &group.SortOrder, &group.IsCollapsed, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &group, nil
}

// Delete 删除分组
func (r *GroupRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 将分组下的网站移到未分组状态
	_, err = tx.Exec("UPDATE websites SET group_id = NULL WHERE group_id = ?", id)
	if err != nil {
		return err
	}

	// 删除分组
	_, err = tx.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Reorder 重新排序分组
func (r *GroupRepository) Reorder(groupIds []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, groupId := range groupIds {
		_, err = tx.Exec("UPDATE groups SET sort_order = ?, updated_at = ? WHERE id = ?", 
			i, time.Now(), groupId)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// MoveWebsiteToGroup 移动网站到分组
func (r *GroupRepository) MoveWebsiteToGroup(websiteId string, groupId *string, position *int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 如果指定了位置，需要更新排序
	if position != nil && groupId != nil {
		// 为该位置及之后的网站排序号加1
		_, err = tx.Exec(`UPDATE websites SET sort_order = sort_order + 1 
		                  WHERE group_id = ? AND sort_order >= ?`, *groupId, *position)
		if err != nil {
			return err
		}

		// 更新目标网站
		_, err = tx.Exec(`UPDATE websites SET group_id = ?, sort_order = ?, updated_at = ? 
		                  WHERE id = ?`, groupId, *position, time.Now(), websiteId)
	} else {
		// 如果没有指定位置，放到最后
		var sortOrder int
		if groupId != nil {
			var maxSort sql.NullInt64
			err = tx.QueryRow("SELECT MAX(sort_order) FROM websites WHERE group_id = ?", *groupId).Scan(&maxSort)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if maxSort.Valid {
				sortOrder = int(maxSort.Int64) + 1
			}
		} else {
			var maxSort sql.NullInt64
			err = tx.QueryRow("SELECT MAX(sort_order) FROM websites WHERE group_id IS NULL").Scan(&maxSort)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if maxSort.Valid {
				sortOrder = int(maxSort.Int64) + 1
			}
		}

		_, err = tx.Exec(`UPDATE websites SET group_id = ?, sort_order = ?, updated_at = ? 
		                  WHERE id = ?`, groupId, sortOrder, time.Now(), websiteId)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}