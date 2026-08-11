// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsNodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsNodesLogic {
	return &OpenapiUnsNodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsNodesLogic) OpenapiUnsNodes(req *types.UnsListReq) (resp *types.Envelope, err error) {
	return NewOpenapiUnsNodeListLogic(l.ctx, l.svcCtx).OpenapiUnsNodeList(req)
}
