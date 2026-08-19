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

type ApiKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiKeyDeleteLogic {
	return &ApiKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiKeyDeleteLogic) ApiKeyDelete(req *types.ApiKeyIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.APIKey.Delete(l.ctx, req.ApiKeyId, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
