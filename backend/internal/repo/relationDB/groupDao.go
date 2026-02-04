package relationDB

import (
	"backend/internal/types"
	"context"

	"gitee.com/unitedrhino/share/errors"
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
func (m *GroupMapper) Save(ctx context.Context, groups []*GroupModel) error {
	if len(groups) == 0 {
		return nil
	}
	db := GetDb(ctx)
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

// UpdateSortById 根据 ID 更新 Group 的 Sort 字段
func (m *GroupMapper) UpdateSortById(db *gorm.DB, id int64, sort int32) error {
	err := db.Model(&GroupModel{}).
		Where("id = ?", id).
		Update("sort", sort).Error
	if err != nil {
		logx.Errorf("failed to update group sort: %v", err)
		return err
	}
	return nil
}

// CleanBizGroupById 清除分组下的所有业务关联（将业务的分组ID设为NULL）
func (m *GroupMapper) CleanBizGroupById(db *gorm.DB, id int64) error {
	// 执行三个更新语句，确保在同一个事务中执行
	err := db.Transaction(func(tx *gorm.DB) error {
		// 更新 uns_dashboard 表
		if err := tx.Exec(`UPDATE "uns_dashboard" SET group_id = NULL WHERE group_id = ?`, id).Error; err != nil {
			logx.Errorf("failed to update uns_dashboard group_id: %v", err)
			return err
		}

		// 更新 supos_node_flows 表
		if err := tx.Exec(`UPDATE "supos_node_flows" SET group_id = NULL WHERE group_id = ?`, id).Error; err != nil {
			logx.Errorf("failed to update supos_node_flows group_id: %v", err)
			return err
		}

		// 更新 supos_event_flows 表
		if err := tx.Exec(`UPDATE "supos_event_flows" SET group_id = NULL WHERE group_id = ?`, id).Error; err != nil {
			logx.Errorf("failed to update supos_event_flows group_id: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		logx.Errorf("failed to clean biz group by id: %v", err)
		return err
	}

	return nil
}

// OperationGroup 操作分组数据（移入移出）
func (m *GroupMapper) OperationGroup(db *gorm.DB, req *types.OperationGroupReq) error {
	// 查询分组信息
	var group GroupModel
	err := db.Where("id = ?", *req.ID).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFind
		}
		logx.Errorf("failed to select group by id: %v", err)
		return err
	}

	// 根据分组类型确定要操作的表
	var tableName string
	var flowType string
	switch *group.Type {
	case 1:
		tableName = "supos_node_flows"
		flowType = " and template = 'node-red'"
	case 2:
		tableName = "supos_node_flows"
		flowType = " and template = 'event-flow'"
	case 3:
		flowType = " "
		tableName = "uns_dashboard"
	default:
		return errors.Parameter
	}

	// 根据状态确定要设置的group_id值
	if *req.Status {
		// 移入分组：set group_id = req.ID
		err = db.Exec(
			"UPDATE "+tableName+" SET group_id = ? WHERE id = ?"+flowType,
			*req.ID, *req.BizId,
		).Error
	} else {
		// 移出分组：set group_id = NULL
		err = db.Exec(
			"UPDATE "+tableName+" SET group_id = NULL WHERE id = ?"+flowType,
			*req.BizId,
		).Error
	}

	if err != nil {
		logx.Errorf("failed to operation group: %v", err)
		return err
	}

	return nil
}

// 根据name查询 Group
func (m *GroupMapper) SelectByName(db *gorm.DB, name string) ([]*GroupModel, error) {
	var groups []*GroupModel
	err := db.Where("name = ?", name).Find(&groups).Error
	if err != nil {
		logx.Errorf("failed to select groups by name: %v", err)
		return nil, err
	}
	return groups, nil
}

// 根据name查询 Group not id equals
func (m *GroupMapper) SelectByNameNotId(db *gorm.DB, id int64, name string) ([]*GroupModel, error) {
	var groups []*GroupModel
	err := db.Where("name = ?", name).Where("id != ?", id).Find(&groups).Error
	if err != nil {
		logx.Errorf("failed to select groups by name: %v", err)
		return nil, err
	}
	return groups, nil
}
