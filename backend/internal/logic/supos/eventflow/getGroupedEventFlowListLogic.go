// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupedEventFlowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页按分组获取event flow列表
func NewGetGroupedEventFlowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupedEventFlowListLogic {
	return &GetGroupedEventFlowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupedEventFlowListLogic) GetGroupedEventFlowList(req *types.GroupPageRequest) (resp *types.EventFlowGroupPageResp, err error) {

	return
}
