package dashboard

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type IsExistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// isExist
func NewIsExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsExistLogic {
	return &IsExistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsExistLogic) IsExist() error {
	// todo: add your logic here and delete this line

	return nil
}
