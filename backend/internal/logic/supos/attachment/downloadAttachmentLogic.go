package attachment

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadAttachmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadAttachmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadAttachmentLogic {
	return &DownloadAttachmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadAttachmentLogic) DownloadAttachment(req *types.DownloadAttachmentRequest, w http.ResponseWriter) error {
	// 1. 根据fileID查找文件信息（这里需要实现数据库查询）
	// fileInfo, err := l.getAttachmentInfo(req.FileID)
	// if err != nil {
	//     return errors.NotFound.WithMsg("文件不存在")
	// }

	// 2. 构建文件路径（这里需要根据实际存储结构实现）
	// filePath := filepath.Join(l.svcCtx.Config.Attachment.StorageRoot, fileInfo.RelativePath)

	// 临时实现：从固定目录查找
	storageRoot := "/app/go-edge/attachments"
	filePath := l.findFileByID(storageRoot, req.FileID)
	if filePath == "" {
		return errors.NotFind.WithMsg("文件不存在")
	}

	// 3. 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.NotFind.WithMsg("文件不存在")
		}
		return errors.NotFind.WithMsg("检查文件失败")
	}

	// 4. 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Failure.WithMsg("打开文件失败")
	}
	defer file.Close()

	// 5. 设置HTTP头
	filename := req.Filename
	if filename == "" {
		filename = filepath.Base(filePath)
	}

	// 设置Content-Disposition头，支持中文文件名
	disposition := fmt.Sprintf("attachment; filename*=UTF-8''%s", filename)
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", req.FileID))

	// 6. 支持断点续传
	rangeHeader := w.Header().Get("Range")
	if rangeHeader != "" {
		l.serveRange(file, fileInfo.Size(), rangeHeader, w)
		return nil
	}

	// 7. 发送文件内容
	if _, err := io.Copy(w, file); err != nil {
		l.Errorf("发送文件内容失败: %v", err)
		return errors.Failure.WithMsg("发送文件失败")
	}

	// 8. 更新下载统计（这里需要实现）
	// l.updateDownloadCount(req.FileID)

	l.Infof("文件下载成功: %s (大小: %d bytes)", filename, fileInfo.Size())
	return nil
}

// findFileByID 根据文件ID查找文件路径
func (l *DownloadAttachmentLogic) findFileByID(rootDir, fileID string) string {
	// 递归查找包含fileID的文件
	var foundPath string
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 检查文件名是否包含fileID
		if strings.Contains(info.Name(), fileID) {
			foundPath = path
			return filepath.SkipAll // 找到后停止遍历
		}

		return nil
	})

	return foundPath
}

// serveRange 处理断点续传
func (l *DownloadAttachmentLogic) serveRange(file *os.File, fileSize int64, rangeHeader string, w http.ResponseWriter) {
	// 解析Range头
	// 格式: bytes=0-499, 1000-1499, 2000-
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	ranges := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(ranges, "-")
	if len(parts) != 2 {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	var start, end int64
	if parts[0] != "" {
		_, err := fmt.Sscanf(parts[0], "%d", &start)
		if err != nil {
			http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	if parts[1] != "" {
		_, err := fmt.Sscanf(parts[1], "%d", &end)
		if err != nil {
			http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
			return
		}
	} else {
		end = fileSize - 1
	}

	// 验证范围
	if start < 0 || end >= fileSize || start > end {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// 设置部分内容响应头
	contentLength := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	w.WriteHeader(http.StatusPartialContent)

	// 定位并发送指定范围的内容
	file.Seek(start, 0)
	io.CopyN(w, file, contentLength)
}

// getAttachmentInfo 获取附件信息（需要实现数据库查询）
// func (l *DownloadAttachmentLogic) getAttachmentInfo(fileID string) (*models.Attachment, error) {
//     // 这里需要实现数据库查询逻辑
//     return nil, nil
// }

// updateDownloadCount 更新下载次数（需要实现）
// func (l *DownloadAttachmentLogic) updateDownloadCount(fileID string) {
//     // 这里需要实现更新下载统计的逻辑
// }
