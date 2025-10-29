package dashboard

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnmarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消置顶
func NewUnmarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarkLogic {
	return &UnmarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnmarkLogic) Unmark() error {
	// todo: add your logic here and delete this line

	return nil
}
