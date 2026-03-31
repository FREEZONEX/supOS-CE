// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/uns/uns/service"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchRemoveByAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据别名集合批量删除文件夹和文件
func NewBatchRemoveByAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRemoveByAliasLogic {
	return &BatchRemoveByAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchRemoveByAliasLogic) BatchRemoveByAlias(req *types.BatchRemoveUnsDto) (resp *types.RemoveResult, err error) {
	// 调用 UnsRemoveService.BatchRemoveResultByAliasList
	result, err := spring.GetBean[*service.UnsRemoveService]().BatchRemoveResultByAliasList(l.ctx, req)
	if err != nil {
		return &types.RemoveResult{
			BaseResult: types.BaseResult{
				Code: 500,
				Msg:  I18nUtils.GetMessageWithCtx(l.ctx, "uns.batch.delete.failed") + ": " + err.Error(),
			},
		}, nil
	}

	return result, nil
}
