package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsHistoryDeleteJobRepo struct {
	db *gorm.DB
}

func NewUnsHistoryDeleteJobRepo(in any) *UnsHistoryDeleteJobRepo {
	return &UnsHistoryDeleteJobRepo{db: stores.GetCommonConn(in)}
}

type UnsHistoryDeleteJobFilter struct {
	//todo 添加过滤字段
}

func (p UnsHistoryDeleteJobRepo) fmtFilter(ctx context.Context, f UnsHistoryDeleteJobFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	//todo 添加条件
	return db
}

func (p UnsHistoryDeleteJobRepo) Insert(ctx context.Context, data *UnsHistoryDeleteJob) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsHistoryDeleteJobRepo) FindOneByFilter(ctx context.Context, f UnsHistoryDeleteJobFilter) (*UnsHistoryDeleteJob, error) {
	var result UnsHistoryDeleteJob
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsHistoryDeleteJobRepo) FindByFilter(ctx context.Context, f UnsHistoryDeleteJobFilter, page *stores.PageInfo) ([]*UnsHistoryDeleteJob, error) {
	var results []*UnsHistoryDeleteJob
	db := p.fmtFilter(ctx, f).Model(&UnsHistoryDeleteJob{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsHistoryDeleteJobRepo) CountByFilter(ctx context.Context, f UnsHistoryDeleteJobFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsHistoryDeleteJob{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsHistoryDeleteJobRepo) Update(ctx context.Context, data *UnsHistoryDeleteJob) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsHistoryDeleteJobRepo) DeleteByFilter(ctx context.Context, f UnsHistoryDeleteJobFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsHistoryDeleteJob{}).Error
	return stores.ErrFmt(err)
}

func (p UnsHistoryDeleteJobRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsHistoryDeleteJob{}).Error
	return stores.ErrFmt(err)
}
func (p UnsHistoryDeleteJobRepo) FindOne(ctx context.Context, id int64) (*UnsHistoryDeleteJob, error) {
	var result UnsHistoryDeleteJob
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsHistoryDeleteJobRepo) MultiInsert(ctx context.Context, data []*UnsHistoryDeleteJob) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsHistoryDeleteJob{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsHistoryDeleteJobRepo) UpdateWithField(ctx context.Context, f UnsHistoryDeleteJobFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsHistoryDeleteJob{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
