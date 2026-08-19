// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"
	"mime/multipart"
	"strings"

	domainasset "backend/internal/domain/asset"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsAttachmentUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsAttachmentUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsAttachmentUploadLogic {
	return &OpenapiUnsAttachmentUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsAttachmentUploadLogic) OpenapiUnsAttachmentUpload(req *types.UnsAttachmentUploadReq, file multipart.File, header *multipart.FileHeader) (resp *types.Envelope, err error) {
	node, err := resolveOpenapiAttachmentNode(l.ctx, req.UnsId, req.Topic)
	if err != nil {
		return nil, logicx.Error(err)
	}
	fileName := strings.TrimSpace(req.FileName)
	if fileName == "" && header != nil {
		fileName = header.Filename
	}
	contentType := ""
	size := int64(0)
	if header != nil {
		contentType = header.Header.Get("Content-Type")
		size = header.Size
	}
	data, err := l.svcCtx.App.Asset.UploadAttachment(l.ctx, domainasset.AttachmentUploadCommand{
		OwnerType:   "unsNode",
		OwnerID:     node.ID,
		UnsID:       node.ID,
		FileName:    fileName,
		ContentType: contentType,
		Size:        size,
		Sha256:      req.Sha256,
		Reader:      file,
		UserID:      logicx.UserID(l.ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
