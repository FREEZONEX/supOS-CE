package attachment

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/common/utils/fileutil"
	"backend/internal/common/utils/idutil"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss/common"
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
	err = l.r.ParseForm()
	if err != nil {
		return nil, errors.Parameter.AddMsg("解析表单失败")
	}

	files := l.r.MultipartForm.File["files"]
	if len(files) == 0 {
		l.Error(l.r.MultipartForm.File)
		return nil, errors.Parameter.AddMsg("没有上传文件")
	}
	repo := dao.NewUnsAttachmentRepo(l.ctx)
	unsRepo := dao.NewUnsNamespaceRepo()
	db := dao.GetDb(l.ctx)

	// 检查 UNS 是否存在
	_, err = unsRepo.GetByAlias(db, req.Alias)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			logx.Errorf("找不到 UNS: %s", req.Alias)
			return nil, err
		}
		return nil, err
	}

	var boList []types.UnsAttachmentBo
	for _, f := range files {
		// 生成附件名称
		attachmentName := fmt.Sprintf("%d", idutil.NextID())
		// 这里需要从实际文件获取原始文件名和扩展名
		// 暂时使用默认值
		originalName := f.Filename
		extensionName := filepath.Ext(originalName)
		if strings.HasPrefix(extensionName, ".") {
			extensionName = extensionName[1:]
		}
		if extensionName != "" {
			attachmentName += "." + extensionName
		}
		filePath := req.Alias + "/" + attachmentName
		fi, err := f.Open()
		if err != nil {
			l.Error(err)
			continue
		}
		defer fi.Close()
		_, err = l.svcCtx.OssClient.PublicBucket().Upload(l.ctx, filePath, fi, common.OptionKv{})
		if err != nil {
			l.Error(err)
			continue
		}
		//// 写入文件（这里需要实际的文件数据）
		//// 实际应该从 multipart.FileHeader 读取
		//if err := os.WriteFile(filePath, req.Files, 0644); err != nil {
		//	return nil, logx.Errorf("保存文件失败: %w", err)
		//}

		attachmentPath := fileutil.GetRelativePath(filePath)

		// 保存到数据库
		attachment := &dao.UnsAttachment{
			ID:             idutil.NextID(),
			UnsAlias:       req.Alias,
			OriginalName:   originalName,
			AttachmentName: attachmentName,
			AttachmentPath: attachmentPath,
			ExtensionName:  strings.ToUpper(extensionName),
			CreateAt:       time.Now(),
		}

		if err := repo.Insert(l.ctx, attachment); err != nil {
			// 如果数据库保存失败，删除已保存的文件
			l.svcCtx.OssClient.PublicBucket().Delete(l.ctx, filePath, common.OptionKv{})
			return nil, err
		}

		bo := types.UnsAttachmentBo{
			Id:             attachment.ID,
			UnsAlias:       attachment.UnsAlias,
			OriginalName:   attachment.OriginalName,
			AttachmentName: attachment.AttachmentName,
			AttachmentPath: attachment.AttachmentPath,
			ExtensionName:  attachment.ExtensionName,
			CreateAt:       attachment.CreateAt.UnixMilli(),
		}
		boList = append(boList, bo)
	}

	//// 更新 UNS flags
	//flags := uns.Flags
	//if flags == nil {
	//	flags = new(int32)
	//}
	//*flags = *flags | constants.UnsFlagWithAttachment
	//uns.Flags = flags

	//if err := unsRepo.Update(db, uns); err != nil {
	//	return nil, err
	//}

	resp = &types.AttachmentUploadResp{
		List: boList,
	}
	return resp, nil
}
