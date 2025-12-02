// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUploadLogic UNS 上传文件
func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.MultipartFile) (resp *types.StringResult, err error) {
	return spring.GetBean[*service.UnsImportExportService]().UploadFile(req)
}
