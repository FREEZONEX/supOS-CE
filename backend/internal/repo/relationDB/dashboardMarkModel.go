package relationDB

const TableNameUnsDashboardMarkTop = "uns_dashboard_mark_top"

// DashboardMarkModel Dashboard 置顶标记模型
type DashboardMarkModel struct {
	ID     string `gorm:"column:id;primaryKey" json:"id"`
	UserID string `gorm:"column:user_id" json:"userId"`
}

// TableName 返回表名
func (DashboardMarkModel) TableName() string {
	return TableNameUnsDashboardMarkTop
}
