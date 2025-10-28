package example

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UninstallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 卸载
func NewUninstallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UninstallLogic {
	return &UninstallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UninstallLogic) Uninstall() error {
	// todo: add your logic here and delete this line

	return nil
}
