package attachment

import (
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/app/util"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAttachmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAttachmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAttachmentLogic {
	return &UploadAttachmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAttachmentLogic) UploadAttachment(req *types.UploadAttachmentRequest) (*types.UploadAttachmentResponse, error) {
	// 1. 验证文件大小（最大2GB）
	if req.FileHeader.Size > 2*1024*1024*1024 { // 2GB
		return nil, errors.Parameter.WithMsg("文件大小超过2GB限制")
	}

	// 2. 验证文件类型
	//if err := l.validateFileType(req.FileHeader); err != nil {
	//	return nil, err
	//}

	// 3. 使用附件管理器上传文件
	attachmentInfo, err := util.UploadAttachment(req.FileHeader, req.Category, req.Description, req.Tags)
	if err != nil {
		l.Errorf("上传附件失败: %v", err)
		return nil, errors.Failure.WithMsg("上传附件失败")
	}

	// 4. 保存附件信息到数据库（这里需要实现数据库存储）
	// fileRecord, err := l.saveAttachmentToDB(attachmentInfo)
	// if err != nil {
	//     // 如果数据库保存失败，尝试删除已上传的文件
	//     os.Remove(attachmentInfo.StoragePath)
	//     return nil, errors.Internal.WithMsg("保存附件信息失败")
	// }

	// 5. 构建响应
	data := types.UploadResult{
		FileID:      attachmentInfo.FileID,
		Filename:    attachmentInfo.Filename,
		FileSize:    attachmentInfo.FileSize,
		FileType:    attachmentInfo.FileType,
		MD5:         attachmentInfo.MD5,
		UploadTime:  attachmentInfo.UploadTime.Format("2006-01-02 15:04:05"),
		StoragePath: attachmentInfo.StoragePath,
	}

	resp := &types.UploadAttachmentResponse{
		Code:    200,
		Message: "upload success",
		Data:    data,
	}
	l.Infof("附件上传成功: %s (大小: %d bytes)", attachmentInfo.Filename, attachmentInfo.FileSize)
	return resp, nil
}

// validateFileType 验证文件类型
func (l *UploadAttachmentLogic) validateFileType(fileHeader *multipart.FileHeader) error {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return errors.Parameter.WithMsg("文件没有扩展名")
	}

	// 定义允许的文件类型
	allowedTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".7z":   "application/x-7z-compressed",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".json": "application/json",
		".xml":  "application/xml",
		".mp3":  "audio/mpeg",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".wmv":  "video/x-ms-wmv",
	}

	// 检查是否在允许的类型中
	if mimeType, ok := allowedTypes[ext]; ok {
		// 可以进一步验证MIME类型
		l.Debugf("文件类型验证通过: %s -> %s", ext, mimeType)
		return nil
	}

	return errors.Parameter.WithMsg(fmt.Sprintf("不支持的文件类型: %s", ext))
}

// saveAttachmentToDB 保存附件信息到数据库（需要实现）
// func (l *UploadAttachmentLogic) saveAttachmentToDB(info *util.AttachmentInfo) (*models.Attachment, error) {
//     // 这里需要实现数据库存储逻辑
//     // 可以使用gorm或其他ORM框架
//     return nil, nil
// }
