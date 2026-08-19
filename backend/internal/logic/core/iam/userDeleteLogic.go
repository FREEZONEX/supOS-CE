// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package iam

import (
	"context"
	"strconv"

	auditdomain "backend/internal/domain/audit"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/repo"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteLogic) UserDelete(req *types.IdReq) (resp *types.Envelope, err error) {
	user, err := repo.NewIAMRepo(l.ctx).GetUserByID(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	if err := l.svcCtx.App.IAM.DeleteUser(l.ctx, req.Id, logicx.UserID(l.ctx)); err != nil {
		return nil, logicx.Error(err)
	}
	resName := user.NickName
	if resName == "" {
		resName = user.UserName
	}
	l.svcCtx.App.Audit.Record(l.ctx, auditdomain.RecordInput{
		ScopeType:    auditdomain.ScopeTypePlatform,
		ResType:      auditdomain.ResTypeUserManagement,
		ResID:        strconv.FormatInt(req.Id, 10),
		ResName:      resName,
		BusinessType: auditdomain.BizDelete,
		Detail:       map[string]any{"userId": req.Id, "userName": user.UserName, "nickName": user.NickName},
	})
	return respx.Envelope(map[string]any{"id": req.Id}), nil
}
