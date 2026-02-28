// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInstalledAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List all installed applications
func NewListInstalledAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInstalledAppsLogic {
	return &ListInstalledAppsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInstalledAppsLogic) ListInstalledApps() (resp *types.InstalledAppsResponse, err error) {
	//err := app.InstallFeature(fc)
	return nil, err
}
