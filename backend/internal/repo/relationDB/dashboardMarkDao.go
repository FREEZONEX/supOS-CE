package relationDB

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

// DashboardMarkedMapper Dashboard 置顶标记数据访问对象
type DashboardMarkedMapper struct {
}

// Insert 插入置顶标记
func (m *DashboardMarkedMapper) Insert(ctx context.Context, mark *DashboardMarkModel) error {
	err := GetDb(ctx).Create(mark).Error
	if err != nil {
		logx.Errorf("failed to insert dashboard mark: %v", err)
		return err
	}
	return nil
}

// Delete 删除置顶标记
func (m *DashboardMarkedMapper) Delete(ctx context.Context, id string, userID string) error {
	err := GetDb(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&DashboardMarkModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboard mark: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 Dashboard ID 删除所有置顶标记
func (m *DashboardMarkedMapper) DeleteById(ctx context.Context, id string) error {
	err := GetDb(ctx).Where("id = ?", id).Delete(&DashboardMarkModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboard mark by id: %v", err)
		return err
	}
	return nil
}
