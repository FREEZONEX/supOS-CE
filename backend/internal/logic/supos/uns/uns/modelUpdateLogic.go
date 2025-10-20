package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ModelUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改文件夹或文件明细
func NewModelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModelUpdateLogic {
	return &ModelUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModelUpdateLogic) ModelUpdate() error {
	// todo: add your logic here and delete this line

	return nil
}
