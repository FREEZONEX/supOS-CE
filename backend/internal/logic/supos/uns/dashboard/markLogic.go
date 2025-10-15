package dashboard

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 置顶
func NewMarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkLogic {
	return &MarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkLogic) Mark() error {
	// todo: add your logic here and delete this line

	return nil
}
