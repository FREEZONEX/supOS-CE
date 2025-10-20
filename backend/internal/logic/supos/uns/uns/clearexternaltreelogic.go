// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearExternalTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清除所有外部topic
func NewClearExternalTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearExternalTreeLogic {
	return &ClearExternalTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearExternalTreeLogic) ClearExternalTree() (resp *types.StringResult, err error) {
	// todo: add your logic here and delete this line

	return
}
