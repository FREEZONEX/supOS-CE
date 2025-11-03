package model

import "time"

// DashboardModel Dashboard 数据库模型
type DashboardModel struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Type        int       `db:"type" json:"type"`          // 1-grafana 2-fuxa
	NeedInit    bool      `db:"need_init" json:"needInit"` // 是否需要初始化
	Description string    `db:"description" json:"description"`
	JsonContent string    `db:"json_content" json:"jsonContent"`
	Creator     string    `db:"creator" json:"creator"`
	UpdateTime  time.Time `db:"update_time" json:"updateTime"`
	CreateTime  time.Time `db:"create_time" json:"createTime"`
	Error       string    `db:"-" json:"error,omitzero"` // 不存储在数据库中
}

// TableName 返回表名
func (DashboardModel) TableName() string {
	return "uns_dashboard"
}
