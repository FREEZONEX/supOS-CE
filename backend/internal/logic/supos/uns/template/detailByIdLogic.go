// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package template

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID查询模板详情
func NewDetailByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByIdLogic {
	return &DetailByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailByIdLogic) DetailById(req *types.WithID) (resp *types.TemplateDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
