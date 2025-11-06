// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Dashboard 详情
func NewGetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDetailLogic {
	return &GetDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDetailLogic) GetDetail(req *types.DetailRequest) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
