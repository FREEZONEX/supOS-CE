// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type TemplateDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 下载模版
func NewTemplateDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TemplateDownloadLogic {
	return &TemplateDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TemplateDownloadLogic) TemplateDownload(req *types.TemplateDownloadReq, r *http.Request, w http.ResponseWriter) error {
	return spring.GetBean[*service.UnsImportExportService]().TemplateDownload(req, r, w)
}
