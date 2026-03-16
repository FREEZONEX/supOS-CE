// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/logic/supos/open/openservice"
	"context"
	"net/http"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type OpenUserPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户列表
func NewOpenUserPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenUserPageListLogic {
	return &OpenUserPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenUserPageListLogic) OpenUserPageList(r *http.Request) (resp *types.OpenUserPageResult, err error) {
	// 解析查询参数
	var req types.OpenUserPageQuery
	if err := httpx.Parse(r, &req); err != nil {
		return nil, err
	}

	// 使用 UserOpenapiService 获取用户列表
	userOpenapiService := openservice.NewUserOpenapiService(l.ctx, l.svcCtx)
	result, err := userOpenapiService.UserManageList(openservice.UserPageQueryDto{
		PageNo:      req.PageNo,
		PageSize:    req.PageSize,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Phone:       req.Phone,
		Enabled:     req.Enabled,
	})
	if err != nil {
		return nil, err
	}

	return &types.OpenUserPageResult{
		Code:     0,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Total:    result.Total,
		Data:     result.Users,
	}, nil
}
