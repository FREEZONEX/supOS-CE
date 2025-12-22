package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsAttachmentRepo struct {
	db *gorm.DB
}

func NewUnsAttachmentRepo(in context.Context) *UnsAttachmentRepo {
	return &UnsAttachmentRepo{db: GetDb(in)}
}

type UnsAttachmentFilter struct {
	UnsAlias       string // 关联模型实例别名
	AttachmentPath string // 附件存储路径
}

func (p UnsAttachmentRepo) fmtFilter(ctx context.Context, f UnsAttachmentFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UnsAlias != "" {
		db = db.Where("uns_alias = ?", f.UnsAlias)
	}
	if f.AttachmentPath != "" {
		db = db.Where("attachment_path = ?", f.AttachmentPath)
	}
	return db
}

func (p UnsAttachmentRepo) Insert(ctx context.Context, data *UnsAttachment) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsAttachmentRepo) FindOneByFilter(ctx context.Context, f UnsAttachmentFilter) (*UnsAttachment, error) {
	var result UnsAttachment
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsAttachmentRepo) FindByFilter(ctx context.Context, f UnsAttachmentFilter, page *stores.PageInfo) ([]*UnsAttachment, error) {
	var results []*UnsAttachment
	db := p.fmtFilter(ctx, f).Model(&UnsAttachment{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsAttachmentRepo) CountByFilter(ctx context.Context, f UnsAttachmentFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&UnsAttachment{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsAttachmentRepo) Update(ctx context.Context, data *UnsAttachment) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsAttachmentRepo) DeleteByFilter(ctx context.Context, f UnsAttachmentFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&UnsAttachment{}).Error
	return stores.ErrFmt(err)
}

func (p UnsAttachmentRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&UnsAttachment{}).Error
	return stores.ErrFmt(err)
}
func (p UnsAttachmentRepo) FindOne(ctx context.Context, id int64) (*UnsAttachment, error) {
	var result UnsAttachment
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsAttachmentRepo) MultiInsert(ctx context.Context, data []*UnsAttachment) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsAttachment{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsAttachmentRepo) UpdateWithField(ctx context.Context, f UnsAttachmentFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&UnsAttachment{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

// FindByUnsAlias 根据模型实例别名查询附件列表
func (p UnsAttachmentRepo) FindByUnsAlias(ctx context.Context, unsAlias string) ([]*UnsAttachment, error) {
	var results []*UnsAttachment
	err := p.db.WithContext(ctx).
		Where("uns_alias = ?", unsAlias).
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

// FindByAttachmentPath 根据附件存储路径查询附件列表
func (p UnsAttachmentRepo) FindByAttachmentPath(ctx context.Context, attachmentPath string) ([]*UnsAttachment, error) {
	var results []*UnsAttachment
	err := p.db.WithContext(ctx).
		Where("attachment_path = ?", attachmentPath).
		Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
