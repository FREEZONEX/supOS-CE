package relationDB

import (
	"backend/internal/types"
	"context"
	"time"

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
func (m *GroupMapper) SelectByType(db *gorm.DB, typ int16, name string) ([]*GroupModel, error) {
	var groups []*GroupModel

	query := db.Model(&GroupModel{})
	query = db.Where("type = ? ", typ)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&groups).Error
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

// UpdateSortById 根据 ID 更新 Group 的 Sort 字段和 MarkTime 字段
func (m *GroupMapper) UpdateSortById(db *gorm.DB, id int64, sort int32) error {
	updateData := map[string]interface{}{
		"sort": sort,
	}

	if sort == 1 {
		// 置顶时，设置 mark_time 为当前系统时间
		updateData["mark_time"] = time.Now()
	} else {
		// 取消置顶时，设置 mark_time 为 null
		updateData["mark_time"] = nil
	}

	err := db.Model(&GroupModel{}).
		Where("id = ?", id).
		Updates(updateData).Error
	if err != nil {
		logx.Errorf("failed to update group sort: %v", err)
		return err
	}
	return nil
}

/**
  清除分组下的所有业务关联
  groupType 1 对应 supos_node_flows 表   deleteSourceFlowLogic.go
  groupType 2 对应 supos_node_flows 表   deleteEventFlowLogic.go
  groupType 3 对应 uns_dashboard 表 deleteLogic.go
  先通过groupType判断对应的业务表   通过groupId查询所有的业务数据，执行对应的delete方法
*/
// CleanBizGroupById 查询分组下的所有业务数据
func (m *GroupMapper) CleanBizGroupById(db *gorm.DB, groupId int64, groupType *int16) (
	sourceFlows []*NoderedSourceFlow,
	eventFlows []*NoderedSourceFlow,
	dashboards []*DashboardModel,
	err error,
) {
	if groupType == nil {
		logx.Infof("groupType is nil, cannot determine business table")
		return
	}

	// 根据 groupType 执行对应的查询操作
	switch *groupType {
	case 1: // groupType 1 对应 supos_node_flows 表 (source flow)
		// 查询该分组下的所有 source flow
		err = db.Where("flow_id is not null AND group_id = ? AND template = ?", groupId, "node-red").Find(&sourceFlows).Error
		if err != nil {
			logx.Errorf("failed to select source flows by group id: %v", err)
			return
		}

	case 2: // groupType 2 对应 supos_node_flows 表 (event flow)
		// 查询该分组下的所有 event flow
		err = db.Where("flow_id is not null AND group_id = ? AND template = ?", groupId, "event-flow").Find(&eventFlows).Error
		if err != nil {
			logx.Errorf("failed to select event flows by group id: %v", err)
			return
		}

	case 3: // groupType 3 对应 uns_dashboard 表 (dashboard)
		// 查询该分组下的所有 dashboard
		err = db.Where("group_id = ?", groupId).Find(&dashboards).Error
		if err != nil {
			logx.Errorf("failed to select dashboards by group id: %v", err)
			return
		}

	default:
		logx.Infof("unknown group type: %d", *groupType)
	}

	return
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
func (m *GroupMapper) SelectByName(db *gorm.DB, name string, typ int16) ([]*GroupModel, error) {
	var groups []*GroupModel
	err := db.Where("name = ? and type = ? ", name, typ).Find(&groups).Error
	if err != nil {
		logx.Errorf("failed to select groups by name: %v", err)
		return nil, err
	}
	return groups, nil
}

// 根据name查询 Group not id equals
func (m *GroupMapper) SelectByNameNotId(db *gorm.DB, id int64, name string , typ int16) ([]*GroupModel, error) {
	var groups []*GroupModel
	err := db.Where("name = ? and type = ?", name, typ).Where("id != ?", id).Find(&groups).Error
	if err != nil {
		logx.Errorf("failed to select groups by name: %v", err)
		return nil, err
	}
	return groups, nil
}
