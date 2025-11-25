package attachment

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAttachmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取模板实例附件列表
func NewListAttachmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAttachmentLogic {
	return &ListAttachmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAttachmentLogic) ListAttachment(req *types.ListAttachmentReq) (resp *types.ListAttachmentResp, err error) {
	repo := relationDB.NewUnsAttachmentRepo(l.ctx)

	list, err := repo.FindByUnsAlias(l.ctx, req.Alias)
	if err != nil {
		return nil, err
	}

	resp = &types.ListAttachmentResp{}
	resp.List = make([]types.UnsAttachmentBo, 0, len(list))
	for _, it := range list {
		bo := types.UnsAttachmentBo{
			Id:             it.ID,
			UnsAlias:       it.UnsAlias,
			OriginalName:   it.OriginalName,
			AttachmentName: it.AttachmentName,
			AttachmentPath: it.AttachmentPath,
			ExtensionName:  it.ExtensionName,
			CreateAt:       it.CreateAt.UnixMilli(),
		}
		resp.List = append(resp.List, bo)
	}
	return resp, nil
}
