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

type ApiKeyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiKeyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiKeyCreateLogic {
	return &ApiKeyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiKeyCreateLogic) ApiKeyCreate(req *types.ApiKeyCreateReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.APIKey.Create(l.ctx, req.Name, logicx.UserID(l.ctx), req.Permission, req.UsageType, req.KeyType, req.ResourceKeys)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
