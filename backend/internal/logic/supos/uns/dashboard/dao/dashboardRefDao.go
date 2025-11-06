package dao

import (
	"backend/internal/logic/supos/uns/dashboard/model"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// DashboardRefMapper Dashboard 引用关系数据访问对象
type DashboardRefMapper struct {
	db     *gorm.DB
	ctx    context.Context
	logger logx.Logger
}

// NewDashboardRefMapper 创建 DashboardRefMapper 实例
func NewDashboardRefMapper(db *gorm.DB, ctx context.Context) *DashboardRefMapper {
	return &DashboardRefMapper{
		db:     db,
		ctx:    ctx,
		logger: logx.WithContext(ctx),
	}
}

// Insert 插入 Dashboard 引用关系
func (m *DashboardRefMapper) Insert(ref *model.DashboardRefModel) error {
	err := m.db.WithContext(m.ctx).Create(ref).Error
	if err != nil {
		m.logger.Errorf("failed to insert dashboard ref: %v", err)
		return err
	}
	return nil
}

// DeleteByDashboardId 根据 Dashboard ID 删除引用关系
func (m *DashboardRefMapper) DeleteByDashboardId(dashboardID string) error {
	err := m.db.WithContext(m.ctx).Where("dashboard_id = ?", dashboardID).Delete(&model.DashboardRefModel{}).Error
	if err != nil {
		m.logger.Errorf("failed to delete dashboard ref: %v", err)
		return err
	}
	return nil
}

// GetByUns 根据 UNS 别名获取 Dashboard
func (m *DashboardRefMapper) GetByUns(unsAlias string) (*model.DashboardModel, error) {
	var dashboard model.DashboardModel
	err := m.db.WithContext(m.ctx).
		Table("uns_dashboard a").
		Select("a.*").
		Joins("LEFT JOIN uns_dashboard_ref b ON a.id = b.dashboard_id").
		Where("b.uns_alias = ?", unsAlias).
		First(&dashboard).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		m.logger.Errorf("failed to get dashboard by uns: %v", err)
		return nil, err
	}
	return &dashboard, nil
}

// SelectByUnsAlias 根据 UNS 别名查询引用关系
func (m *DashboardRefMapper) SelectByUnsAlias(unsAlias string) (*model.DashboardRefModel, error) {
	var ref model.DashboardRefModel
	err := m.db.WithContext(m.ctx).Where("uns_alias = ?", unsAlias).First(&ref).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		m.logger.Errorf("failed to select dashboard ref: %v", err)
		return nil, err
	}
	return &ref, nil
}

// SelectByUnsAliases selects dashboard references by a list of UNS aliases.
func (m *DashboardRefMapper) SelectByUnsAliases(aliases []string) ([]*model.DashboardRefModel, error) {
	if len(aliases) == 0 {
		return []*model.DashboardRefModel{}, nil
	}

	var refs []*model.DashboardRefModel
	err := m.db.WithContext(m.ctx).Where("uns_alias IN ?", aliases).Find(&refs).Error
	if err != nil {
		m.logger.Errorf("failed to select dashboard refs by aliases: %v", err)
		return nil, err
	}
	return refs, nil
}
