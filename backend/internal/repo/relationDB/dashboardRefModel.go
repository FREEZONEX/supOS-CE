package relationDB

import "time"

const TableNameUnsDashboardRef = "uns_dashboard_ref"

// DashboardRefModel Dashboard 引用关系模型
type DashboardRefModel struct {
	DashboardID string    `gorm:"column:dashboard_id;not null" json:"dashboard_id"`
	UnsAlias    string    `gorm:"column:uns_alias;not null" json:"uns_alias"`
	CreateAt    time.Time `gorm:"column:create_at;default:now()" json:"create_at"`
}

// TableName 返回表名
func (DashboardRefModel) TableName() string {
	return TableNameUnsDashboardRef
}
