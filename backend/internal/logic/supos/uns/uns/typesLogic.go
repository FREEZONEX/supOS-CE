package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 枚举数据类型
func NewTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TypesLogic {
	return &TypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TypesLogic) Types() error {
	// todo: add your logic here and delete this line

	return nil
}
