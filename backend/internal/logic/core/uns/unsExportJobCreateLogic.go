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

type UnsExportJobCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsExportJobCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsExportJobCreateLogic {
	return &UnsExportJobCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsExportJobCreateLogic) UnsExportJobCreate(req *types.UnsExportJobReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.CreateExportJob(l.ctx, req.RootNodeId, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
