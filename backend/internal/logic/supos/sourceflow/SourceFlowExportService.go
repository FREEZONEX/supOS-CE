package sourceflow

import (
	"backend/internal/common"
	"backend/internal/common/event"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type SourceFlowExportService struct {
	apiHost string
}

func init() {
	spring.RegisterBean(&SourceFlowExportService{})
}

func (s *SourceFlowExportService) FileName() string {
	return "source_flow.json"
}
func (s *SourceFlowExportService) Order() int {
	return 5000
}

func (s *SourceFlowExportService) ExportStream(ctx context.Context, arg *types.GlobalExportParam) (exporter func(writer io.Writer)) {
	req := arg.SrcFlowExportParam
	if req == nil {
		return
	}
	return NodeRedFlowExport(ctx, req.GroupIds, req.FlowIds, true, s.apiHost)
}
func (s *SourceFlowExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {
	logx.WithContext(ctx).Info("ImportStream start: ", fileName)
	NodeRedImport(ctx, s.apiHost, true, reader, statusConsumer)
}
func (s *SourceFlowExportService) OnEventContextRefreshedEvent(evt *event.ContextRefreshedEvent) error {
	cf := evt.SvcContext.Config.NodeRed.Source
	apiHost := fmt.Sprintf("%s:%v", cf.Host, cf.Port)
	if !strings.HasPrefix(apiHost, "http") {
		apiHost = "http://" + apiHost
	}
	s.apiHost = apiHost
	return nil
}
