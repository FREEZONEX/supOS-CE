// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// UNS 导出,返回文件路径
func NewExportPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportPathLogic {
	return &ExportPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportPathLogic) ExportPath(req *types.ExportReq) (resp *types.ExportResp, err error) {
	return spring.GetBean[*service.UnsImportExportService]().ExportPath(l.ctx, req)
}
