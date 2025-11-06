package relationDB

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// DashboardMarkedMapper Dashboard 置顶标记数据访问对象
type DashboardMarkedMapper struct {
	db     *gorm.DB
	ctx    context.Context
	logger logx.Logger
}

// NewDashboardMarkedMapper 创建 DashboardMarkedMapper 实例
func NewDashboardMarkedMapper(db *gorm.DB, ctx context.Context) *DashboardMarkedMapper {
	return &DashboardMarkedMapper{
		db:     db,
		ctx:    ctx,
		logger: logx.WithContext(ctx),
	}
}

// Insert 插入置顶标记
func (m *DashboardMarkedMapper) Insert(mark *DashboardMarkModel) error {
	err := m.db.WithContext(m.ctx).Create(mark).Error
	if err != nil {
		m.logger.Errorf("failed to insert dashboard mark: %v", err)
		return err
	}
	return nil
}

// Delete 删除置顶标记
func (m *DashboardMarkedMapper) Delete(id string, userID string) error {
	err := m.db.WithContext(m.ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&DashboardMarkModel{}).Error
	if err != nil {
		m.logger.Errorf("failed to delete dashboard mark: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 Dashboard ID 删除所有置顶标记
func (m *DashboardMarkedMapper) DeleteById(id string) error {
	err := m.db.WithContext(m.ctx).Where("id = ?", id).Delete(&DashboardMarkModel{}).Error
	if err != nil {
		m.logger.Errorf("failed to delete dashboard mark by id: %v", err)
		return err
	}
	return nil
}
