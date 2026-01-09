package serviceApi

import "backend/internal/types"

type UnsData struct {
	Uns  *types.CreateTopicDto
	Data []map[string]string
}
type DataSourceProperties struct {
	Url      string
	UserName string
	Password string
	DbName   string
	HostPort string
	Schema   string
}

// IPersistentService 数据持久化服务
type IPersistentService interface {
	// Persistent 持久化数据
	Persistent(unsData []UnsData)
	// GetDataSrcId 返回数据源ID
	GetDataSrcId() types.SrcJdbcType

	GetDataSourceProperties() DataSourceProperties

	// FillLastRecord 填充最后一条记录
	FillLastRecord(uns *types.CreateTopicDto)

	Save(creates []types.UnsInfo) error
	Remove(topics []types.UnsInfo) error
}
