package service

import (
	"backend/internal/common"
	"backend/internal/types"
	"context"
	"io"
)

type ExportService interface {
	FileName() string
	ExportStream(ctx context.Context, arg *types.ExportReq) (exporter func(writer io.Writer))
	ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus))
	Order() int // 优先级
}
