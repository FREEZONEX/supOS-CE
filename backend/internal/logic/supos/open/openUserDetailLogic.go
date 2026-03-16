// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/open/openservice"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户详情
func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(username string) (resp *types.UserDetailResult, err error) {
	if l.svcCtx.Keycloak == nil {
		return &types.UserDetailResult{
			Code: 404,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "user.not.exist"),
		}, nil
	}

	// 使用 UserOpenapiService 通过用户名获取用户信息
	userOpenapiService := openservice.NewUserOpenapiService(l.ctx, l.svcCtx)
	result, err := userOpenapiService.UserManageList(openservice.UserPageQueryDto{
		ExactUsername: username,
	})
	if err != nil {
		l.Errorf("Failed to fetch user %s: %v", username, err)
		return nil, err
	}

	if len(result.Users) == 0 {
		return &types.UserDetailResult{
			Code: 404,
			Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "user.not.exist"),
		}, nil
	}

	return &types.UserDetailResult{
		Code: 0,
		Msg:  "success",
		Data: result.Users[0],
	}, nil
}
