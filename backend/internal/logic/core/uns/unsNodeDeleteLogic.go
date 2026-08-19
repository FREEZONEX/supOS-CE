// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	auditdomain "backend/internal/domain/audit"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsNodeDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsNodeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsNodeDeleteLogic {
	return &UnsNodeDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsNodeDeleteLogic) UnsNodeDelete(req *types.NodeIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Delete(l.ctx, req.NodeId, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordUNSAudit(l.ctx, l.svcCtx, data, auditdomain.BizDelete)
	return respx.Envelope(data), nil
}
