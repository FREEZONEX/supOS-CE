package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type NameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改文件夹或文件名称
func NewNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NameLogic {
	return &NameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NameLogic) Name() error {
	// todo: add your logic here and delete this line

	return nil
}
