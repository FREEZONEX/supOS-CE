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

type SubscribeModelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件或文件夹修改订阅
func NewSubscribeModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeModelLogic {
	return &SubscribeModelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeModelLogic) SubscribeModel(req *types.SubscribeModelReq) (resp *types.ResultVO, err error) {
	resp, err = spring.GetBean[*service.UnsUpdateService]().SubscribeModel(l.ctx, req)
	return
}
