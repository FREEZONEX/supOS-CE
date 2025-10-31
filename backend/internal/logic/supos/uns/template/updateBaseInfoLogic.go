// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package template

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改模板基本信息
func NewUpdateBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseInfoLogic {
	return &UpdateBaseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBaseInfoLogic) UpdateBaseInfo(req *types.UpdateTemplateBaseInfoReq) (resp *types.BaseResult, err error) {
	// todo: add your logic here and delete this line

	return
}
