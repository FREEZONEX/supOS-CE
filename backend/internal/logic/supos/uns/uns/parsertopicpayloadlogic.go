// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParserTopicPayloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部topic payload解析
func NewParserTopicPayloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParserTopicPayloadLogic {
	return &ParserTopicPayloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParserTopicPayloadLogic) ParserTopicPayload(req *types.ParserTopicPayloadReq) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
