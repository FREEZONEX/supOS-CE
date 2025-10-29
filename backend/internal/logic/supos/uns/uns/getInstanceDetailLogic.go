// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInstanceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文件详情
func NewGetInstanceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInstanceDetailLogic {
	return &GetInstanceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInstanceDetailLogic) GetInstanceDetail(req *types.InstanceDetailReq) (resp *types.InstanceDetailResp, err error) {
	resp, err = spring.GetBean[*service.UnsQueryService]().GetInstanceDetail(l.ctx, req, "")
	return
}
