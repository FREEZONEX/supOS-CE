package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseJson2unsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部JSON定义转uns字段定义
func NewParseJson2unsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseJson2unsLogic {
	return &ParseJson2unsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseJson2unsLogic) ParseJson2uns(req *types.JsonBodyReq) (resp *types.ParseJson2UnsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
