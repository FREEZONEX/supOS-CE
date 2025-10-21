package relationDB

import (
	"time"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnsLabelRepo struct {
}

func NewUnsLabelRepo() UnsLabelRepo {
	return UnsLabelRepo{}
}

type UnsLabelFilter struct {
	//todo 添加过滤字段
	LabelName string
}

func (p UnsLabelRepo) fmtFilter(db *gorm.DB, f UnsLabelFilter) *gorm.DB {
	//todo 添加条件
	if f.LabelName != "" {
		db = db.Where("label_name = ?", f.LabelName)
	}
	return db
}

func (p UnsLabelRepo) Insert(db *gorm.DB, data *UnsLabel) error {
	result := db.Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UnsLabelRepo) FindOneByFilter(db *gorm.DB, f UnsLabelFilter) (*UnsLabel, error) {
	var result UnsLabel
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UnsLabelRepo) FindByFilter(db *gorm.DB, f UnsLabelFilter, page *stores.PageInfo) ([]*UnsLabel, error) {
	var results []*UnsLabel
	db = db.Model(&UnsLabel{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsLabelRepo) FindByNames(db *gorm.DB, names []string) ([]*UnsLabel, error) {
	var results []*UnsLabel
	db = db.Where("label_name in ? ", names).Model(&UnsLabel{})
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsLabelRepo) FindByName(db *gorm.DB, name string) (*UnsLabel, error) {
	var label UnsLabel
	db = db.Where("label_name = ? ", name).Model(&UnsLabel{})
	err := db.Find(&label).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &label, nil
}

func (p UnsLabelRepo) CountByFilter(db *gorm.DB, f UnsLabelFilter) (size int64, err error) {
	err = db.Model(&UnsLabel{}).Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsLabelRepo) Update(db *gorm.DB, data *UnsLabel) error {
	err := db.Model(&UnsLabel{}).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

//func (p UnsLabelRepo) DeleteByFilter(db *gorm.DB, f UnsLabelFilter) error {
//	err := db.Model(&UnsLabel{}).Delete(&UnsLabel{}).Error
//	return stores.ErrFmt(err)
//}

func (p UnsLabelRepo) Delete(db *gorm.DB, id int64) error {
	err := db.Model(&UnsLabel{}).Where("id = ?", id).Delete(&UnsLabel{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRepo) DeleteRefByLabelId(db *gorm.DB, id int64) error {
	err := db.Model(&UnsLabel{}).Where("label_id = ?", id).Delete(&UnsLabel{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRepo) DeleteRefByUnsId(db *gorm.DB, unsId int64) error {
	err := db.Model(&UnsLabel{}).Where("uns_id = ?", unsId).Delete(&UnsLabel{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRepo) FindOne(db *gorm.DB, id int64) (*UnsLabel, error) {
	var result UnsLabel
	err := db.Model(&UnsLabel{}).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UnsLabelRepo) MultiInsert(db *gorm.DB, data []*UnsLabel) error {
	err := db.Model(&UnsLabel{}).Clauses(clause.OnConflict{UpdateAll: true}).Model(&UnsLabel{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d UnsLabelRepo) UpdateWithField(db *gorm.DB, f UnsLabelFilter, updates map[string]any) error {
	err := db.Model(&UnsLabel{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

// GORM hooks
// AfterUpdate: touch update_at to current time to ensure timestamp consistency
func (u *UnsLabel) AfterUpdate(tx *gorm.DB) (err error) {
	if u == nil || u.ID == 0 {
		return nil
	}
	// Skip hooks to avoid recursion
	if err = tx.Session(&gorm.Session{SkipHooks: true}).Model(&UnsLabel{}).
		Where("id = ?", u.ID).
		Update("update_at", time.Now()).Error; err != nil {
		return stores.ErrFmt(err)
	}
	return nil
}

// AfterDelete: cascade delete label refs to avoid orphaned rows
func (u *UnsLabel) AfterDelete(tx *gorm.DB) (err error) {
	if u == nil || u.ID == 0 {
		return nil
	}
	if err = tx.Where("label_id = ?", u.ID).Delete(&UnsLabelRef{}).Error; err != nil {
		return stores.ErrFmt(err)
	}
	return nil
}
