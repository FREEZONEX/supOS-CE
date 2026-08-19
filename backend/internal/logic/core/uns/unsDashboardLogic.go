// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	domainuns "backend/internal/domain/uns"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsDashboardLogic {
	return &UnsDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsDashboardLogic) UnsDashboard(req *types.UnsDashboardReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Dashboard(l.ctx, domainuns.DashboardQuery{
		NodeID:    req.NodeId,
		TimeStart: req.TimeStart,
		TimeEnd:   req.TimeEnd,
		Limit:     req.Limit,
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
