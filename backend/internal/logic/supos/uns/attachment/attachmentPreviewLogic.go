package attachment

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"backend/internal/common/utils/fileutil"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模板实例附件预览
func NewAttachmentPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentPreviewLogic {
	return &AttachmentPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentPreviewLogic) AttachmentPreview(req *types.AttachmentPreviewReq) (resp *types.AttachmentPreviewResp, err error) {
	repo := relationDB.NewUnsAttachmentRepo(l.ctx)

	attachments, err := repo.FindByAttachmentPath(l.ctx, req.ObjectName)
	if err != nil {
		return nil, err
	}
	if len(attachments) == 0 {
		logx.Errorf("附件不存在: %s", req.ObjectName)
		return nil, nil
	}

	attachment := attachments[0]
	filePath := filepath.Join(fileutil.GetFileRootPath(), attachment.AttachmentPath)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		logx.Errorf("读取文件失败: %w", err)
		return nil, nil
	}

	// SVG 文件特殊处理
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(attachment.OriginalName), ".svg") {
		contentType = "image/svg+xml;charset=UTF-8"
	}

	resp = &types.AttachmentPreviewResp{
		FileName: attachment.OriginalName,
		Content:  base64.StdEncoding.EncodeToString(fileData),
	}
	_ = contentType // 预览时可能需要设置 Content-Type，但这里返回的是 base64
	return resp, nil
}
