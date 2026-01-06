package serviceApi

import "backend/internal/types"

type TopicMessage struct {
	UnsId     int64
	DataSrcId types.SrcJdbcType
	Data      []map[string]any
}

// IDataSinkService 数据下沉服务
type IDataSinkService interface {
	// Sink 下沉数据
	Sink(unsData []TopicMessage)
}
