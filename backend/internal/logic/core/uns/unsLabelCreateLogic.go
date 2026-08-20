// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsLabelCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsLabelCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsLabelCreateLogic {
	return &UnsLabelCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsLabelCreateLogic) UnsLabelCreate(req *types.UnsLabelReq) (resp *types.Envelope, err error) {
	// todo: add your logic here and delete this line

	return
}
