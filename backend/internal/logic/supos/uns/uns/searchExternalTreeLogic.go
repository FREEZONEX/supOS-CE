// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchExternalTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索外部topic主题树，默认整个树
func NewSearchExternalTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchExternalTreeLogic {
	return &SearchExternalTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchExternalTreeLogic) SearchExternalTree(req *types.SearchTreeReq) (resp *types.SearchTreeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
