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

type ApiKeyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiKeyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiKeyListLogic {
	return &ApiKeyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiKeyListLogic) ApiKeyList(req *types.PageReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.APIKey.List(l.ctx, logicx.UserID(l.ctx), req.KeyType, req.Keyword)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
