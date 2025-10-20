package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type Json2fsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部JSON定义转uns字段定义
func NewJson2fsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Json2fsLogic {
	return &Json2fsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Json2fsLogic) Json2fs() error {
	// todo: add your logic here and delete this line

	return nil
}
