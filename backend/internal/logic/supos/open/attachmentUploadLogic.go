// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"net/http"

	"backend/internal/svc"
	"backend/internal/types"

	"backend/internal/logic/supos/uns/attachment"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

// 模板实例附件上传
func NewAttachmentUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *AttachmentUploadLogic {
	return &AttachmentUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *AttachmentUploadLogic) AttachmentUpload(req *types.AttachmentUploadReq) (resp *types.AttachmentUploadResp, err error) {
	// 调用 uns/attachment 包中的 AttachmentUpload 方法
	unsLogic := attachment.NewAttachmentUploadLogic(l.ctx, l.svcCtx, l.r)
	return unsLogic.AttachmentUpload(req)
}
