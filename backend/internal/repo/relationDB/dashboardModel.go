package relationDB

import "time"

const TableNameUnsDashboard = "uns_dashboard"

// DashboardModel Dashboard 数据库模型
type DashboardModel struct {
	ID          string            `gorm:"column:id;primaryKey" json:"id"`
	Name        string            `gorm:"column:name" json:"name"`
	Type        int               `gorm:"column:type" json:"type,omitzero"` // 来源：1-grafana 2-fuxa
	GroupId     *int64            `gorm:"column:group_id;comment:分组ID" json:"-"`
	ExportType  string            `gorm:"-"  json:"exportType,omitempty"`            // 导出类型：group-分组 默认：dashboard-面板
	NeedInit    bool              `gorm:"column:need_init" json:"needInit,omitzero"` // 是否需要初始化
	Description string            `gorm:"column:description" json:"description,omitempty"`
	JsonContent string            `gorm:"column:json_content" json:"jsonContent,omitempty"`
	Creator     string            `gorm:"column:creator" json:"creator,omitempty"`
	UpdateTime  time.Time         `gorm:"column:update_time;autoUpdateTime" json:"updateTime,omitzero"`
	CreateTime  time.Time         `gorm:"column:create_time;autoCreateTime" json:"createTime,omitzero"`
	Error       string            `gorm:"-" json:"error,omitzero"` // 不存储在数据库中
	Children    []*DashboardModel `gorm:"-" json:"children,omitempty"`
	Sort        int32             `gorm:"-" json:"sort,omitzero"`
}

// TableName 返回表名
func (DashboardModel) TableName() string {
	return TableNameUnsDashboard
}
