// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package apikey

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiKeyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiKeyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiKeyUpdateLogic {
	return &ApiKeyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiKeyUpdateLogic) ApiKeyUpdate(req *types.ApiKeyUpdateReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.APIKey.Update(l.ctx, req.ApiKeyId, logicx.UserID(l.ctx), req.Name, req.Permission)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
