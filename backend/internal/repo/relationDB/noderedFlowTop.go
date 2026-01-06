package relationDB

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将NoderedFlowTop全局替换为模型的表名
2. 完善todo
*/

type NoderedFlowTopRepo struct {
	db *gorm.DB
}

func NewNoderedFlowTopRepo(in context.Context) *NoderedFlowTopRepo {
	return &NoderedFlowTopRepo{db: GetDb(in)}
}

type NoderedFlowTopFilter struct {
	//todo 添加过滤字段
}

func (p NoderedFlowTopRepo) fmtFilter(ctx context.Context, f NoderedFlowTopFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p NoderedFlowTopRepo) Insert(ctx context.Context, data *NoderedFlowTop) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NoderedFlowTopRepo) FindOneByFilter(ctx context.Context, f NoderedFlowTopFilter) (*NoderedFlowTop, error) {
	var result NoderedFlowTop
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NoderedFlowTopRepo) FindByFilter(ctx context.Context, f NoderedFlowTopFilter, page *stores.PageInfo) ([]*NoderedFlowTop, error) {
	var results []*NoderedFlowTop
	db := p.fmtFilter(ctx, f).Model(&NoderedFlowTop{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NoderedFlowTopRepo) CountByFilter(ctx context.Context, f NoderedFlowTopFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&NoderedFlowTop{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NoderedFlowTopRepo) Update(ctx context.Context, data *NoderedFlowTop) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p NoderedFlowTopRepo) DeleteByFilter(ctx context.Context, f NoderedFlowTopFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&NoderedFlowTop{}).Error
	return stores.ErrFmt(err)
}

func (p NoderedFlowTopRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&NoderedFlowTop{}).Error
	return stores.ErrFmt(err)
}
func (p NoderedFlowTopRepo) FindOne(ctx context.Context, id int64) (*NoderedFlowTop, error) {
	var result NoderedFlowTop
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NoderedFlowTopRepo) MultiInsert(ctx context.Context, data []*NoderedFlowTop) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&NoderedFlowTop{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d NoderedFlowTopRepo) UpdateWithField(ctx context.Context, f NoderedFlowTopFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&NoderedFlowTop{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

// Upsert marks the flow as pinned for the specified user.
func (r NoderedFlowTopRepo) Upsert(ctx context.Context, flowID int64, userID string, mark int) error {
	if flowID <= 0 {
		return stores.ErrFmt(gorm.ErrInvalidData)
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return stores.ErrFmt(gorm.ErrInvalidData)
	}
	top := &NoderedFlowTop{
		ID:     flowID,
		UserID: userID,
		Mark:   mark,
	}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{"mark": mark, "mark_time": gorm.Expr("CURRENT_TIMESTAMP"), "update_time": gorm.Expr("CURRENT_TIMESTAMP")}),
		}).
		Create(top).Error
	return stores.ErrFmt(err)
}

// DeleteByUser removes the pinned record for the given user.
func (r NoderedFlowTopRepo) DeleteByUser(ctx context.Context, flowID int64, userID string) error {
	userID = strings.TrimSpace(userID)
	if flowID <= 0 || userID == "" {
		return stores.ErrFmt(gorm.ErrInvalidData)
	}
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", flowID, userID).Delete(&NoderedFlowTop{}).Error
	return stores.ErrFmt(err)
}
