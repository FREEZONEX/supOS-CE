package serviceApi

import (
	"backend/internal/types"
	"context"
)

type TopicMessage struct {
	UnsId     int64
	DataSrcId types.SrcJdbcType
	Data      []map[string]string
}

// IDataSinkService 数据下沉服务
type IDataSinkService interface {
	// Sink 下沉数据
	Sink(ctx context.Context, unsData []TopicMessage)
}
