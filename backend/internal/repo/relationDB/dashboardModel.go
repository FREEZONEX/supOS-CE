package relationDB

import "time"

const TableNameUnsDashboard = "uns_dashboard"

// DashboardModel Dashboard 数据库模型
type DashboardModel struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Type        int       `gorm:"column:type" json:"type"`          // 1-grafana 2-fuxa
	NeedInit    bool      `gorm:"column:need_init" json:"needInit"` // 是否需要初始化
	Description string    `gorm:"column:description" json:"description"`
	JsonContent string    `gorm:"column:json_content" json:"jsonContent"`
	Creator     string    `gorm:"column:creator" json:"creator"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	Error       string    `gorm:"-" json:"error,omitzero"` // 不存储在数据库中
}

// TableName 返回表名
func (DashboardModel) TableName() string {
	return TableNameUnsDashboard
}
