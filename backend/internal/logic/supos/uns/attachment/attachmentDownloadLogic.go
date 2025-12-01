package attachment

import (
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模板实例附件下载
func NewAttachmentDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDownloadLogic {
	return &AttachmentDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDownloadLogic) AttachmentDownload(req *types.AttachmentDownloadReq, w http.ResponseWriter, r *http.Request) error {
	repo := relationDB.NewUnsAttachmentRepo(l.ctx)

	attachment, err := repo.FindOneByFilter(l.ctx, relationDB.UnsAttachmentFilter{AttachmentPath: req.ObjectName})
	if err != nil {
		return err
	}
	tmpFile := fmt.Sprintf("/tmp/%s", req.ObjectName)
	err = l.svcCtx.OssClient.PublicBucket().GetObjectLocal(l.ctx, req.ObjectName, tmpFile)
	if err != nil {
		return err
	}
	// 打开文件
	file, err := os.Open(tmpFile)
	if err != nil {
		return errors.System.AddMsgf("文件不存在").AddDetail(err)
	}
	defer file.Close()
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return errors.System.AddMsgf("无法获取文件信息").AddDetail(err)
	}
	// 设置响应头
	// Content-Disposition 用于指定文件名
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachment.AttachmentName))
	c := mime.TypeByExtension(path.Ext(attachment.AttachmentName))
	if c != "" {
		w.Header().Set("Content-Type", c)
	} else {
		// 只需要读取前 512 个字节来判断内容类型
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return errors.System.AddMsgf("读取文件内容失败").AddDetail(err)
		}
		// 重置文件读取位置，以便后续操作
		_, err = file.Seek(0, 0)
		if err != nil {
			return errors.System.AddMsgf("重置文件指针失败").AddDetail(err)
		}
		// Content-Type 可以根据文件类型设置，这里使用通用的二进制类型
		w.Header().Set("Content-Type", http.DetectContentType(buffer))
	}

	// Content-Length 用于指定文件大小
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	_, err = io.Copy(w, file)
	if err != nil {
		return errors.System.AddMsgf("文件下载出错").AddDetail(err)
	}
	return nil
}
