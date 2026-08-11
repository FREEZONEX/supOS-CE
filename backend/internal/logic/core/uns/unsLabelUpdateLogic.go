// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsLabelUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsLabelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsLabelUpdateLogic {
	return &UnsLabelUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsLabelUpdateLogic) UnsLabelUpdate(req *types.UnsLabelReq) (resp *types.Envelope, err error) {
	// todo: add your logic here and delete this line

	return
}
