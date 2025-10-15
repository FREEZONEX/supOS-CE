package dashboard

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BindUnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// bindUns
func NewBindUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUnsLogic {
	return &BindUnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindUnsLogic) BindUns() error {
	// todo: add your logic here and delete this line

	return nil
}
