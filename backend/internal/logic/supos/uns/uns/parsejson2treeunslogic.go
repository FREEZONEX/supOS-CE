package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseJson2TreeUnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部JSON定义转树结构uns字段定义
func NewParseJson2TreeUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseJson2TreeUnsLogic {
	return &ParseJson2TreeUnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseJson2TreeUnsLogic) ParseJson2TreeUns(req *types.JsonBodyReq) (resp *types.ParseJson2TreeUnsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
