package relationDB

import "time"

// suposResource maps to table supos_resource.
type SuposResource struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ParentID        *int64    `gorm:"column:parent_id"`
	Type            int       `gorm:"column:type"`
	Source          *string   `gorm:"column:source"`
	Code            string    `gorm:"column:code"`
	NameCode        *string   `gorm:"column:name_code"`
	RouteSource     *int      `gorm:"column:route_source"`
	URL             *string   `gorm:"column:url"`
	URLType         *int      `gorm:"column:url_type"`
	OpenType        *int      `gorm:"column:open_type"`
	Icon            *string   `gorm:"column:icon"`
	DescriptionCode *string   `gorm:"column:description_code"`
	Sort            *int      `gorm:"column:sort"`
	EditEnable      *bool     `gorm:"column:edit_enable"`
	HomeEnable      *bool     `gorm:"column:home_enable"`
	Fixed           *bool     `gorm:"column:fixed"`
	Enable          *bool     `gorm:"column:enable"`
	UpdateAt        time.Time `gorm:"column:update_at"`
	CreateAt        time.Time `gorm:"column:create_at"`
}

func (SuposResource) TableName() string {
	return "supos.supos_resource"
}
