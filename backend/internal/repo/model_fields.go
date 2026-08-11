package repo

import (
	"time"

	"backend/internal/contextx"

	"gorm.io/gorm"
)

type CreateTime struct {
	CreatedTime time.Time `gorm:"column:created_time;index;sort:desc;autoCreateTime" json:"createdTime"`
}

type OnlyTime struct {
	CreateTime
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime" json:"updatedTime"`
}

type DeletedTimeField struct {
	DeletedTime DeletedTime `gorm:"column:deleted_time;default:0" json:"deletedTime"`
}

type NoDelTime struct {
	CreatedBy int64 `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy int64 `gorm:"column:updated_by" json:"updatedBy"`
	OnlyTime
}

type SoftTime struct {
	NoDelTime
	DeletedBy int64 `gorm:"column:deleted_by" json:"deletedBy"`
	DeletedTimeField
}

type SoftNoDelByTime struct {
	NoDelTime
	DeletedTimeField
}

type SoftOnlyTime struct {
	OnlyTime
	DeletedTimeField
}

type CreatorTime struct {
	CreatedBy int64 `gorm:"column:created_by" json:"createdBy"`
	OnlyTime
}

type CreatorSoftTime struct {
	CreatorTime
	DeletedTimeField
}

type CreatorCreateTime struct {
	CreatedBy int64 `gorm:"column:created_by" json:"createdBy"`
	CreateTime
}

type CreatorCreateSoftTime struct {
	CreatorCreateTime
	DeletedTimeField
}

func (m *NoDelTime) BeforeCreate(tx *gorm.DB) error {
	fillCreateActor(tx, &m.CreatedBy, &m.UpdatedBy)
	return nil
}

func (m *NoDelTime) BeforeUpdate(tx *gorm.DB) error {
	fillUpdateActor(tx, &m.UpdatedBy)
	return nil
}

func (m *SoftTime) BeforeCreate(tx *gorm.DB) error {
	return m.NoDelTime.BeforeCreate(tx)
}

func (m *SoftTime) BeforeUpdate(tx *gorm.DB) error {
	return m.NoDelTime.BeforeUpdate(tx)
}

func (m *SoftNoDelByTime) BeforeCreate(tx *gorm.DB) error {
	return m.NoDelTime.BeforeCreate(tx)
}

func (m *SoftNoDelByTime) BeforeUpdate(tx *gorm.DB) error {
	return m.NoDelTime.BeforeUpdate(tx)
}

func (m *CreatorTime) BeforeCreate(tx *gorm.DB) error {
	fillCreateOnlyActor(tx, &m.CreatedBy)
	return nil
}

func (m *CreatorSoftTime) BeforeCreate(tx *gorm.DB) error {
	return m.CreatorTime.BeforeCreate(tx)
}

func (m *CreatorCreateTime) BeforeCreate(tx *gorm.DB) error {
	fillCreateOnlyActor(tx, &m.CreatedBy)
	return nil
}

func (m *CreatorCreateSoftTime) BeforeCreate(tx *gorm.DB) error {
	return m.CreatorCreateTime.BeforeCreate(tx)
}

func repoNowMilli() int64 {
	return time.Now().UTC().UnixMilli()
}

func repoNowTime() time.Time {
	return time.Now().UTC()
}

func repoTimeFromMilli(value int64) time.Time {
	if value <= 0 {
		return repoNowTime()
	}
	return time.UnixMilli(value).UTC()
}

func repoDeleteTimeFromMilli(value int64) int64 {
	if value <= 0 {
		return time.Now().UTC().Unix()
	}
	return time.UnixMilli(value).UTC().Unix()
}

func actorIDFromDB(tx *gorm.DB) int64 {
	if tx == nil || tx.Statement == nil || tx.Statement.Context == nil {
		return 0
	}
	subject, ok := contextx.SubjectFrom(tx.Statement.Context)
	if !ok {
		return 0
	}
	return subject.UserID
}

func fillCreateActor(tx *gorm.DB, createdBy, updatedBy *int64) {
	actorID := actorIDFromDB(tx)
	if createdBy != nil && *createdBy == 0 {
		*createdBy = actorID
	}
	if updatedBy != nil && *updatedBy == 0 {
		if createdBy != nil && *createdBy != 0 {
			*updatedBy = *createdBy
		} else {
			*updatedBy = actorID
		}
	}
}

func fillCreateOnlyActor(tx *gorm.DB, createdBy *int64) {
	if createdBy != nil && *createdBy == 0 {
		*createdBy = actorIDFromDB(tx)
	}
}

func fillUpdateActor(tx *gorm.DB, updatedBy *int64) {
	if updatedBy == nil || *updatedBy != 0 {
		return
	}
	actorID := actorIDFromDB(tx)
	if actorID == 0 {
		return
	}
	*updatedBy = actorID
	if tx != nil && tx.Statement != nil {
		tx.Statement.SetColumn("updated_by", actorID)
	}
}

func touchValues(values map[string]any, now int64) map[string]any {
	if values == nil {
		values = map[string]any{}
	}
	if _, ok := values["updated_time"]; !ok {
		values["updated_time"] = repoTimeFromMilli(now)
	}
	return values
}

func touchByValues(values map[string]any, userID, now int64) map[string]any {
	values = touchValues(values, now)
	if userID != 0 {
		if _, ok := values["updated_by"]; !ok {
			values["updated_by"] = userID
		}
	}
	return values
}

func softDeleteValues(userID, now int64) map[string]any {
	values := touchByValues(map[string]any{}, userID, now)
	values["deleted_time"] = repoDeleteTimeFromMilli(now)
	if userID != 0 {
		values["deleted_by"] = userID
	}
	return values
}

func softDeleteNoDelByValues(userID, now int64) map[string]any {
	values := touchByValues(map[string]any{}, userID, now)
	values["deleted_time"] = repoDeleteTimeFromMilli(now)
	return values
}
