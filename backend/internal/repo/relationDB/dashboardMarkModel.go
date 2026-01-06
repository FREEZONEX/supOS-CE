package relationDB

import "time"

const TableNameUnsDashboardMarkTop = "uns_dashboard_top_recodes"

// DashboardMarkModel Dashboard 置顶标记模型
type DashboardMarkModel struct {
	ID         string    `gorm:"column:id;uniqueIndex:udx_dashboard_id_user;not null" json:"id"`
	UserID     string    `gorm:"column:user_id;uniqueIndex:udx_dashboard_id_user;not null" json:"user_id"`
	Mark       int16     `gorm:"column:mark;default:1" json:"mark"`
	MarkTime   time.Time `gorm:"column:mark_time;default:CURRENT_TIMESTAMP" json:"mark_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// TableName 返回表名
func (DashboardMarkModel) TableName() string {
	return TableNameUnsDashboardMarkTop
}
