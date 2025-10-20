package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type Ds2fsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部数据源表的字段定义转uns字段定义
func NewDs2fsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Ds2fsLogic {
	return &Ds2fsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Ds2fsLogic) Ds2fs() error {
	// todo: add your logic here and delete this line

	return nil
}
