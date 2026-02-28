package app

import (
	"context"
	"time"

	"backend/internal/common/I18nUtils"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/app"
	"backend/share/app/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type InstallAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewInstallAppLogic 创建 InstallAppLogic 实例
func NewInstallAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstallAppLogic {
	return &InstallAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// InstallApp 安装应用
func (l *InstallAppLogic) InstallApp(req *types.InstallAppRequest) (resp *types.InstallAppResponse, err error) {
	logx.Info("[InstallApp] 开始安装应用:", req.Name)

	// 将请求转换为 NewFeatureModel
	featureModel := &model.NewFeatureModel{
		Name:        req.Name,
		Description: req.Description,
		ImagePath:   req.ImagePath,
		ImageUrl:    req.ImageUrl,
		Menu:        req.Menu,
		ComposeYaml: req.ComposeYaml,
	}

	// 调用安装功能
	if err := app.InstallFeature(featureModel); err != nil {
		logx.Errorf("[InstallApp] 安装应用失败: %v", err)
		return nil, err
	}

	// 构建成功响应
	installTime := time.Now().Format("2006-01-02 15:04:05")
	responseData := types.InstallResult{
		Name:        req.Name,
		Success:     true,
		Message:     I18nUtils.GetMessage("app.install.success"),
		InstallTime: installTime,
	}

	logx.Info("[InstallApp] 应用安装成功:", req.Name)
	return &types.InstallAppResponse{
		Code:    200,
		Message: I18nUtils.GetMessage("app.install.success"),
		Data:    responseData,
	}, nil
}
