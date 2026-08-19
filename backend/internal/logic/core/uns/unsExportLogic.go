// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsExportLogic {
	return &UnsExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsExportLogic) UnsExport(req *types.UnsExportJobReq) (resp *types.Envelope, err error) {
	// todo: add your logic here and delete this line

	return
}
