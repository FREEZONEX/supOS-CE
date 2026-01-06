package relationDB

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DashboardRefMapper Dashboard 引用关系数据访问对象
type DashboardRefMapper struct {
}

// Insert 插入 Dashboard 引用关系
func (m DashboardRefMapper) Insert(db *gorm.DB, ref *DashboardRefModel) error {
	err := db.Model(&DashboardRefModel{}).Clauses(clause.OnConflict{DoNothing: true}).Create(ref).Error
	if err != nil {
		logx.Errorf("failed to insert dashboard ref: %v", err)
		return err
	}
	return nil
}
func (m DashboardRefMapper) SaveBatch(db *gorm.DB, refers []*DashboardRefModel) error {
	err := db.Model(&DashboardRefModel{}).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(refers, 1000).Error
	if err != nil {
		logx.Errorf("failed to insert dashboard ref: %v", err)
		return err
	}
	return nil
}
func (m DashboardRefMapper) DeleteByUnsAlias(db *gorm.DB, unsAlias string) error {
	err := db.Where("uns_alias = ?", unsAlias).Delete(&DashboardRefModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboard ref: %v, uns=%s", err, unsAlias)
		return err
	}
	return nil
}
func (m DashboardRefMapper) DeleteByUnsAliasList(db *gorm.DB, unsAlias []string) error {
	err := db.Where("uns_alias IN ?", unsAlias).Delete(&DashboardRefModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboard ref: %v, uns=%s", err, unsAlias)
		return err
	}
	return nil
}

// DeleteByDashboardId 根据 Dashboard ID 删除引用关系
func (m DashboardRefMapper) DeleteByDashboardId(db *gorm.DB, dashboardID string) error {
	err := db.Where("dashboard_id = ?", dashboardID).Delete(&DashboardRefModel{}).Error
	if err != nil {
		logx.Errorf("failed to delete dashboard ref: %v", err)
		return err
	}
	return nil
}

// GetByUns 根据 UNS 别名获取 Dashboard
func (m DashboardRefMapper) GetByUns(db *gorm.DB, unsAlias string) (*DashboardModel, error) {
	var dashboard DashboardModel
	err := db.
		Table("uns_dashboard a").
		Select("a.*").
		Joins("LEFT JOIN uns_dashboard_ref b ON a.id = b.dashboard_id").
		Where("b.uns_alias = ?", unsAlias).
		First(&dashboard).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logx.Errorf("failed to get dashboard by uns: %v", err)
		return nil, err
	}
	return &dashboard, nil
}

// SelectByUnsAlias 根据 UNS 别名查询引用关系
func (m DashboardRefMapper) SelectByUnsAlias(db *gorm.DB, unsAlias string) (*DashboardRefModel, error) {
	var ref DashboardRefModel
	err := db.Where("uns_alias = ?", unsAlias).First(&ref).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logx.Errorf("failed to select dashboard ref: %v", err)
		return nil, err
	}
	return &ref, nil
}

// SelectByUnsAliases selects dashboard references by a list of UNS aliases.
func (m DashboardRefMapper) SelectByUnsAliases(db *gorm.DB, aliases []string) ([]*DashboardRefModel, error) {
	if len(aliases) == 0 {
		return []*DashboardRefModel{}, nil
	}

	var refs []*DashboardRefModel
	err := db.Where("uns_alias IN ?", aliases).Find(&refs).Error
	if err != nil {
		logx.Errorf("failed to select dashboard refs by aliases: %v", err)
		return nil, err
	}
	return refs, nil
}
