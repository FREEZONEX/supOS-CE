// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package template

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageListUnsByTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取模板下的文件列表
func NewPageListUnsByTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageListUnsByTemplateLogic {
	return &PageListUnsByTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageListUnsByTemplateLogic) PageListUnsByTemplate(req *types.PageListUnsByTemplateReq) (resp *types.PageListUnsByTemplateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
