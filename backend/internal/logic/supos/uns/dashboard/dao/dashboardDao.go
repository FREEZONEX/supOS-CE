package dao

import (
	"backend/internal/logic/supos/uns/dashboard/model"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DashboardMapper Dashboard 数据访问对象
type DashboardMapper struct {
	db     *gorm.DB
	ctx    context.Context
	logger logx.Logger
}

// NewDashboardMapper 创建 DashboardMapper 实例
func NewDashboardMapper(db *gorm.DB, ctx context.Context) *DashboardMapper {
	return &DashboardMapper{
		db:     db,
		ctx:    ctx,
		logger: logx.WithContext(ctx),
	}
}

// SelectById 根据 ID 查询 Dashboard
func (m *DashboardMapper) SelectById(id string) (*model.DashboardModel, error) {
	var dashboard model.DashboardModel
	err := m.db.WithContext(m.ctx).Where("id = ?", id).First(&dashboard).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		m.logger.Errorf("failed to select dashboard by id: %v", err)
		return nil, err
	}
	return &dashboard, nil
}

// Insert 插入 Dashboard
func (m *DashboardMapper) Insert(dashboard *model.DashboardModel) error {
	err := m.db.WithContext(m.ctx).Create(dashboard).Error
	if err != nil {
		m.logger.Errorf("failed to insert dashboard: %v", err)
		return err
	}
	return nil
}

// UpdateById 根据 ID 更新 Dashboard
func (m *DashboardMapper) UpdateById(dashboard *model.DashboardModel) error {
	// 使用 map 更新非零值字段，避免gorm默认的“忽略零值”行为
	// 这里假设所有字段都需要更新
	err := m.db.WithContext(m.ctx).Model(&model.DashboardModel{}).Where("id = ?", dashboard.ID).Updates(dashboard).Error
	if err != nil {
		m.logger.Errorf("failed to update dashboard: %v", err)
		return err
	}
	return nil
}

// DeleteById 根据 ID 删除 Dashboard
func (m *DashboardMapper) DeleteById(id string) error {
	err := m.db.WithContext(m.ctx).Where("id = ?", id).Delete(&model.DashboardModel{}).Error
	if err != nil {
		m.logger.Errorf("failed to delete dashboard: %v", err)
		return err
	}
	return nil
}

// SelectByFlowNames 根据名称列表查询 Dashboard
func (m *DashboardMapper) SelectByFlowNames(names []string) ([]*model.DashboardModel, error) {
	if len(names) == 0 {
		return []*model.DashboardModel{}, nil
	}
	var dashboards []*model.DashboardModel
	err := m.db.WithContext(m.ctx).Where("name IN ?", names).Find(&dashboards).Error
	if err != nil {
		m.logger.Errorf("failed to select dashboards by names: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SaveOrIgnoreBatch 批量保存或忽略
func (m *DashboardMapper) SaveOrIgnoreBatch(dashboards []*model.DashboardModel) error {
	if len(dashboards) == 0 {
		return nil
	}
	// GORM v2 的 Clauses(clause.OnConflict{DoNothing: true}) 提供了优雅的方式
	err := m.db.WithContext(m.ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&dashboards).Error
	if err != nil {
		m.logger.Errorf("failed to batch save or ignore dashboards: %v", err)
		return err
	}
	return nil
}

// DashboardExtends Dashboard 扩展信息（包含置顶标记）
type DashboardExtends struct {
	model.DashboardModel
	Mark     *int   `db:"mark" json:"mark,omitzero"`          // 置顶标记
	MarkTime *int64 `db:"mark_time" json:"markTime,omitzero"` // 置顶时间
}

// SelectDashboard 分页查询 Dashboard（包含置顶信息）
func (m *DashboardMapper) SelectDashboard(
	userID string,
	fuzzyName string,
	typ *int,
	orderCode string,
	descOrAsc string,
	pageNo int64,
	pageSize int64,
) ([]*DashboardExtends, error) {
	var dashboards []*DashboardExtends
	query := m.db.WithContext(m.ctx).
		Table("uns_dashboard a").
		Select("a.*, b.mark, b.mark_time").
		Joins("LEFT JOIN uns_dashboard_top_recodes b ON a.id = b.id AND b.user_id = ?", userID)

	if fuzzyName != "" {
		searchPattern := "%" + fuzzyName + "%"
		query = query.Where("a.name LIKE ? OR a.description LIKE ?", searchPattern, searchPattern)
	}

	if typ != nil {
		query = query.Where("a.type = ?", *typ)
	}

	// 排序
	if orderCode == "" {
		query = query.Order("b.mark ASC, b.mark_time DESC, a.create_time DESC")
	} else {
		query = query.Order(fmt.Sprintf("%s %s", orderCode, descOrAsc))
	}

	// 分页
	offset := (pageNo - 1) * pageSize
	query = query.Limit(int(pageSize)).Offset(int(offset))

	err := query.Find(&dashboards).Error
	if err != nil {
		m.logger.Errorf("failed to select dashboards: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SelectDashboardCount 查询 Dashboard 总数
func (m *DashboardMapper) SelectDashboardCount(fuzzyName string, typ *int) (int64, error) {
	var count int64
	query := m.db.WithContext(m.ctx).Model(&model.DashboardModel{})

	if fuzzyName != "" {
		searchPattern := "%" + fuzzyName + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	if typ != nil {
		query = query.Where("type = ?", *typ)
	}

	err := query.Count(&count).Error
	if err != nil {
		m.logger.Errorf("failed to count dashboards: %v", err)
		return 0, err
	}
	return count, nil
}

// SelectAll selects all DashboardModel from the database.
func (m *DashboardMapper) SelectAll() ([]*model.DashboardModel, error) {
	var dashboards []*model.DashboardModel
	err := m.db.WithContext(m.ctx).Find(&dashboards).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.DashboardModel{}, nil
		}
		m.logger.Errorf("failed to select all dashboards: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// SelectByIds selects multiple DashboardModel from the database by their IDs.
func (m *DashboardMapper) SelectByIds(ids []string) ([]*model.DashboardModel, error) {
	if len(ids) == 0 {
		return []*model.DashboardModel{}, nil
	}
	var dashboards []*model.DashboardModel
	err := m.db.WithContext(m.ctx).Where("id IN ?", ids).Find(&dashboards).Error
	if err != nil {
		m.logger.Errorf("failed to select dashboards by ids: %v", err)
		return nil, err
	}
	return dashboards, nil
}

// DeleteBatchIds deletes multiple dashboards from the database by their IDs.
func (m *DashboardMapper) DeleteBatchIds(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	err := m.db.WithContext(m.ctx).Where("id IN ?", ids).Delete(&model.DashboardModel{}).Error
	if err != nil {
		m.logger.Errorf("failed to delete dashboards by ids: %v", err)
		return err
	}
	return nil
}

// SelectDashboardsToInit selects dashboards that need to be initialized.
func (m *DashboardMapper) SelectDashboardsToInit() ([]*model.DashboardModel, error) {
	var dashboards []*model.DashboardModel
	err := m.db.WithContext(m.ctx).
		Where("need_init = ? AND type = ? AND json_content IS NOT NULL AND json_content != ?", true, 1, "").
		Find(&dashboards).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.DashboardModel{}, nil
		}
		m.logger.Errorf("failed to select dashboards to init: %v", err)
		return nil, err
	}
	return dashboards, nil
}
