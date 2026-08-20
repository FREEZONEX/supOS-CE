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

type AssetDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetDownloadLogic {
	return &AssetDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetDownloadLogic) AssetDownload(req *types.FileIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Asset.Download(l.ctx, req.FileId)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
