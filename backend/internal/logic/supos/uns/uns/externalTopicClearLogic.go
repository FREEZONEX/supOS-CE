package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExternalTopicClearLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清除所有外部topic
func NewExternalTopicClearLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExternalTopicClearLogic {
	return &ExternalTopicClearLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExternalTopicClearLogic) ExternalTopicClear() error {
	// todo: add your logic here and delete this line

	return nil
}
