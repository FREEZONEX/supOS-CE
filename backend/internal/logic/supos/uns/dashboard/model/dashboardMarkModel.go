package model

// DashboardMarkModel Dashboard 置顶标记模型
type DashboardMarkModel struct {
	ID     string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"userId"`
}

// TableName 返回表名
func (DashboardMarkModel) TableName() string {
	return "uns_dashboard_mark_top"
}
