// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"context"
	"strconv"

	auditdomain "backend/internal/domain/audit"
	domainiam "backend/internal/domain/iam"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentUserPasswordUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentUserPasswordUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserPasswordUpdateLogic {
	return &CurrentUserPasswordUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentUserPasswordUpdateLogic) CurrentUserPasswordUpdate(req *types.CurrentUserPasswordReq) (resp *types.Envelope, err error) {
	oldPassword := req.OldPassword
	if oldPassword == "" {
		oldPassword = req.Password
	}
	if err := l.svcCtx.App.IAM.ChangeCurrentUserPassword(l.ctx, domainiam.CurrentUserPasswordCommand{
		UserID:      logicx.UserID(l.ctx),
		OldPassword: oldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		return nil, logicx.Error(err)
	}
	userID := logicx.UserID(l.ctx)
	l.svcCtx.App.Audit.Record(l.ctx, auditdomain.RecordInput{
		ScopeType:    auditdomain.ScopeTypePlatform,
		ResType:      auditdomain.ResTypeAuth,
		ResID:        strconv.FormatInt(userID, 10),
		ResName:      strconv.FormatInt(userID, 10),
		BusinessType: auditdomain.BizResetPassword,
		Detail:       map[string]any{"userId": userID},
	})
	return respx.Envelope(map[string]any{"updated": true}), nil
}
