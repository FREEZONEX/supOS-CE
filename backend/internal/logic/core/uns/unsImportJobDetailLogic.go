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

type UnsImportJobDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsImportJobDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsImportJobDetailLogic {
	return &UnsImportJobDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsImportJobDetailLogic) UnsImportJobDetail(req *types.JobIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Job(l.ctx, req.JobId)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
