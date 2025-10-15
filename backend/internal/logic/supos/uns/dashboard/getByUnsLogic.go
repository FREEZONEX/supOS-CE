package dashboard

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetByUnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// getByUns
func NewGetByUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByUnsLogic {
	return &GetByUnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetByUnsLogic) GetByUns() error {
	// todo: add your logic here and delete this line

	return nil
}
