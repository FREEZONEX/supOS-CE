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

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件下载
func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadReq, r *http.Request, w http.ResponseWriter) error {
	return spring.GetBean[*service.UnsImportExportService]().FileDownload(req, r, w)
}
