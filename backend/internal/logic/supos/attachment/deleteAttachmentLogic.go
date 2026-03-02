// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAttachmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete attachment
func NewDeleteAttachmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAttachmentLogic {
	return &DeleteAttachmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAttachmentLogic) DeleteAttachment(req *types.DeleteAttachmentRequest) (resp *types.DeleteAttachmentResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
