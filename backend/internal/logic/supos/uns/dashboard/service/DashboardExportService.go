package service

import (
	"backend/internal/common"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"io"
)

type DashboardExportService struct {
}

func init() {
	spring.RegisterBean(&DashboardExportService{})
}

func (*DashboardExportService) FileName() string {
	return "dashboard.json"
}
func (*DashboardExportService) Order() int {
	return 1000
}

func (*DashboardExportService) ExportStream(ctx context.Context, arg *types.ExportReq) (exporter func(writer io.Writer)) {
	return
}
func (*DashboardExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {

}
