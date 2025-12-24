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
func (p UnsLabelRepo) SelectById(db *gorm.DB, id int64) (*UnsLabel, error) {
	var result UnsLabel
	err := db.Model(&UnsLabel{}).Where("id=?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	} else if result.ID == 0 {
		return nil, nil
	}
	return &result, nil
}
func (p UnsLabelRepo) ListByIds(db *gorm.DB, ids []int64) ([]*UnsLabel, error) {
	var results []*UnsLabel
	db = db.Model(&UnsLabel{}).Where("id in ? ", ids)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsLabelRepo) ListAll(db *gorm.DB, page, pageSize int) ([]*UnsLabel, error) {
	var results []*UnsLabel
	if page < 1 {
		page = 1
	}
	if pageSize > 1000 {
		pageSize = 1000
	} else if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	db = db.Model(&UnsLabel{}).Order("id").Offset(offset).Limit(pageSize)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
func (p UnsLabelRepo) ListByUnsId(db *gorm.DB, unsId int64) (result []*UnsLabel, err error) {
	err = db.Model(&UnsLabel{}).Raw(`
       select ul.id ,ul.label_name from uns_label ul join  uns_label_ref rf on ul.id = rf.label_id where rf.uns_id =? 
       `, unsId).
		Find(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return
}
func (p UnsLabelRepo) Insert(db *gorm.DB, data *UnsLabel) error {
	result := db.Model(&UnsLabel{}).Create(data)
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
func (p UnsLabelRepo) LikeName(db *gorm.DB, key string) ([]*UnsLabel, error) {
	var results []*UnsLabel
	db = db.Model(&UnsLabel{})
	if len(key) > 0 {
		db = db.Where("label_name like '%" + escapeLikePattern(escapeSQL(key)) + "%' ")
	}
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
	err := db.First(&label).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &label, nil
}
func (p UnsLabelRepo) CountByName(db *gorm.DB, name string) (count int64, err error) {
	db = db.Where("label_name = ? ", name).Model(&UnsLabel{})
	err = db.Count(&count).Error
	if err != nil {
		return count, stores.ErrFmt(err)
	}
	return
}
func (p UnsLabelRepo) CountByFilter(db *gorm.DB, f UnsLabelFilter) (size int64, err error) {
	err = db.Model(&UnsLabel{}).Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UnsLabelRepo) Update(db *gorm.DB, data *UnsLabel) error {
	err := db.Model(&UnsLabel{}).Where("id = ?", data.ID).Omit("id").Updates(data).Error
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
	err := db.Model(&UnsLabelRef{}).Where("label_id = ?", id).Delete(&UnsLabelRef{}).Error
	return stores.ErrFmt(err)
}
func (p UnsLabelRepo) DeleteRefByUnsId(db *gorm.DB, unsId int64) error {
	err := db.Model(&UnsLabelRef{}).Where("uns_id = ?", unsId).Delete(&UnsLabelRef{}).Error
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

func (p UnsLabelRepo) UpdateWithField(db *gorm.DB, f UnsLabelFilter, updates map[string]any) error {
	err := db.Model(&UnsLabel{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

// GORM hooks
// AfterUpdate: touch update_at to current time to ensure timestamp consistency
func (l *UnsLabel) AfterUpdate(tx *gorm.DB) (err error) {
	if l == nil || l.ID == 0 {
		return nil
	}
	// Skip hooks to avoid recursion
	if err = tx.Session(&gorm.Session{SkipHooks: true}).Model(&UnsLabel{}).
		Where("id = ?", l.ID).
		Update("update_at", time.Now()).Error; err != nil {
		return stores.ErrFmt(err)
	}
	return nil
}

// AfterDelete: cascade delete label refs to avoid orphaned rows
func (l *UnsLabel) AfterDelete(tx *gorm.DB) (err error) {
	if l == nil || l.ID == 0 {
		return nil
	}
	if err = tx.Where("label_id = ?", l.ID).Delete(&UnsLabelRef{}).Error; err != nil {
		return stores.ErrFmt(err)
	}
	return nil
}
