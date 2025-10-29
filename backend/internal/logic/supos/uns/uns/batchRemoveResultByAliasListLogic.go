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

type BatchRemoveResultByAliasListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据别名集合批量删除文件夹和文件
func NewBatchRemoveResultByAliasListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRemoveResultByAliasListLogic {
	return &BatchRemoveResultByAliasListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchRemoveResultByAliasListLogic) BatchRemoveResultByAliasList(req *types.BatchRemoveUnsDto) (resp *types.RemoveResult, err error) {
	resp, err = spring.GetBean[*service.UnsRemoveService]().BatchRemoveResultByAliasList(l.ctx, req)
	return
}
