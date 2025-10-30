package relationDB

import (
	"backend/share/base"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsLabelRefRepo struct {
}

func NewUnsLabelRefRepo() UnsLabelRefRepo {
	return UnsLabelRefRepo{}
}

type UnsLabelRefFilter struct {
	//todo 添加过滤字段
}

func (p UnsLabelRefRepo) model(db *gorm.DB) *gorm.DB {
	return db.Model(&UnsLabelRef{})
}
func (p UnsLabelRefRepo) Insert(db *gorm.DB, data *UnsLabelRef) error {
	result := p.model(db).Create(data)
	return stores.ErrFmt(result.Error)
}
func (p UnsLabelRefRepo) ListUnsIds(db *gorm.DB, labelId int64) (unsIds []int64, err error) {
	var result []UnsLabelRef
	err = p.model(db).Select("uns_id").Where("label_id=?", labelId).Find(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	} else if len(result) == 0 {
		return
	}
	unsIds = base.Map[UnsLabelRef, int64](result, func(e UnsLabelRef) int64 {
		return e.UnsID
	})
	return
}
func (p UnsLabelRefRepo) ListByUnsId(db *gorm.DB, unsId int64) (result []*UnsLabelRef, err error) {
	err = p.model(db).Where("uns_id=?", unsId).Find(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsLabelRefRepo) DeleteByUnsIds(db *gorm.DB, unsIds []int64) error {
	err := p.model(db).Where("uns_id in ?", unsIds).Delete(&UnsLabelRef{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRefRepo) DeleteByLabelIds(db *gorm.DB, labelIds []int64) error {
	err := p.model(db).Where("label_id in ?", labelIds).Delete(&UnsLabelRef{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRefRepo) DeleteByUnsIdAndLabelIds(db *gorm.DB, unsId int64, labelIds []int64) error {
	err := p.model(db).Where("uns_id = ?", unsId).Where("label_id in ?", labelIds).Delete(&UnsLabelRef{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRefRepo) FindOneByFilter(db *gorm.DB, f UnsLabelRefFilter) (*UnsLabelRef, error) {
	var result UnsLabelRef
	err := p.model(db).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsLabelRefRepo) FindByFilter(db *gorm.DB, f UnsLabelRefFilter, page *stores.PageInfo) ([]*UnsLabelRef, error) {
	var results []*UnsLabelRef
	db = p.model(db)
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UnsLabelRefRepo) CountByFilter(db *gorm.DB, f UnsLabelRefFilter) (size int64, err error) {
	err = p.model(db).Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsLabelRefRepo) Update(db *gorm.DB, data *UnsLabelRef) error {
	// 组合主键，直接Save
	err := p.model(db).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UnsLabelRefRepo) DeleteByFilter(db *gorm.DB, f UnsLabelRefFilter) error {
	err := p.model(db).Delete(&UnsLabelRef{}).Error
	return stores.ErrFmt(err)
}

func (p UnsLabelRefRepo) FindOne(db *gorm.DB, labelID int64, unsID int64) (*UnsLabelRef, error) {
	var result UnsLabelRef
	err := p.model(db).Where("label_id = ? AND uns_id = ?", labelID, unsID).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入记录
func (p UnsLabelRefRepo) SaveOrIgnore(db *gorm.DB, data []*UnsLabelRef) error {
	err := p.model(db).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, 1000).Error
	return stores.ErrFmt(err)
}
