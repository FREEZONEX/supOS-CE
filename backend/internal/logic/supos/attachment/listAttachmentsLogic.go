package attachment

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAttachmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAttachmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAttachmentsLogic {
	return &ListAttachmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAttachmentsLogic) ListAttachments(req *types.ListAttachmentsRequest) (*types.ListAttachmentsResponse, error) {
	// 1. 验证分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 2. 从数据库查询附件列表（这里需要实现）
	// attachments, total, err := l.getAttachmentsFromDB(req.Category, req.Tag, req.Page, req.PageSize)
	// if err != nil {
	//     return nil, errors.Internal.WithMsg("查询附件列表失败")
	// }

	// 临时实现：从文件系统获取
	storageRoot := "./data/attachments"
	attachments, total := l.getAttachmentsFromFS(storageRoot, req.Category, req.Tag, req.Page, req.PageSize)

	// 3. 转换为响应格式
	items := make([]types.AttachmentInfo, 0, len(attachments))
	for _, att := range attachments {
		items = append(items, types.AttachmentInfo{
			FileID:      att.FileID,
			Filename:    att.Filename,
			FileSize:    att.FileSize,
			FileType:    att.FileType,
			Category:    att.Category,
			Description: att.Description,
			Tags:        att.Tags,
			MD5:         att.MD5,
			Uploader:    att.Uploader,
		})
	}

	// 4. 计算总页数
	totalPages := total / req.PageSize
	if total%req.PageSize > 0 {
		totalPages++
	}

	// 5. 构建响应
	resp := &types.ListAttachmentsResponse{
		Code:    200,
		Message: "查询成功",
		Data: types.ListResult{
			Items:      items,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: totalPages,
		},
		Total: total,
	}

	l.Debugf("查询附件列表成功: 第%d页，每页%d条，共%d条", req.Page, req.PageSize, total)
	return resp, nil
}

// getAttachmentsFromFS 从文件系统获取附件列表（临时实现）
func (l *ListAttachmentsLogic) getAttachmentsFromFS(rootDir, category, tag string, page, pageSize int) ([]*types.AttachmentInfo, int) {
	var allAttachments []*types.AttachmentInfo

	// 递归遍历目录
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 跳过隐藏文件和临时文件
		filename := info.Name()
		if filename[0] == '.' || strings.HasSuffix(filename, "_thumb") {
			return nil
		}

		// 构建附件信息
		att := &types.AttachmentInfo{
			FileID:        strings.TrimSuffix(filename, filepath.Ext(filename)),
			Filename:      filename,
			FileSize:      info.Size(),
			FileType:      filepath.Ext(filename),
			Category:      "",         // 需要从数据库获取
			Description:   "",         // 需要从数据库获取
			Tags:          []string{}, // 需要从数据库获取
			MD5:           "",         // 需要计算
			UploadTime:    info.ModTime().Format("2006-01-02 15:04:05"),
			Uploader:      "", // 需要从数据库获取
			DownloadCount: 0,  // 需要从数据库获取
		}

		// 过滤条件
		if category != "" && att.Category != category {
			return nil
		}

		if tag != "" && !containsTag(att.Tags, tag) {
			return nil
		}

		allAttachments = append(allAttachments, att)
		return nil
	})

	// 按上传时间倒序排序
	sort.Slice(allAttachments, func(i, j int) bool {
		timeI, _ := time.Parse("2006-01-02 15:04:05", allAttachments[i].UploadTime)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", allAttachments[j].UploadTime)
		return timeI.After(timeJ)
	})

	// 分页
	total := len(allAttachments)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []*types.AttachmentInfo{}, total
	}

	if end > total {
		end = total
	}

	return allAttachments[start:end], total
}

// containsTag 检查是否包含指定标签
func containsTag(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

// getAttachmentsFromDB 从数据库获取附件列表（需要实现）
// func (l *ListAttachmentsLogic) getAttachmentsFromDB(category, tag string, page, pageSize int) ([]*models.Attachment, int, error) {
//     // 这里需要实现数据库查询逻辑
//     return nil, 0, nil
// }
