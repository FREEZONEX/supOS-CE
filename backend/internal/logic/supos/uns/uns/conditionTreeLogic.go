package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConditionTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 多条件分页查询树结构
func NewConditionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConditionTreeLogic {
	return &ConditionTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConditionTreeLogic) ConditionTree() error {
	// todo: add your logic here and delete this line

	return nil
}
