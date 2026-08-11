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

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() (resp *types.Envelope, err error) {
	if subject, ok := contextx.SubjectFrom(l.ctx); ok {
		l.svcCtx.App.Audit.Record(l.ctx, auditdomain.RecordInput{
			ScopeType:    auditdomain.ScopeTypePlatform,
			ResType:      auditdomain.ResTypeAuth,
			ResID:        strconv.FormatInt(subject.UserID, 10),
			ResName:      subject.UserName,
			BusinessType: auditdomain.BizLogout,
			Detail:       map[string]any{"username": subject.UserName},
		})
	}
	return respx.Empty(), nil
}
