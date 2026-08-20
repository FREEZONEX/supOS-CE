// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package iam

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.PageReq) (resp *types.Envelope, err error) {
	users, err := l.svcCtx.App.IAM.Users(l.ctx)
	if err != nil {
		return nil, err
	}
	return respx.Envelope(map[string]any{"list": users, "total": len(users)}), nil
}
