package attachment

import (
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"

	"gitee.com/unitedrhino/share/oss/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type AttachmentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模板实例附件删除
func NewAttachmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttachmentDeleteLogic {
	return &AttachmentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AttachmentDeleteLogic) AttachmentDelete(req *types.AttachmentDeleteReq) error {
	repo := dao.NewUnsAttachmentRepo(l.ctx)
	attachment, err := repo.FindOneByFilter(l.ctx, dao.UnsAttachmentFilter{AttachmentPath: req.ObjectName})
	if err != nil {
		return err
	}
	err = l.svcCtx.OssClient.Delete(l.ctx, attachment.AttachmentPath, common.OptionKv{})
	if err != nil {
		return err
	}
	// 删除数据库记录
	if err := repo.Delete(l.ctx, attachment.ID); err != nil {
		return err
	}

	return nil
}
