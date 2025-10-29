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

type SearchPagedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页搜索主题
func NewSearchPagedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPagedLogic {
	return &SearchPagedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchPagedLogic) SearchPaged(req *types.SearchPagedReq) (resp *types.TopicPaginationSearchResult, err error) {
	resp, err = spring.GetBean[*service.UnsQueryService]().SearchPaged(l.ctx, req)
	return
}
