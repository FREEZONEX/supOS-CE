// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package iam

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

type UserPasswordResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPasswordResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPasswordResetLogic {
	return &UserPasswordResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPasswordResetLogic) UserPasswordReset(req *types.UserPasswordResetReq) (resp *types.Envelope, err error) {
	if err := l.svcCtx.App.IAM.ResetUserPassword(l.ctx, domainiam.UserPasswordResetCommand{
		UserID:   req.Id,
		Password: req.Password,
		ActorID:  logicx.UserID(l.ctx),
	}); err != nil {
		return nil, logicx.Error(err)
	}
	l.svcCtx.App.Audit.Record(l.ctx, auditdomain.RecordInput{
		ScopeType:    auditdomain.ScopeTypePlatform,
		ResType:      auditdomain.ResTypeUserManagement,
		ResID:        strconv.FormatInt(req.Id, 10),
		ResName:      strconv.FormatInt(req.Id, 10),
		BusinessType: auditdomain.BizResetPassword,
		Detail:       map[string]any{"userId": req.Id},
	})
	return respx.Envelope(map[string]any{"id": req.Id}), nil
}
