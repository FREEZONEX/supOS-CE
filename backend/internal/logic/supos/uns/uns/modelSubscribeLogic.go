package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ModelSubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件或文件夹修改订阅
func NewModelSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModelSubscribeLogic {
	return &ModelSubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModelSubscribeLogic) ModelSubscribe() error {
	// todo: add your logic here and delete this line

	return nil
}
