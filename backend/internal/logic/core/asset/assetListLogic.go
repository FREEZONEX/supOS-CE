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

type AssetListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetListLogic {
	return &AssetListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetListLogic) AssetList(req *types.AssetListReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Asset.List(l.ctx, req.OwnerType, req.OwnerId)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
