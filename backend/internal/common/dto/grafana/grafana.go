package grafana

import (
	"encoding/base64"
	"fmt"
	"time"
)

type CreateDashboardDto struct {
	UID     string `json:"uid"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type DashboardDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        int       `json:"type"` // 1-grafana 2-fuxa
	Description string    `json:"description"`
	Mark        int       `json:"mark"` // 1 means pinned
	Creator     string    `json:"creator"`
	UpdateTime  time.Time `json:"updateTime"`
	CreateTime  time.Time `json:"createTime"`
}

type DashboardRefDto struct {
	DashboardID string `json:"dashboardId"`
	UnsAlias    string `json:"unsAlias"`
}

// GrafanaDashboardParam is the parameter for creating a Grafana dashboard.
type GrafanaDashboardParam struct {
	UID              string `json:"uid"`
	Title            string `json:"title"`
	DataSourceType   string `json:"dataSourceType"`
	DataSourceUID    string `json:"dataSourceUid"`
	Schema           string `json:"schema"`
	TableName        string `json:"tableName"`
	TagNameCondition string `json:"tagNameCondition"` // vqt 模式，位号名称sql条件 "tag='{tbValue}'"
	Columns          string `json:"columns"`
	Version          int64  `json:"version"`
}

type PgDashboardParam struct {
	GrafanaDashboardParam
}

type TdDashboardParam struct {
	GrafanaDashboardParam
}

type GrafanaDataSourceDto struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	User      string `json:"user"`
	Password  string `json:"password"`
	BasicAuth string `json:"basicAuth"`
}

// CreateBasicAuth generates the BasicAuth string from User and Password.
func (ds *GrafanaDataSourceDto) CreateBasicAuth() {
	auth := fmt.Sprintf("%s:%s", ds.User, ds.Password)
	ds.BasicAuth = base64.StdEncoding.EncodeToString([]byte(auth))
}

type GrafanaFolderDto struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
}

type TdPanelParam struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	DataSourceUID string `json:"dataSourceUid"`
	Columns       string `json:"columns"`
	Schema        string `json:"schema"`
	TableName     string `json:"tableName"`
	GridPosX      int    `json:"gridPosX"`
}
