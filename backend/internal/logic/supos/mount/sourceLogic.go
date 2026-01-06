package mount

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取挂载数据源
func NewSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SourceLogic {
	return &SourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SourceLogic) Source() error {
	// todo: add your logic here and delete this line

	return nil
}
