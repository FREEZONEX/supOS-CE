package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ModelDetectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 预先判断是否有属性关联
func NewModelDetectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModelDetectLogic {
	return &ModelDetectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModelDetectLogic) ModelDetect() error {
	// todo: add your logic here and delete this line

	return nil
}
