// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"context"

	"backend/internal/contextx"
	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserLogic {
	return &CurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentUserLogic) CurrentUser() (resp *types.Envelope, err error) {
	subject, ok := contextx.SubjectFrom(l.ctx)
	if !ok {
		return nil, respx.NewHTTPError(401, "unauthorized")
	}
	user, keys, uiActions, forceChangePassword, err := l.svcCtx.App.Auth.CurrentUser(l.ctx, subject.UserID)
	if err != nil {
		return nil, err
	}
	return respx.Envelope(types.CurrentUserResp{
		UserId:              user.ID,
		UserName:            user.UserName,
		NickName:            user.NickName,
		Email:               user.Email,
		Phone:               user.Phone,
		RoleCode:            l.svcCtx.App.IAM.UserRoleCode(l.ctx, subject.UserID),
		ForceChangePassword: forceChangePassword,
		ResourceKeys:        keys,
		UiActions:           uiActions,
	}), nil
}
