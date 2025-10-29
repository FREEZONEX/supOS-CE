// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索主题树，默认整个树
func NewSearchTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchTreeLogic {
	return &SearchTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchTreeLogic) SearchTree(req *types.SearchTreeReq) (resp *types.SearchTreeResp, err error) {
	resp, err = spring.GetBean[*service.UnsQueryService]().SearchTree(l.ctx, req)
	return
}
