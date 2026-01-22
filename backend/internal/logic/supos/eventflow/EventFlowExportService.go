package eventflow

import (
	"backend/internal/common"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"io"
)

type EventFlowExportService struct {
}

func init() {
	spring.RegisterBean(&EventFlowExportService{})
}

func (*EventFlowExportService) FileName() string {
	return "event_flow.json"
}
func (*EventFlowExportService) Order() int {
	return 3000
}

func (*EventFlowExportService) ExportStream(ctx context.Context, arg *types.ExportReq) (exporter func(writer io.Writer)) {
	return
}
func (*EventFlowExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {

}
