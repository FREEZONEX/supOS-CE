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

type RemoveModelOrInstanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除指定路径下的所有文件夹和文件
func NewRemoveModelOrInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveModelOrInstanceLogic {
	return &RemoveModelOrInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveModelOrInstanceLogic) RemoveModelOrInstance(req *types.RemoveReq) (resp *types.RemoveResult, err error) {
	resp, err = spring.GetBean[*service.UnsRemoveService]().RemoveModelOrInstance(l.ctx, req)
	return
}
