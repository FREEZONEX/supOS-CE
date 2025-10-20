package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExternalTopic2UnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部topic转UNS
func NewExternalTopic2UnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExternalTopic2UnsLogic {
	return &ExternalTopic2UnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExternalTopic2UnsLogic) ExternalTopic2Uns() error {
	// todo: add your logic here and delete this line

	return nil
}
