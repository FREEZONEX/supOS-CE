package relationDB

import "time"

const TableNameGroup = "resource_group"

// GroupModel
type GroupModel struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	Type        *int16     `gorm:"column:type;not null" json:"type"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	Description string     `gorm:"column:description;not null" json:"description"`
	Sort        int32      `gorm:"column:sort;not null" json:"sort"`
	MarkTime    *time.Time `gorm:"column:mark_time" json:"mark_time"`
	UpdateAt    time.Time  `gorm:"column:update_at;default:now()" json:"update_at"`
	CreateAt    time.Time  `gorm:"column:create_at;default:now()" json:"create_at"`
	Creator     string     `gorm:"column:creator" json:"creator"`
}

// TableName 返回表名
func (GroupModel) TableName() string {
	return TableNameGroup
}
