// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExternalTopicAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部topic转UNS
func NewExternalTopicAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExternalTopicAddLogic {
	return &ExternalTopicAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExternalTopicAddLogic) ExternalTopicAdd(req *types.CreateFileDto) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
