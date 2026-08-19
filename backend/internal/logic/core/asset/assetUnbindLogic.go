// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package asset

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssetUnbindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetUnbindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetUnbindLogic {
	return &AssetUnbindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetUnbindLogic) AssetUnbind(req *types.IdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Asset.Unbind(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
