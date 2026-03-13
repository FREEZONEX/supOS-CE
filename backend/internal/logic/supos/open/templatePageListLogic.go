// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TemplatePageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询模板列表
func NewTemplatePageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TemplatePageListLogic {
	return &TemplatePageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TemplatePageListLogic) TemplatePageList(req *types.TemplateQueryVo) (resp *types.TemplatePageResult, err error) {
	// todo: add your logic here and delete this line

	return
}
