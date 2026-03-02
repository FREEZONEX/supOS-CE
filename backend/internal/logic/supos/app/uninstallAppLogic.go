// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UninstallAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Uninstall an application
func NewUninstallAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UninstallAppLogic {
	return &UninstallAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UninstallAppLogic) UninstallApp(req *types.UninstallAppRequest) (resp *types.UninstallAppResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
