// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取最新消息
func NewGetLastMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastMsgLogic {
	return &GetLastMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLastMsgLogic) GetLastMsg(req *types.GetLastMsgReq) (resp *types.GetLastMsgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
