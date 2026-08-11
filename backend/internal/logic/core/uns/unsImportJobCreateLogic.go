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

type UnsImportJobCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsImportJobCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsImportJobCreateLogic {
	return &UnsImportJobCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsImportJobCreateLogic) UnsImportJobCreate(req *types.UnsImportJobReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.CreateImportJob(l.ctx, req.DryRun, req.SourceFileId, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
