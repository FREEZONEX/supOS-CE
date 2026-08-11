// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsExportJobDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsExportJobDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsExportJobDetailLogic {
	return &UnsExportJobDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsExportJobDetailLogic) UnsExportJobDetail(req *types.JobIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Job(l.ctx, req.JobId)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
