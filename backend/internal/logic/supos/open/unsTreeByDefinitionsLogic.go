// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsTreeByDefinitionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 多条件分页查询树结构
func NewUnsTreeByDefinitionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsTreeByDefinitionsLogic {
	return &UnsTreeByDefinitionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsTreeByDefinitionsLogic) UnsTreeByDefinitions(req *types.UnsTreeCondition) (resp *types.TreePageResult, err error) {
	// todo: add your logic here and delete this line

	return
}
