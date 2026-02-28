package app

import (
	"backend/share/app/model"
	"context"

	"backend/internal/svc"
	"backend/share/app"
)

type AppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppsLogic {
	return &AppsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppsLogic) ListInstalledApps() ([]model.NewFeatureModel, error) {
	// 调用之前实现的 ListInstalledFeatures 函数
	features, err := app.ListInstalledFeatures()
	return features, err
}

func (l *AppsLogic) InstallApp(fc *model.NewFeatureModel) error {
	// 调用之前实现的 InstallFeature 函数
	err := app.InstallFeature(l.ctx, fc)
	return err
}
