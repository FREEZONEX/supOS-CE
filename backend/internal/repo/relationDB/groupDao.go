package relationDB

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// GroupMapper resource_group 数据访问对象
type GroupMapper struct {
}

// SelectById 根据 ID 查询 Group
func (m *GroupMapper) SelectById(db *gorm.DB, id int64) (*GroupModel, error) {
	var group GroupModel
	err := db.Where("id = ?", id).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logx.Errorf("failed to select group by id: %v", err)
		return nil, err
	}
	return &group, nil
}

// Insert 插入 Group
func (m *GroupMapper) Insert(db *gorm.DB, group *GroupModel) error {
	err := db.Create(group).Error
	if err != nil {
		logx.Errorf("failed to insert group: %v", err)
		return err
	}
	return nil
}

// SaveBatch 批量插入 Group
func (m *GroupMapper) SaveBatch(db *gorm.DB, groups []*GroupModel) error {
	if len(groups) == 0 {
		return nil
	}
	err := db.Create(groups).Error
	if err != nil {
		logx.Errorf("failed to batch insert groups: %v", err)
		return err
	}
	return nil
}

// UpdateById 根据 ID 更新 Group
func (m *GroupMapper) UpdateById(db *gorm.DB, group *GroupModel) error {
	err := db.Model(&GroupModel{}).
		Where("id = ?", group.ID).
		Updates(group).Error
	if err != nil {
		logx.Errorf("failed to update group: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 ID 删除 Group
func (m *GroupMapper) DeleteById(db *gorm.DB, id int64) error {
	err := db.Where("id = ?", id).Delete(&GroupModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete group by id: %v", err)
		return err
	}
	return nil
}

// SelectAll 查询所有 Group
func (m *GroupMapper) SelectAll(db *gorm.DB) ([]*GroupModel, error) {
	var groups []*GroupModel
	err := db.Find(&groups).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*GroupModel{}, nil
		}
		logx.Errorf("failed to select all groups: %v", err)
		return nil, err
	}
	return groups, nil
}

// SelectByType 根据 type 查询 Group
func (m *GroupMapper) SelectByType(db *gorm.DB, typ int16) ([]*GroupModel, error) {
	var groups []*GroupModel
	err := db.Where("type = ?", typ).Find(&groups).Error
	if err != nil {
		logx.Errorf("failed to select groups by type: %v", err)
		return nil, err
	}
	return groups, nil
}

// SelectByIds 根据 ID 列表查询 Group
func (m *GroupMapper) SelectByIds(db *gorm.DB, ids []int64) ([]*GroupModel, error) {
	if len(ids) == 0 {
		return []*GroupModel{}, nil
	}
	var groups []*GroupModel
	err := db.Where("id IN ?", ids).Find(&groups).Error
	if err != nil {
		logx.Errorf("failed to select groups by ids: %v", err)
		return nil, err
	}
	return groups, nil
}

// DeleteBatchIds 根据 ID 列表批量删除 Group
func (m *GroupMapper) DeleteBatchIds(db *gorm.DB, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	err := db.Where("id IN ?", ids).Delete(&GroupModel{}).Error
	if err != nil {
		logx.Errorf("failed to batch delete groups: %v", err)
		return err
	}
	return nil
}
