package service

import (
	"fmt"
	"nav-panel-backend/internal/model"
	"nav-panel-backend/internal/repository"
)

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

// GetAllGroups 获取所有分组
func (s *GroupService) GetAllGroups() ([]model.Group, error) {
	return s.repo.GetAll()
}

// GetAllGroupsWithWebsites 获取所有分组及其网站
func (s *GroupService) GetAllGroupsWithWebsites() ([]model.GroupWithWebsites, error) {
	return s.repo.GetAllWithWebsites()
}

// CreateGroup 创建分组
func (s *GroupService) CreateGroup(req model.CreateGroupRequest) (*model.Group, error) {
	// 业务逻辑验证
	if req.Name == "" {
		return nil, fmt.Errorf("分组名称不能为空")
	}

	return s.repo.Create(req)
}

// UpdateGroup 更新分组
func (s *GroupService) UpdateGroup(id string, req model.UpdateGroupRequest) error {
	// 业务逻辑验证
	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("分组名称不能为空")
	}

	return s.repo.Update(id, req)
}

// DeleteGroup 删除分组
func (s *GroupService) DeleteGroup(id string) error {
	return s.repo.Delete(id)
}

// ReorderGroups 重新排序分组
func (s *GroupService) ReorderGroups(groupIds []string) error {
	if len(groupIds) == 0 {
		return fmt.Errorf("分组ID列表不能为空")
	}

	return s.repo.Reorder(groupIds)
}

// MoveWebsiteToGroup 移动网站到分组
func (s *GroupService) MoveWebsiteToGroup(websiteId string, groupId *string, position *int) error {
	if websiteId == "" {
		return fmt.Errorf("网站ID不能为空")
	}

	return s.repo.MoveWebsiteToGroup(websiteId, groupId, position)
}