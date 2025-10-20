package uns

import (
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
	// todo: add your logic here and delete this line

	return
}
