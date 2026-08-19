// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsAttachmentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsAttachmentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsAttachmentListLogic {
	return &OpenapiUnsAttachmentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsAttachmentListLogic) OpenapiUnsAttachmentList(req *types.UnsAttachmentListReq) (resp *types.Envelope, err error) {
	node, err := resolveOpenapiAttachmentNode(l.ctx, req.UnsId, req.Topic)
	if err != nil {
		return nil, logicx.Error(err)
	}
	includeFileURL := true
	if req.IncludeFileUrl != nil {
		includeFileURL = *req.IncludeFileUrl
	}
	data, err := l.svcCtx.App.Asset.ListAttachments(l.ctx, "unsNode", node.ID, node.ID, req.Page, req.Size, includeFileURL)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
