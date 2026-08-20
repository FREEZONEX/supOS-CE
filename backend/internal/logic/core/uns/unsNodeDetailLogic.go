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

type UnsNodeDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsNodeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsNodeDetailLogic {
	return &UnsNodeDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsNodeDetailLogic) UnsNodeDetail(req *types.NodeIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Detail(l.ctx, req.NodeId, false)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
