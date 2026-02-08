package app

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"backend/internal/types"
	"backend/internal/logic/supos/app"
	"backend/internal/svc"
	"backend/share/app/model"
)

// ListInstalledAppsHandler 获取已安装应用列表
func ListInstalledAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := app.NewAppsLogic(r.Context(), svcCtx)
		resp, err := l.ListInstalledApps()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// InstallAppsHandler 安装应用
func InstallAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.InstallAppRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 转换为 NewFeatureModel
		featureModel := &model.NewFeatureModel{
			Name:        req.Name,
			Description: req.Description,
			ImagePath:   req.ImagePath,
			ImageUrl:    req.ImageUrl,
			IconPath:    req.IconPath,
			MenuUrl:     req.MenuUrl,
			ComposeYaml: req.ComposeYaml,
			RouterTrim:  req.RouterTrim,
			KeepHost:    req.KeepHost,
		}

		// 验证模型
		if err := featureModel.Validate(); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := app.NewAppsLogic(r.Context(), svcCtx)
		err := l.InstallApp(featureModel)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
				"code": 200,
				"message": "安装成功",
				"data": map[string]interface{}{
					"name": featureModel.Name,
					"success": true,
				},
			})
		}
	}
}