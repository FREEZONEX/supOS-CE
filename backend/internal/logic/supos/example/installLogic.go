package example

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InstallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 安装接口
func NewInstallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstallLogic {
	return &InstallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstallLogic) Install() error {
	// todo: add your logic here and delete this line

	return nil
}
