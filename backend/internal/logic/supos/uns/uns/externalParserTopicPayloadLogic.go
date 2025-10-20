package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExternalParserTopicPayloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部topic payload解析
func NewExternalParserTopicPayloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExternalParserTopicPayloadLogic {
	return &ExternalParserTopicPayloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExternalParserTopicPayloadLogic) ExternalParserTopicPayload() error {
	// todo: add your logic here and delete this line

	return nil
}
