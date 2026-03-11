package eventflow

import (
	"backend/internal/common"
	"backend/internal/common/event"
	"backend/internal/logic/supos/sourceflow"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type EventFlowExportService struct {
	apiHost string
}

func init() {

	spring.RegisterBean(&EventFlowExportService{})
}

func (e *EventFlowExportService) FileName() string {
	return "event_flow.json"
}
func (e *EventFlowExportService) Order() int {
	return 3000
}

func (e *EventFlowExportService) ExportStream(ctx context.Context, arg *types.GlobalExportParam) (exporter func(writer io.Writer)) {
	req := arg.EventFlowExportParam
	if req == nil {
		return
	}
	return sourceflow.NodeRedFlowExport(ctx, req.GroupIds, req.FlowIds, false, e.apiHost)
}
func (e *EventFlowExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {
	logx.WithContext(ctx).Info("ImportStream start: ", fileName)
	sourceflow.NodeRedImport(ctx, e.apiHost, false, reader, statusConsumer)
}
func (e *EventFlowExportService) OnEventContextRefreshedEvent(evt *event.ContextRefreshedEvent) error {
	cf := evt.SvcContext.Config.NodeRed.Event
	apiHost := fmt.Sprintf("%s:%v", cf.Host, cf.Port)
	if !strings.HasPrefix(apiHost, "http") {
		apiHost = "http://" + apiHost
	}
	e.apiHost = apiHost
	return nil
}
