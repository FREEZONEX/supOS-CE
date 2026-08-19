// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"context"
	"strconv"

	"backend/internal/contextx"
	auditdomain "backend/internal/domain/audit"
	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.Envelope, err error) {
	result, err := l.svcCtx.App.Auth.Login(l.ctx, req.Username, req.Password)
	if err != nil {
		l.svcCtx.App.Audit.Record(l.ctx, auditdomain.RecordInput{
			ScopeType:    auditdomain.ScopeTypePlatform,
			ResType:      auditdomain.ResTypeAuth,
			ResName:      req.Username,
			BusinessType: auditdomain.BizLogin,
			Detail:       map[string]any{"username": req.Username},
		})
		return nil, respx.NewHTTPError(401, "invalid username or password")
	}
	subject := contextx.Subject{
		UserID:   result.User.ID,
		UserName: result.User.UserName,
		Email:    result.User.Email,
		AuthType: "session",
	}
	l.svcCtx.App.Audit.BindSubject(l.ctx, subject)
	l.svcCtx.App.Audit.Record(l.ctx, auditdomain.RecordInput{
		ScopeType:    auditdomain.ScopeTypePlatform,
		ResType:      auditdomain.ResTypeAuth,
		ResID:        strconv.FormatInt(result.User.ID, 10),
		ResName:      result.User.UserName,
		BusinessType: auditdomain.BizLogin,
		Detail:       map[string]any{"username": result.User.UserName},
	})
	return respx.Envelope(types.LoginResp{
		Token:               result.Token,
		UserId:              result.User.ID,
		UserName:            result.User.UserName,
		NickName:            result.User.NickName,
		Email:               result.User.Email,
		Phone:               result.User.Phone,
		RoleCode:            l.svcCtx.App.IAM.UserRoleCode(l.ctx, result.User.ID),
		IsRandomPwd:         result.User.IsRandomPwd,
		ForceChangePassword: result.ForceChangePassword,
		ResourceKeys:        result.ResourceKeys,
		UiActions:           result.UIActions,
	}), nil
}
