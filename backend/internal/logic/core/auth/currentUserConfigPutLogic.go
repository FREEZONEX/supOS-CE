// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentUserConfigPutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentUserConfigPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserConfigPutLogic {
	return &CurrentUserConfigPutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentUserConfigPutLogic) CurrentUserConfigPut(req *types.CurrentUserConfigReq) (resp *types.Envelope, err error) {
	return NewCurrentUserConfigUpdateLogic(l.ctx, l.svcCtx).CurrentUserConfigUpdate(req)
}
