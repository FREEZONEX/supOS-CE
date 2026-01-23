// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupedDashboardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页按分组获取dashboard列表
func NewGetGroupedDashboardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupedDashboardListLogic {
	return &GetGroupedDashboardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupedDashboardListLogic) GetGroupedDashboardList(req *types.GroupListQuery) (resp *types.DashboardGroupPageResp, err error) {
	// todo: add your logic here and delete this line

	return
}
