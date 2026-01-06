package mount

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手动挂载
func NewMountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MountLogic {
	return &MountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MountLogic) Mount() error {
	// todo: add your logic here and delete this line

	return nil
}
