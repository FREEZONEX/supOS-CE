// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppByNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get application details by name
func NewGetAppByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppByNameLogic {
	return &GetAppByNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppByNameLogic) GetAppByName(req *types.GetAppRequest) (resp *types.AppDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
