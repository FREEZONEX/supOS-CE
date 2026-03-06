// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAppsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Search installed applications
func NewSearchAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAppsLogic {
	return &SearchAppsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchAppsLogic) SearchApps(req *types.SearchAppsRequest) (resp *types.InstalledAppsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
