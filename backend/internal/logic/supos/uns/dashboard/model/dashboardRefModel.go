package model

import "time"

// DashboardRefModel Dashboard 引用关系模型
type DashboardRefModel struct {
	DashboardID string    `db:"dashboard_id" json:"dashboardId"`
	UnsAlias    string    `db:"uns_alias" json:"unsAlias"`
	CreateAt    time.Time `db:"create_at" json:"createAt"`
}

// TableName 返回表名
func (DashboardRefModel) TableName() string {
	return "uns_dashboard_ref"
}
