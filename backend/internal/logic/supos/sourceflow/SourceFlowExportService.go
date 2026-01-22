package sourceflow

import (
	"backend/internal/common"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"io"
)

type SourceFlowExportService struct {
}

func init() {
	spring.RegisterBean(&SourceFlowExportService{})
}

func (*SourceFlowExportService) FileName() string {
	return "source_flow.json"
}
func (*SourceFlowExportService) Order() int {
	return 5000
}

func (*SourceFlowExportService) ExportStream(ctx context.Context, arg *types.ExportReq) (exporter func(writer io.Writer)) {
	return
}
func (*SourceFlowExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {

}
